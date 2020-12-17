package service

import (
	"cblog/config"
	"cblog/utils/V"
	logger "github.com/caitinggui/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"

	"cblog/models"
	"cblog/utils"
	"cblog/utils/e"
)

type indexContext struct {
	Cates      []models.Category
	Tags       []models.Tag
	Visitors   []models.Visitor
	Links      []models.Link
	VisitorSum interface{}
}

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
	mc := NewAdvancedGinContext(c)
	err = c.ShouldBind(&form)
	logger.Debug("创建文章: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	if form.ID != 0 {
		var tags []models.Tag
		for _, tagId := range form.TagsId {
			tag := models.Tag{}
			tag.ID = tagId
			tags = append(tags, tag)
		}
		// TODO 这里不一定要用ReplaceTags，可以直接delete已去除的tag，然后update会全部insert
		form.ReplaceTags(tags)
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
	mc := NewAdvancedGinContext(c)
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
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	logger.Info("get article : ", id)
	article, err := models.GetFullArticleById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	visitors, err := models.GetVisitorsByArticle(id)
	for k, _ := range visitors {
		visitors[k].IP = utils.FormatIP(visitors[k].IP)
	}
	if mc.CheckGormErr(err) != nil {
		return
	}
	comments, err := models.GetCommentByArticleId(article.ID)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.Res = gin.H{
		"Article":       article,
		"Visitors":      visitors,
		"Comments":      comments,
		"CommentsNum":   len(comments),
		"IsCommentOpen": config.Config.IsCommentOpen,
	}
	mc.SuccessHtml("blog/detail.html", mc.Res)
}

/**
* @api {get} /v1/article/:id/download?fileId= 下载文章附件
* @apiGroup Article
* @apiVersion 0.1.0
*
 */
func PraiseArticle(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	arti, err := models.GetArticleById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	arti.Likes += 1
	err = arti.Update()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, arti.Likes)
}

/**
* @api {get} /v1/article/:id/download?fileId= 下载文章附件
* @apiGroup Article
* @apiVersion 0.1.0
*
 */
func DownloadArticleAttachment(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	fileId := c.Query("fileId")
	arti, err := models.GetArticleById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	fileIdInt := utils.StrToInt64(fileId)
	attachments := strings.Split(arti.AttachmentUrl, V.AttachmentSeparator)
	if len(attachments) < int(fileIdInt) {
		logger.Warnf("用户%s在下载不存在的附件%s,文章为:%s", mc.GetCurrentUser(), fileId, id)
		mc.Redirect("/blog/article/" + id)
		return
	}
	mc.ServeFile(path.Join(V.AttachmentDirectory, id), attachments[fileIdInt])
}

/**
* @api {get} /v1/article/:id/upload 上传文章附件
* @apiGroup Article
* @apiVersion 0.1.0
*
 */
