package service

import (
	"net/http"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"
	//"cblog/models"
)

func GetArticle(c *gin.Context) {
	logger.Info("get article:")
	c.JSON(http.StatusOK, gin.H{"hello": "article"})
}
