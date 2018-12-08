package service

import (
	"net/http"

	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils/e"
)

func CreateCategory(c *gin.Context) {
	mc := Gin{C: c}
	name := c.PostForm("name")
	if name == "" {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		//c.AbortWithStatusJSON(http.StatusOK, gin.H{"errMsg": "名字不能为空"})
		return
	}
	logger.Info("find if exist ", name, " in database")
	cateNum, err := models.CountCategoryByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errMsg": "数据库异常"})
		return
	}
	if cateNum != 0 {
		c.JSON(http.StatusOK, gin.H{"errMsg": "该类型已存在"})
		return
	}
	cate := models.Category{Name: name}
	err = cate.Insert()
	// TODO 要设计怎么返回
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": "创建失败", "reason": err})
	}
	c.JSON(http.StatusOK, gin.H{"errMsg": "创建成功", "id": cate.ID})
}

func GetCategories(c *gin.Context) {
	cates, err := models.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": "数据库异常"})
	}
	c.HTML(http.StatusOK, "admin/category.html", gin.H{"Cates": cates})
}

func UpdateCategory(c *gin.Context) {
	name := c.PostForm("name")
	id := c.PostForm("id")
	if id == "" || name == "" {
		c.JSON(http.StatusOK, gin.H{"errMsg": "参数错误"})
		return
	}
	cate, err := models.GetCategoryById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errMsg": err})
		return
	}
	cate.Name = name
	err = cate.UpdateAllField()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 4001, "errMsg": err, "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": gin.H{"name": cate.Name}})
}
