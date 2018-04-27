package store

import (
	"testing"

	"strconv"

	"github.com/open-fightcoder/oj-judger/common/g"
)

func TestList(t *testing.T) {
	g.LoadConfig("../../cfg/cfg.toml.debug")
	InitRedis()
	RedisClient.ZRank("person_week_rank", strconv.FormatInt(11111, 10))

}
