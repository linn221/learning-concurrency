to limit the active number of go routines / work in progress.

by creating a buffered channel of length N, and writing to it before doing the work, the write will block if the buffer is full (which means, the processing go routines have exceeded the limit), and once we pass the write, we immediately read the ch in defer.
sending to the channel check the capacity, and block till there is available space.

we don't limit the number of go routines, we are setting the maximum active works to be done at once.


> [!NOTE] 
> A counting semiphore of 1 is the same concept as [mutex](mutex.md)
