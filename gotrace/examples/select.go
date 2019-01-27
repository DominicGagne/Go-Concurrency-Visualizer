package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

func main() {
	trace.Start(os.Stderr)
	ch, ch1 := make(chan int), make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()

	for i := 0; i < 20; i++ {
		select {
		case v := <-ch:
			fmt.Println("Recv ", v)
		case v := <-ch1:
			fmt.Println("Recv ", v)
		}
	}
	trace.Stop()
}
