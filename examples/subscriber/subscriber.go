package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

type notifier struct {
	events            chan int
	subscribers       map[int]chan int
	subscribeInternal chan chan int
}

func NewNotifier() *notifier {
	return &notifier{events: make(chan int), subscribers: make(map[int]chan int), subscribeInternal: make(chan chan int)}
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
			for id, sub := range n.subscribers {
				select {
				case sub <- event:
					// subscriber has received the message, carry on to the next one
				case <-time.After(time.Millisecond * 500):
					// subscriber has not responded in 500ms, kill that consumer
					delete(n.subscribers, id)
				}
			}
		case sub := <-n.subscribeInternal:
			n.subscribers[len(n.subscribers)] = sub
		}
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
			if i > 4 {
				break
			}
			i++
		}
	}()

	time.Sleep(time.Second * 10)

	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}
