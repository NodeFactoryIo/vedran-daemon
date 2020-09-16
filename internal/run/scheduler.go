package run

import "github.com/go-co-op/gocron"

type Scheduler interface {
	Every(interval uint64) *gocron.Scheduler
	StartBlocking()
}

type scheduler struct {
	scheduler *gocron.Scheduler
}
