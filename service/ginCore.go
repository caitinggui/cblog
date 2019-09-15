package service

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"cblog/utils"
	"cblog/utils/e"
)

type Gin struct {
	C *gin.Context
	Res gin.H
}

// new a Gin struct
func NewAdvancedGinContext(c *gin.Context) *Gin {
	mc := Gin{C: c, Res: gin.H{}}
	return &mc
}

// 临时重定向
func (self *Gin) Redirect(url string) {
	self.C.Redirect(http.StatusMovedPermanently, url)
}

// 返回html code 为200的json response
func (self *Gin) WebJson(code int, data interface{}) {
	logger.Debug("json response: ", code, e.GetMsg(code), data)
	self.C.JSON(http.StatusOK, gin.H{"errCode": code, "errMsg": e.GetMsg(code), "data": data})
}

// 返回html code为200的html response
func (self *Gin) SuccessHtml(templateName string, data interface{}) {
	//logger.Debug("html response data: ", data)
	self.C.HTML(http.StatusOK, templateName, data)
}

// 返回html code为400的html response
func (self *Gin) BadHtml(templateName string, data interface{}) {
	logger.Warn("bad html response data", data)
	self.C.HTML(http.StatusBadRequest, templateName, data)
}

// ErrRecordNotFound时返回无数据的response
func (self *Gin) CheckGormErr(err error) error {
	if err == nil {
		return nil
	}
	if gorm.IsRecordNotFoundError(err) {
		logger.Warn("数据库无此数据")
		self.WebJson(e.ERR_NO_DATA, nil)
		return err
	}
	logger.Error("sql error: ", err.Error())
	self.WebJson(e.ERR_SQL, err.Error())
	return err
}

func (self *Gin) CheckBindErr(err error) error {
	if err != nil {
		logger.Info("解析参数失败:", err)
		self.WebJson(e.ERR_INVALID_PARAM, err)
		return err
	}
	return nil
}

// load template directory
func loadTemplateDir(r multitemplate.Renderer, funcMap template.FuncMap, templatesDir string, module string) {
	logger.Infof("loac template module: %s/%s", templatesDir, module)
	moduleBase, err := filepath.Glob(fmt.Sprint(templatesDir, "/layouts/", module, "-base.html"))
	utils.PanicErr(err)
	moduleHtmls, err := filepath.Glob(fmt.Sprint(templatesDir , "/", module, "/*.html"))
	utils.PanicErr(err)
	for _, adminHtml := range moduleHtmls {
		layoutCopy := make([]string, len(moduleBase))
		copy(layoutCopy, moduleBase)
		files := append(layoutCopy, adminHtml)
		relativePath, err := filepath.Rel(templatesDir, adminHtml)
		utils.PanicErr(err)
		logger.Info("template name: ", relativePath)
		r.AddFromFilesFuncs(relativePath, funcMap, files...)
	}
}

// 加载模板
func LoadTemplates(templatesDir string) multitemplate.Renderer {
	//var relativePath string
	r := multitemplate.NewRenderer()
	// 定义模板函数
	funcMap := template.FuncMap{
		"FormatAsDate": FormatAsDate,
		"ToStr":        utils.ToStr,
		"Split":        SplitSring,
	}

	loadTemplateDir(r, funcMap, templatesDir, "admin")
	loadTemplateDir(r, funcMap, templatesDir, "blog")
	// login.html模板不使用base.html渲染
	r.AddFromFilesFuncs("login.html", funcMap, templatesDir+"/login.html")

	r.AddFromFilesFuncs("blog/hello.html", funcMap, templatesDir+"/blog/hello.html")
	return r

}

// html模板函数，时间转为字符串格式
func FormatAsDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// template function
func SplitSring(s string, sep string) []string {
	return strings.Split(s, sep)
}

// 用来logger记录gin框架的log
type GinLog struct{}

func (self GinLog) Write(p []byte) (n int, err error) {
	logger.Trace(string(p))
	return len(p), nil
}
