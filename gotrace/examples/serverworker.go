package main

import (
	"fmt"
	"net"
	"os"
	"runtime/trace"
	"time"
)

func handler(c net.Conn, ch chan int) {
	time.Sleep(50 * time.Millisecond)
	ch <- 0
	c.Write([]byte("ok"))
	c.Close()
}

func worker(wch chan int, results chan int) {
	for {
		data := <-wch
		data++
		results <- data
	}
}

func parse(results chan int) {
	for {
		<-results
	}
}

func pool(ch chan int, n int) {
	wch := make(chan int)
	results := make(chan int)
	for i := 0; i < n; i++ {
		go worker(wch, results)
		time.Sleep(1 * time.Millisecond)
	}
	go parse(results)
	for {
		val := <-ch
		wch <- val
	}
}

func server(l net.Listener, ch chan int) {
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go handler(c, ch)
	}
}

func main() {
	trace.Start(os.Stderr)

	fmt.Println("Listening on :5000. Send something using nc: echo hello | nc localhost 5000")
	fmt.Println("Exiting in 2 seconds...")
	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	ch := make(chan int)
	go pool(ch, 36)
	go server(l, ch)
	time.Sleep(2 * time.Second)
	trace.Stop()
}
