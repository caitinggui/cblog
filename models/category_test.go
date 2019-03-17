package models

import (
	"testing"
	"time"
)

func TestCategoryInsert(t *testing.T) {
	var zeroTime time.Time
	cate := Category{Name: "TestCategoryInsert"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("测试插入失败：", err)
	}
	defer cate.Delete()
	cate2 := cate
	cate2.Name = "TestCategoryInsert2" // 不会改变cate的值
	cate2.UpdatedAt = zeroTime
	cate2.CreatedAt = zeroTime
	err = cate2.Insert()
	t.Log("重复插入的结果:", err)
	if err != ERR_EXIST_ID {
		t.Fatal("重复id插如有误: ", err)
	}
}

func TestCategoryUpdate(t *testing.T) {
	var zeroTime time.Time
	cate := Category{Name: "TestCategoryInsert"}
	err := cate.Insert()
	defer cate.Delete()
	cate, err = GetCategoryByName("TestCategoryInsert")
	if err != nil {
		t.Fatal("测试查询失败: ", err)
	}
	cate.Name = "TestCategoryUpdate"
	cate.UpdatedAt = zeroTime
	cate.CreatedAt = zeroTime
	err = cate.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}

	cate2 := Category{Name: "TestCategoryUpdate2"}
	cate2.UpdatedAt = zeroTime
	cate2.CreatedAt = zeroTime
	err = cate2.Update()
	if err == nil {
		t.Fatal("空id不允许被更新")
	}
}

func TestCategoryDelete(t *testing.T) {
	cate1 := Category{Name: "TestCategoryDelete1"}
	if err := cate1.Delete(); err != ERR_EMPTY_ID {
		t.Fatal("空id不允许被删除")
	}
	if err := cate1.Insert(); err != nil {
		t.Fatal(err)
	}
	cate1.Delete()
	cate, err := GetCategoryByName("TestCategoryDelete1")
	if cate.Name != "" {
		t.Fatal("delete category fail: ", err)
	}
}

func TestCategoryUnique(t *testing.T) {
	cate1 := Category{Name: "TestCategoryUnique"}
	err := cate1.Insert()
	if err != nil {
		t.Fatal("类型插入失败：", err)
	}
	//defer cate1.Delete()
	cate2 := Category{Name: "TestCategoryUnique"}
	err = cate2.Insert()
	if err == nil {
		t.Fatal("重复类型可插入, 唯一索引未生效：")
	}
	//defer cate2.Delete()
	err = cate1.Delete() // 要删除cate1, 因为cate2没有成功插入
	t.Log("cate1删除后: ", cate1)
	if err != nil {
		t.Fatal("文章类型失败：", err)
	}
	cate3 := Category{Name: "TestCategoryUnique"}
	err = cate3.Insert()
	if err != nil {
		t.Fatal("已删除类型插入失败, 唯一索引设置有问题：", err)
	}
	defer cate3.Delete()
}
