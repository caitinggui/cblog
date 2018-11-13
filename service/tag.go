package service

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	//"cblog/models"
)

func CreateTag(c *gin.Context) {
	logger.Info("create tag:")
	c.JSON(http.StatusOK, gin.H{"hello": "article"})
}
