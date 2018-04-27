package sandbox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

var (
	OutOfTimeError   = errors.New("out of time")
	OutOfMemoryError = errors.New("out of memory")
)

type Sandbox struct {
	Bin         string    // binary path
	Args        []string  // arguments
	Input       io.Reader // standard input
	TimeLimit   int64     // time limit in ms
	MemoryLimit int64     // memory limit in kb
}

func NewSandbox(bin string, args []string, input io.Reader, timeLimit int64, memoryLimit int64) *Sandbox {
	sandbox := new(Sandbox)
	sandbox.Bin = bin
	sandbox.Args = args
	sandbox.Input = input
	sandbox.TimeLimit = timeLimit
	sandbox.MemoryLimit = memoryLimit
	return sandbox
}

func (s *Sandbox) Run() (output []byte, errput []byte, timeUse int64, memoryUse int64, err error) {
	cmd := exec.Command(s.Bin, s.Args...)

	errBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	cmd.Stderr = errBuf
	cmd.Stdout = outBuf
	cmd.Stdin = s.Input

	if err = cmd.Start(); err != nil {
		return
	}
	defer cmd.Process.Kill()

	errCh := make(chan error)
	defer close(errCh)

	go func() {
		var rusage syscall.Rusage
		var wStatus syscall.WaitStatus

		_, err = syscall.Wait4(cmd.Process.Pid, &wStatus, syscall.WUNTRACED, &rusage)
		if err != nil {
			fmt.Println("wait error", err)
			errCh <- err
		}
		if wStatus.Signaled() {
			sig := wStatus.Signal()
			errCh <- fmt.Errorf("get signal %s", sig)
		}
		errCh <- nil
	}()

	fmt.Println(cmd.Process.Pid)
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		ok, vm, rss, runningTime, cpuTime := GetResourceUsage(cmd.Process.Pid)
		if !ok {
			fmt.Println("cmd 退出")
			break
		}
		fmt.Println(vm, rss, runningTime, cpuTime)
		timeUse = cpuTime
		memoryUse = rss * 3 / 2

		if cpuTime > s.TimeLimit ||
			runningTime*2 > 3*s.TimeLimit {
			err = OutOfTimeError
			fmt.Println("cpu limit")
			break

		}

		if rss*3 > s.MemoryLimit*2 ||
			vm > s.MemoryLimit*10 {
			err = OutOfMemoryError
			fmt.Println("rss limit")
			break
		}
	}

	err = <-errCh
	fmt.Println("err = ", err)

	output = outBuf.Bytes()
	errput = errBuf.Bytes()
	return
}
