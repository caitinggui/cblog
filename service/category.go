package service

import (
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils/e"
)

/**
* @api {post} /v1/category 创建博客类别
* @apiGroup 类别
* @apiVersion 0.1.0
*
* @apiParam {string} name 类别名称
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": 21            // 类别id
*    }
 */
func CreateCategory(c *gin.Context) {
	mc := Gin{C: c}
	name := c.PostForm("name")
	if name == "" || len(name) > 20 {
		logger.Warn("create category param error: ", name)
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		return
	}
	logger.Info("find if exist ", name, " in database")
	ifExist := models.CheckIsExistCategoryByName(name)
	// 找到了
	if ifExist {
		mc.WebJson(e.ERR_SQL_DATA_DUPLICATED, nil)
		return
	}
	cate := models.Category{Name: name}
	err := cate.Insert()
	if err != nil {
		mc.WebJson(e.ERR_SQL, err)
	}
	mc.WebJson(e.SUCCESS, cate.ID)
}

func GetCategories(c *gin.Context) {
	mc := Gin{C: c}
	cates, err := models.GetAllCategories()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.SuccessHtml("admin/category.html", gin.H{"Cates": cates})
}

func UpdateCategory(c *gin.Context) {
	mc := Gin{C: c}
	name := c.PostForm("name")
	id := c.PostForm("id")
	if id == "" || name == "" || len(name) > 20 {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		return
	}
	ifExist := models.CheckIsExistCategoryByName(name)
	// 找到了
	if ifExist {
		mc.WebJson(e.ERR_SQL_DATA_DUPLICATED, nil)
		return
	}
	cate, err := models.GetCategoryById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	cate.Name = name
	err = cate.UpdateAllField()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, cate)
}

func DeleteCategory(c *gin.Context) {
	mc := Gin{C: c}
	id := c.Param("id")
	logger.Info("try to delete category: ", id)
	err := models.DeleteCategoryById(id)
	if mc.CheckGormErr(err) != nil {
		logger.Error("delete category error: ", err)
		return
	}
	mc.WebJson(e.SUCCESS, nil)

}
