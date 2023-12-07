package job

import (
	"time"

	"github.com/back-end-labs/ruok/pkg/alerting/models"
	"github.com/back-end-labs/ruok/pkg/cronParser"
	"github.com/rs/zerolog/log"
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
	Id              int                      `json:"id"`
	CronExp         cronParser.CronExpresion `json:"-"`
	CronExpString   string                   `json:"cronexp"`
	LastExecution   time.Time                `json:"lastExecution"`
	ShouldExecuteAt time.Time                `json:"shouldExecuteAt"`
	LastResponseAt  time.Time                `json:"lastResponseAt"`
	LastMessage     string                   `json:"lastMessage"`
	LastStatusCode  int                      `json:"lastStatusCode"`
	MaxRetries      int                      `json:"maxRetries"`
	Endpoint        string                   `json:"endpoint"`
	HttpMethod      string                   `json:"httpmethod"`
	Headers         []Header                 `json:"headers"`
	SuccessStatuses []int                    `json:"successStatuses"`
	Succeeded       string                   `json:"succeeded"`
	Status          string                   `json:"status"`
	ClaimedBy       string                   `json:"claimedBy"`
	CreatedAt       int                      `json:"createdAt"`
	AlertStrategy   string                   `json:"alertStrategy"`
	AlertMethod     string                   `json:"alertMethod"`
	AlertEndpoint   string                   `json:"alertEndpoint"`
	AlertPayload    string                   `json:"alertPayload"`
	AlertHeaders    []Header                 `json:"alertHeaders"`
	TLSClientCert   string                   `json:"-"`
	Scheduled       bool                     `json:"-"`
	AbortChannel    chan struct{}            `json:"-"`

	Doer     `json:"-"`
	Handlers Handlers `json:"-"`
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
	Succeeded       string    `json:"succeeded"`
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
		log.Error().Err(err).Msgf("error while parsing expresion for job %v", j.Id)
		return err
	}
	j.CronExp = expr
	return nil
}

func (j *Job) Schedule(notifier chan int) string {
	now := time.Now()
	nextExecution := j.CronExp.Next(now)
	log.Info().Msgf("next execution of job %v will be at %q", j.Id, nextExecution.String())
	timer := time.After(nextExecution.Sub(now))
	select {
	case <-j.AbortChannel:
		return "aborted"

	case executionTime := <-timer:
		result := j.Execute()
		j.LastResponseAt = result.ResponseTime
		j.LastExecution = executionTime
		j.LastMessage = result.Message
		j.LastStatusCode = result.Status
		j.Succeeded = "error"
		if j.IsSuccess(result.Status) {
			j.Succeeded = "ok"
			j.OnSuccess()
		} else {
			j.Succeeded = "error"
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

func (j *Job) AlertingInput() models.AlertInput {
	return models.AlertInput{
		AlertStrategy:  j.AlertStrategy,
		Url:            j.AlertEndpoint,
		Method:         j.AlertMethod,
		Payload:        j.AlertPayload,
		ExpectedStatus: 200,
	}
}
