package scheduler

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestJobsList_AvailableSpace(t *testing.T) {
	jl := NewJobList(config.MaxJobs())

	// Simulate adding 3 jobs to the list
	job1 := &job.Job{}
	job2 := &job.Job{}
	job3 := &job.Job{}
	jl.list[1] = job1
	jl.list[2] = job2
	jl.list[3] = job3

	expectedSpace := config.MaxJobs() - len(jl.list)
	if space := jl.AvailableSpace(); space != expectedSpace {
		t.Errorf("Expected available space: %d, got: %d", expectedSpace, space)
	}
}

func TestScheduler_DumpToFile(t *testing.T) {

	sched := &Scheduler{
		l: NewJobList(config.MaxJobs()),
	}
	sched.l.list[1] = &job.Job{Id: 1, Handlers: job.Handlers{}}

	tempFile, err := os.CreateTemp("", "dump_test_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	err = sched.DumpToFile(tempFile)
	if err != nil {
		t.Fatalf("Error dumping to file: %v", err)
	}

	dumpedData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Error reading dumped file: %v", err)
	}

	// Unmarshal JSON and check if it contains the expected job
	var dumped struct {
		Jobs []job.Job `json:"jobs"`
	}
	err = json.Unmarshal(dumpedData, &dumped)
	if err != nil {
		t.Fatalf("Error unmarshaling dumped data: %v", err)
	}

	// Verify that the dumped data contains the expected job
	if len(dumped.Jobs) != 1 || dumped.Jobs[0].Id != 1 {
		t.Errorf("Dumped data does not contain the expected job")
	}
}

// Mock implementation of the Storage interface for testing
type mockStorage struct{}

var mockedJobList *JobsList

func (ms *mockStorage) GetAvailableJobs(space int) []*job.Job {
	// Mock implementation, returns an empty list for simplicity
	return []*job.Job{
		{Id: 1, CronExpString: "10 * * * *"},
		{Id: 2, CronExpString: "10 * * * *"},
		{Id: 3, CronExpString: "10 * * * *"},
		{Id: 4, CronExpString: "10 * * * *"},
		{Id: 5, CronExpString: "10 * * * *"},
		{Id: 6, CronExpString: "10 * * * *"},
		{Id: 7, CronExpString: "10 * * * *"},
		{Id: 8, CronExpString: "10 * * * *"},
		{Id: 9, CronExpString: "10 * * * *"},
		{Id: 10, CronExpString: "10 * * * *"},
	}
}

var releasedJobs []*job.Job

func (ms *mockStorage) ReleaseAll(jobs []*job.Job) error {
	releasedJobs = jobs
	return nil
}

func (ms *mockStorage) GetClient() *pgxpool.Pool {
	// Mock implementation, does nothing for simplicity
	return nil
}

func (ms *mockStorage) RegisterSelf() {
	// Mock implementation, does nothing for simplicity

}

func (ms *mockStorage) WriteDone(j *job.Job) error {
	return nil
}

func (ms *mockStorage) GetClaimedJobs(limit int, offset int) []*job.Job {
	return nil
}

func (ms *mockStorage) GetClaimedJobsExecutions(jobId int, limit int, offset int) []*job.JobExecution {
	return nil
}

func TestScheduler_Start(t *testing.T) {
	releasedJobs = []*job.Job{}
	mockedJobList = NewJobList(config.MaxJobs())
	sched := NewScheduler(&mockStorage{}, mockedJobList)

	exitCodeCh := make(chan int, 1)
	signalCh := make(chan os.Signal, 1)
	go func() {
		exitCode := sched.Start(signalCh)
		exitCodeCh <- exitCode
	}()

	for _, j := range sched.l.list {
		t.Log(j.Id, "scheduled?", j.Scheduled)
		assert.True(t, j.Scheduled)
	}
	signalCh <- os.Interrupt
	exitCode := <-exitCodeCh

	if exitCode != 0 {
		t.Errorf("expected exit code 0 in a normal flow. Instead got %d", exitCode)
	}

	for i := 1; i <= 10; i++ {
		j, ok := mockedJobList.list[i]
		assert.True(t, ok)
		assert.False(t, j.Scheduled)
		releasedJobsId := []int{}
		for _, rj := range releasedJobs {
			releasedJobsId = append(releasedJobsId, rj.Id)
		}
		assert.Contains(t, releasedJobsId, j.Id)
	}

}
