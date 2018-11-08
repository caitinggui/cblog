package main

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"cblog/models"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "blog/hello.tmpl", gin.H{"title": models.MyName()})
}

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
		admin.GET("/blog", Hello)
	}

	router.POST("/login", Login)
	router.GET("/health", Health)
	router.GET("/hello", Hello)
	router.Run("0.0.0.0:8089")
}
