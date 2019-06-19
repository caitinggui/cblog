package service

import (
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/e"
)

/**
* @api {post} /v1/article 创建文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiParam {string} title 文章标题
* @apiParam {string} body 文章内容
* @apiParam {string} abstract 摘要
* @apiParam {int=-1, 1} status -1表示未发表，1表示已发表
* @apiParam {int=-1, 1} topped -1表示不置顶，1表示置顶
* @apiParam {int} category_id 类别id
* @apiParam {int[]} [tags_id] 标签id
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*    }
 */
func CreateOrUpdateArticle(c *gin.Context) {
	var (
		form models.Article
		err  error
	)
	mc := Gin{C: c}
	err = c.ShouldBind(&form)
	logger.Debug("创建文章: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	if form.ID != 0 {
		form.Update()
	} else {
		err = form.Insert()
	}
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.Redirect("/admin/article")
}

/**
* @api {put} /v1/article 修改文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiParam {string} title 文章标题
* @apiParam {string} body 文章内容
* @apiParam {string} abstract 摘要
* @apiParam {int=-1, 1} status -1表示未发表，1表示已发表
* @apiParam {int=-1, 1} topped -1表示不置顶，1表示置顶
* @apiParam {int} category_id 类别id
* @apiParam {int[]} [tags_id] 标签id
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*    }
 */
func UpdateArticle(c *gin.Context) {
	var (
		form models.Article
		err  error
	)
	mc := Gin{C: c}
	err = c.ShouldBind(&form)
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
* @api {get} /v1/article/edit/:id 创建、修改文章
* @apiGroup Article
* @apiVersion 0.1.0
*
* @apiSuccessExample {json} Success-Response:
*   {
*     "errCode": "0",
*     "errMsg": "请求成功",
*     "data": Object
*}
**/
func EditArticle(c *gin.Context) {
	mc := Gin{C: c}
	id := c.Query("id")
	res := gin.H{
		"Article": models.Article{},
	}
	if id != "" {
		article, err := models.GetFullArticleById(id)
		logger.Debug("get article: ", article)
		if mc.CheckGormErr(err) != nil {
			return
		}
		tags := make([]string, len(article.Tags))
		for _, tag := range article.Tags {
			tags = append(tags, utils.ToStr(tag.ID))
		}
		res["Article"] = article
		res["ExistTags"] = tags
	}
	cates, err2 := models.GetAllCategories()
	if mc.CheckGormErr(err2) != nil {
		return
	}
	res["Cates"] = cates
	tags, err3 := models.GetAllTags()
	if mc.CheckGormErr(err3) != nil {
		return
	}
	res["Tags"] = tags
	mc.SuccessHtml("admin/article-edit.html", res)
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
	err := models.DeleteArticleById(intId)
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
	mc.SuccessHtml("admin/article-list.html", gin.H{"Article": articles})
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
