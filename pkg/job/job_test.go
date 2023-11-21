package job

import (
	"testing"
	"time"

	"github.com/back-end-labs/ruok/pkg/cronParser"
)

type returnNow struct{}

func (rn *returnNow) Next(t time.Time) time.Time {
	_ = t
	return time.Now().Add(time.Microsecond * 100)
}

var scheduleRightNowFn cronParser.ParseFn = func(cronline string) (cronParser.CronExpresion, error) {
	_ = cronline
	return &returnNow{}, nil
}

func TestContains(t *testing.T) {
	testCases := []struct {
		array      []int
		target     int
		shouldPass bool
	}{
		{[]int{1, 2, 3, 4, 5, 6}, 1, true},
		{[]int{2, 3, 4, 5, 6, 7, 8, 9, 1}, 1, true},
		{[]int{2, 3, 4, 5, 6}, 1, false},
		{[]int{1}, 1, true},
		{[]int{}, 1, false},
		{[]int{11, 2, 3, 4, 5, 6}, 1, false},
		{[]int{1, 1, 1, 2, 3, 4, 5, 6}, 1, true},
		{[]int{11, 1, 1, 2, 3, 4, 5, 6}, 1, true},
	}
	t.Run("Test Contains function", func(t *testing.T) {
		for _, test := range testCases {
			if Contains(test.target, test.array) != test.shouldPass {
				t.Errorf("mismatch between expected value and result. shouldPass=%v target=%d array=%v",
					test.shouldPass, test.target, test.array)
			}
		}
	})
}

func TestInitExpression(t *testing.T) {
	tests := []struct {
		expression  string
		shouldError bool
	}{
		{"", true},
		{"a", true},
		{"123", true},
		{"1 2 3", true},
		{"* * * * *", false},
		{"1 2 3 4 5 6", true},
		{"* * * * * * *", false},
		{"17-43/5 * * * *", false},
		{"15-30/4,55 * * * *", false},
		// Testing a valid expression for every 5 minutes
		{"*/5 * * * *", false},
		// Testing a valid expression for midnight every day
		{"0 0 * * *", false},
		// Testing a valid expression for 3:15 AM every weekday
		{"15 3 * * MON-FRI", false},
		// Testing a valid expression for noon every other month on the 1st
		{"0 12 1 */2 *", false},
		// Testing a valid expression for every 15 minutes
		{"*/15 * * * *", false},
		// Testing an invalid expression with too many fields
		{"*/5 * * * * *", false},
		// Testing an invalid expression with too few fields
		{"30 12 * *", true},
		// Testing an invalid expression with non-numeric value
		{"abc * * * *", true},
		// Testing an invalid expression with an out-of-range value
		{"61 * * * *", true},
		// Testing an invalid expression with a range exceeding 0-59
		{"0-60/5 * * * *", true},
	}
	t.Run("Check if cron library integrates ok", func(t *testing.T) {
		for _, test := range tests {
			j := Job{}
			j.CronExpString = test.expression
			err := j.InitExpression(cronParser.Parse)
			gotError := err != nil
			if gotError != test.shouldError {
				t.Errorf("expecting different result. expression=%q shouldError=%v error=%q", test.expression, test.shouldError, err)
			}

		}
	})
}

func TestScheduleHook(t *testing.T) {
	tests := []struct {
		OKs                  []int
		status               int
		shouldTriggerError   bool
		shouldTriggerSuccess bool
	}{
		{[]int{200, 201}, 200, false, true},
		{[]int{200, 201}, 400, true, false},
		{[]int{201}, 400, true, false},
		{[]int{1, 2, 3}, 7, true, false},
		{[]int{1, 2, 3}, 7, true, false},
		{[]int{}, 7, true, false},
	}

	for _, test := range tests {
		errorTriggered, successTriggered := false, false

		executionFn := func(j *Job) ExecutionResult {
			return ExecutionResult{
				Status: test.status,
			}
		}

		onErrorFn := func(j *Job) {
			errorTriggered = true
		}
		onSuccessFn := func(j *Job) {
			successTriggered = true
		}
		ch := make(chan int)
		j := &Job{
			SuccessStatuses: test.OKs,
			Id:              1,
			Handlers: Handlers{
				OnErrorFn:   onErrorFn,
				OnSuccessFn: onSuccessFn,
				ExecuteFn:   executionFn,
			},
		}
		j.InitExpression(scheduleRightNowFn)
		go j.Schedule(ch)
		msg := <-ch
		if msg != j.Id {
			t.Errorf("expected job id %q, got %q", j.Id, msg)
		}
		if errorTriggered != test.shouldTriggerError {
			t.Errorf("expected error triggered %v, got %v. oks=%v status=%d",
				test.shouldTriggerError, errorTriggered, test.OKs, test.status)
		}
		if successTriggered != test.shouldTriggerSuccess {
			t.Errorf("expected success triggered %v, got %v.  oks=%v status=%d",
				test.shouldTriggerSuccess, successTriggered, test.OKs, test.status)
		}

	}

}
