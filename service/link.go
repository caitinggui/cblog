package service

import (
	"errors"
	"net/http"

	"cblog/models"
	"cblog/utils"
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func GetLinks(c *gin.Context) {
	links, err := models.GetAllLinks()
	logger.Info("get links result: ", err)
	c.JSON(http.StatusOK, gin.H{"data": links, "errMsg": err})
}

func GetLink(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("get link by id: ", id)
	link, err := models.GetLinkById(id)
	logger.Info("get link result: ", err)
	c.JSON(http.StatusOK, gin.H{"data": link, "errMsg": err})
}

// 为避免字段改变影响到业务，所以不用c.PostForm来获取参数，统一用c.Bind
func CreateLink(c *gin.Context) {
	form := models.Link{}
	// bind会优先json，xml，然后匹配不到才找form
	err := c.Bind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err = checkLinkForm(&form); err != nil {
		logger.Info("bad request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"errMsg": err.Error()})
		return
	}

	// 防止被恶意修改id
	form.ID = utils.V.EmptyIntId
	logger.Info("create form: ", form)
	err = models.CreateLink(&form)
	logger.Info("create link result: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err, "data": form.ID})
}

func UpdateLink(c *gin.Context) {
	var form, link models.Link
	err := c.Bind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err = checkLinkForm(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errMsg": err})
		return
	}
	if form.ID == utils.V.EmptyIntId {
		c.JSON(http.StatusOK, gin.H{"errMsg": err, "data": "empty id"})
		return
	}
	link, err = models.GetLinkById(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": err})
		return
	}
	err = link.UpdateNonzero(form)
	logger.Info("update resule: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}

func DeleteLink(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("delete link by id: ", id)
	err := models.DeleteLinkById(id)
	logger.Info("delete link result: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}

func checkLinkForm(form *models.Link) error {
	if len(form.Name) > 128 || len(form.Url) > 512 || len(form.Desc) > 512 {
		return errors.New("parameter too long")
	}
	if form.Url == "" {
		return errors.New("empty url")
	}
	return nil
}
