package scheduler

import "github.com/go-co-op/gocron"

type Scheduler interface {
	Every(interval uint64) *gocron.Scheduler
	Do(jobFun interface{}, params ...interface{}) (*gocron.Job, error)
	StartBlocking()
}
