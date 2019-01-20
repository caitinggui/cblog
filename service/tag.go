package service

import (
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/e"
)

func GetTags(c *gin.Context) {
	mc := Gin{C: c}
	tags, err := models.GetAllTags()
	logger.Info("get tags result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, tags)
}

func GetTag(c *gin.Context) {
	mc := Gin{C: c}
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("get tag by id: ", id)
	tag, err := models.GetTagById(id)
	logger.Info("get tag result: ", err)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, tag)
}

// 为避免字段改变影响到业务，所以不用c.PostForm来获取参数，统一用c.Bind
func CreateTag(c *gin.Context) {
	var tagNum int64
	mc := Gin{C: c}
	form := &models.Tag{}
	// bind会优先json，xml，然后匹配不到才找form
	err := c.Bind(form)
	logger.Info("origin form: ", form, " err: ", err)
	// 防止被恶意修改id
	if err != nil || len(form.Name) > 20 || form.Name == "" || form.ID != 0 {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)

		return
	}
	tagNum, err = models.CountTagByName(form.Name)
	if tagNum != 0 || err != nil {
		logger.Info("exist tag: ", form.Name, " err: ", err)
		mc.WebJson(e.ERR_INVALID_PARAM, "exist tag name")
		return
	}

	logger.Info("create form: ", form)
	err = models.CreateTag(form)
	logger.Info("create tag result: ", err)
	mc.WebJson(e.SUCCESS, form)
}

func UpdateTag(c *gin.Context) {
	var form, tag models.Tag
	mc := Gin{C: c}
	err := c.Bind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err != nil || form.ID == utils.V.EmptyIntId {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		return
	}
	tag, err = models.GetTagById(form.ID)
	if mc.CheckGormErr(err) != nil {
		logger.Warn("get tag failed: ", tag)
		return
	}
	err = tag.UpdateNoneZero(form)
	if mc.CheckGormErr(err) != nil {
		logger.Warn("update tag failed: ", err)
		return
	}
	logger.Info("update result: ", err)
	mc.WebJson(e.SUCCESS, form)
}

func DeleteTag(c *gin.Context) {
	mc := Gin{C: c}
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("delete tag by id: ", id)
	err := models.DeleteTagById(id)
	if mc.CheckGormErr(err) != nil {
		logger.Warn("delete tag failed: ", err)
		return
	}
	logger.Info("delete tag result: ", err)
	mc.WebJson(e.SUCCESS, nil)
}
