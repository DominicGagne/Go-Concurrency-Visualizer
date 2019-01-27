package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func producer(ch chan int, d time.Duration) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(d)
	}
}

func reader(out chan int) {
	for i := 0; i < 20; i++ {
		x := <-out
		fmt.Println(x)
	}
}

func main() {
	trace.Start(os.Stderr)
	ch := make(chan int)
	out := make(chan int)
	go producer(ch, 10*time.Millisecond)
	go producer(ch, 25*time.Millisecond)
	go reader(out)
	for i := 0; i < 20; i++ {
		i := <-ch
		out <- i
	}
	trace.Stop()
}
