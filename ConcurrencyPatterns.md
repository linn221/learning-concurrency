## Generator pattern
![](attachments/Pasted%20image%2020261112165143.png)

a function that returns a channel you can read from, abstracting the creation of channel under it.

## FanIn
collect multiple channels to a single channel
a func that takes multiple channels, and return one read channel, constructing an inner channel and writing to it inside the func.

![](attachments/Pasted%20image%2020261112165919.png)
![](attachments/Pasted%20image%2020261112172135.png)

## FanOut
from one channel to multiple channels