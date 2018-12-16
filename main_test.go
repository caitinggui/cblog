package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/health", Health)

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("health status code error: ", w.Code)
	}
}

func BenchmarkHealth(b *testing.B) {
	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能
	b.StopTimer() //调用该函数停止压力测试的时间计数

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/health", Health)

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
