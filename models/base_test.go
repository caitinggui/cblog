package models

import (
	"testing"

	"cblog/config"
	"cblog/utils"
)

// 初始化数据库连接
func TestMain(m *testing.M) {
	db := InitDB()
	utils.InitUniqueId(config.Config.UniqueId.WorkerId, config.Config.UniqueId.ReserveId)
	defer db.Close()
	m.Run()
}

func TestTableName(t *testing.T) {
	table := IntIdModelWithoutDeletedAt{ID: 12}
	if table.TableName() != "IntIdModelWithoutDeleteAt" {
		t.Fatalf("table error, expectd: %v, get: %v ", "IntIdModelWithoutDeletedAt", table.TableName())
	}
}
