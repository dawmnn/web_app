package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis/v8"
)

func getIDsFormKey(ctx context.Context, key string, page, size int64) ([]string, error) {
	//2.确定查询的索引起始点
	start := (page - 1) * size
	end := start*size - 1
	//3.ZRevRange 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(ctx, key, start, end).Result()
}

func GetPostIDsInOrder(ctx context.Context, p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引起始点
	return getIDsFormKey(ctx, key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投票赞成的的数据
func GetPostVoteData(ctx context.Context, ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	//查找key中分数是1的元素的数量
	//	v := client.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, v)
	//}

	//使用pipeline异常发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {

		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(ctx context.Context, p *models.ParamPostList) ([]string, error) {

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	fmt.Println(cKey)
	//利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(ctx, key).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Aggregate: "MAX",
			Keys:      []string{"cKey", "orderKey"},
		}) //ZInterStore 计算
		pipeline.Expire(ctx, key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	//存在的话就直接根据key查询ids
	return getIDsFormKey(ctx, key, p.Page, p.Size)
}
