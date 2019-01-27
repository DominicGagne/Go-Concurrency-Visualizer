package main

import (
	"time"
)

// This example shows blocked state and some CPU usage visualization
func main() {
	results := make(chan int)

	go func() {
		for {
			val := <-results
			time.Sleep(time.Second)
			results <- val + 1
		}
	}()

	go func() {
		for {
			val := <-results
			time.Sleep(time.Second)
			results <- val + 1
		}
	}()

	results <- 1

	time.Sleep(time.Second * 3)

	// sleep for some time before stopping trace, because
	// calling stop too soon seems to result in an incomplete visualization
	time.Sleep(time.Millisecond * 100)
}
