package logic

import (
	"context"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(ctx context.Context, p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	//3.返回
	err = redis.CreatePost(ctx, p.ID, p.CommunityID)
	return err
}
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合我们接口想用的数据
	//data=new(models.ApiPostDetail)
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed",
			zap.Int64("author_id", post.CommunityID),
			zap.Error(err))
		return
	}
	//接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	//data.AuthorName=user.Username
	//data.Community=community
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//3.根据id去Mysql数据库查询帖子详细信息
	//返回的数据还有按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好查询的每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}
	zap.L().Info("redis.GetPostVoteData", zap.Any("voteData", voteData))
	//将帖子的作者及分区信息查询出来填充到帖子中
	for idex, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			VoteNUm:         voteData[idex],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	//2.去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))
	//3.根据id去Mysql数据库查询帖子详细信息
	//返回的数据还有按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好查询的每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	//将帖子的作者及分区信息查询出来填充到帖子中
	for idex, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("author_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			VoteNUm:         voteData[idex],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个查询帖子列表合二为一的函数
func GetPostListNew(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		//按时间
		data, err = GetPostList2(ctx, p)
	} else {
		//根据社区id查询
		//按分数
		data, err = GetCommunityPostList(ctx, p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
