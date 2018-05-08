package redis

import (
	"testing"

	"fmt"

	"github.com/open-fightcoder/oj-web/common/g"
	"github.com/open-fightcoder/oj-web/common/store"
)

func TestProblemCountSet(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	for i := 1; i <= 2000; i++ {
		ProblemCountSet(int64(i), "{\"ac_num\":0,\"total_num\":1}")
	}
}

func TestProblemCountGet(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	aa, err := ProblemCountGet(13)
	fmt.Println(aa, err)
}
