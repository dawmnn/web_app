package controllers

import (
	"errors"
	"fmt"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数

// SignUpHandler 用户注册接口
// @Summary 用户注册接口
// @Description 注册用户账户
// @Tags 用户相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param ParamSignUp body models.ParamSignUp true "用户注册参数"
// @Security ApiKeyAuth
// @Success 200 {object} models.ResponseSuccess "成功响应"
// @Success 400 {object} models.ResponseError "请求参数错误"
// @Success 409 {object} models.ResponseError "用户名已存在"
// @Success 500 {object} models.ResponseError "服务器错误"
// @Router /api/v1/signup [post]
func SignUpHandler(c *gin.Context) {
	//1、获取参数和参数效验
	p := new(models.ParamSignUp)
	if err := c.ShouldBind(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors 类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//fmt.Println("aaaaaaaa")
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		return
	}
	//手动对请求参数进行详细的业务规则效验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数错误",
	//	})
	//	return
	//}

	//fmt.Println(p)
	//fmt.Println("??????")
	//2、业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3、返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success。。。",
	//})
	ResponseSuccess(c, nil)
}

// 登录

// LoginHandler 用户登录接口
// @Summary 用户登录接口
// @Description 登录用户账户
// @Tags 用户相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param ParamLogin body models.ParamLogin true "用户登录参数"
// @Security ApiKeyAuth
// @Success 200 {object} models.ResponseSuccess "成功响应"
// @Success 400 {object} models.ResponseError "响应错误"
// @Success 500 {object} models.ResponseError "服务器错误"
// @Router /api/v1/login [post]
func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数效验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator。ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译错误
		//})
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或密码错误",
		//})
		ResponseError(c, CodeInvalidPassword)
		return
	}
	//3、返回响应0
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "查询成功",
	//})
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //id值大于1<<53-1  int64类型的最大值1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}
