package service

import (
	"net/http"

	"cblog/models"
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func CreateTag(c *gin.Context) {
	var form models.Tag
	logger.Info("create tag:")
	logger.Info("postform", c.Request)
	// bind会优先json，xml，然后匹配不到才找form
	err := c.Bind(&form)
	logger.Info("form: ", form)
	if err != nil {
		logger.Warn(err)
		return
	}
	models.DB.Save(&form)
	c.JSON(http.StatusOK, gin.H{"hello": "article"})
}
