package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)

	c := gen(2, 3, 4, 5)
	outBound := square(square(square(square(square(c)))))

	for val := range outBound {
		fmt.Printf("received val: %d\n", val)
	}

	// small sleeps are required before stopping the trace
	// to ensure the output has been collected
	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}

// copied from https://blog.golang.org/pipelines
// converts variadic or slice of ints to distinct values in a pipeline
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			time.Sleep(time.Millisecond * 100)
			out <- n
		}
		// small sleep before terminating the goroutine to ensure
		// trace output is collected
		time.Sleep(time.Millisecond * 50)
		// inform the pipeline that's all we've got
		close(out)
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			time.Sleep(time.Millisecond * 100)
			out <- n * n
		}
		// small sleep before terminating the goroutine to ensure
		// trace output is collected
		time.Sleep(time.Millisecond * 50)
		close(out)
	}()
	return out
}
