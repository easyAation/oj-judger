package judge

import (
	"bufio"
	"os"

	"github.com/open-fightcoder/oj-judger/sandbox"
)

type JudgePy struct {
}

func (this *JudgePy) Compile(workDir string, codeFile string) Result {
	// 编译成功
	// 将code名改为bin名，后缀可在run里添加.py即可
	return Result{
		ResultCode: Normal,
		ResultDes:  "",
	}
}

func (this *JudgePy) Run(bin string, inputFile string, outputFile string, timeLimit int64, memoryLimit int64) Result {
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

	sd := sandbox.NewSandbox("python",
		[]string{""},
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
