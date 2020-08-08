package main

import (
	"fmt"
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
	quit  chan bool
	wg    sync.WaitGroup
	pwg  *sync.WaitGroup
}

type Job func()

func NewJobScheduler() *JobScheduler {
	js := JobScheduler{
		quit: make(chan bool),
	}
	js.pwg = &js.wg
	return &js
}

func (js *JobScheduler) addJob(job Job, startTime time.Time) {
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

func (js *JobScheduler) addRecurrentJob(job Job, startTime time.Time, interval time.Duration) {
	go func(quit chan bool) {
		defer js.pwg.Done()
		js.pwg.Add(1)
		select {
		case <- quit:
			return
		case <- time.After(startTime.Sub(time.Now())):
			fmt.Println(time.Now(), "Start recurrent job cycle")
		}
		ticker := time.NewTicker(interval)
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
	close(js.quit)
	js.pwg.Wait()
}

func main() {
	fmt.Println(time.Now(), "Started")
	job := func() {
		fmt.Println(time.Now(), "Job done")
	}
	rjob := func() {
		fmt.Println(time.Now(), "Recurrent Job done")
	}
	js := NewJobScheduler()
	js.addJob(job, time.Now().Add(time.Second * 1))
	js.addJob(job, time.Now().Add(time.Second * 3))
	js.addRecurrentJob(rjob, time.Now().Add(time.Second), time.Second)
	time.Sleep(time.Second * 10)
	js.Stop()
	fmt.Println(time.Now(), "All done")
}

