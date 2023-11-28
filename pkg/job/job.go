package job

import (
	"log"
	"time"

	"github.com/back-end-labs/ruok/pkg/cronParser"
)

type Doer interface {
	Schedule()
	Execute() ExecutionResult
	OnSuccess()
	OnError()
	OnUnidentified()
	Log()
}

func Contains(x int, arr []int) bool {
	var i int
	for i = 0; i < len(arr); i++ {
		if arr[i] == x {
			return true
		}
	}
	return false

}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Handlers struct {
	ExecuteFn   func(*Job) ExecutionResult
	OnErrorFn   func(*Job)
	OnSuccessFn func(*Job)
}
type Job struct {
	CronExp         cronParser.CronExpresion
	CronExpString   string    `json:"cronexp"`
	LastExecution   time.Time `json:"lastExecution"`
	ShouldExecuteAt time.Time `json:"shouldExecuteAt"`
	LastResponseAt  time.Time `json:"lastResponseAt"`
	LastMessage     string    `json:"lastMessage"`
	LastStatusCode  int       `json:"lastStatusCode"`
	Id              int       `json:"id"`
	MaxRetries      int       `json:"maxRetries"`
	Endpoint        string    `json:"endpoint"`
	HttpMethod      string    `json:"httpmethod"`
	Headers         []Header  `json:"headers"`
	SuccessStatuses []int     `json:"successStatuses"`
	Status          string    `json:"status"`
	ClaimedBy       string    `json:"claimedBy"`
	CreatedAt       int       `json:"createdAt"`
	TLSClientCert   string
	AbortChannel    chan struct{} `json:"-"`
	Scheduled       bool          `json:"-"`

	Handlers Handlers `json:"-"`

	Doer
}

type JobExecution struct {
	Id              int       `json:"id"`
	JobId           int       `json:"jobId"`
	CronExpString   string    `json:"cronexp"`
	LastExecution   time.Time `json:"lastExecution"`
	ShouldExecuteAt time.Time `json:"shouldExecuteAt"`
	LastResponseAt  time.Time `json:"lastResponseAt"`
	LastMessage     string    `json:"lastMessage"`
	LastStatusCode  int       `json:"lastStatusCode"`
	Endpoint        string    `json:"endpoint"`
	HttpMethod      string    `json:"httpmethod"`
	Headers         []Header  `json:"headers"`
	SuccessStatuses []int     `json:"successStatuses"`
	Status          string    `json:"status"`
	ClaimedBy       string    `json:"claimedBy"`
	CreatedAt       int       `json:"createdAt"`
	DeletedAt       int       `json:"deletedAt,omitempty"`
}

func (j *Job) IsSuccess(x int) bool {
	return Contains(x, j.SuccessStatuses)
}

func (j *Job) InitExpression(parsefn cronParser.ParseFn) error {
	expr, err := parsefn(j.CronExpString)
	if err != nil {
		log.Println("error while parsing expresion: ", err)
		return err
	}
	j.CronExp = expr
	return nil
}

func (j *Job) Schedule(notifier chan int) string {
	now := time.Now()
	nextExecution := j.CronExp.Next(now)
	log.Printf("next execution will be at: %q", nextExecution.String())
	timer := time.After(nextExecution.Sub(now))
	select {
	case <-j.AbortChannel:
		j.Scheduled = false
		return "aborted"
	case executionTime := <-timer:
		result := j.Execute()
		j.LastResponseAt = result.ResponseTime
		j.LastExecution = executionTime
		j.LastMessage = result.Message
		j.LastStatusCode = result.Status

		if j.IsSuccess(result.Status) {
			j.OnSuccess()
		} else {
			j.OnError()
		}
		notifier <- j.Id
	}
	return "re-schedule"

}

type ExecutionResult struct {
	Status         int       `json:"status"`
	Message        string    `json:"message"`
	ResponseTime   time.Time `json:"responseTime"`
	SchedulerError string    `json:"schedulerError"`
}

func (j *Job) Execute() ExecutionResult {
	return j.Handlers.ExecuteFn(j)
}

func (j *Job) OnError() {
	j.Handlers.OnErrorFn(j)
}
func (j *Job) OnSuccess() {
	j.Handlers.OnSuccessFn(j)
}
