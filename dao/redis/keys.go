package redis

//redis key

//redis key 尽量注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix              = "web_app:"
	KeyPostTimeZSet        = "post:time"   //zset:帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset:帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" //zset:记录用户及投票类型;参数是post id

	KeyCommunitySetPF = "community:" //set;保存每个分区下的帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
