package judge

type Judge interface {
	Compile(codeFile string) Result
	Run(inputFile string, outputFile string) Result
}

func newJudge(language string) Judge {
	switch language {
	case "c++":
		return &JudgeCpp{}
	default:
		panic("No such judge: " + language)
	}
}
