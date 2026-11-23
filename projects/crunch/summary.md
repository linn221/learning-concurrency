> crunch will generate a wordlist for bruteforcing passwords, provided the charset (ie, 0-9) and length (min, max), using recursion

**sequential** will create the output file, and if find the word, use bufio to write the words directly for faster performance, the recursive function will re-use a []byte with maximum length in each call for efficiency. it runs about **1.6** second to generate 0-9 digits of length 8, close to the actual crunch program performance. in extreme case, it will generate a length of 5 for character sets of (0-9,a-z,A-Z), and took **14** seconds.
concurrent programs will employ multiple go routines, each one already generating first character, in case of 8 length digit wordlist, we will have 10 go routines, first one will start with 00000000, then 00000001, second one will start with 100000000, 10000001, etc. *each go routine does not depend on another go routine for the word generation.* however, writing to the output file have to be synchornized, since bufio's Write is not concurrent safe.

**channels version** will use wordsCh (chan []byte), for each word, recursive function will write to the channel, main go routine will read it and write to the file. it has the overhead of sending and receiving from a channel, which slows down the program. it is the slowest.

**mutex version** will use a mutex for exclusive access to file's write. it is faster than channels, but still slower than the sequential program, taking **2 minute** in extreme case.

**noBlocking** is the fastest, taking only **8.5** second on extreme case, although it is slower than **sequential** on small character set. it remove the need for synchronization, each go routine will write to its own separate file therefore no need of an exclusive lock. after all the go routines have finished, main will merge those files into a single output file. it introuce an extra workload (which is why it is slower on small data set) but runs faster when generating big set.
