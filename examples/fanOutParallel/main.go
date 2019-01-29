package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	trace.Start(os.Stderr)

	ch1 := make(chan int)
	ch2 := make(chan int)

	wg := &sync.WaitGroup{}

	go produceData(ch1)
	go produceData(ch2)

	// add two consumers to our waitGroup, do not proceed until results come in
	wg.Add(2)

	go monitorResults(ch1, wg)
	go monitorResults(ch2, wg)

	// block until results have come in
	wg.Wait()

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

func monitorResults(ch chan int, wg *sync.WaitGroup) {
	for val := range ch {
		fmt.Printf("got value: %d\n", val)
	}
	// inform the waitGroup that we're done
	wg.Done()
}
