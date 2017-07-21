/**
source code are from https://github.com/go-cmd/cmd
#author: daniel-nichter
**/

package exec

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type KolaCmd struct {
	Name string
	Args []string

	*sync.Mutex
	started   bool
	stopped   bool
	done      bool
	final     bool
	startTime time.Time
	stdout    *output
	stderr    *output
	status    Status
	doneChan  chan Status
}

type Status struct {
	Cmd       string
	PID       int
	Completed bool
	Exit      int
	Error     error
	StartTs   int64
	StopTs    int64
	Runtime   float64
	Stdout    []string
	Stderr    []string
}

type output struct {
	buffer *bytes.Buffer
	lines  []string
	*sync.Mutex
}

func newOutput() *output {
	return &output{
		buffer: &bytes.Buffer{},
		lines:  []string{},
		Mutex:  &sync.Mutex{},
	}
}

func NewCmd(name string, args ...string) *KolaCmd {
	return &KolaCmd{
		Name: name,
		Args: args,

		Mutex: &sync.Mutex{},
		status: Status{
			Cmd:       name,
			PID:       0,
			Completed: false,
			Exit:      -1,
			Error:     nil,
			Runtime:   0,
		},
	}
}

func (c *KolaCmd) run() {
	defer func() {
		c.doneChan <- c.Status()
	}()
	//setup command
	cmd := exec.Command(c.Name, c.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.stdout = newOutput()
	c.stderr = newOutput()
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	//start the command base on exect run()
	now := time.Now()
	if err := cmd.Start(); err != nil {
		c.Lock()
		c.status.Error = err
		c.status.StartTs = now.UnixNano()
		c.status.StopTs = time.Now().UnixNano()
		c.done = true
		c.Unlock()
		return
	}

	//no error, set the initial status
	c.Lock()
	c.startTime = now
	c.status.PID = cmd.Process.Pid
	c.status.StartTs = now.UnixNano()
	c.started = true
	c.Unlock()

	//wait for the command to finish or be killed
	err := cmd.Wait()

	//Get the exit code of the command
	exitCode := 0
	signaled := false
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			err = nil
			if waitStatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus() // -i if signaled
				// if the command was termainated by a signal, then exiterr.Error()
				// is a string like "signal: terminated"
				if waitStatus.Signaled() {
					signaled = true
					err = errors.New(exiterr.Error())
				}
			}
		}
	}

	// Set the final status for the command
	c.Lock()
	if !c.stopped && !signaled {
		c.status.Completed = true
	}
	c.status.Runtime = time.Now().Sub(c.startTime).Seconds()
	c.status.StopTs = time.Now().UnixNano()
	c.status.Exit = exitCode
	c.status.Error = err
	c.done = true
	c.Unlock()
}

func (c *KolaCmd) Status() Status {
	c.Lock()
	defer c.Unlock()
	if c.doneChan == nil || !c.started {
		return c.status
	}

	if c.done {
		// No longer running
		if !c.final {
			c.status.Stdout = c.stdout.dumpLines()
			c.status.Stderr = c.stderr.dumpLines()
			c.stdout = nil // release buffer
			c.stderr = nil

			c.final = true
		}
	} else {
		// still running
		c.status.Runtime = time.Now().Sub(c.startTime).Seconds()
		c.status.Stdout = c.stdout.dumpLines()
		c.status.Stderr = c.stderr.dumpLines()
	}

	return c.status
}

func (c *KolaCmd) Start() <-chan Status {
	c.Lock()
	defer c.Unlock()
	if c.doneChan != nil {
		return c.doneChan
	}
	//avoid the race condition
	// this done channel actually works like a activation lock
	c.doneChan = make(chan Status, 1)
	go c.run()
	return c.doneChan
}

func (c *KolaCmd) Stop() error {
	c.Lock()
	defer c.Unlock()

	if c.doneChan == nil || !c.started || !c.done {
		return nil
	}
	c.stopped = true
	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return syscall.Kill(-c.status.PID, syscall.SIGTERM)
}

func (o *output) dumpLines() []string {
	o.Lock()
	defer o.Unlock()
	s := bufio.NewScanner(o.buffer)
	for s.Scan() {
		o.lines = append(o.lines, s.Text())
	}
	return o.lines
}

func (o *output) Write(p []byte) (int, error) {
	o.Lock()
	defer o.Unlock()
	return o.buffer.Write(p)
}
