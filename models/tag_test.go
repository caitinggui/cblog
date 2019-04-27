package models

import (
	"testing"
	"time"
)

func TestTagInsert(t *testing.T) {
	var zeroTime time.Time
	tag := Tag{Name: "TestTagInsert"}
	err := tag.Insert()
	if err != nil {
		t.Fatal("测试插入失败：", err)
	}
	defer tag.Delete()
	tag2 := tag
	tag2.Name = "TestTagInsert2" // 不会改变tag的值
	tag2.UpdatedAt = zeroTime
	tag2.CreatedAt = zeroTime
	err = tag2.Insert()
	t.Log("重复插入的结果:", err)
	if err != ERR_EXIST_ID {
		t.Fatal("重复id插如有误: ", err)
	}
}

func TestTagUpdate(t *testing.T) {
	var zeroTime time.Time
	tag := Tag{Name: "TestTagInsert"}
	err := tag.Insert()
	defer tag.Delete()
	tag, err = GetTagByName("TestTagInsert")
	if err != nil {
		t.Fatal("测试查询失败: ", err)
	}
	tag.Name = "TestTagUpdate"
	tag.UpdatedAt = zeroTime
	tag.CreatedAt = zeroTime
	err = tag.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}

	tag2 := Tag{Name: "TestTagUpdate2"}
	tag2.UpdatedAt = zeroTime
	tag2.CreatedAt = zeroTime
	err = tag2.Update()
	if err == nil {
		t.Fatal("空id不允许被更新")
	}
}

func TestTagDelete(t *testing.T) {
	tag1 := Tag{Name: "TestTagDelete1"}
	if err := tag1.Delete(); err != ERR_EMPTY_ID {
		t.Fatal("空id不允许被删除")
	}
	if err := tag1.Insert(); err != nil {
		t.Fatal(err)
	}
	tag1.Delete()
	tag, err := GetTagByName("TestTagDelete1")
	if tag.Name != "" {
		t.Fatal("delete Tag fail: ", err)
	}
}

func TestTagUnique(t *testing.T) {
	tag1 := Tag{Name: "TestTagUnique"}
	err := tag1.Insert()
	if err != nil {
		t.Fatal("类型插入失败：", err)
	}
	//defer tag1.Delete()
	tag2 := Tag{Name: "TestTagUnique"}
	err = tag2.Insert()
	if err == nil {
		t.Fatal("重复类型可插入, 唯一索引未生效：")
	}
	//defer tag2.Delete()
	err = tag1.Delete() // 要删除tag1, 因为tag2没有成功插入
	t.Log("tag1删除后: ", tag1)
	if err != nil {
		t.Fatal("文章类型失败：", err)
	}
	tag3 := Tag{Name: "TestTagUnique"}
	err = tag3.Insert()
	if err != nil {
		t.Fatal("已删除类型插入失败, 唯一索引设置有问题：", err)
	}
	defer tag3.Delete()
}
