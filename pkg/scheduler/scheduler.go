package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"

	"os"
	"sync"
	"time"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/cronParser"
	jobs "github.com/back-end-labs/ruok/pkg/job"
	jobhandler "github.com/back-end-labs/ruok/pkg/jobHandler"
	"github.com/back-end-labs/ruok/pkg/storage"
)

type JobsList struct {
	lock *sync.Mutex
	list map[int]*jobs.Job
}

func (jl *JobsList) AvailableSpace() int {
	jl.lock.Lock()
	defer jl.lock.Unlock()
	return config.MaxJobs() - len(jl.list)
}

func NewJobList(maxJobs int) *JobsList {
	return &JobsList{
		lock: &sync.Mutex{},
		list: make(map[int]*jobs.Job, maxJobs),
	}
}

type Scheduler struct {
	l        *JobsList
	storage  storage.SchedulerStorage
	parser   cronParser.ParseFn
	notifier chan int
}

func NewScheduler(s storage.SchedulerStorage, jobList *JobsList) *Scheduler {
	return &Scheduler{l: jobList, storage: s, parser: cronParser.Parse}
}

func (sched *Scheduler) initJobList(j []*jobs.Job, notifier chan int) {
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	for _, job := range j {
		err := job.InitExpression(sched.parser)
		if err != nil {
			// TODO: log error in db and continue without this one
			log.Info().Msgf("skipping job %v because we couldn't init cron expression %q", job.Id, job.CronExpString)
			continue
		}
		job.AbortChannel = make(chan struct{})
		job.Handlers.ExecuteFn = jobhandler.HTTPExecutor
		job.Handlers.OnSuccessFn = jobhandler.OnSuccessHandler(sched.storage)
		job.Handlers.OnErrorFn = jobhandler.OnErrorHandler(sched.storage)
		sched.l.list[job.Id] = job
		go job.Schedule(notifier)
		job.Scheduled = true
	}
	config.AppStats.ClaimedJobs = len(sched.l.list)
}

func (sched *Scheduler) checkForNewJobs(notifier chan int) {
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	freeSpace := config.MaxJobs() - len(sched.l.list)
	if freeSpace == 0 {
		log.Info().Msg("There is no more space for new jobs")
		return
	}
	j := sched.storage.GetAvailableJobs(freeSpace)
	sched.initJobList(j, notifier)
}

func (sched *Scheduler) Drain() error {
	releaseList := []*jobs.Job{}
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	for _, v := range sched.l.list {
		v.Scheduled = false
		v.AbortChannel <- struct{}{}
		close(v.AbortChannel)
		releaseList = append(releaseList, v)
	}
	err := sched.storage.ReleaseAll(releaseList)
	if err != nil {
		log.Error().Err(err).Msg("could not release claimed jobs. They may still be marked as claimed by this application in the db")
		return err
	}
	return nil
}

func (sched *Scheduler) DumpToFile(w io.Writer) error {

	releaseList := []jobs.Job{}
	for _, v := range sched.l.list {
		releaseList = append(releaseList, *v)
	}

	err := json.NewEncoder(w).Encode(
		&struct {
			Jobs []jobs.Job `json:"jobs"`
		}{
			releaseList,
		},
	)
	if err != nil {
		fmt.Printf("couldnt create json from jobs. error=%q", err.Error())
	}

	return nil
}

func (sched *Scheduler) Start(signalsCh chan os.Signal) int {
	log.Info().Msg("about to get available jobs to start working :)")

	j := sched.storage.GetAvailableJobs(sched.l.AvailableSpace())
	log.Info().Msgf("got %d jobs", len(j))

	sched.notifier = make(chan int, len(sched.l.list))

	log.Info().Msg("About to init all jobs")
	sched.initJobList(j, sched.notifier)

	log.Info().Msg("About to spawn 'listen for job updates' gorutine")
	updatedJobsNotificationsch := make(chan int, 100)
	updatesListenerCtx, cancelUpdateListener := context.WithCancel(context.Background())
	sched.storage.ListenForChanges(updatedJobsNotificationsch, updatesListenerCtx)

	log.Info().Msg("starting new ticker for poller")
	pollSignal := time.NewTicker(config.PollingInterval())

	for {
		select {
		case <-pollSignal.C:
			log.Info().Msg("Tick! time for polling")
			sched.checkForNewJobs(sched.notifier)

		case doneJobId := <-sched.notifier:
			sched.reschedule(doneJobId)

		case updatedJobId := <-updatedJobsNotificationsch:
			log.Info().Msgf("received signal to re-schedule job %d", updatedJobId)
			sched.refreshJob(updatedJobId)

		case <-signalsCh:
			// TODO: if we couldn't release we should push an alert to some channel
			return sched.shutDown(pollSignal, cancelUpdateListener, signalsCh)
		}
	}
}

// 1. Closes the poll signal, triggers the cancel for the notifications listener, closes the job done notifier.
//
// 2. Sends a message to the db to unlisten.
//
// 3. Releases all the jobs to the db or write them down to a file if  the db doesn't respond.
func (sched *Scheduler) shutDown(
	pollSignal *time.Ticker,
	cancelUpdateListener context.CancelFunc,
	signalsCh chan os.Signal,
) int {

	log.Info().Msg("About to drain because we got a signal")
	pollSignal.Stop()
	sched.storage.StopListeningForChanges()
	cancelUpdateListener()
	close(sched.notifier)
	close(signalsCh)

	err := sched.Drain()
	if err != nil {

		log.Error().Err(err).Msg("there was a problem with the database, trying to dump jobs into a file")
		f, err := os.Create("./dump.json")
		if err != nil {
			log.Error().Err(err).Msg("could not create file to write jobs as json")
			return 1
		}
		err = sched.DumpToFile(f)
		if err != nil {
			log.Error().Err(err).Msg("could not write jobs into a file")
			return 1
		}
		log.Info().Msg("Dumped all jobs into a file")
		return 1

	}
	log.Info().Msg("Drain operation succeeded")
	return 0
}

func (sched *Scheduler) reschedule(doneJobId int) {
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	log.Info().Msgf("job %v done!", doneJobId)
	job, ok := sched.l.list[doneJobId]
	if !ok {
		log.Error().Msgf("jod %v marked as done but can't reschedule because it is not on our job list", doneJobId)
		return
	}
	log.Info().Msgf("rescheduling job %v", doneJobId)
	go job.Schedule(sched.notifier)
}

func (sched *Scheduler) refreshJob(jobId int) {
	// Lock for a simple check
	sched.l.lock.Lock()
	j, ok := sched.l.list[jobId]
	sched.l.lock.Unlock()
	if !ok {
		log.Error().Msgf("couldn't find and update job %d in our list", jobId)
		return
	}
	updates := sched.storage.GetJobUpdates(jobId)
	if updates == nil {
		log.Error().Msg("Received empty updates. Keeping the old job.")
		return
	}

	// lock for the update
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	j.Scheduled = false
	j.AbortChannel <- struct{}{}
	j.Endpoint = updates.Endpoint
	j.HttpMethod = updates.Endpoint
	j.MaxRetries = updates.Max_retries
	j.SuccessStatuses = updates.Success_statuses
	if j.CronExpString != updates.Cron_exp_string {
		oldExpr := j.CronExpString
		j.CronExpString = updates.Cron_exp_string
		err := j.InitExpression(sched.parser)
		if err != nil {
			log.Error().Err(err).Msgf("could not init expression for job %d due to invalid expression", jobId)
			j.CronExpString = oldExpr
			j.InitExpression(sched.parser)
		}
	}
	j.Scheduled = true
	go j.Schedule(sched.notifier)

}
