package service

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

// 获取登录页面
func GetLogin(c *gin.Context) {
	logger.Info("get login page")
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
