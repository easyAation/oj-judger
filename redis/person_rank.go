package redis

import (
	"strconv"

	"errors"

	"github.com/go-redis/redis"
	. "github.com/open-fightcoder/oj-web/common/store"
)

func PersonWeekRankAdd(userId int64) error {
	res := RedisClient.ZAdd("person_week_rank", redis.Z{0, userId})
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func PersonWeekRankUpdate(increment int, userId int64) error {
	res := RedisClient.ZIncrBy("person_week_rank", float64(increment), strconv.FormatInt(userId, 10))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func PersonWeekRankGet(userId int64) ([]map[string]interface{}, error) {
	idStr := strconv.FormatInt(userId, 10)
	isExitRet := RedisClient.ZScore("person_week_rank", idStr)
	if isExitRet.Val() > 0 {
		sizeRet := RedisClient.ZCard("person_week_rank")
		if sizeRet.Err() != nil {
			return nil, errors.New("获取失败")
		}
		size := sizeRet.Val()
		var start int64
		var end int64
		if size <= 5 {
			start = 0
			end = 4
		} else {
			res := RedisClient.ZRevRank("person_week_rank", idStr)
			if res.Err() != nil {
				return nil, errors.New("获取失败")
			}
			index := res.Val()
			if index < 2 {
				start = 0
				end = 4
			} else if index > size-3 {
				start = size - 5
				end = size - 1
			} else {
				start = index - 2
				end = index + 2
			}
		}
		result := RedisClient.ZRevRange("person_week_rank", start, end)
		if result.Err() != nil {
			return nil, errors.New("获取失败")
		}
		var rankLists []map[string]interface{}
		for i, v := range result.Val() {
			projects := make(map[string]interface{})
			scoreRes := RedisClient.ZScore("person_week_rank", v)
			projects["rank_num"] = i + 1
			projects["user_id"] = v
			projects["ac_num"] = scoreRes.Val()
			rankLists = append(rankLists, projects)
		}
		return rankLists, nil
	} else {
		return nil, nil
	}
}

func PersonMonthRankAdd(userId int64) error {
	res := RedisClient.ZAdd("person_month_rank", redis.Z{0, userId})
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func PersonMonthRankUpdate(increment int, userId int64) error {
	res := RedisClient.ZIncrBy("person_month_rank", float64(increment), strconv.FormatInt(userId, 10))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func PersonMonthRankGet(userId int64) ([]map[string]interface{}, error) {
	idStr := strconv.FormatInt(userId, 10)
	isExitRet := RedisClient.ZScore("person_month_rank", idStr)
	if isExitRet.Val() > 0 {
		sizeRet := RedisClient.ZCard("person_month_rank")
		if sizeRet.Err() != nil {
			return nil, errors.New("获取失败")
		}
		size := sizeRet.Val()
		var start int64
		var end int64
		if size <= 5 {
			start = 0
			end = 4
		} else {
			res := RedisClient.ZRevRank("person_month_rank", idStr)
			if res.Err() != nil {
				return nil, errors.New("获取失败")
			}
			index := res.Val()
			if index < 2 {
				start = 0
				end = 4
			} else if index > size-3 {
				start = size - 5
				end = size - 1
			} else {
				start = index - 2
				end = index + 2
			}
		}
		result := RedisClient.ZRevRange("person_month_rank", start, end)
		if result.Err() != nil {
			return nil, errors.New("获取失败")
		}
		var rankLists []map[string]interface{}
		for _, v := range result.Val() {
			projects := make(map[string]interface{})
			scoreRes := RedisClient.ZScore("person_month_rank", v)
			rankId := RedisClient.ZRevRank("person_month_rank", v)
			projects["rank_num"] = rankId.Val() + 1
			projects["user_id"] = v
			projects["ac_num"] = scoreRes.Val()
			rankLists = append(rankLists, projects)
		}
		return rankLists, nil
	} else {
		return nil, nil
	}
}
