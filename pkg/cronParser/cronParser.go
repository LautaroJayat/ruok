package cronParser

import (
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

type CronExpresion interface {
	Next(time.Time) time.Time
}

type ParseFn func(cronLine string) (CronExpresion, error)

func Parse(cronLine string) (CronExpresion, error) {
	return cronexpr.Parse(cronLine)
}

func IsValidExpression(cronLine string) bool {
	_, err := cronexpr.Parse(cronLine)
	return err != nil
}
