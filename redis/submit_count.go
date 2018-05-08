package redis

import (
	"strconv"

	"fmt"

	. "github.com/open-fightcoder/oj-web/common/store"
)

func SubmitCountGet(userId int64) (string, error) {
	res := RedisClient.HMGet("submit_count", strconv.FormatInt(userId, 10))
	if res.Err() != nil {
		return "", res.Err()
	}
	if res.Val()[0] == nil {
		return "", nil
	}
	return fmt.Sprint(res.Val()[0]), nil
}

func SubmitCountSet(userId int64, mess string) bool {
	res := RedisClient.HMSet("submit_count", map[string]interface{}{strconv.FormatInt(userId, 10): mess})
	str, _ := res.Result()
	if str == "OK" {
		return true
	}
	return false
}
