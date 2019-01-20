package service

import (
	logger "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils/e"
)

/**
* @api {post} /v1/category 创建博客类别
* @apiGroup Category
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

/**
* @api {get} /v1/category 获取所有博客类别
* @apiGroup Category
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": [object, object]
*    }
 */
func GetCategories(c *gin.Context) {
	mc := Gin{C: c}
	cates, err := models.GetAllCategories()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.SuccessHtml("admin/category.html", gin.H{"Cates": cates})
}

/**
* @api {put} /v1/category 更新某个博客类别
* @apiGroup Category
* @apiVersion 0.1.0

* @apiParam {string} name 类别名称
* @apiParam {string} id 类别id
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": object   // object 为类别详情对象
*
}
*/
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

/**
* @api {delete} /v1/category/:id 删除某个博客类别
* @apiGroup Category
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": null
*
}
*/
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
