package redis

import (
	"fmt"
	"strings"

	. "github.com/open-fightcoder/oj-judger/common/store"
)

func ProblemNumGet() (string, error) {
	res := RedisClient.Get("problem_num")
	if res.Err() != nil {
		if strings.Contains(fmt.Sprint(res.Err()), "redis: nil") {
			return "", nil
		}
		return "", res.Err()
	}
	return res.Val(), nil
}

func ProblemNumUpdate() error {
	res := RedisClient.Incr("problem_num")
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
