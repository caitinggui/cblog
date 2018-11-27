package main

import (
	"net/http"

	logger "github.com/cihub/seelog"
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
	c.HTML(http.StatusOK, "blog/hello.tmpl", gin.H{"title": "ctg"})
}

// 登陆
func Login(c *gin.Context) {
	username := c.PostForm("username")
	passwd := c.PostForm("passwd")
	session := sessions.Default(c)
	if username != "test" || passwd != "test" {
		// 密码错误则清空session, 一定要Save，否则前端不响应.本质上是通过Set-Cookie
		//这个http header生效
		session.Clear()
		session.Save()
		logger.Info("error user try to login: ", username, "@", passwd)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号或者密码错误"})

	} else {
		logger.Info("set login user:", username)
		// 设置cookie
		session.Set("user", username)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"msg": "login success"})
	}
}

// 登出
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"msg": "logout success"})
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 检查cookie
		user := session.Get("user")
		if user == nil || user == "" {
			logger.Warnf("User %s not authorized to visit %s", user, c.Request.RequestURI)
			// 清空session
			session.Clear()
			session.Save()
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Login required"})
			c.Abort()

		} else {
			logger.Info("user logined:", user)
			c.Next()
		}
	}
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
	db := models.InitDB()
	defer db.Close()
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	store := cookie.NewStore([]byte("I'am a very secert string"))
	// 前端的document.cookie无法获取我们设置的session值
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   3600 * 12, // 3600*12 = 12h

	})
	router.Use(sessions.Sessions("cblog", store))

	// 可以注册根目录，不影响router在根目录继续添加路由
	admin := router.Group("/admin")
	admin.Use(LoginRequired())
	{
		admin.GET("/hello", Hello)
	}

	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.GET("/health", Health)
	//router.GET("/hello", Hello)
	router.POST("/category", service.CreateCategory)
	router.GET("/category", service.GetCategories)
	router.PUT("/category", service.UpdateCategory)

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

	//router.GET("/get", GetCategory)
	router.Run("0.0.0.0:8089")
}
