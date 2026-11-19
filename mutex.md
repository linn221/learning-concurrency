to synchronize access to a resource, by locking.

mutexs are typically used as struct field, and then passing the pointer to the struct to functions, to prevent the mutex from being copied.

### RWMutex
RWMutex can be handly when you don't want 2 reads to block each other, since it should be safe to read a value by multiple go routines as long as the value has not been modified.
you can use .RLock(Unlock) for the reading purpose, which will only block while the value is being modified.
for writing, use .Lock(Unlock). write will be blocked while the value has been read.