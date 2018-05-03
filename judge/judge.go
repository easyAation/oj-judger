package judge

import (
	log "github.com/sirupsen/logrus"
)

type Judge interface {
	Compile(workDir string, codeFile string) Result
	Run(bin string, inputFile string, outputFile string, timeLimit int64, memoryLimit int64) Result
}

func NewJudge(language string) Judge {
	log.Debugf("new judge from language [%s]", language)
	switch language {
	case "c":
		return &JudgeCpp{}
	case "c++":
		return &JudgeCpp{}
	case "golang":
		return &JudgeGo{}
	case "java":
		return &JudgeJava{}
	case "python":
		return &JudgePy{}
	default:
		panic("No such judge: " + language)
	}
}
