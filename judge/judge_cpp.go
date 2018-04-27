package judge

import (
	"fmt"

	"bufio"
	"os"

	"github.com/open-fightcoder/oj-judger/sandbox"
)

type JudgeCpp struct {
}

func (this *JudgeCpp) Compile(codeFile string) Result {
	sd := sandbox.NewSandbox("g++",
		[]string{codeFile, "-fmax-errors=200", "-o", "user.bin"},
		nil, 5000, 100000)
	output, errput, _, _, err := sd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
	fmt.Println(string(errput))
	return Result{}
}

func (this *JudgeCpp) Run(inputFile string, outputFile string) Result {
	input, err := os.Open(inputFile)
	sd := sandbox.NewSandbox("./user.bin",
		[]string{},
		bufio.NewReader(input), 5000, 100000)
	output, errput, timeUse, memoryUse, err := sd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("output [%s]\n", string(output))
	fmt.Printf("errput [%s]\n", string(errput))
	fmt.Println(timeUse, memoryUse)
	return Result{}
}
