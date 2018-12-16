package config

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type mysql struct {
	Server  string
	MaxIdle int
	MaxOpen int
	MaxLife time.Duration
	LogMode bool
}

type config struct {
	Listen    string
	Mysql     mysql
	CacheFile string
	Secret    string
}

var Config config
var LoggerConfig []byte

// 根据GIN_MODE环境变量，获取相应的logger的配置和应用配置
func init() {
	// 处理logger
	var (
		content []byte
		err     error
	)
	if os.Getenv("GIN_MODE") == gin.ReleaseMode {
		// 打包会合并到可执行文件中，所以目录以可执行文件为基准
		content, err = ioutil.ReadFile("config/pro/log.xml")

	} else {
		content, err = ioutil.ReadFile("config/dev/log.xml")

	}
	if err != nil {
		panic(err)

	}
	LoggerConfig = content
	// 处理配置
	if os.Getenv("GIN_MODE") == gin.ReleaseMode {
		// 打包会合并到可执行文件中，所以目录以可执行文件为基准
		content, err = ioutil.ReadFile("config/pro/config.yaml")

	} else {
		content, err = ioutil.ReadFile("config/dev/config.yaml")

	}
	if err != nil {
		panic(err)

	}
	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		panic(err)

	}
	return

}
