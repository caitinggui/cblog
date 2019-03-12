package service

import (
	"cblog/models"
	"cblog/utils/e"
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"
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

func GetArticleNames(c *gin.Context) {
	mc := Gin{C: c}
	logger.Info("get articles name")
	names, err := models.GetArticleNames()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, names)
}
