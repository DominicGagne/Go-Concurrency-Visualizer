package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)

	ch1 := make(chan int)
	ch2 := make(chan int)

	go produceData(ch1)
	go produceData(ch2)

	for val := range ch1 {
		fmt.Printf("got value: %d\n", val)
	}

	for val := range ch2 {
		fmt.Printf("got value: %d\n", val)
	}

	// small sleeps are required before stopping the trace
	// to ensure the output has been collected
	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}

// push data through a channel for a while and then close that channel
func produceData(ch chan int) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500-100)+100))
		ch <- i
	}
	// tell whoever is listening that no more values are coming
	close(ch)
}
