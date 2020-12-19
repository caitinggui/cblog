package service

import (
	"cblog/config"
	"net/http"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"cblog/utils/V"
)

// 登陆
func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	redirect_uri := c.DefaultQuery("redirect_uri", "/")
	logger.Info("redirect_uri: ", redirect_uri)
	session := sessions.Default(c)
	if username != config.Config.Admin.Name || password != config.Config.Admin.Password {
		// 密码错误则清空session, 一定要Save，否则前端不响应.本质上是通过Set-Cookie
		//这个http header生效
		session.Clear()
		session.Save()
		logger.Info("error user try to login: ", username, " @", password)
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		logger.Info("set login user:", username)
		// 设置cookie
		session.Set(V.CurrentUser, username)
		session.Save()
		logger.Info("登陆成功，重定向到: ", redirect_uri)
		c.Redirect(http.StatusMovedPermanently, redirect_uri)
	}
}

// 登出
func Logout(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	logger.Info("-----", c.Writer.Header())
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	logger.Warnf("logout: %v %+v", err, c.Writer.Header())
	mc.Redirect("/")
}

// 获取登录页面
func GetLogin(c *gin.Context) {
	logger.Info("get login page")
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
