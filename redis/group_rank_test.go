package redis

import (
	"testing"

	"fmt"

	"github.com/open-fightcoder/oj-web/common/g"
	"github.com/open-fightcoder/oj-web/common/store"
)

func TestGroupRankAdd(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	err := GroupRankAdd(8)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGroupRankUpdate(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	err := GroupRankIncr(1, 1)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGroupRankGet(t *testing.T) {
	g.LoadConfig("../cfg/cfg.toml.debug")
	store.InitRedis()

	res, err := GroupRankGet(1, 5)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}
