package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

const (
	WORKERS    = 6
	SUBWORKERS = 12
	TASKS      = 12
	SUBTASKS   = 12
)

func subworker(subtasks chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-subtasks
		if !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
		fmt.Println(task)
	}
}

func worker(tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}

		var wg2 sync.WaitGroup
		wg2.Add(SUBWORKERS)
		subtasks := make(chan int)
		for i := 0; i < SUBWORKERS; i++ {
			time.Sleep(1 * time.Millisecond)
			go subworker(subtasks, &wg2)
		}
		for i := 0; i < SUBTASKS; i++ {
			task1 := task * i
			subtasks <- task1
		}
		close(subtasks)
		wg2.Wait()
	}
}

func main() {
	trace.Start(os.Stderr)
	var wg sync.WaitGroup
	wg.Add(WORKERS)
	tasks := make(chan int)

	for i := 0; i < WORKERS; i++ {
		time.Sleep(1 * time.Millisecond)
		go worker(tasks, &wg)
	}

	for i := 0; i < TASKS; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
	trace.Stop()
}
