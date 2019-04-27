package service

import (
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/e"
)

func GetLinks(c *gin.Context) {
	mc := Gin{C: c}
	links, err := models.GetAllLinks()
	logger.Info("get links result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, links)
}

func GetLink(c *gin.Context) {
	mc := Gin{C: c}
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("get link by id: ", id)
	link, err := models.GetLinkById(id)
	logger.Info("get link result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, link)
}

// 为避免字段改变影响到业务，所以不用c.PostForm来获取参数，统一用c.ShouldBind
func CreateLink(c *gin.Context) {
	mc := Gin{C: c}
	form := models.Link{}
	// bind会优先json，xml，然后匹配不到才找form
	err := c.ShouldBind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if mc.CheckBindErr(err) != nil {
		return
	}

	err = form.Insert()
	logger.Info("create link result: ", err)
	if err != nil {
		logger.Error("创建link失败: ", err)
		mc.WebJson(e.ERR_SQL, err)
		return
	}
	mc.WebJson(e.SUCCESS, form)
}

func UpdateLink(c *gin.Context) {
	var form models.Link
	mc := Gin{C: c}
	err := c.ShouldBind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if mc.CheckBindErr(err) != nil {
		return
	}
	err = form.Update()
	logger.Info("update link result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, form)
}

func DeleteLink(c *gin.Context) {
	mc := Gin{C: c}
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("delete link by id: ", id)
	err := models.DeleteLinkById(id)
	logger.Info("delete link result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, id)
}
