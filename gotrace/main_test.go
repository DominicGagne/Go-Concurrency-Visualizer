package main

import (
	"path/filepath"
	"testing"
)

func TestExamples(t *testing.T) {
	tests := []struct {
		path        string
		cmdCount    int
		createCount int
		stopCount   int
		sendCount   int
	}{
		{"hello.go", 4, 2, 2, 1},
		{"pingpong01.go", 16, 3, 1, 12},
		{"pingpong02.go", 18, 4, 1, 13},
		{"pingpong03.go", 84, 37, 1, 46},
		{"fanin1.go", 46, 4, 2, 40},
		{"workers1.go", 126, 38, 38, 50},
		{"workers2.go", 346, 66, 60, 220},
		{"server1.go", 5, 3, 2, 0},
		{"primesieve.go", 188, 12, 1, 175},
		//{"select.go", 24, 3, 3, 0},
	}

	for _, test := range tests {
		path := filepath.Join("examples", test.path)
		t.Log("Testing", path)
		src := NewNativeRun(path)
		events, err := src.Events()
		if err != nil {
			t.Fatal(path, err)
		}
		commands, err := ConvertEvents(events)
		if err != nil {
			t.Fatal(path, err)
		}

		if commands.Count() != test.cmdCount {
			t.Log(commands.String())
			t.Fatalf("Wrong number of commands: %s: expecting %d, but got %d", path, test.cmdCount, commands.Count())
		}
		if commands.CountCreateGoroutine() != test.createCount {
			t.Fatalf("Wrong number of Create commands: %s: expecting %d, but got %d", path, test.createCount, commands.CountCreateGoroutine())
		}
		/*
			if commands.CountStopGoroutine() != test.stopCount {
				t.Fatalf("Wrong number of Stop commands: %s: expecting %d, but got %d", path, test.stopCount, commands.CountStopGoroutine())
			}
		*/
		if commands.CountSendToChannel() != test.sendCount {
			t.Fatalf("Wrong number of Send commands: %s: expecting %d, but got %d", path, test.sendCount, commands.CountSendToChannel())
		}
	}
}
