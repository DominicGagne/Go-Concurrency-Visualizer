package main

import (
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)

	table := make(chan int)

	// launch two goroutines to pass a resource back and forth
	go play(table)
	go play(table)

	// start the match by passing a resource onto the channel
	table <- 1

	// let the match go on for a while, then end it by grabbing the 'ball'
	time.Sleep(time.Second * 3)
	<-table

	// small sleeps are required before stopping the trace
	// to ensure the output has been collected
	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}

func play(table chan int) {
	for {
		ball := <-table
		// simulate doing some work on this resource
		time.Sleep(time.Millisecond * 500)

		// to enrich visualization, increment by one
		ball++
		table <- ball
	}
}
