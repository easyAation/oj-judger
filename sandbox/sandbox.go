package sandbox

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	OutOfTimeError   = errors.New("out of time")
	OutOfMemoryError = errors.New("out of memory")
)

type Sandbox struct {
	Bin         string
	Args        []string
	Input       *bufio.Reader
	Output      *bufio.Writer
	TimeLimit   int64
	MemoryLimit int64
}

func NewSandbox(bin string, args []string, input *bufio.Reader, output *bufio.Writer, timeLimit int64, memoryLimit int64) *Sandbox {
	sandbox := new(Sandbox)
	sandbox.Bin = bin
	sandbox.Args = args
	sandbox.Input = input
	sandbox.Output = output
	sandbox.TimeLimit = timeLimit
	sandbox.MemoryLimit = memoryLimit
	return sandbox
}

func (s *Sandbox) Run() (timeUse int, memoryUse int, err error) {
	cmd := exec.Command(s.Bin, s.Args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Debugf("cmd.StderrPipe: %s", err.Error())
		return
	}
	if s.Output != nil {
		cmd.Stdout = s.Output
		defer s.Output.Flush()
	}
	if s.Input != nil {
		cmd.Stdin = s.Input
	}
	if err = cmd.Start(); err != nil {
		log.Debugf("cmd.Start: %s", err.Error())
		return
	}
	defer cmd.Process.Kill()
	errCh := make(chan error)
	isEnd := false

	go func() {
		var rusage syscall.Rusage
		var wStatus syscall.WaitStatus

		_, err := syscall.Wait4(cmd.Process.Pid, &wStatus, syscall.WUNTRACED, &rusage)
		isEnd = true
		if err != nil {
			log.Debugf("syscall.Wait4: %s", err.Error())
			errCh <- fmt.Errorf("syscall wait %s", err.Error())
		}
		if wStatus.Signaled() {
			sig := wStatus.Signal()
			log.Debugf("wStatus.Signal: %s", err.Error())
			errCh <- fmt.Errorf("get signal %s", sig)
		}
		errCh <- nil
	}()

	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		if isEnd {
			break
		}
		ok, vm, rss, runningTime, cpuTime := GetResourceUsage(cmd.Process.Pid)

		if !ok {
			log.Debugf("cmd 退出")
			break
		}
		//fmt.Println(vm, rss, runningTime, cpuTime)
		if timeUse < int(cpuTime) {
			timeUse = int(cpuTime)
		}
		if memoryUse < int(rss*3/2) {
			memoryUse = int(rss * 3 / 2)
		}

		if cpuTime > s.TimeLimit ||
			runningTime > 5*s.TimeLimit {
			err = OutOfTimeError
			log.Debug("cpu limit: ", runningTime, cpuTime, s.TimeLimit)
			return
		}

		if rss*3 > s.MemoryLimit*2 ||
			vm > s.MemoryLimit*10 {
			err = OutOfMemoryError
			log.Debug("rss limit: ", vm, rss, s.MemoryLimit)
			return
		}
	}

	if err != nil {
		log.Debug("A", err.Error())
	}

	err = <-errCh
	if err != nil {
		return
	}

	var errput []byte
	errput, err = ioutil.ReadAll(stderr)
	if err != nil {
		return
	}
	if string(errput) != "" {
		err = errors.New(string(errput))
	}

	return
}
