package max_jobs_all_http_alert

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/back-end-labs/ruok/e2e"
	"github.com/back-end-labs/ruok/pkg/config"
	"github.com/back-end-labs/ruok/pkg/storage"
	"github.com/gin-gonic/gin"
)

type mapping struct {
	received     map[string]int
	alerts       map[string]int
	receivedLock sync.Mutex
	alertsLock   sync.Mutex
}

func (m *mapping) addReceived(key string) {
	m.receivedLock.Lock()
	defer m.receivedLock.Unlock()
	v, ok := m.received[key]
	if !ok {
		m.received[key] = 1
		return
	}
	m.received[key] = v + 1
}
func (m *mapping) addAlert(key string) {
	m.alertsLock.Lock()
	defer m.alertsLock.Unlock()
	v, ok := m.alerts[key]
	if !ok {
		m.alerts[key] = 1
		return
	}
	m.alerts[key] = v + 1
}

func makeTestHandler(results *mapping) gin.HandlerFunc {
	return func(c *gin.Context) {
		results.addReceived(fmt.Sprintf("http://%s%s", c.Request.Host, c.Request.URL.String()))
		c.JSON(http.StatusNotFound, gin.H{})
	}
}

func makeAlertHandler(results *mapping) gin.HandlerFunc {
	return func(c *gin.Context) {
		results.addAlert(fmt.Sprintf("http://%s%s", c.Request.Host, c.Request.URL.String()))
		c.JSON(http.StatusOK, gin.H{})
	}
}

func setupTestServer(results *mapping) *httptest.Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/test", makeTestHandler(results))
	router.POST("/alert", makeAlertHandler(results))

	thirdPartyServer := httptest.NewServer(router)
	return thirdPartyServer
}

func createJobInput(id int, jobName int, host string) storage.CreateJobInput {
	return storage.CreateJobInput{
		Name:            fmt.Sprintf("job %d", id),
		CronExpString:   "*/5 * * * * * *",
		MaxRetries:      1,
		Endpoint:        e2e.MakeTestURL(host, jobName, id),
		HttpMethod:      "GET",
		SuccessStatuses: []int{200},
		AlertStrategy:   config.ALERT_HTTP,
		AlertMethod:     "POST",
		AlertEndpoint:   e2e.MakeAlertUrl(host, jobName, id),
		AlertHeaders:    map[string]string{"Authorization": "bearer jwt"},
		AlertPayload:    "",
	}
}

/*
It will mount a testing server that echoes 404 for each job request.

Then it will try to create users*jobsPeruser number of jobs.

# It will register all executions and alerts from the scheduler

Then it will compare results the ones in the DB
*/
func TestCreateMaxJobs_all_http_alerts(t *testing.T) {
	storage.Drop()
	defer storage.Drop()
	users := 100
	jobsPerUser := 101
	results := mapping{
		alerts:       map[string]int{},
		received:     map[string]int{},
		alertsLock:   sync.Mutex{},
		receivedLock: sync.Mutex{},
	}
	jobList := []storage.CreateJobInput{}

	toBeMonitored := setupTestServer(&results)
	defer toBeMonitored.CloseClientConnections()
	defer toBeMonitored.Close()

	for i := 0; i < users; i++ {
		for j := 0; j < jobsPerUser; j++ {
			jobList = append(jobList, createJobInput(j, i, toBeMonitored.URL))
		}
	}
	_, currentFile, _, _ := runtime.Caller(0)
	base := path.Dir(currentFile)
	base = path.Join(base, "..", "..")
	base = path.Clean(base)
	cmd := exec.Command(base+"/ruok", "start")

	cmd.Env = []string{fmt.Sprintf("%s=%d\n", config.POLL_INTERVAL_SECONDS, 5)}

	f, err := os.Create("./e2etest.log")

	if err != nil {
		t.Errorf("cant create temp file for logs, %q\n", err.Error())
	} else {
		defer f.Close()
		cmd.Stderr = f
		cmd.Stdout = f
	}

	ready := make(chan struct{})
	done := make(chan struct{})
	go func() {
		err := cmd.Start()
		if err != nil {
			log.Fatalf("aborting, couldn't start scheduler: %q\n", err.Error())
		}
		ready <- struct{}{}
		err = cmd.Wait()
		if err != nil {
			log.Fatalf("there was an error while shutting down the scheduler: %q\n", err.Error())
		}
		done <- struct{}{}
	}()

	<-ready
	tryAgain := 5
	dbConnected := e2e.ServerUp(t)

	for !dbConnected {
		tryAgain--
		if tryAgain < 0 {
			t.Fatal("didn't receive dbUp from scheduler")
		}
		time.Sleep(time.Millisecond * 100)
		dbConnected = e2e.ServerUp(t)
	}

	for i, input := range jobList {
		if ok := e2e.CreateJob(t, input); !ok {
			t.Errorf("couldn't create job number %d", i)
			cmd.Process.Signal(os.Interrupt)
			<-done
			t.Fatal("aborting")

		}
	}
	time.Sleep(time.Second * 30)
	jobsCount, err := e2e.ClaimedJobs(t)
	if err != nil {
		t.Errorf("error while querying claimed jobs: %q", err.Error())
		cmd.Process.Signal(os.Interrupt)
		<-done
		t.Fatal("aborting")
	}
	if jobsCount != users*jobsPerUser && jobsCount != config.MaxJobs() {
		t.Errorf("expected %d or %d jobs but got %d", users*jobsPerUser, config.MaxJobs(), jobsCount)
		cmd.Process.Signal(os.Interrupt)
		<-done
		t.Fatal("aborting")
	}
	time.Sleep(time.Second * 15)
	cmd.Process.Signal(os.Interrupt)
	<-done
	cfg := config.FromEnvs()

	s, closeDb := storage.NewStorage(&cfg)
	defer closeDb()

	rows, err := s.GetClient().Query(
		context.Background(),
		`	select alert_endpoint, count(r.id) num 
			from ruok.job_results as r
			join ruok.jobs as j 
			on r.job_id = j.id
			where r.succeeded = 'error' 
			group by j.alert_endpoint`,
	)
	if err != nil {
		t.Errorf("expected nil error to get job executions: %q", err.Error())
		return
	}

	gottenAcc := 0
	dbAcc := 0
	for rows.Next() {
		var endpoint string
		var count int
		err := rows.Scan(&endpoint, &count)
		if err != nil {
			t.Errorf("expected no error while scanning row: %q", err.Error())

		}
		gotten, ok := results.alerts[endpoint]
		if !ok {
			t.Errorf("expected %q to exist in collected results", endpoint)
			continue
		}
		gottenAcc += gotten
		dbAcc += count
		if gotten != count {
			t.Errorf("expected collected count to be equal to value from db. colleted=%d in db=%d", gotten, count)
		}
	}
	if gottenAcc != dbAcc {
		t.Errorf("expected total collected to be equal to total executions in db. colleted=%d in db=%d", gottenAcc, dbAcc)

	}

	row := s.GetClient().QueryRow(
		context.Background(),
		"select count(id) num from ruok.jobs where status != 'pending to be claimed'",
	)
	var otherThanClaimedJobs int
	err = row.Scan(&otherThanClaimedJobs)
	if err != nil {
		t.Errorf("expected non ni error while checking how many rows kept claimed: %q", err.Error())
	}
	if otherThanClaimedJobs != 0 {
		t.Errorf("expected 0 rows with a status different to 'pending to be claimed'. got=%d", otherThanClaimedJobs)

	}
}
