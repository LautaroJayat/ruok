package scheduler

import (
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
	l       *JobsList
	storage storage.SchedulerStorage
	parser  cronParser.ParseFn
}

func NewScheduler(s storage.SchedulerStorage, jobList *JobsList) *Scheduler {
	return &Scheduler{l: jobList, storage: s, parser: cronParser.Parse}
}

func (sched *Scheduler) initJobList(j []*jobs.Job, notifier chan int) {
	for _, job := range j {
		err := job.InitExpression(sched.parser)
		if err != nil {
			// TODO: log error in db and continue without this one
			log.Info().Msgf("skipping job %v because we couldn't init cron expression %q", job.Id, job.CronExpString)
			continue
		}
		job.Scheduled = true
		job.AbortChannel = make(chan struct{})
		job.Handlers.ExecuteFn = jobhandler.HTTPExecutor
		job.Handlers.OnSuccessFn = jobhandler.OnSuccessHandler(sched.storage)
		job.Handlers.OnErrorFn = jobhandler.OnErrorHanler(sched.storage)
		sched.l.list[job.Id] = job
		go job.Schedule(notifier)
	}
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
	for _, v := range sched.l.list {
		v.AbortChannel <- struct{}{}
		close(v.AbortChannel)
		releaseList = append(releaseList, v)
	}
	for {
		wait := false
		for _, v := range sched.l.list {
			if v.Scheduled {
				wait = true
			}
		}
		if !wait {
			break
		}
		time.Sleep(time.Microsecond * 100)
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

	notifier := make(chan int, len(sched.l.list))

	log.Info().Msg("About to init all jobs")
	sched.initJobList(j, notifier)

	log.Info().Msg("starting new ticker for poller")
	pollSignal := time.NewTicker(config.PollingInterval())

	for {
		select {
		case <-pollSignal.C:
			log.Info().Msg("Tick! time for polling")
			sched.checkForNewJobs(notifier)
		case doneJobId := <-notifier:
			log.Info().Msgf("job %v done!", doneJobId)
			job, ok := sched.l.list[doneJobId]
			if !ok {
				log.Error().Msgf("jodb %v marked as done but can't reschedule because it is not on our job list", doneJobId)
				continue
			}
			log.Info().Msgf("rescheduling job %v", doneJobId)
			go job.Schedule(notifier)

		case <-signalsCh:
			log.Info().Msg("About to drain because we got a signal")
			pollSignal.Stop()
			close(notifier)
			close(signalsCh)

			err := sched.Drain()
			// if we couldn't release...
			if err != nil {
				// should push an alert to some channel
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
	}
}
