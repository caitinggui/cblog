package service

import (
	"net/http"
	"strings"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils/V"
)

// 从cookie检查是否已登录
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 检查cookie
		uid := session.Get(V.CurrentUser)
		if uid == nil || uid == "" {
			logger.Warnf("user %s not authorized to visit %s", uid, c.Request.RequestURI)
			// 清空session
			session.Clear()
			session.Save()
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()

		} else {
			logger.Info("user logined:", uid)
			c.Next()
		}
	}
}

// 判断是否为管理员
// 一定要在LoginRequired后面，否则始终都是非管理员
func AdminRequierd() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid, _ := session.Get(V.CurrentUser).(string)
		logger.Info("check if ", uid, " is admin")
		isAdmin := models.IsAdminByUid(uid)
		if !isAdmin {
			logger.Warn(uid, " is not admin")
			c.String(http.StatusBadRequest, "非管理员")
			c.Abort()
			return
		}
		logger.Info(uid, " is admin")
		c.Next()
	}
}

func RecordClientIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var article_id string
		clientIp := c.ClientIP()
		url := c.Request.URL.String()
		if strings.HasPrefix(url, "/article/") {
			article_id = strings.Split(url, "/")[2]
		}
		logger.Debug("request url: ", c.Request.URL, " client Ip: ", clientIp)
		visitor := models.Visitor{
			IP:        clientIp,
			Referer:   c.Request.Referer(),
			ArticleId: article_id,
		}
		logger.Info("visitor: ", visitor)
		go func() {
			err := visitor.PraseIp()
			if err != nil {
				logger.Error("Prase visitor Ip failed: ", visitor, err)
				return
			}
			err = visitor.Insert()
			if err != nil {
				logger.Error("Save visitor Ip failed: ", visitor, err)
			}
		}()
	}
}
