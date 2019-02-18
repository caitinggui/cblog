package main

import (
	"net/http"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"cblog/config"
	"cblog/models"
	"cblog/service"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "blog/hello.html", gin.H{"title": "ctg"})
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index.html", gin.H{"test": "test"})
}

func main() {
	log, err := logger.LoggerFromConfigAsBytes(config.LoggerConfig)
	if err != nil {
		panic(err)
	}
	err = logger.ReplaceLogger(log)
	if err != nil {
		panic(err)
	}
	logger.Info("start cblog...")
	defer logger.Flush()
	db := models.InitDB()
	defer db.Close()

	// 用logger 的trace记录gin框架的日志
	var lg service.GinLog
	gin.DisableConsoleColor()
	gin.DefaultWriter = lg
	gin.DefaultErrorWriter = lg
	router := gin.Default()

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

	// 可以注册根目录，不影响router在根目录继续添加路由
	admin := router.Group("/admin")
	admin.Use(service.LoginRequired())
	{
		admin.GET("/hello", Hello)
		admin.GET("", Index)
	}

	router.POST("/login", service.PostLogin)
	router.GET("/login", service.GetLogin)
	router.GET("/logout", service.Logout)

	router.GET("/health", Health)
	//router.GET("/hello", Hello)
	router.POST("/category", service.CreateCategory)
	router.GET("/category", service.GetCategories)
	router.PUT("/category", service.UpdateCategory)
	router.DELETE("/category/:id", service.DeleteCategory)

	router.GET("/tag/:id", service.GetTag)
	router.GET("/tag", service.GetTags)
	router.POST("/tag", service.CreateTag)
	router.PUT("/tag", service.UpdateTag)
	router.DELETE("/tag/:id", service.DeleteTag)

	router.GET("/link/:id", service.GetLink)
	router.GET("/link", service.GetLinks)
	router.POST("/link", service.CreateLink)
	router.PUT("/link", service.UpdateLink)
	router.DELETE("/link/:id", service.DeleteLink)

	router.GET("/visitor", service.GetVisitors)

	router.GET("/", Hello)

	err = models.InitCache(config.Config.CacheFile)
	defer func() {
		logger.Info("start dump cache")
		models.DumpCache(config.Config.CacheFile)
	}()
	if err != nil {
		logger.Warn("load cache file error: ", err)
	} else {
		logger.Info("load cache file success")
	}
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
	err = router.Run("0.0.0.0:8089")
	logger.Errorf("stop server: %2v", err)
}
