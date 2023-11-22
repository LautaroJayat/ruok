package scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
			log.Printf("erro=%q", err.Error())
			continue
		}
		sched.l.list[job.Id] = job
		job.Handlers.ExecuteFn = jobhandler.HTTPExecutor
		job.Handlers.OnSuccessFn = jobhandler.OnSuccessHandler(sched.storage)
		job.Handlers.OnErrorFn = jobhandler.OnErrorHanler(sched.storage)
		go job.Schedule(notifier)
	}
}

func (sched *Scheduler) checkForNewJobs(notifier chan int) {
	sched.l.lock.Lock()
	defer sched.l.lock.Unlock()
	freeSpace := config.MaxJobs() - len(sched.l.list)
	if freeSpace == 0 {
		log.Println("no more free space for")
		return
	}
	j := sched.storage.GetAvailableJobs(freeSpace)
	sched.initJobList(j, notifier)
}

func (sched *Scheduler) Drain() error {
	releaseList := []*jobs.Job{}
	for _, v := range sched.l.list {
		releaseList = append(releaseList, v)
	}
	err := sched.storage.ReleaseAll(releaseList)
	if err != nil {
		fmt.Printf("could not release all jobs. err=%q.", err.Error())
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
	log.Println("about to get jobs")
	j := sched.storage.GetAvailableJobs(sched.l.AvailableSpace())
	log.Printf("we got %d jobs\n", len(j))
	notifier := make(chan int, len(sched.l.list))

	log.Println("about init jobs")
	sched.initJobList(j, notifier)

	log.Println("starting new ticker for poller")
	pollSignal := time.NewTicker(config.PollingInterval())

	for {
		select {
		case <-pollSignal.C:
			log.Println("Tick! time for polling")
			sched.checkForNewJobs(notifier)

		case doneJobId := <-notifier:
			log.Println("job done!")
			job, ok := sched.l.list[doneJobId]
			if !ok {
				fmt.Println("error=job not found")
			} else {
				log.Println("scheduling it again")
				go job.Schedule(notifier)
			}
		case <-signalsCh:
			fmt.Println("About to drain because we got a signal")
			pollSignal.Stop()
			close(notifier)
			close(signalsCh)

			err := sched.Drain()
			// if we couldn't release...
			if err != nil {
				// should push an alert to some channel
				f, err := os.Create("./dump.json")
				if err != nil {
					fmt.Printf("could not create file to write jobs as json. error=%q", err.Error())
					return 1
				}
				sched.DumpToFile(f)
				return 1
			}
			return 0
		}
	}
}
