package routes

import (
	"net/http"
	"web_app/controllers"
	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs
	"web_app/logger"
	"web_app/middlewares"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.SetMode(gin.ReleaseMode) //gin设置成发布模式

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//登录业务路由`
	v1.POST("/login", controllers.LoginHandler)
	//注册
	v1.POST("/signup", controllers.SignUpHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.CetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)
		//根据时间或分数获取帖子列表
		v1.GET("/posts2", controllers.GetPostListHandler2)

		//投票
		v1.POST("/vote", controllers.PostVoteController)
	}
	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//如果是登录的用户,判断请求头中是否有有效的JWT
	//isLogin := true
	//c.Request.Header.Get("Authorization")
	//if isLogin {
	//	c.String(http.StatusOK, "pong")
	//} else {
	//否则就直接返回请登录
	//	c.String(http.StatusOK, "请登录")

	//})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	r.Run(":8080")
	return r
}
