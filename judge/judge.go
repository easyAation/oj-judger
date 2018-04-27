package judge

type Judge interface {
	Compile(workDir string, codeFile string) Result
	Run(bin string, inputFile string, outputFile string) Result
}

func NewJudge(language string) Judge {
	switch language {
	case "c++":
		return &JudgeCpp{}
	default:
		panic("No such judge: " + language)
	}
}
