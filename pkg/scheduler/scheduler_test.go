package scheduler

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/back-end-labs/ruok/pkg/alerting"
	"github.com/back-end-labs/ruok/pkg/alerting/models"
	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/job"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestJobsList_AvailableSpace(t *testing.T) {

	jl := NewJobList(config.MaxJobs())
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	id3, _ := uuid.NewV7()

	// Simulate adding 3 jobs to the list
	job1 := &job.Job{}
	job2 := &job.Job{}
	job3 := &job.Job{}
	jl.list[id1] = job1
	jl.list[id2] = job2
	jl.list[id3] = job3

	expectedSpace := config.MaxJobs() - len(jl.list)
	if space := jl.AvailableSpace(); space != expectedSpace {
		t.Errorf("Expected available space: %d, got: %d", expectedSpace, space)
	}
}

func TestScheduler_DumpToFile(t *testing.T) {
	id, _ := uuid.NewV7()

	sched := &Scheduler{
		l: NewJobList(config.MaxJobs()),
	}
	sched.l.list[id] = &job.Job{Id: id, Handlers: job.Handlers{}}

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
	if len(dumped.Jobs) != 1 || dumped.Jobs[0].Id.String() != id.String() {
		t.Errorf("Dumped data does not contain the expected job")
	}
}

type mockStorage struct {
	JobUpdatesCh chan uuid.UUID
}

func NewMockStorage() *mockStorage {
	return &mockStorage{
		JobUpdatesCh: make(chan uuid.UUID, 1),
	}
}

var mockedJobList *JobsList
var gotAvailableJobs = false

var id1, _ = uuid.NewV7()
var id2, _ = uuid.NewV7()
var id3, _ = uuid.NewV7()
var id4, _ = uuid.NewV7()
var id5, _ = uuid.NewV7()
var id6, _ = uuid.NewV7()
var id7, _ = uuid.NewV7()
var id8, _ = uuid.NewV7()
var id9, _ = uuid.NewV7()
var id10, _ = uuid.NewV7()

func (ms *mockStorage) GetAvailableJobs(space int) []*job.Job {

	gotAvailableJobs = true
	return []*job.Job{
		{Id: id1, CronExpString: "10 * * * *"},
		{Id: id2, CronExpString: "10 * * * *"},
		{Id: id3, CronExpString: "10 * * * *"},
		{Id: id4, CronExpString: "10 * * * *"},
		{Id: id5, CronExpString: "10 * * * *"},
		{Id: id6, CronExpString: "10 * * * *"},
		{Id: id7, CronExpString: "10 * * * *"},
		{Id: id8, CronExpString: "10 * * * *"},
		{Id: id9, CronExpString: "10 * * * *"},
		{Id: id10, CronExpString: "10 * * * *"},
	}
}

var releasedJobs []*job.Job

func (ms *mockStorage) ReleaseAll(jobs []*job.Job) error {
	releasedJobs = jobs
	return nil
}

func (ms *mockStorage) GetClient() *pgxpool.Pool {
	return nil
}

func (ms *mockStorage) RegisterSelf() {

}

func (ms *mockStorage) WriteDone(j *job.Job) error {
	return nil
}

func (ms *mockStorage) GetClaimedJobs(limit int, offset int) []*job.Job {
	return nil
}

func (ms *mockStorage) GetClaimedJobsExecutions(jobId uuid.UUID, limit int, offset int) []*job.JobExecution {
	return nil
}
func (ms *mockStorage) ListenForChanges(jobIDUpdatedCh chan uuid.UUID, ctx context.Context) {
	// Simulate sending updates to the provided channel
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("done listening for notifications")
				return
			case jobID := <-ms.JobUpdatesCh:
				jobIDUpdatedCh <- jobID
			}
		}
	}()
}

func (ms *mockStorage) GetJobUpdates(jobId uuid.UUID) *storage.JobUpdates {
	return &storage.JobUpdates{
		Cron_exp_string:  "*/5 * * * *",
		Endpoint:         "/updated",
		Httpmethod:       "PUT",
		Max_retries:      2,
		Headers_string:   `{"Content-Type": "application/json"}`,
		Success_statuses: []int{200, 201},
		Tls_client_cert:  "updated_cert",
		Updated_at:       time.Now().Unix(),
	}
}

func (ms *mockStorage) StopListeningForChanges() error {
	return nil
}

func TestScheduler_Start_HappyPath(t *testing.T) {
	dummyfn := func(i models.AlertInput) (string, error) {
		_ = i
		return "", nil

	}

	mockAlertingManager := alerting.CreateAlertManager(
		[]string{"http"},
		models.PluginList{
			func() (string, models.AlertFunc) {
				return "http", dummyfn
			},
		},
	)

	releasedJobs = []*job.Job{}
	mockedJobList = NewJobList(config.MaxJobs())
	mockedStorage := &mockStorage{
		JobUpdatesCh: make(chan uuid.UUID, 1),
	}

	sched := NewScheduler(mockedStorage, mockAlertingManager, mockedJobList)

	exitCodeCh := make(chan int, 1)
	signalCh := make(chan os.Signal, 1)

	go func() {
		exitCode := sched.Start(signalCh)
		exitCodeCh <- exitCode
	}()

	for !gotAvailableJobs {
		time.Sleep(time.Millisecond * 10)
	}
	sched.l.lock.Lock()
	for _, j := range sched.l.list {
		assert.True(t, j.Scheduled)
	}

	j := sched.l.list[id1]
	oldJob := *j
	sched.l.lock.Unlock()

	updatedJobID := id1
	mockedStorage.JobUpdatesCh <- updatedJobID
	time.Sleep(10 * time.Millisecond)

	sched.l.lock.Lock()
	j, ok := sched.l.list[updatedJobID]
	updatedJob := *j
	sched.l.lock.Unlock()

	assert.True(t, ok, "Refreshed job should exist in the mocked list")
	assert.True(t, updatedJob.Scheduled, "Refreshed job should not be scheduled after stopping")
	assert.NotEqual(t, oldJob.CronExpString, updatedJob.CronExpString, "Unexpected CronExpString")
	assert.NotEqual(t, oldJob.Endpoint, updatedJob.Endpoint, "Unexpected Endpoint")
	assert.NotEqual(t, oldJob.HttpMethod, updatedJob.HttpMethod, "Unexpected HttpMethod")
	assert.NotEqual(t, oldJob.MaxRetries, updatedJob.MaxRetries, "Unexpected MaxRetries")
	assert.NotEqual(t, oldJob.SuccessStatuses, updatedJob.SuccessStatuses, "Unexpected SuccessStatuses")

	signalCh <- os.Interrupt
	exitCode := <-exitCodeCh

	if exitCode != 0 {
		t.Errorf("expected exit code 0 in a normal flow. Instead got %d", exitCode)
	}

	for _, id := range []uuid.UUID{
		id1,
		id2,
		id3,
		id4,
		id5,
		id6,
		id7,
		id8,
		id9,
		id10} {
		j, ok := mockedJobList.list[id]
		assert.True(t, ok)
		assert.False(t, j.Scheduled)
		releasedJobsId := []uuid.UUID{}
		for _, rj := range releasedJobs {
			releasedJobsId = append(releasedJobsId, rj.Id)
		}
		assert.Contains(t, releasedJobsId, j.Id)
	}

}
