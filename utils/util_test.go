package utils

import (
	"testing"
)

func TestUidEncypt(t *testing.T) {
	data := "123456"
	str1, err := UidEncrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	if str1 != "40d868bb2517dc94" {
		t.Fatal("加密后值不对: ", str1, " 期望值: 40d868bb2517dc94")
	}
	t.Log("加密后: ", str1)
}

func TestUidDecrypt(t *testing.T) {
	data := "40d868bb2517dc94"
	str1 := UidDecrypt(data)
	if str1 != "123456" {
		t.Fatal("解密后值不对: ", str1, " 期望值: 123456")
	}
	t.Log("解密后: ", str1)
}
