
> [!NOTE]
> particularly useful for synchronizing of go routines, when you don't know the number of them.
> if you do, you could use a done channel, write to it in each of go routines' exit, then, reading from it the number of go routines in main function.

`wg := new(sync.WaitGroup)` is preferred and conventionally wait groups are passed to funcs by pointers.
