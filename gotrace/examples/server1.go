package main

import (
	"fmt"
	"net"
	"os"
	"runtime/trace"
	"time"
)

func handler(c net.Conn) {
	c.Write([]byte("ok"))
	c.Close()
}

func main() {
	trace.Start(os.Stderr)
	fmt.Println("Listening on :5000. Send something using nc: echo hello | nc localhost 5000")
	fmt.Println("Exiting in 1 second...")

	go time.AfterFunc(1*time.Second, func() { trace.Stop(); os.Exit(0) })

	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go handler(c)
	}
}
