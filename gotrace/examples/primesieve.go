package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		time.Sleep(10 * time.Millisecond)
		ch <- i
	}
}

func Filter(ch <-chan int, out chan<- int, prime int) {
	for {
		i := <-ch
		if i%prime != 0 {
			out <- i
		}
	}
}

func main() {
	trace.Start(os.Stderr)
	ch := make(chan int)
	go Generate(ch)
	for i := 0; i < 10; i++ {
		prime := <-ch
		fmt.Println(prime)
		out := make(chan int)
		go Filter(ch, out, prime)
		ch = out
	}
	trace.Stop()
}
