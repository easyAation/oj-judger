package judge

import (
	"bufio"
	"fmt"
	"os"

	"github.com/open-fightcoder/oj-judger/sandbox"
)

type JudgeCpp struct {
}

func (this *JudgeCpp) Compile(workDir string, codeFile string) Result {
	sd := sandbox.NewSandbox("g++",
		[]string{workDir + "/" + codeFile, "-fmax-errors=200", "-w", "-o", workDir + "/user.bin"},
		nil, nil,
		5000, 100000)
	_, _, err := sd.Run()
	if err != nil {
		return Result{
			ResultCode:    CompilationError,
			ResultDes:     err.Error(),
			RunningTime:   -1,
			RunningMemory: -1,
		}
	}

	// 编译成功
	return Result{
		ResultCode:    Normal,
		ResultDes:     "",
		RunningTime:   -1,
		RunningMemory: -1,
	}
}

func (this *JudgeCpp) Run(inputFile string, outputFile string) Result {
	fmt.Println(inputFile, outputFile)
	input, err := os.OpenFile(inputFile, os.O_RDWR, 0777)
	if err != nil {
		return Result{
			ResultCode: SystemError,
			ResultDes:  err.Error(),
		}
	}
	defer input.Close()
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return Result{
			ResultCode: SystemError,
			ResultDes:  err.Error(),
		}
	}
	defer output.Close()

	sd := sandbox.NewSandbox("./user.bin",
		[]string{},
		bufio.NewReader(input), bufio.NewWriter(output), 5000, 100000)
	timeUse, memoryUse, err := sd.Run()
	if err != nil {
		if err == sandbox.OutOfMemoryError {
			return Result{
				ResultCode:    MemoryLimitExceeded,
				RunningTime:   timeUse,
				RunningMemory: memoryUse,
				ResultDes:     err.Error(),
			}
		}
		if err == sandbox.OutOfTimeError {
			return Result{
				ResultCode:    TimeLimitExceeded,
				RunningTime:   timeUse,
				RunningMemory: memoryUse,
				ResultDes:     err.Error(),
			}
		}
		return Result{
			ResultCode:    RuntimeError,
			RunningMemory: memoryUse,
			RunningTime:   timeUse,
			ResultDes:     err.Error(),
		}
	}

	fmt.Println(timeUse, memoryUse)
	return Result{
		ResultCode:    Normal,
		RunningTime:   timeUse,
		RunningMemory: memoryUse,
		ResultDes:     "",
	}
}
