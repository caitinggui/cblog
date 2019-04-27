package config

import (
	"io/ioutil"
	"os"
	"time"

	logger "github.com/caitinggui/seelog"
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
	UniqueId  uniqueId
	CacheFile string
	Secret    string
}

type uniqueId struct {
	WorkerId  uint16
	ReserveId uint8
}

var Config config
var LoggerConfig []byte

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 根据GIN_MODE环境变量，获取相应的logger的配置和应用配置, 也可以用-configPath指定绝对目录
func init() {
	var (
		content    []byte
		err        error
		configPath string
	)

	// 先定义，后Parse
	//flag.StringVar(&configPath, "configPath", "", "set configuration file path")
	//flag.Parse()
	configPath = os.Getenv("CBLOG_CONFIG_PATH")

	// 如果指定了目录，就直接读取目录下配置，否则根据环境读取相应配置*/
	if configPath != "" {

	} else if os.Getenv("GIN_MODE") == gin.ReleaseMode {
		// 打包会合并到可执行文件中，所以目录以可执行文件为基准
		configPath = "config/pro/"
	} else {
		configPath = "config/dev/"
	}
	logger.Info("configPath: ", configPath)
	LoggerConfig, err = ioutil.ReadFile(configPath + "log.xml")
	PanicErr(err)
	content, err = ioutil.ReadFile(configPath + "config.yaml")
	PanicErr(err)

	err = yaml.Unmarshal(content, &Config)
	PanicErr(err)
	logger.Info("config.yaml: ", Config)

	//password := utils.UidDecrypt(Config.Mysql.Password)
	//if password == "" {
	//panic("Db password error")
	//}
	//Config.Mysql.Server = strings.Replace(Config.Mysql.Server, "{password}", password, -1)
	return
}
