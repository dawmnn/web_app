package logic

import (
	"context"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

//基于用户投票功能：
//1.用户投票的数据
//

func VoteForPost(ctx context.Context, userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(ctx, strconv.Itoa(int(userID)), p.PostID, float64(p.Direction)) //整数转字符串

}
