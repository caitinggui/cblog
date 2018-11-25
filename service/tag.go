package service

import (
	"net/http"

	"cblog/models"
	"cblog/utils"
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	tags, err := models.GetAllTags()
	logger.Info("get tags result: ", err)
	c.JSON(http.StatusOK, gin.H{"data": tags, "errMsg": err})
}

func GetTag(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("get tag by id: ", id)
	tag, err := models.GetTagById(id)
	logger.Info("get tag result: ", err)
	c.JSON(http.StatusOK, gin.H{"data": tag, "errMsg": err})
}

// 为避免字段改变影响到业务，所以不用c.PostForm来获取参数，统一用c.Bind
func CreateTag(c *gin.Context) {
	var tagNum int64
	form := &models.Tag{}
	// bind会优先json，xml，然后匹配不到才找form
	err := c.Bind(form)
	logger.Info("origin form: ", form, " err: ", err)
	if err != nil || len(form.Name) == 20 || form.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	tagNum, err = models.CountTagByName(form.Name)
	if tagNum != 0 || err != nil {
		logger.Info("exist tag: ", form.Name, " err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"errMsg": err, "data": "exist tag name"})
		return
	}

	// 防止被恶意修改id
	form.ID = utils.V.EmptyIntId
	logger.Info("create form: ", form)
	err = models.CreateTag(form)
	logger.Info("create tag result: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err, "data": form.ID})
}

func UpdateTag(c *gin.Context) {
	var form, tag models.Tag
	err := c.Bind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err != nil || form.ID == utils.V.EmptyIntId {
		c.JSON(http.StatusOK, gin.H{"errMsg": "参数错误", "data": err})
		return
	}
	tag, err = models.GetTagById(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": err})
		return
	}
	err = tag.UpdateNonzero(form)
	logger.Info("update resule: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}

func DeleteTag(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("delete tag by id: ", id)
	err := models.DeleteTagById(id)
	logger.Info("delete tag result: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}
