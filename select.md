listens on more than one channel at a time.

whatever channel return first on the case statement will be read.


functions that return channels
- time.After(time.Second * 5)
- time.NewTicker(time.Second * 5).C
listening to those channels is a way to control the execution.

[context](context.md) has a Done that returns a signaling channel.