## Workerpool

```go
// chatGPT supplied code
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int)
    results := make(chan int)

    // Start workers
    for i := 0; i < 3; i++ {
        go worker(jobs, results)
    }

    // Send jobs
    go func() {
        for i := 0; i < 10; i++ {
            jobs <- i
        }
        close(jobs)
    }()

    // Collect results
    for i := 0; i < 10; i++ {
        fmt.Println(<-results)
    }
}
```


A **worker** is simply a function that do stuffs with the channels (mostly a read and write one).
when the function is runned via a go routine multiple times via a for loop, it becomes a **worker pool**, the channel will be read and written by multiple go routines at the same time.
instead of spawing a goroutine for every job, we divide the job into a limited number of goroutine(worker) and let them have it.

*each worker must know when to exit, after the read channel is closed. it is conventional to use a for range for reading from the channel*