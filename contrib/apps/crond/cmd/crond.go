package main

import (
	"github.com/oceanho/gw/logger"
	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	var job Job
	c := cron.New(
		cron.WithSeconds())
	spec := "* * * * * ?"
	c.Start()
	defer c.Stop()
	_, e := c.AddJob(spec, job)
	if e != nil {
		panic(e)
	}
	var ch = make(chan bool, 1)
	_ = <-ch
}

type Job struct {
}

func (j Job) Run() {
	//
	logger.Info("now, %v", time.Now())
}
