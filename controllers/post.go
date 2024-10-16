package controllers

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的处理函数

// CreatePostHandler 创建帖子的处理函数
// @Summary 创建新帖子
// @Description 根据请求体中的参数创建一个新的帖子
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.Post true "帖子内容" // 这里说明了请求体参数 post 的类型和是否必需
// @Success 201 {object} models.ResponseSuccess // 成功时返回的状态和对象类型
// @Failure 400 {object} models.ResponseError // 参数错误时返回的对象类型
// @Failure 401 {object} models.ResponseError // 未登录时返回的对象类型
// @Failure 500 {object} models.ResponseError // 服务器错误时返回的对象类型
// @Router /posts [post]
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的效验

	//c.ShouldBindJSON() validator -->binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到发请求的的用户的ID
	//从c取到当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(c, p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

// CreatePostDetailHandler  获取帖子详情的处理函数

// CetPostDetailHandler 获取帖子详情的处理函数
// @Summary 获取帖子详情
// @Description 根据帖子ID获取特定帖子的详细信息
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "帖子ID" // 从URL路径中获取的帖子ID，必需参数
// @Success 200 {object} models.Post // 成功时返回帖子详情
// @Failure 400 {object} models.ResponseError // 参数错误时返回的对象类型
// @Failure 404 {object} models.ResponseError // 找不到帖子时返回的对象类型
// @Failure 500 {object} models.ResponseError // 服务器错误时返回的对象类型
// @Router /posts/{id} [get]
func CetPostDetailHandler(c *gin.Context) {

	//1.获取参数（从URL获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据帖子id取出数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

//CetPostListHandler 获取帖子列表的处理函数

// GetPostListHandler 获取帖子列表的处理函数
// @Summary 获取帖子列表
// @Description 根据分页参数获取帖子列表
// @Tags post
// @Accept json
// @Produce json
// @Param page query int false "页码" // 可选，默认值为 1
// @Param size query int false "每页数量" // 可选，默认值为 10
// @Success 200 {object} models.ApiPostDetail // 成功时返回的帖子列表对象
// @Failure 500 {object} models.ResponseError // 服务器错误时返回的对象类型
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	//返回响应
	ResponseSuccess(c, data)
}

//CetPostListHandler2 获取帖子列表的处理函数
//根据前端传来的参数（）动态的获取帖子列表
//按创建时间排序 或者 按照 分数排序
//1.获取参数
//2.去redis查询id列表
//3.根据id去数据库查询帖子详细信息

// GetPostListHandler2 升级帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口（api分组展示使用的）
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	//GET请求参数/api/v1/posts2??page=1&size=10&order=time       query string
	//获取分页参数
	//初始化结构体时指定初始化参数
	p := &models.ParamPostList{
		CommunityID: 0,
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
	}
	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON()  如果请求中携带的是JSON格式的数据，才能用这个方法取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetPostListNew(c, p) //更新：合二为一
	//page, size := getPageInfo(c)
	//获取数据
	//data, err := logic.GetPostList2(c, p)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	//返回响应
	ResponseSuccess(c, data)
}

// 根据社区去查询帖子列表

//func GetCommunityListHandler(c *gin.Context) {
//	//GET请求参数/api/v1/posts2??page=1&size=10&order=time       query string
//	//获取分页参数
//	//初始化结构体时指定初始化参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
//	//c.ShouldBindJSON()  如果请求中携带的是JSON格式的数据，才能用这个方法取到数据
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//page, size := getPageInfo(c)
//	//获取数据
//	data, err := logic.GetCommunityPostList(c, p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//	}
//	//返回响应
//	ResponseSuccess(c, data)
//}
