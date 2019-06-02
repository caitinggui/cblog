package service

import (
	"net/http"

	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/V"
)

func GetVisitors(c *gin.Context) {
	mc := Gin{C: c}
	page := utils.StrToUnit64(c.Query("page"))
	pageSize := utils.StrToUnit64(c.Query("pageSize"))
	if pageSize == 0 || pageSize > V.MaxPageSize {
		pageSize = V.MaxPageSize
	}
	logger.Info("get visitorys page, pageSize: ", page, pageSize)
	visitors, err := models.GetVisitors(page, pageSize)
	logger.Info("get visitors result err: ", err)
	mc.SuccessHtml("admin/visitor.html", gin.H{"Visitors": visitors})
}

func GetVisitor(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("get visitor by id: ", id)
	visitor, err := models.GetVisitorById(id)
	logger.Info("get visitor result: ", err)
	c.JSON(http.StatusOK, gin.H{"data": visitor, "errMsg": err})
}

// 在数据库添加访问者ip
func CreateVisitor(visitor *models.Visitor) {
	err := models.CreateVisitor(visitor)
	if err != nil {
		logger.Error("add visitor error: ", err)
	} else {
		logger.Info("add visitor success: ", visitor.IP)
	}
}

func UpdateVisitor(c *gin.Context) {
	var form, visitor models.Visitor
	err := c.Bind(&form)
	logger.Info("origin form: ", form, " err: ", err)
	if err != nil || form.ID == V.EmptyIntId {
		c.JSON(http.StatusOK, gin.H{"errMsg": err, "data": "参数错误"})
		return
	}
	visitor, err = models.GetVisitorById(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": err})
		return
	}
	err = visitor.UpdateNonzero(form)
	logger.Info("update resule: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}

func DeleteVisitor(c *gin.Context) {
	id := utils.StrToUnit64(c.Param("id"))
	logger.Info("delete visitor by id: ", id)
	err := models.DeleteVisitorById(id)
	logger.Info("delete visitor result: ", err)
	c.JSON(http.StatusOK, gin.H{"errMsg": err})
}
