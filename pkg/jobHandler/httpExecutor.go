package jobhandler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog/log"

	"github.com/back-end-labs/ruok/pkg/job"
)

func HTTPExecutor(j *job.Job) job.ExecutionResult {
	r, err := http.NewRequest(j.HttpMethod, j.Endpoint, nil)

	if err != nil {
		log.Error().Err(err).Msgf("could not create request for job %v", j.Id)
		return job.ExecutionResult{}
	}

	for i := 0; i < len(j.Headers); i++ {
		r.Header.Set(j.Headers[i].Name, j.Headers[i].Name)
	}

	result := job.ExecutionResult{}
	client := http.Client{}
	res, err := client.Do(r)
	result.ResponseTime = time.Now()

	if err != nil {

		log.Error().Err(err).Msgf("there was an error while sending the request for job %v", j.Id)
		result.SchedulerError = err.Error()
		return result
	}

	result.Status = res.StatusCode

	if res.Body != nil {
		body, err := io.ReadAll(res.Body)
		stringBody := string(body)

		if err != nil {
			result.SchedulerError = fmt.Sprintf("could not read body from request. error=%q\n", err)
			result.SchedulerError += "\n"
			result.SchedulerError += err.Error() + "\n"
		}

		if !utf8.ValidString(stringBody) {
			log.Info().Msgf("Converting service response to valid UTF8 for job %v", j.Id)
			stringBody = strings.ToValidUTF8(stringBody, "")
		}

		result.Message = stringBody
	}

	return result
}
