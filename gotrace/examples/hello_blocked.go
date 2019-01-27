package main

import (
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)
	// create new channel of type int
	ch := make(chan int)

	// start new anonymous goroutine
	go func() {
		// send 42 to channel
		time.Sleep(10 * time.Millisecond)
		ch <- 42
	}()
	// read from channel
	<-ch
	trace.Stop()
}
