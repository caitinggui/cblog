package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"cblog/config"
	"cblog/models"
	"cblog/service"
	"cblog/utils"
)

func Health(c *gin.Context) {
	err := models.Ping()
	if err == nil {
		c.String(http.StatusOK, "success")
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func BindRoute(router *gin.Engine) {
	// 可以注册根目录，不影响router在根目录继续添加路由
	admin := router.Group("/admin")
	admin.Use(service.LoginRequired(), service.AdminRequierd())
	{
		admin.GET("", service.AdminIndex)
		admin.GET("/article", service.GetArticles)
		admin.GET("/article-edit", service.EditArticle)
		admin.GET("/article/:id", service.GetArticle)
		admin.POST("/article", service.CreateOrUpdateArticle)
		admin.PUT("/article", service.UpdateArticle)
		admin.DELETE("/article/:id", service.DeleteArticle)

		admin.GET("/category", service.GetCategories)
		admin.POST("/category", service.CreateCategory)
		admin.PUT("/category", service.UpdateCategory)
		admin.DELETE("/category/:id", service.DeleteCategory)

		admin.GET("/tag", service.GetTags)
		admin.POST("/tag", service.CreateTag)
		admin.PUT("/tag", service.UpdateTag)
		admin.DELETE("/tag/:id", service.DeleteTag)

		admin.GET("/link", service.GetLinks)
		admin.POST("/link", service.CreateLink)
		admin.PUT("/link", service.UpdateLink)
		admin.DELETE("/link/:id", service.DeleteLink)

		admin.GET("/visitor", service.GetVisitors)
	}

	blog := router.Group("/blog")
	{
		blog.GET("/article/:id", service.GetArticle)
		blog.GET("/link/:id", service.GetLink)
		blog.GET("/search", service.SearchArticles)
	}

	router.POST("/login", service.PostLogin)
	router.GET("/login", service.GetLogin)
	router.GET("/logout", service.Logout)

	router.GET("/health", Health)
	//router.GET("/hello", Hello)

	router.GET("/", service.GetArticleIndex)

	router.GET("/testadd", func(c *gin.Context) {
		models.SetCache("test", 100, 0)
		models.DumpCache(config.Config.CacheFile)
		c.String(200, "ok")
	})
	router.GET("/testget", func(c *gin.Context) {
		data, ok := models.GetCache("test")
		c.JSON(200, gin.H{"data": data, "ok": ok})
	})

	//router.GET("/admin", Index)

	//router.GET("/get", GetCategory)
}

// 平滑结束服务，避免链接突然全部断开，结束的超时时间为10s
func ListenAndServeGrace(listen string, router http.Handler) error {
	srv := http.Server{
		Addr:    listen,
		Handler: router,
	}
	go func() {
		err := srv.ListenAndServe()
		// 判断是否启动时端口占用等问题
		if err == http.ErrServerClosed {
			logger.Info("正常结束服务: ", err)
		} else {
			logger.Critical("服务结束异常: ", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("收到结束信号量: ", <-stop)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := srv.Shutdown(ctx)
	return err
}

func InitRouterAndDb() (router *gin.Engine, db *gorm.DB) {
	logger.Info("start cblog...")

	utils.InitUniqueId(config.Config.UniqueId.WorkerId, config.Config.UniqueId.ReserveId)
	logger.Info("初始化唯一id生成器成功: ", utils.GenerateId())

	db = models.InitDB()

	router = gin.New()
	var lg service.GinLog
	gin.DefaultErrorWriter = lg // 错误日志也要记录

	ginLoggerConfig := gin.LoggerConfig{
		Output: lg,
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%v |%3d| %8v %8vkB | %15s |%-5s %s | %s\n%s",
				param.TimeStamp.Format("2006/01/02-15:04:05"),
				param.StatusCode,
				param.Latency, param.BodySize/1024,
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Referer(),
				param.ErrorMessage,
			)
		},
	}
	router.Use(gin.LoggerWithConfig(ginLoggerConfig))
	router.Use(gin.Recovery())

	router.HTMLRender = service.LoadTemplates("templates")
	// router.Static相当于用router.Group为静态链接的请求建立了路由，所以/static就是路由地址"./static"就是指当前目录的static/目录
	router.Static("/static", "static")

	store := cookie.NewStore([]byte(config.Config.Secret))
	// 前端的document.cookie无法获取我们设置的session值
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   3600 * 12, // 3600*12 = 12h

	})
	router.Use(sessions.Sessions("cblog", store))
	router.Use(service.AbortClientCache(), service.RecordClientIp())
	BindRoute(router)

	err := models.InitCache(config.Config.CacheFile)
	if err != nil {
		logger.Warn("load cache file error: ", err)
	} else {
		logger.Info("load cache file success")
	}
	return
}

func main() {
	config.InitConfig()
	log, err := logger.LoggerFromConfigAsBytes(config.LoggerConfig)
	if err != nil {
		panic(err)
	}
	err = logger.ReplaceLogger(log)
	if err != nil {
		panic(err)
	}
	defer logger.Flush()

	router, db := InitRouterAndDb()
	defer db.Close()
	defer models.DumpCache(config.Config.CacheFile)

	//models.Search()
	se := models.InitIndex()
	//se := models.RestoreIndex()
	// will stop exit
	defer se.Close()

	err = ListenAndServeGrace(config.Config.Listen, router)
	logger.Errorf("stop server: %v", err)
}
