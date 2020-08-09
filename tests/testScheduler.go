package main

import (
	"fmt"
	"github.com/zywangzy/JobScheduler"
	"time"
)
func main() {
	fmt.Println(time.Now(), "Started")
	job := func(params ...interface{}) {
		params = append([]interface{}{time.Now(), "Job done"}, params...)
		fmt.Println(params...)
	}
	recurJob := func(...interface{}) {
		fmt.Println(time.Now(), "Recurrent Job done")
	}
	js := scheduler.NewJobScheduler()
	js.Start()
	js.AddJob(job, time.Now().Add(time.Second * 1))
	js.AddJob(job, time.Now().Add(time.Second * 3), "3 seconds later")
	js.AddRecurrentJob(recurJob, time.Now().Add(time.Second * 3), time.Second)
	time.Sleep(time.Second * 5)
	js.Stop()
	fmt.Println(time.Now(), "All done")
}
