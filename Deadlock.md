## Deadlock
when a channel/mutex is waiting but immposible to finish
egs:
- writing to an unbufferred channel twice (the sender code is waiting for the channel to be read but not allowed it to)
```go
	  ch := make(chan int)
	  ch <- 3
	  ch <- 4 // Deadlock!
	  go func() {
		for range 2{
		v := <-ch
		fmt.Println(v)
		}
	  }()
```
- locking a mutex twice