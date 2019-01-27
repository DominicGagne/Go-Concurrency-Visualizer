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
	ch := make(chan int)
	ch2 := make(chan int)

	aggregateResults := merge(ch, ch2)

	go work(ch, 1)
	go work(ch2, 2)

	for val := range aggregateResults {
		fmt.Printf("received a value: %d\n", val)
	}

	time.Sleep(time.Second)

	trace.Stop()
}

func work(ch chan int, i int) {
	max := 8
	min := 1
	time.Sleep(time.Second * time.Duration(rand.Intn(max-min)+min))
	ch <- i
	time.Sleep(time.Millisecond * 100)
	close(ch)
}

func merge(channels ...chan int) chan int {
	outBound := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	sendValuesToOutbound := func(channel chan int) {
		for value := range channel {
			outBound <- value
		}
		// once this channel has closed, inform the waitgroup
		wg.Done()
	}

	for _, channel := range channels {
		go sendValuesToOutbound(channel)
	}

	go func() {
		// once all channels have closed, close our outBound channel
		wg.Wait()
		close(outBound)
	}()

	return outBound
}
