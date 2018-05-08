package redis

import (
	"strconv"

	"fmt"

	. "github.com/open-fightcoder/oj-web/common/store"
)

func ProblemCountGet(problemId int64) (string, error) {
	res := RedisClient.HMGet("problem_count", strconv.FormatInt(problemId, 10))
	if res.Err() != nil {
		return "", res.Err()
	}
	if res.Val()[0] == nil {
		return "", nil
	}
	return fmt.Sprint(res.Val()[0]), nil
}

func ProblemCountSet(problemId int64, mess string) bool {
	res := RedisClient.HMSet("problem_count", map[string]interface{}{strconv.FormatInt(problemId, 10): mess})
	str, _ := res.Result()
	if str == "OK" {
		return true
	}
	return false
}
