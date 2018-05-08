package redis

import (
	"strconv"

	"github.com/go-redis/redis"
	. "github.com/open-fightcoder/oj-web/common/store"
)

func GroupRankAdd(groupId int64) error {
	res := RedisClient.ZAdd("group_rank", redis.Z{0, groupId})
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func GroupRankIncr(increment int, groupId int64) error {
	res := RedisClient.ZIncrBy("group_rank", float64(increment), strconv.FormatInt(groupId, 10))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func GroupRankGet(currentPage int, perPage int) ([]map[string]interface{}, error) {
	res := RedisClient.ZRevRange("group_rank", int64((currentPage-1)*perPage), int64(currentPage*perPage-1))
	if res.Err() != nil {
		return nil, res.Err()
	}
	rankLists := make([]map[string]interface{}, 0)
	for i, v := range res.Val() {
		projects := make(map[string]interface{})
		scoreRes := RedisClient.ZScore("group_rank", v)
		projects["rank_num"] = i + 1
		projects["group_id"] = v
		projects["ac_num"] = scoreRes.Val()
		rankLists = append(rankLists, projects)
	}
	return rankLists, nil
}
