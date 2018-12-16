package service

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
)

// 从cookie检查是否已登录
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 检查cookie
		uid := session.Get(utils.V.CurrentUser)
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
		uid, _ := session.Get(utils.V.CurrentUser).(string)
		logger.Info("check if ", uid, " is admin")
		ifExist := models.IsAdminExistByUid(uid)
		if !ifExist {
			logger.Warn(uid, " is not admin")
			c.String(http.StatusBadRequest, "非管理员")
			c.Abort()
			return
		}
		logger.Info(uid, " is admin")
		c.Next()
	}
}
