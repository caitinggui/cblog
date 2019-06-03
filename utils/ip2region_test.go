package utils

import (
	"testing"
)

func TestIp2Region(t *testing.T) {
	ip2Region := Ip2Region{
		IP: "119.123.179.47",
	}
	err := ip2Region.PraseIp()
	if err != nil {
		t.Fatal("解析ip出错: ", err)
	}
	t.Log("ip2Region: ", ip2Region)
	ip2Region = Ip2Region{IP: "119.123.179.47"}
	err = ip2Region.PraseIp()
	if err != nil {
		t.Fatal("解析ip出错: ", err)
	}
	t.Log("ip2Region: ", ip2Region)
}
