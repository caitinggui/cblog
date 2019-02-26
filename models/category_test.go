package models

import (
	"testing"
)

func TestCategoryInsert(t *testing.T) {
	cate := Category{Name: "TestCategoryInsert"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("测试插入失败：", err)
	}
}

func TestCategoryUpdate(t *testing.T) {
	cate, err := GetCategoryByName("TestCategoryInsert")
	if err != nil {
		t.Fatal("测试查询失败: ", err)
	}
	cate.Name = "TestCategoryUpdate"
	err = cate.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}
}

func TestCategoryUnique(t *testing.T) {
	cate := Category{Name: "TestCategoryUnique"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("文章插入失败：", err)
	}
	cate = Category{Name: "TestCategoryUnique"}
	err = cate.Insert()
	if err == nil {
		t.Fatal("重复文章可插入, 唯一索引未生效：")
	}
	err = cate.Delete()
	if err != nil {
		t.Fatal("文章删除失败：", err)
	}
	cate = Category{Name: "TestCategoryUnique"}
	err = cate.Insert()
	if err != nil {
		t.Fatal("已删除文章插入失败, 唯一索引设置有问题：", err)
	}
}
