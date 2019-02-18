package service

import (
	"errors"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/V"
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
	if err = checkLinkForm(&form); err != nil {
		logger.Info("bad request: ", err)
		mc.WebJson(e.ERR_INVALID_PARAM, err)
		return
	}

	logger.Info("create form: ", form)
	err = models.CreateLink(&form)
	logger.Info("create link result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, form)
}

func UpdateLink(c *gin.Context) {
	var form, link models.Link
	mc := Gin{C: c}
	err := c.ShouldBind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err != nil {
		mc.WebJson(e.ERR_INVALID_PARAM, err)
		return
	}
	link, err = models.GetLinkById(form.ID)
	if mc.CheckGormErr(err) != nil {
		return
	}
	err = link.UpdateNoneZero(form)
	logger.Info("update resule: ", err)
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

// 检查Link的表单,要和models对应
func checkLinkForm(form *models.Link) error {
	if len(form.Name) > 128 || len(form.Url) > 512 || len(form.Desc) > 512 {
		return errors.New("parameter too long")
	}
	if form.Url == "" {
		return errors.New("empty url")
	}
	if form.ID != V.EmptyIntId {
		return errors.New("id error")
	}
	return nil
}
