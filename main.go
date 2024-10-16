package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controllers"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

//gin框架

// @title web_app 项目接口文档
// @version 1.0
// @description Go web开发web_app项目

// @contact.name 程雨涵

// @host 127.0.0.1:8080
// @BasePath /api/v1
func main() {
	//1、加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v", err)
		return
	}
	//2、初始化日志
	//fmt.Println(settings.Conf.Mode)
	//fmt.Println(settings.Conf.LogConfig.Mode)
	if err := logger.Init(settings.Conf.LogConfig.Mode); err != nil {
		fmt.Printf("init logger failed,err:%v", err)
		return
	}
	defer zap.L().Sync()
	//zap.L().Debug("logger init success...")
	//3、初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed,err:%v", err)
		return
	}
	defer mysql.Close()
	//fmt.Println(settings.Conf.StartTime)
	//fmt.Println(settings.Conf.MachineID)
	if err := snowflake.Init(viper.GetString("start_time"), int64(viper.GetInt("machine_id"))); err != nil {
		fmt.Printf("init snowflake failed,err:%v\n", err)
		return
	}
	//4、初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed,err:%v", err)
		return
	}
	defer redis.Close()

	//初始化gin框架内置的效验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed,err:%v\n", err)
		return
	}

	//5、注册路由
	r := routes.Setup(settings.Conf.LogConfig.Mode)
	//6、启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen:%s\n", err)
		}
	}()
	//等待中断信号来优雅的关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	//kill默认会发送syscall.SIGTERM信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
