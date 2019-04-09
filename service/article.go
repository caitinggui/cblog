package service

import (
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/e"
)

/**
* @api {get} /v1/article/:id 获取单篇文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*}
 */
func GetArticle(c *gin.Context) {
	mc := Gin{C: c}
	id := c.Param("id")
	logger.Info("get article : ", id)
	article, err := models.GetArticleById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, article)
}

/**
* @api {post} /v1/article 创建文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiParam {string} title 文章标题
* @apiParam {string} body 文章正文
* @apiParam {int} status -1表示未发布，1表示已发布，默认为-1
* @apiParam {string} abstract 文章摘要，为空就为body的前128字
* @apiParam {int} topped -1表示不置顶，1表示置顶，默认为-1
* @apiParam {int} category_id 类别id
* @apiParam {list} tags_id 标签id，int型的列表
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*}
 */
func PostArticle(c *gin.Context) {
	var (
		err  error
		form models.Article
	)
	mc := Gin{C: c}
	err = c.ShouldBind(&form)
	logger.Info("form: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	err = form.Insert()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, form)
}

/**
* @api {post} /v1/article 修改文章
* apiDescription id不对时不报错
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiParam {int} id 文章id
* @apiParam {string} title 文章标题
* @apiParam {string} body 文章正文
* @apiParam {int} status -1表示未发布，1表示已发布，默认为-1
* @apiParam {string} abstract 文章摘要，为空就为body的前128字
* @apiParam {int} topped -1表示不置顶，1表示置顶，默认为-1
* @apiParam {int} category_id 类别id
* @apiParam {list} tags_id 标签id，int型的列表
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*}
 */
func PutArticle(c *gin.Context) {
	var (
		err  error
		form models.Article
	)
	mc := Gin{C: c}
	err = c.ShouldBind(&form)
	logger.Info("form: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	err = form.Update()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, form)
}

/**
* @api {delete} /v1/article/:id 删除某个文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": null
*}
 */
func DeleteArticle(c *gin.Context) {
	mc := Gin{C: c}
	id := c.Param("id")
	logger.Info("try to delete article: ", id)
	intId := utils.StrToUnit64(id)
	if intId == 0 {
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
		return
	}
	err := models.DeleteArticleById(utils.StrToUnit64(id))
	if mc.CheckGormErr(err) != nil {
		logger.Error("delete category error: ", err)
		return
	}
	mc.WebJson(e.SUCCESS, nil)

}

/**
* @api {get} /v1/article 获取文章列表
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": [Object,]
*}
 */
func GetArticles(c *gin.Context) {
	var (
		form     models.ArticleListParam
		articles []*models.Article
		err      error
	)
	mc := Gin{C: c}
	err = c.ShouldBindQuery(&form)
	logger.Info("form: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	articles, err = models.GetArticleInfos(form)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, articles)
}

func GetArticleNames(c *gin.Context) {
	mc := Gin{C: c}
	logger.Info("get articles name")
	names, err := models.GetAllArticleNames()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, names)
}
