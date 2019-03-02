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
	cate.Delete()
}

func TestCategoryUpdate(t *testing.T) {
	cate := Category{Name: "TestCategoryInsert"}
	err := cate.Insert()
	cate, err = GetCategoryByName("TestCategoryInsert")
	if err != nil {
		t.Fatal("测试查询失败: ", err)
	}
	cate.Name = "TestCategoryUpdate"
	err = cate.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}
	cate.Delete()
}

func TestCategoryUnique(t *testing.T) {
	cate1 := Category{Name: "TestCategoryUnique"}
	err := cate1.Insert()
	if err != nil {
		t.Fatal("类型插入失败：", err)
	}
	cate2 := Category{Name: "TestCategoryUnique"}
	err = cate2.Insert()
	if err == nil {
		t.Fatal("重复类型可插入, 唯一索引未生效：")
	}
	err = cate1.Delete() // 要删除cate1, 因为cate2没有成功插入
	if err != nil {
		t.Fatal("文章类型失败：", err)
	}
	cate3 := Category{Name: "TestCategoryUnique"}
	err = cate3.Insert()
	if err != nil {
		t.Fatal("已删除类型插入失败, 唯一索引设置有问题：", err)
	}
}
