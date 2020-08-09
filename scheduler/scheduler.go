package scheduler

import (
	"sync"
	"time"
)

/* Usage:
 * var jobScheduler JobScheduler
 * jobScheduler.start()
 * var job func()
 * scheduledTime := time.Now() + time.Duration(10)
 * jobScheduler.addJob(job, scheduledTime)
 * interval := time.Duration(10) * time.Minute
 * jobScheduler.addRecurrentJob(job, interval)
 * jobScheduler.stop()
 */

type JobScheduler struct {
	running bool
	quit    chan bool
	wg      sync.WaitGroup
	pwg     *sync.WaitGroup
}

type Job func()

func NewJobScheduler() *JobScheduler {
	js := JobScheduler{
		running: false,
		quit:    make(chan bool),
	}
	js.pwg = &js.wg
	return &js
}

func (js *JobScheduler) Start() {
	if js.running {
		return
	}
	js.running = true
}

func (js *JobScheduler) AddJob(job Job, startTime time.Time) {
	if !js.running {
		return
	}
	go func(quit chan bool) {
		defer js.pwg.Done()
		js.pwg.Add(1)
		select {
		case <- quit:
			return
		case <- time.After(startTime.Sub(time.Now())):
			job()
		}
	}(js.quit)
}

func (js *JobScheduler) AddRecurrentJob(job Job, startTime time.Time, interval time.Duration) {
	if !js.running {
		return
	}
	go func(quit chan bool) {
		defer js.pwg.Done()
		js.pwg.Add(1)
		select {
		case <- quit:
			return
		case <- time.After(startTime.Sub(time.Now())):
		}
		ticker := time.NewTicker(interval)
		job()
		for {
			select {
			case <-quit:
				ticker.Stop()
				return
			case <-ticker.C:
				job()
			default:
				time.Sleep(time.Second)
			}
		}
	}(js.quit)
}

func (js *JobScheduler) Stop() {
	if !js.running {
		return
	}
	js.running = false
	close(js.quit)
	js.pwg.Wait()
}



