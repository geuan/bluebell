package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title bluebell
// @version 1.0
// @description bluebell
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path

func main() {
	//1、加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return // 如果返回错误则直接结束执行
	}
	/*
	调试打印配置文件中的参数,并判断是否读物到了配置文件
	fmt.Println(settings.Conf)
	fmt.Println(settings.Conf.LogConfig == nil)
	*/
	//2、初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync() //把缓存区的日志追加到日志文件中
	zap.L().Debug("logger init success...")

	//3、初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()
	//4、初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()


	if err := snowflake.Init(settings.Conf.StartTime,settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed,err:%v\n",err)
		return
	}

	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed,err:%v\n",err)
		return
	}

	//5、注册路由
	r := routes.Setup(settings.Conf.Mode)

	//6、启用服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()

	//等待终端信号来优雅地关闭服务器，为关闭服务器操作设置一个5s的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	//kill 默认会发送 syscall。SIGTERM 信号
	// kill -2 发送 syscall.SIGINT信号，我们常用的Ctrl+C就是触发SIGINY信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT 或 syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会被阻塞
	<-quit                                               //阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server...")
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5秒内优雅关闭服务(将未处理完的请求处理完再关闭服务)，超过5s就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
