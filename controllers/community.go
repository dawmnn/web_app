package controllers

import (
	"strconv"
	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- 和社区相关的----

// CommunityHandler 主页接口
// @Summary 得到社区信息的接口
// @Description 得到社区信息
// @Tags 社区相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param ParamSignUp body models.Community true "社区参数"
// @Security ApiKeyAuth
// @Success 200 {object} models.ResponseSuccess "成功响应"
// @Success 400 {object} models.ResponseError "请求参数错误"
// @Success 409 {object} models.ResponseError "用户名已存在"
// @Success 500 {object} models.ResponseError "服务器错误"
// @Router /api/v1/community [get]
func CommunityHandler(c *gin.Context) {
	//查询到所有的社区 （community_id,community_name）以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandle 社区分类详情

// CommunityDetailHandler 社区分类详情
// @Summary 获取社区详情
// @Description 根据社区ID获取社区的详细信息
// @Tags community
// @Accept json
// @Produce json
// @Param id path int true "社区ID" // 这里说明了路径参数 id 的类型和是否必需
// @Success 200 {object} models.ResponseSuccess "成功响应"
// @Success 400 {object} models.ResponseError "请求参数错误"
// @Success 409 {object} models.ResponseError "用户名已存在"
// @Success 500 {object} models.ResponseError "服务器错误"
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	communityIDStr := c.Param("id")
	id, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询到所有的社区 （community_id,community_name）以列表形式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
