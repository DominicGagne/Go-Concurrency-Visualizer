package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	// gotrace can take either .go file or already generated
	// trace file as an argument.
	// - for .go file, it will use "Native Run" mode - instrument
	// the code with tracing, run it and collect the trace.
	// - for trace file, it will proceed directly with its content.
	var src EventSource
	if len(args) == 1 {
		if strings.HasSuffix(args[0], ".go") {
			src = NewNativeRun(args[0])
		} else {
			// assuming all other types of file are trace files
			src = NewTraceSource(args[0])
		}
	} else {
		Usage()
		os.Exit(1)
	}

	events, err := src.Events()
	if err != nil {
		panic(err)
	}

	commands, err := ConvertEvents(events)
	if err != nil {
		panic(err)
	}

	ProcessCommands(commands)
}

// ProcessCommands processes command list.
func ProcessCommands(cmds Commands) {
	params := GuessParams(cmds)
	data := cmds.toJSON()

	StartServer(":2000", data, params)
}

// Usage prints usage information, overriding default one.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [trace.out] or [main.go]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       (if you pass .go file to gotrace, it will modify code on the fly,\n")
	fmt.Fprintf(os.Stderr, "       adding tracing, run it and collect the trace automagically)\n")
	flag.PrintDefaults()
}
