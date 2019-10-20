package config

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"

	"cblog/utils"
)

type mysql struct {
	Server   string
	Password string
	MaxIdle  int
	MaxOpen  int
	MaxLife  time.Duration
	LogMode  bool
}

type searcher struct {
	IsTest       bool
	DictoryPath  string
	StopWordPath string
}

type praseIp struct {
	Interval time.Duration
	Capacity int
	IsOpen   bool
}

type config struct {
	Listen    string
	Mysql     mysql
	UniqueId  uniqueId
	CacheFile string
	Secret    string
	Searcher  searcher
	PraseIp   praseIp
}

type uniqueId struct {
	WorkerId  uint16
	ReserveId uint8
}

var Config config
var LoggerConfig []byte

// 定时更新配置(5min)，但是数据库链接等更新了也没法使用
func UpdateConfigFrequency(configPath string) {
	tick := time.Tick(5 * time.Minute)
	for {
		select {
		case <-tick:
			content, err := ioutil.ReadFile(configPath + "config.yaml")
			if err != nil {
				logger.Warn("定时读取配置文件失败: ", err)

			}
			err = yaml.Unmarshal(content, &Config)
			if err != nil {
				logger.Warn("定时解析配置失败: ", err)

			}
			logger.Info("定时刷新配置成功")

		}

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
	utils.PanicErr(err)
	content, err = ioutil.ReadFile(configPath + "config.yaml")
	utils.PanicErr(err)

	err = yaml.Unmarshal(content, &Config)
	utils.PanicErr(err)
	logger.Info("config.yaml: ", Config)

	password := utils.UidDecrypt(Config.Mysql.Password)
	if password == "" {
		panic("Db password error")
	}
	Config.Mysql.Server = strings.Replace(Config.Mysql.Server, "{password}", password, -1)
	go UpdateConfigFrequency(configPath)
	return
}
