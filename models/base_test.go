package models

import (
	"testing"
)

// 初始化数据库连接
func TestMain(m *testing.M) {
	db := InitDB()
	defer db.Close()
	m.Run()
}
