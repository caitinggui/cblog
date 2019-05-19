package utils

import (
	//"io/ioutil"
	//"net/http"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Ip2Region struct {
	IP       string        `json:"ip"`
	Country  string        `json:"country"`
	Province string        `json:"region"`
	City     string        `json:"city"`
	Isp      string        `json:"isp"`
	Timeout  time.Duration // 请求的超时时间
}

// 用淘宝的接口去解析
func (self *Ip2Region) Name() string {
	return "淘宝"
}

func (self *Ip2Region) Url() string {
	return "http://ip.taobao.com/service/getIpInfo.php?ip="
}

func (self *Ip2Region) PraseIp() error {
	url := self.Url() + self.IP
	//client := http.Client{
	//Timeout: self.Timeout,
	//}
	//resp, err := client.Get(url)
	body, err := HttpRetryGet(url)
	if err != nil {
		logger.Error("ip2Region request ip error: ", err)
		return err
	}
	//defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//logger.Error("ip2Region prase ip error: ", err)
	//return err
	//}

	json.Unmarshal([]byte(jsoniter.Get(body, "data").ToString()), &self)
	return nil
}
