package main

import (
	"net/http"
	"path/filepath"

	logger "github.com/cihub/seelog"
	"github.com/gin-contrib/multitemplate"
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

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	var relativePath string
	r := multitemplate.NewRenderer()

	adminBase, err := filepath.Glob(templatesDir + "/layouts/admin-base.html")
	if err != nil {
		panic(err.Error())

	}

	adminHtmls, err := filepath.Glob(templatesDir + "/admin/*.html")
	if err != nil {
		panic(err.Error())

	}
	logger.Info("adminBase: ", adminBase)
	logger.Info("adminHtmls: ", adminHtmls)

	// Generate our templates map from our adminBase/ and admin/ directories
	for _, adminHtml := range adminHtmls {
		layoutCopy := make([]string, len(adminBase))
		copy(layoutCopy, adminBase)
		files := append(layoutCopy, adminHtml)
		relativePath, err = filepath.Rel(templatesDir, adminHtml)
		if err != nil {
			panic(err)
		}
		logger.Info("template name: ", relativePath)
		r.AddFromFiles(relativePath, files...)

	}
	// login.html模板不使用base.html渲染
	r.AddFromFiles("login.html", templatesDir+"/login.html")

	r.AddFromFiles("blog/hello.html", templatesDir+"/blog/hello.html")
	return r

}

func GetLogin(c *gin.Context) {
	logger.Info("get login page")
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// 登陆
func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	session := sessions.Default(c)
	if username != "test" || password != "test" {
		// 密码错误则清空session, 一定要Save，否则前端不响应.本质上是通过Set-Cookie
		//这个http header生效
		session.Clear()
		session.Save()
		logger.Info("error user try to login: ", username, " @", password)
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		logger.Info("set login user:", username)
		// 设置cookie
		session.Set("user", username)
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

// 登出
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/")
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
			c.Redirect(http.StatusMovedPermanently, "/login")
			//c.JSON(http.StatusUnauthorized, gin.H{"msg": "Login required"})
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
	defer logger.Flush()
	db := models.InitDB()
	defer db.Close()
	router := gin.Default()

	router.HTMLRender = LoadTemplates("templates")
	// router.Static相当于用router.Group为静态链接的请求建立了路由，所以/static就是路由地址"./static"就是指当前目录的static/目录
	router.Static("/static", "static")

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
		admin.GET("/", Index)
	}

	router.POST("/login", PostLogin)
	router.GET("/login", GetLogin)
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

	router.GET("/", Hello)

	err = models.InitCache("/home/ctg/go/src/cblog/cache.dump")
	defer func() {
		logger.Info("start dump cache")
		err = models.DumpCache("/home/ctg/go/src/cblog/cache.dump")
		logger.Info("dump cache result: ", err)
	}()
	if err != nil {
		logger.Warn("load cache file error: ", err)
	} else {
		logger.Info("load cache file success")
	}
	router.GET("/testadd", func(c *gin.Context) {
		models.SetCache("test", 100, 0)
		models.DumpCache("cache.dump")
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
