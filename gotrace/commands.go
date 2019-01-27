package main

import (
	"encoding/json"
	"fmt"
)

type Commands struct {
	cmds []*Command
	gd   map[string]int //goroutines depth map
}

func NewCommands() Commands {
	return Commands{
		gd: make(map[string]int),
	}
}

const (
	CmdCreate  = "create goroutine"
	CmdStop    = "stop goroutine"
	CmdSend    = "send to channel"
	CmdBlock   = "block goroutine"
	CmdUnblock = "unblock goroutine"
	CmdSleep   = "sleep goroutine"
)

// Command is a common structure for all
// types of supported events (aka 'commands').
// It's main purpose to handle JSON marshalling.
type Command struct {
	Time     int64       "json:\"t\""
	Command  string      "json:\"command\""
	Name     string      "json:\"name,omitempty\""
	Parent   string      "json:\"parent,omitempty\""
	Channels []string    "json:\"channels,omitempty\""
	From     string      "json:\"from,omitempty\""
	To       string      "json:\"to,omitempty\""
	Channel  string      "json:\"ch,omitempty\""
	Value    interface{} "json:\"value,omitempty\""
	EventID  string      "json:\"eid,omitempty\""
	Duration int64       "json:\"duration,omitempty\""
	Depth    int         "json:\"depth,omitempty\""
}

func (c *Commands) toJSON() []byte {
	data, err := json.MarshalIndent(c.cmds, "", "  ")
	if err != nil {
		panic(err)
	}

	return data
}

func (c *Commands) StartGoroutine(ts int64, gname string, gid, pid uint64) {
	// TODO: use gname
	name := fmt.Sprintf("#%d", gid)
	parent := fmt.Sprintf("#%d", pid)

	// ignore parent for 'main()' which has pid 0
	if pid == 0 {
		parent = ""
	}

	c.gd[name] = 0
	if parent != "" {
		c.gd[name] = c.gd[parent] + 1
	}

	cmd := &Command{
		Time:    ts,
		Command: CmdCreate,
		Name:    name,
		Parent:  parent,
		Depth:   c.gd[name],
	}
	c.cmds = append(c.cmds, cmd)
}

func (c *Commands) StopGoroutine(ts int64, name string, gid uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdStop,
		Name:    fmt.Sprintf("#%d", gid),
	}
	c.cmds = append(c.cmds, cmd)
}

func (c *Commands) ChanSend(send_ts, recv_ts int64, cid, fgid, tgid, val uint64) {
	cmd := &Command{
		Time:     recv_ts,
		Command:  CmdSend,
		From:     fmt.Sprintf("#%d", fgid),
		To:       fmt.Sprintf("#%d", tgid),
		Channel:  fmt.Sprintf("#%d", cid),
		Value:    fmt.Sprintf("%d", val),
		Duration: recv_ts - send_ts,
	}
	c.cmds = append(c.cmds, cmd)
}

func (c *Commands) BlockGoroutine(ts int64, gid uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdBlock,
		Name:    fmt.Sprintf("#%d", gid),
	}
	c.cmds = append(c.cmds, cmd)
}

func (c *Commands) UnblockGoroutine(ts int64, gid uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdUnblock,
		Name:    fmt.Sprintf("#%d", gid),
	}
	c.cmds = append(c.cmds, cmd)
}

func (c *Commands) SleepGoroutine(ts int64, gid uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdSleep,
		Name:    fmt.Sprintf("#%d", gid),
	}
	c.cmds = append(c.cmds, cmd)
}

//ByTimestamp implements sort.Interface for sorting command by timestamp.
type ByTimestamp []*Command

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Time < a[j].Time }

// Count counts total number of commands, which is not block/unblock ones.
func (c Commands) Count() int {
	var count int
	for _, c := range c.cmds {
		if c.Command != CmdUnblock && c.Command != CmdBlock && c.Command != CmdSleep {
			count++
		}
	}
	return count
}

// CountCreateGoroutine counts total number of CreateGoroutine commands.
func (c Commands) CountCreateGoroutine() int {
	var count int
	for _, cmd := range c.cmds {
		if cmd.Command == CmdCreate {
			count++
		}
	}
	return count
}

// CountStopGoroutine counts total number of StopGoroutine commands.
func (c Commands) CountStopGoroutine() int {
	var count int
	for _, cmd := range c.cmds {
		if cmd.Command == CmdStop {
			count++
		}
	}
	return count
}

// CountSendToChannel counts total number of SendToChannel commands.
func (c Commands) CountSendToChannel() int {
	var count int
	for _, cmd := range c.cmds {
		if cmd.Command == CmdSend {
			count++
		}
	}
	return count
}

// String implements Stringer inteface for Commands.
func (c Commands) String() string {
	var out string
	for _, cmd := range c.cmds {
		if cmd.Command != CmdUnblock && cmd.Command != CmdBlock && cmd.Command != CmdSleep {
			out = fmt.Sprintf("%s%s %v %v (%v -> %v)\n", out, cmd.Command, cmd.Parent, cmd.Name, cmd.From, cmd.To)
		}
	}
	return out
}
