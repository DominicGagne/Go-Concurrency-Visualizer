package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func worker(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-ch
		if !ok {
			return
		}
		time.Sleep(1 * time.Millisecond)
		fmt.Println("processing task", task)
	}
}

func pool(wg *sync.WaitGroup, workers, tasks int) {
	ch := make(chan int)

	for i := 0; i < workers; i++ {
		time.Sleep(1 * time.Millisecond)
		go worker(ch, wg)
	}

	for i := 0; i < tasks; i++ {
		time.Sleep(10 * time.Millisecond)
		ch <- i
	}

	close(ch)
}

func main() {
	trace.Start(os.Stderr)
	var wg sync.WaitGroup
	wg.Add(36)
	go pool(&wg, 36, 36)
	wg.Wait()
	trace.Stop()
}
