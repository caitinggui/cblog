package service

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/utils/e"
)

type Gin struct {
	C *gin.Context
}

// 返回html code 为200的json response
func (self *Gin) WebJson(code int, data interface{}) {
	logger.Info("json response: ", code, data)
	self.C.JSON(http.StatusOK, gin.H{"errCode": code, "ErrMsg": e.GetMsg(code), "data": data})
}

// 返回html code为200的html response
func (self *Gin) SuccessHtml(templateName string, data interface{}) {
	logger.Info("html response data: ", data)
	self.C.HTML(http.StatusOK, templateName, data)
}

// 返回html code为400的html response
func (self *Gin) BadHtml(templateName string, data interface{}) {
	logger.Warn("bad html response data", data)
	self.C.HTML(http.StatusBadRequest, templateName, data)
}
