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

	// chanOne := make(chan int)
	// chanTwo := make(chan int)

	var channels []chan int

	for i := 0; i < 20; i++ {
		ch := make(chan int)
		go produceData(ch)
		channels = append(channels, ch)
	}

	aggregateData := merge(channels...)

	// go produceData(chanOne)
	// go produceData(chanTwo)

	results := make(chan int)

	go func() {
		for val := range aggregateData {
			fmt.Printf("got value: %d\n", val)
		}
		results <- 1
		time.Sleep(time.Millisecond * 50)
	}()

	<-results

	// small sleeps are required before stopping the trace
	// to ensure the output has been collected
	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}

// receive a variadic number of channels, merge all their values
// onto a single channel, and pass that channel back as read-only
func merge(channels ...chan int) <-chan int {
	// aggregate is where all output of these channels
	// gets funneled into
	outBound := make(chan int)

	wg := &sync.WaitGroup{}

	wg.Add(len(channels))

	output := func(dataSource chan int) {
		for val := range dataSource {
			outBound <- val
		}
		// as soon as the channel is closed, inform the waitgroup
		// and terminate this gouroutine
		wg.Done()
	}

	for _, channel := range channels {
		// start up a goroutine for each of our channels, monitor
		// it's output and pass it along to the aggregate outBound channel
		go output(channel)
	}

	go func() {
		// wait for all data channels to be closed and supervisor goroutines
		// to be killed, and then close the aggregate outBound channel
		wg.Wait()
		close(outBound)
	}()

	return outBound
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
