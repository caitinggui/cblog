package models

import (
	"testing"
	"time"
)

func TestLinkInsert(t *testing.T) {
	link := Link{
		Name: "link1",
		Url:  "http://url1.com",
	}
	err := link.Insert()
	if err != nil {
		t.Fatal("测试插入失败：", err)
	}
	defer link.Delete()
}

func TestLinkUpdate(t *testing.T) {
	var zeroTime time.Time
	link := Link{
		Name: "link1",
		Url:  "http://url1.com",
	}
	err := link.Insert()
	defer link.Delete()
	link.Name = "TestLinkUpdate"
	err = link.Update()
	if err == nil {
		t.Fatal("UpdateAt不为空时不允许更新")
	}
	link.UpdatedAt = zeroTime
	link.CreatedAt = zeroTime
	err = link.Update()
	if err != nil {
		t.Fatal("测试更新失败: ", err)
	}

	link2 := Link{
		Name: "link2",
		Url:  "http://url2.com",
	}
	link2.UpdatedAt = zeroTime
	link2.CreatedAt = zeroTime
	err = link2.Update()
	if err == nil {
		t.Fatal("空id不允许被更新")
	}
}

func TestLinkDelete(t *testing.T) {
	link1 := Link{
		Name: "link1",
		Url:  "http://url1.com",
	}
	if err := link1.Delete(); err != ERR_EMPTY_ID {
		t.Fatal("空id不允许被删除")
	}
	if err := link1.Insert(); err != nil {
		t.Fatal(err)
	}
	link1.Delete()
	link, err := GetLinkById(link1.ID)
	if link.Name != "" {
		t.Fatal("delete link fail: ", err)
	}
}
