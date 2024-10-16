package controllers

import (
	"web_app/logic"
	"web_app/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//投票

// PostVoteController 处理用户对帖子的投票请求
// @Summary 用户对帖子进行投票
// @Description 接收投票数据并处理投票逻辑
// @Tags vote
// @Accept json
// @Produce json
// @Param voteData body models.ParamVoteData true "投票数据"
// @Success 200 {object} models.ResponseSuccess // 成功时返回的响应
// @Failure 400 {object} models.ResponseError // 参数验证错误时返回
// @Failure 500 {object} models.ResponseError // 服务器错误时返回
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	//参数效验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译去除错误提示中的结构体
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	//获取当前请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
	}
	//具体投票的业务逻辑
	if err = logic.VoteForPost(c, userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
