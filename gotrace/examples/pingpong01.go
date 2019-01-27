package main

import (
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)
	var Ball int
	table := make(chan int)

	go player(table)
	go player(table)

	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
	trace.Stop()
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
