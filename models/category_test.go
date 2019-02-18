package models

import (
	"testing"
)

func TestCategoryInsert(t *testing.T) {
	cate := Category{Name: "test"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("测试插入失败：", err)
	}
}

func TestCategoryUpdate(t *testing.T) {
	cate, err := GetCategoryByName("test")
	if err != nil {
		t.Fatal("测试查询失败: ", err)
	}
	cate.Name = "test2"
	err = cate.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}
}
