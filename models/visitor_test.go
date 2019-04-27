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
	visitor2, _ := GetVisitorById(visitor.ID)
	if visitor2.Country == "" {
		t.Fatal("插入的ip数据没有国家: ", visitor2)
	}
}
