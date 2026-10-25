A primary way of passing data between go routines.

They are preferred over using [mutex](mutex.md) in *Don't communicate by sharing memory. **Share memory by communicating** *, which essentially mean communication as a means to an end, which is sharing memory, not the other way around. instead of 2 go routines communicating by sharing a resource (ie, counter, using mutex to lock the access), the goroutine send the result to the channel, while the other goroutine read the data via the channel. 
using channels make the role explicit, writing end or receiving end. it makes the code easier to reason. shared memory vs ownership

memorization of channel syntax is required first.

```golang
ch := make(chan int) // creating a channel

go func() {
	ch <- 10 // write to it
}

v := <-ch // read from it

// passing channels to functions
func writeToCh(ch chan<- int) 
func readFromCh(ch <-chan int)
```


> [!TIP] memonic
> to differentiate between a read/write channel variable, a **writing** channel has arrow on the **right**. `chan<- type`


> [!NOTE] 
> It's conventional to have a goroutine that write to a channel, and a different goroutine reading from it.

Incrementing a counter can be done in safe manner by using channels by counter channel (chan int), a background go function (ch chan<- int) that write to it by incrementing in a for loop, while the handler for getting the counter value read from it (ch <-chan int). it is safe since the counter is incremented only after it is read, and when it read, it is after the counter is written as *a channel read will block execution until it is written, and a channel write will block until the channel is read*


`v, ok := <-ch`
ok return if the channel is close.

Channels can also serve as a signal, not just for trasferring values, and are commonly used with [select](select.md).


you can read a channel through a for loop
```golang
for job := range jobsCh {
}
```