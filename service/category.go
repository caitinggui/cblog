package service

import (
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
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
*     "data": Object           // 类别id
*    }
 */
func CreateCategory(c *gin.Context) {
	var (
		form models.Category
		err  error
	)
	mc := NewAdvancedGinContext(c)
	err = c.ShouldBind(&form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	logger.Info("find if exist ", form.Name, " in database")
	ifExist := models.CheckIsExistCategoryByName(form.Name)
	// 找到了
	if ifExist {
		mc.WebJson(e.ERR_PARAMETER_DUPLICATED, nil)
		return
	}
	err = form.Insert()
	if err != nil {
		mc.WebJson(e.ERR_SQL, err)
		return
	}
	mc.WebJson(e.SUCCESS, form)
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
	mc := NewAdvancedGinContext(c)
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
*}
 */
func UpdateCategory(c *gin.Context) {
	var (
		form models.Category
		err  error
	)
	mc := NewAdvancedGinContext(c)
	err = c.ShouldBind(&form)
	logger.Info("UpdateCategory form: ", form)
	if err != nil || form.ID == 0 {
		logger.Info("参数异常: ", err)
		mc.WebJson(e.ERR_INVALID_PARAM, err)
		return
	}

	err = form.Update()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, form)
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
*}
 */
func DeleteCategory(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	logger.Info("try to delete category: ", id)
	intId := utils.StrToUint64(id)
	if intId == 0 {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		return
	}
	err := models.DeleteCategoryById(intId)
	if mc.CheckGormErr(err) != nil {
		logger.Error("delete category error: ", err)
		return
	}
	mc.WebJson(e.SUCCESS, nil)

}
