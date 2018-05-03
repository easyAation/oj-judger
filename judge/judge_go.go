package judge

import (
	"bufio"
	"os"

	"github.com/open-fightcoder/oj-judger/sandbox"
)

type JudgeGo struct {
}

func (this *JudgeGo) Compile(workDir string, codeFile string) Result {
	sd := sandbox.NewSandbox("go",
		[]string{"build", workDir + "/" + codeFile},
		nil, nil,
		5000, 100000)
	_, _, err := sd.Run()
	if err != nil {
		return Result{
			ResultCode: CompilationError,
			ResultDes:  err.Error(),
		}
	}

	err = os.Rename(string([]byte(codeFile)[:len(codeFile)-3]), workDir+"/user.bin")
	if err != nil {
		return Result{
			ResultCode: SystemError,
			ResultDes:  err.Error(),
		}
	}

	// 编译成功
	return Result{
		ResultCode: Normal,
		ResultDes:  "",
	}
}

func (this *JudgeGo) Run(bin string, inputFile string, outputFile string, timeLimit int64, memoryLimit int64) Result {
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