func UploadArticleAttachment(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	attachment, err := c.FormFile("uploadfile")
	if mc.CheckBindErr(err) != nil {
		return
	}
	if len(attachment.Filename) > 256 {
		mc.WebJson(e.ERR_INVALID_PARAM, "附件文件名过长")
	}
	arti, err := models.GetArticleById(id)
	if mc.CheckGormErr(err) != nil {
		return
	}
	attachmentDirectory := path.Join(V.AttachmentDirectory, id)
	os.MkdirAll(attachmentDirectory, os.ModePerm) // 如果目录存在，函数会返回nil
	err = c.SaveUploadedFile(attachment, path.Join(attachmentDirectory, attachment.Filename))
	if err != nil {
		logger.Error("upload file failed: ", err)
		mc.WebJson(e.ERROR, err)
		return
	}
	arti.AttachmentUrl = arti.AttachmentUrl + V.AttachmentSeparator + attachment.Filename
	err = arti.Update()
	if mc.CheckGormErr(err) != nil {
		return
	}
	logger.Infof("success upload article %d attachment: %s", arti.ID, attachment.Filename)
	mc.Redirect("/blog/article/" + id)
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
	mc := NewAdvancedGinContext(c)
	id := c.Query("id")
	res := gin.H{
		"Article": models.Article{},
	}
	if id != "" {
		article, err := models.GetFullArticleById(id)
		logger.Debugf("get article: %+v", article)
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
	mc := NewAdvancedGinContext(c)
	id := c.Param("id")
	logger.Info("try to delete article: ", id)
	intId := utils.StrToUint64(id)
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
	mc := NewAdvancedGinContext(c)
	err = c.ShouldBindQuery(&form)
	logger.Info("form: ", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	articles, _, err = models.GetArticleInfos(form, false)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.SuccessHtml("admin/article-list.html", gin.H{"Article": articles})
}

func SearchArticles(c *gin.Context) {
	var (
		form models.ArticleSearchParam
		err  error
	)
	mc := NewAdvancedGinContext(c)
	err = c.ShouldBindQuery(&form)
	logger.Infof("search article:  %v", form.Text)
	if mc.CheckBindErr(err) != nil {
		return
	}
	articles, articleNum := models.SearchFullArticle(form.Text, form.Page, form.PageSize)
	pages := Paginator(int(form.Page), int(form.PageSize), int(articleNum))

	res, err := getIndexContext(mc)
	if err != nil {
		return
	}
	mc.Res = map[string]interface{}{
		"Articles":   articles,
		"Cates":      res.Cates,
		"Tags":       res.Tags,
		"Visitors":   res.Visitors,
		"VisitorSum": res.VisitorSum,
		"Paginator":  pages,
		"IsQuery":    true,
		"Links":      res.Links,
	}
	logger.Debugf("res: %+v", mc.Res)
	mc.SuccessHtml("blog/index.html", mc.Res)
}

// index page
func GetArticleIndex(c *gin.Context) {
	var (
		form models.ArticleListParam
	)
	mc := NewAdvancedGinContext(c)
	err := c.ShouldBindQuery(&form)
	logger.Infof("get article index, form: %+v", form)
	if mc.CheckBindErr(err) != nil {
		return
	}
	articles, articleNum, err := models.GetArticleInfos(form, true)
	if mc.CheckGormErr(err) != nil {
		return
	}
	articleByMonth, err := models.CountArticleByMonth()
	if mc.CheckGormErr(err) != nil {
		return
	}
	hotArticle, err := models.GetArticleInfoByWeight(10)
	if mc.CheckGormErr(err) != nil {
		return
	}
	pages := Paginator(int(form.Page), int(form.PageSize), articleNum)
	res, err := getIndexContext(mc)
	if err != nil {
		return
	}
	comment, err := models.GetCommentsByCreatedAt(10)
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.Res = map[string]interface{}{
		"Articles":      articles,
		"Cates":         res.Cates,
		"Tags":          res.Tags,
		"Visitors":      res.Visitors,
		"VisitorSum":    res.VisitorSum,
		"Links":         res.Links,
		"Paginator":     pages,
		"DateArchive":   articleByMonth,
		"HotArticle":    hotArticle,
		"RecentComment": comment,
	}
	mc.SuccessHtml("blog/index.html", mc.Res)
}

func GetArticleNames(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	logger.Info("get articles name")
	names, err := models.GetAllArticleNames()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.WebJson(e.SUCCESS, names)
}

// admin index
func AdminIndex(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	updateBlogTemplateContext(mc)
	c.HTML(http.StatusOK, "admin/index.html", mc.Res)
}

func PostComment(c *gin.Context) {
	mc := NewAdvancedGinContext(c)
	articleId := c.Param("id")
	body := c.PostForm("body")
	if len(body) == 0 || utils.StrToUint64(articleId) == 0 {
		mc.Redirect("/")
		return
	}
	if !config.Config.IsCommentOpen {
		logger.Warn("评论已关闭")
		mc.Redirect("/blog/article/" + articleId)
		return
	}
	form := models.Comment{
		ArticleId: utils.StrToUint64(articleId),
		Body:      body,
		Name:      utils.RandomName(),
	}
	err := form.Insert()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.Redirect("/blog/article/" + articleId)
}

func updateBlogTemplateContext(mc *Gin) {
	if value, ok := mc.Res["Links"]; ok {
		logger.Errorf("template context exist same key: Links, the value: %v", value)
		return
	}
	links, err := models.GetAllLinks()
	if mc.CheckGormErr(err) != nil {
		return
	}
	mc.Res["Links"] = links
	return
}

func getIndexContext(mc *Gin) (res indexContext, err error) {
	cates, err := models.GetAllCategories()
	if mc.CheckGormErr(err) != nil {
		return
	}
	tags, err := models.GetAllTags()
	if mc.CheckGormErr(err) != nil {
		return
	}
	visitors, err := models.GetVisitors(0, V.DefaultPageSize)
	if mc.CheckGormErr(err) != nil {
		return
	}
	for k, _ := range visitors {
		visitors[k].IP = utils.FormatIP(visitors[k].IP)
	}
	visitorSum, _ := models.GetCache(V.VisitorSum)
	links, err := models.GetAllLinks()
	if mc.CheckGormErr(err) != nil {
		return
	}
	res = indexContext{
		Cates:      cates,
		Tags:       tags,
		Visitors:   visitors,
		VisitorSum: visitorSum,
		Links:      links,
	}
	return res, err
}
