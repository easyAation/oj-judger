package redis

import (
	"testing"

	"fmt"

	"github.com/open-fightcoder/oj-judger/common/g"
	"github.com/open-fightcoder/oj-judger/common/store"
)

func TestProblemStatusSet(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	//for i := 1; i <= 2300; i++ {
	fmt.Println(ProblemStatusSet(10, 7, 1))
	//}
}

func TestProblemStatusGet(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	fmt.Println(ProblemStatusGet(1, 3))
}
