// Package scheduler provides class JobScheduler to schedule jobs to run at given times or
// recurrent jobs to run at given interval starting at given time.
package scheduler

import (
	"sync"
	"time"
)

// Class JobScheduler
// Usage:
// Create and start the job scheduler
// jobScheduler := NewJobScheduler()
// jobScheduler.Start()
// Then just add jobs to the scheduler with scheduled time for execution:
// job := func() { do your stuff }
// scheduledTime := time.Now() + time.Duration(10)
// jobScheduler.AddJob(job, scheduledTime)
// You can also schedule a job with fixed interval and scheduled time for first execution:
// interval := time.Minute * 10
// jobScheduler.AddRecurrentJob(job, interval)
// jobScheduler.Stop()
type JobScheduler struct {
	running  bool
	quit     chan bool
	wg       sync.WaitGroup
	pwg     *sync.WaitGroup
}

// Class Job to be executed by job scheduler. Notice that the Job instances don't return
// values explicitly.
type Job func(params ...interface{})

// Function to create new JobScheduler instance, which returns a pointer pointing to the
// newly created instance.
func NewJobScheduler() *JobScheduler {
	js := JobScheduler{
		running: false,
		quit:    make(chan bool),
	}
	js.pwg = &js.wg
	return &js
}

// Start the JobScheduler. The JobScheduler needs to be started before you add any jobs to
// it, otherwise the jobs would be ignored by JobScheduler.
func (js *JobScheduler) Start() {
	if js.running {
		return
	}
	js.running = true
}

// Add a function object as job to the job scheduler. The job would be executed exactly
// once at `startTime`. If startTime is equal or earlier than `time.Now()`, the job
// would be executed immediately after the function is called.
func (js *JobScheduler) AddJob(job Job, startTime time.Time, jobParams ...interface{}) {
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
			job(jobParams...)
		}
	}(js.quit)
}

// Add a function object as recurrent job to job scheduler. The job would be executed
// recurrently with an interval of `interval`. The first execution would happen at
// `startTime`. If startTime is equal or earlier than `time.Now()`, the first job
// execution would happen immediately and then executes with interval.
func (js *JobScheduler) AddRecurrentJob(job Job, startTime time.Time, interval time.Duration, jobParams ...interface{}) {
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
				job(jobParams...)
			default:
				time.Sleep(time.Second)
			}
		}
	}(js.quit)
}

// Stop the job scheduler instance. This is a blocking call to stop all the jobs that
// have been added to this scheduler and wait for all the goroutines to stop. It's
// required to call Stop from the goroutine where JobScheduler instance is created
// before it terminates, otherwise the program might have a panic or crash.
func (js *JobScheduler) Stop() {
	if !js.running {
		return
	}
	js.running = false
	close(js.quit)
	js.pwg.Wait()
}

