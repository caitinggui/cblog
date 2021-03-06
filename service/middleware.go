package service

import (
	"net/http"
	"strings"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"cblog/config"
	"cblog/models"
	"cblog/utils"
	"cblog/utils/V"
)

var lm = utils.NewRateLimter(config.Config.PraseIp.Interval, config.Config.PraseIp.Capacity)
var ipRateLimter = map[string]*utils.RateLimiter{}

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
		var articleId uint64
		clientIp := c.ClientIP()
		url := c.Request.URL.String()
		if strings.Contains(url, "/article/") {
			path := strings.Split(url, "/")
			for k, v := range path {
				if v == "article" && k+1 < len(path) {
					articleId = utils.StrToUint64(path[k+1])
					break
				}
			}
		}
		//logger.Debug("request url: ", c.Request.URL, " client Ip: ", clientIp, " articleId: ", articleId)
		visitor := models.Visitor{
			IP:        clientIp,
			Referer:   c.Request.Referer(),
			ArticleId: articleId,
		}
		go func() {
			// 先创建，后更新，主要是为了保证createAt是正确的时间
			// insert 会自动计数
			err := visitor.Insert()
			if err != nil {
				logger.Error("Save visitor Ip failed: ", visitor, err)
			}
			if config.Config.PraseIp.IsOpen {
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
			}
		}()
	}
}

// 根据ip限制访问频率
func RateLimted() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		ipRate, ok := ipRateLimter[clientIp]
		if !ok {
			ipRate = utils.NewRateLimter(time.Minute, 2)
			ipRateLimter[clientIp] = ipRate
		}
		if !ipRate.Allow() {
			logger.Warnf("%s visite %s too quick, abort!", clientIp, c.Request.URL.String())
			c.String(http.StatusTooManyRequests, "访问太频繁，请稍后重试")
			c.Abort()
			return
		}
		c.Next()
	}
}

// 清除浏览器缓存
func AbortClientCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		c.Header("Cache-Control", "no-cache")
	}
}
