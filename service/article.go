package service

import (
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
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
	if err != nil {
		logger.Info("解析参数失败:", err)
		mc.WebJson(e.ERR_INVALID_PARAM, nil)
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
