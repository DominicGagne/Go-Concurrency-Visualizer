package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

type notifier struct {
	events            chan int
	subscribers       []chan int
	subscribeInternal chan chan int
}

func NewNotifier() *notifier {
	return &notifier{events: make(chan int), subscribeInternal: make(chan chan int)}
}

func (n *notifier) subscribe() <-chan int {
	sub := make(chan int)
	n.subscribeInternal <- sub
	return sub
}

func (n *notifier) run() {
	for {
		select {
		case event := <-n.events:
			fmt.Printf("received a new event: %+v\n", event)

			// inform subscribers
			for _, sub := range n.subscribers {
				sub <- event
			}

		case sub := <-n.subscribeInternal:
			n.subscribers = append(n.subscribers, sub)
		}
		// TODO: heartbeat for events
	}
}

func main() {
	trace.Start(os.Stderr)

	hub := NewNotifier()

	go hub.run()

	go func() {
		var i int
		for {
			time.Sleep(time.Second)
			hub.events <- i
			i++
		}
	}()

	subOne := hub.subscribe()

	go func() {
		var i int
		for val := range subOne {
			fmt.Printf("got value: %d\n", val)
			if i > 1 {
				time.Sleep(time.Millisecond * 50)
				break
			}
			i++
		}
	}()

	subTwo := hub.subscribe()

	go func() {
		for val := range subTwo {
			fmt.Printf("got value: %d\n", val)
		}
	}()

	time.Sleep(time.Second * 10)

	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}
