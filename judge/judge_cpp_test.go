package judge

import "testing"

var workDir = "/Users/shiyi/project/go/src/github.com/open-fightcoder/oj-judger"

func TestJudgeCpp_Compile(t *testing.T) {
	judgeCpp := new(JudgeCpp)
	judgeCpp.Compile("./scripts/test/cpp_ajb.cpp")
}
