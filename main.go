package main

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func main() {
	logger.Info("start cblog...")
	router := gin.Default()
	router.GET("/health", Health)
	router.Run("0.0.0.0:8089")
}
