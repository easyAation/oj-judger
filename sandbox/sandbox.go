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
		log.Println("A")
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
		log.Println(err)
		fmt.Println("B")
		return
	}
	defer cmd.Process.Kill()

	errCh := make(chan error)

	go func() {
		var rusage syscall.Rusage
		var wStatus syscall.WaitStatus

		_, err := syscall.Wait4(cmd.Process.Pid, &wStatus, syscall.WUNTRACED, &rusage)
		if err != nil {
			log.Println("wait error", err)
			errCh <- err
		}
		if wStatus.Signaled() {
			sig := wStatus.Signal()
			log.Println("get signal", sig)
			errCh <- fmt.Errorf("get signal %s", sig)
		}
		errCh <- nil
	}()

	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		ok, vm, rss, runningTime, cpuTime := GetResourceUsage(cmd.Process.Pid)
		if !ok {
			log.Println("cmd 退出")
			break
		}
		//fmt.Println(vm, rss, runningTime, cpuTime)
		timeUse = int(cpuTime)
		memoryUse = int(rss * 3 / 2)

		if cpuTime > s.TimeLimit ||
			runningTime > 5*s.TimeLimit {
			err = OutOfTimeError
			log.Println("cpu limit: ", runningTime, cpuTime, s.TimeLimit)
			return
		}

		if rss*3 > s.MemoryLimit*2 ||
			vm > s.MemoryLimit*10 {
			err = OutOfMemoryError
			log.Println("rss limit: ", vm, rss, s.MemoryLimit)
			return
		}
	}

	if err != nil {
		log.Println("A", err.Error())
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
