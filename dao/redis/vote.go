package redis

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

//投票的几种情况：
//direction=1时，有两种情况：
//1.之前没有投过票，现在投赞成票 --> 更新分数和投票记录  差值的绝对值：1   +432
//2.之前投反对票，现在改投赞成票 --> 更新分数和投票记录  差值的绝对值：2   +432*2

//direction=0时，有两种情况：
//1.之前投过赞成票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  -432
//2.之前投过反对票，现在要取消投票  --> 更新分数和投票记录   差值的绝对值：1  +432

//direction=-1时，有两种情况
//1.之前没有投过票，现在投反对票  --> 更新分数和投票记录   差值的绝对值：1   -432
//2.之前投过赞成票，现在改投反对票  --> 更新分数和投票记录   差值的绝对值：2   -432*2

//本项目使用简化版的投票分数
//投一票加432分   86400/200 -> 需要200张赞成票可以给你的帖子续一天  ->《redis实战》

//VoteForPost 为帖子投票的函数

//投票的限制：每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了
//1，到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
//2.到期之后删除KeyPostVotedZSetPrefix

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

func CreatePost(ctx context.Context, postID int64, communityID int64) error {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(423),
		Member: postID,
	})
	//
	//更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

func VoteForPost(ctx context.Context, userID, postID string, value float64) error {
	// 1.判断 投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2和3需要放到一个pipeline事务中操作
	//2.更新帖子分数
	//先查当前用户给当前的帖子的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	//如果这一次投票的值和上一次投票的值一样，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepested
	}
	var op float64

	if value > ov {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(ov - value) //计算两次投票的差值（取绝对值）
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投过票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), userID)

	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPrefix+postID), &redis.Z{Score: value,
			Member: userID})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
