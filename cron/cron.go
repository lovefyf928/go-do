package cron

import (
	"github.com/robfig/cron"
	"reflect"
)

var cj *cron.Cron

var Schedule *Scheduled

func init() {
	cj = cron.New()
	cj.Start()
	Schedule = &Scheduled{}
}

type Scheduled struct{}

func (*Scheduled) AddJob(job cron.Job) {
	s := reflect.TypeOf(job).Elem()
	if s.NumField() > 1 {

	}
	cronStr := s.Field(0).Tag.Get("corn")
	err := cj.AddJob(cronStr, job)
	if err != nil {
		panic(err)
	}
}

func NewSchedule(cronStr string, fn func()) {
	err := cj.AddFunc(cronStr, fn)
	if err != nil {
		panic(err)
	}
}
