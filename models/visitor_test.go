package models

import (
	"testing"
)

func TestVisitorInsert(t *testing.T) {
	visitor := Visitor{
		IP:      "119.123.179.47",
		Referer: "localhost",
	}
	visitor.PraseIp()
	visitor.Insert()
	defer visitor.Delete()
	visitor2, _ := GetVisitorById(visitor.ID)
	if visitor2.Country == "" {
		t.Fatal("插入的ip数据没有国家: ", visitor2)
	}

}

func TestCountVisitor(t *testing.T) {
	visitor := Visitor{
		IP:      "119.123.179.47",
		Referer: "localhost",
	}
	visitor.Insert()
	defer visitor.Delete()
	n, err := CountVisitor()
	if err != nil || n !=1 {
		t.Fatalf("TestCountVisitor failed, count: %v, err: %v", n, err)
	}

}
