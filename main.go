package main

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "blog/hello.tmpl", gin.H{"title": "world"})
}

func main() {
	logger.Info("start cblog...")
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/health", Health)
	router.GET("/hello", Hello)
	router.Run("0.0.0.0:8089")
}
