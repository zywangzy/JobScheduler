# JobScheduler
Lightweight implementation of Golang job scheduler

This library provides an lightweight implementation of job scheduler in Golang. User can use this library, in the main goroutine, to schedule one-time jobs to be executed at given time, or schedule recurrent jobs to be executed with a given interval between start times of two consecutive calls.

## Usage example
Usage of this library is as simple as:
```go
import (
    "fmt"
    "github.com/zywangzy/JobScheduler"
    "time"
)

func main() {
    oneTimeJob := func(...interface{}) {
        fmt.Println("Hello world")
    }
    recurrentJob := func(...interface{}) {
        fmt.Println("Recurring...")
    }
    jobScheduler := scheduler.NewJobScheduler()
    jobScheduler.Start()

    jobScheduler.AddJob(oneTimeJob, time.Now() + time.Second)
    jobScheduler.AddRecurrentJob(recurrentJob, time.Now() + time.Second * 2, time.Second * 2)

    time.Sleep(time.Second * 10)
    jobScheduler.Stop()
}
