package main

import (
	"fmt"
	"github.com/zywangzy/JobScheduler/scheduler"
	"time"
)
func main() {
	fmt.Println(time.Now(), "Started")
	job := func() {
		fmt.Println(time.Now(), "Job done")
	}
	rjob := func() {
		fmt.Println(time.Now(), "Recurrent Job done")
	}
	js := scheduler.NewJobScheduler()
	js.Start()
	js.AddJob(job, time.Now().Add(time.Second * 1))
	js.AddJob(job, time.Now().Add(time.Second * 3))
	js.AddRecurrentJob(rjob, time.Now().Add(time.Second * 3), time.Second)
	time.Sleep(time.Second * 5)
	js.Stop()
	fmt.Println(time.Now(), "All done")
}
