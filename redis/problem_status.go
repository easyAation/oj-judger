package redis

import (
	"strconv"

	"fmt"

	"strings"

	. "github.com/open-fightcoder/oj-judger/common/store"
)

func ProblemStatusGet(userId int64, problemId int64) (int, error) {
	keyStr := strconv.FormatInt(userId, 10) + "_" + strconv.FormatInt(problemId, 10)
	res := RedisClient.Get(keyStr)
	if res.Err() != nil {
		if strings.Contains(fmt.Sprint(res.Err()), "redis: nil") {
			return 0, nil
		}
		return 0, res.Err()
	}
	status, _ := strconv.Atoi(res.Val())
	return status, nil
}

func ProblemStatusSet(userId int64, problemId int64, status int) error {
	keyStr := strconv.FormatInt(userId, 10) + "_" + strconv.FormatInt(problemId, 10)
	res := RedisClient.Get(keyStr)
	if res.Err() != nil && !strings.Contains(fmt.Sprint(res.Err()), "redis: nil") {
		return res.Err()
	}
	setRes := RedisClient.Set(keyStr, status, 0)
	if setRes.Err() != nil {
		return res.Err()
	}
	return nil
}
