package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	_ "bluebell/docs"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)  //gin设置为发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(mode),middlewares.RateLimitMiddleware(time.Second*2 ,1))
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK,"pong")
	})


	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware())   // 应用JWT认证中间件
	{
		v1.GET("/community",controller.CommunityHandler)
		v1.GET("/community/:id",controller.CommunityDetailHandler)
		v1.GET("/post",controller.CreatePostHandler)
		v1.GET("/post/:id",controller.GetPostDetailHandler)
		v1.GET("/posts",controller.GetPostListHandler)
		//根据时间或分数获取帖子列表
		v1.GET("/posts2",controller.GetPostListHandler2)


		// 投票
		v1.POST("/vote",controller.PostVoteController)
	}

	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"404",
		})

	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r

}

