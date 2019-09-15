package service

import (
	"net/http"
	"strings"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"cblog/config"
	"cblog/models"
	"cblog/utils"
	"cblog/utils/V"
)

var lm = utils.NewRateLimter(config.Config.PraseIp.Interval, config.Config.PraseIp.Capacity)

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
			logger.Info(c.Request.RequestURI, " 未登录")
			c.Redirect(http.StatusMovedPermanently, "/login?redirect_uri="+c.Request.RequestURI)
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
			Logout(c)
			c.Abort()
			return
		}
		logger.Info(uid, " is admin")
		c.Next() // next并不是非得调用，在next之前的为handle处理的步骤，在next之后的就是hander处理完之后，middleware可以继续处理
	}
}

// 记录访问者ip
func RecordClientIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var article_id string
		// update visitor sum
		if _, err := models.IncrUint(V.VisitorSum); err != nil {
			logger.Warnf("there doesn't exist %s in cache", V.VisitorSum)
			visitorSum, err := models.CountVisitor()
			if err != nil {
				logger.Error("count visitor failed: ", err)
			} else {
				models.SetCache(V.VisitorSum, visitorSum, 0)
			}
		}
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
			// 先创建，后更新，主要是为了保证createAt是正确的时间
			err := visitor.Insert()
			if err != nil {
				logger.Error("Save visitor Ip failed: ", visitor, err)
			}
			lm.Wait()
			err = visitor.PraseIp()
			if err != nil {
				logger.Error("Prase visitor Ip failed: ", visitor, err)
				return
			}
			err = visitor.Update()
			if err != nil {
				logger.Error("Update visitor Ip failed: ", visitor, err)
			}
		}()
	}
}

// 紧张浏览器缓存
func AbortClientCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Debug("设置Cache-Control: no-cache")
		c.Header("Cache-Control", "no-cache")
	}
}
