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
		[]string{workDir + "/" + codeFile, "-fmax-errors=200", "-O2", "-static", "-w", "-o", workDir + "/user.bin"},
		nil, nil,
		5000, 100000)
	_, _, err := sd.Run()
	if err != nil {
		return Result{
			ResultCode: CompilationError,
			ResultDes:  err.Error(),
		}
	}

	// 编译成功
	return Result{
		ResultCode: Normal,
		ResultDes:  "",
	}
}

func (this *JudgeCpp) Run(bin string, inputFile string, outputFile string, timeLimit int64, memoryLimit int64) Result {
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

	sd := sandbox.NewSandbox(bin,
		[]string{},
		bufio.NewReader(input), bufio.NewWriter(output), timeLimit, memoryLimit)
	timeUse, memoryUse, err := sd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(timeUse, memoryUse)
	if err != nil {
		if err == sandbox.OutOfMemoryError {
			// 超时
			return Result{
				ResultCode:    MemoryLimitExceeded,
				RunningTime:   timeUse,
				RunningMemory: memoryUse,
				ResultDes:     "",
			}
		}
		if err == sandbox.OutOfTimeError {
			// 超内存
			return Result{
				ResultCode:    TimeLimitExceeded,
				RunningTime:   timeUse,
				RunningMemory: memoryUse,
				ResultDes:     "",
			}
		}
		// 运行异常
		return Result{
			ResultCode:    RuntimeError,
			RunningMemory: memoryUse,
			RunningTime:   timeUse,
			ResultDes:     err.Error(),
		}
	}

	// 正常
	return Result{
		ResultCode:    Normal,
		RunningTime:   timeUse,
		RunningMemory: memoryUse,
		ResultDes:     "",
	}
}
