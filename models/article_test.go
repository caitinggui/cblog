package models

import (
	"testing"

	"cblog/utils"
)

func TestInsertArticle(t *testing.T) {
	article := Article{Title: "TestInsertArticle"}
	err := article.Insert()
	if err != nil {
		t.Fatal("TestInsertArticle error: ", err)
	}
	article.Delete()
	arti2 := Article{Title: "TestInsertArticle2"}
	arti2.Delete()
}

func TestReplaceTags(t *testing.T) {
	tag1 := Tag{Name: "tag1"}
	tag2 := Tag{Name: "tag2"}
	tag3 := Tag{Name: "tag3"}
	tag4 := Tag{Name: "tag4"}
	tag1.Insert()
	defer tag1.Delete()
	tag2.Insert()
	defer tag2.Delete()
	tag3.Insert()
	defer tag3.Delete()
	tag4.Insert()
	defer tag4.Delete()
	arti := Article{
		Title: "arti",
		Body:  "body",
		Tags:  []Tag{tag1, tag2},
	}
	arti.Insert()
	defer arti.Delete()
	err := arti.ReplaceTags([]Tag{tag3, tag4})
	if err != nil {
		t.Fatal("替换tags失败: ", err)
	}
	arti2, _ := GetFullArticleById(utils.ToStr(arti.ID))
	t.Log("arti2 in ReplaceTags: ", arti2)
	if arti2.Tags[0].ID != tag3.ID && arti2.Tags[0].ID != tag4.ID {
		t.Fatal("查询替换后的tags不符合预期")
	}
}

func TestGetFullArticleById(t *testing.T) {
	cate := Category{Name: "TestCategory"}
	cate.Insert()
	defer cate.Delete()
	tag := Tag{Name: "TestTag"}
	tag.Insert()
	defer tag.Delete()
	arti := Article{
		Title:      "标题1",
		Body:       "正文1",
		CategoryId: cate.ID,
		TagsId:     []uint64{tag.ID},
	}
	arti.Insert()
	defer arti.Delete()

	arti2, err := GetFullArticleById(utils.ToStr(arti.ID))
	if err != nil {
		t.Fatal("获取文章失败")
	}
	t.Log("获取全文章信息: ", arti2)
	if len(arti2.Tags) != 1 {
		t.Fatal("文章的tag个数不对")
	}
}

func TestGetArticleByCategory(t *testing.T) {
	var articles []Article
	cate := Category{Name: "TestGetArticleByCategory"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("插入类型失败: ", err)
	}
	defer cate.Delete()
	article := Article{
		Title:      "标题1",
		Body:       "正文1",
		CategoryId: cate.ID,
	}
	err = article.Insert()
	if err != nil {
		t.Fatal("创建文章失败: ", err)
	}
	defer article.Delete()
	article2 := Article{
		Title:    "标题2",
		Body:     "正文2",
		Category: cate,
	}
	err = article2.Insert()
	defer article2.Delete()
	t.Log("创建文章2的结果: ", err)
	articles, err = GetArticlesByCategory("TestGetArticleByCategory")
	if len(articles) != 2 {
		t.Fatal("根据类型查询文章失败: ", articles)
	}
}

func TestGetArticleByTag(t *testing.T) {
	tag := Tag{Name: "tag1"}
	tag.Insert()
	defer tag.Delete()
	article1 := Article{
		Title: "TestGetArticleByTag",
		Tags: []Tag{
			tag,
			{Name: "tag2"},
		},
	}
	article1.Insert()
	defer func() {
		article1.Delete()
		tag2, _ := GetTagByName("tag2")
		tag2.Delete()
	}()
	article2 := Article{
		Title: "TestGetArticleByTag2",
		Tags:  []Tag{tag},
	}
	article2.Insert()
	defer article2.Delete()
	articles, err := GetArticleByTag("tag1")
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 2 {
		t.Fatal("获取的文章数少于2")
	}
}

func TestGetArticleNames(t *testing.T) {
	article1 := Article{
		Title: "article1",
		Body:  "body1",
	}
	article1.Insert()
	defer article1.Delete()
	article2 := Article{
		Title: "article2",
		Body:  "body2",
	}
	article2.Insert()
	defer article2.Delete()
	names, _ := GetAllArticleNames()
	t.Log("文章名: ", names[0])
	if len(names) != 2 {
		t.Fatal("获取文章标题失败:", names)
	}
}

func TestGetArticleInfos(t *testing.T) {
	article1 := Article{
		Title: "article1",
		Body:  "body1",
	}
	article1.Insert()
	defer article1.Delete()
	article2 := Article{
		Title:  "article2",
		Body:   "body2",
		Topped: 1,
	}
	article2.Insert()
	defer article2.Delete()
	article3 := Article{
		Title:  "article3",
		Body:   "body3",
		Topped: -1,
	}
	article3.Insert()
	defer article3.Delete()

	articles, _ := GetAllArticleInfos()
	// 返回是有顺序的
	if len(articles) != 3 || articles[0].Abstract != "body2" {
		t.Log(len(articles), articles[0].Abstract)
		t.Fatal("获取文章信息失败:", articles[0], articles[1], articles[2])
	}

	form := ArticleListParam{
		Offset: 1,
		Size:   10,
	}
	articles, _ = GetArticleInfos(form)
	// 返回是有顺序的
	if len(articles) != 2 || articles[0].Abstract != "body3" {
		t.Fatal("获取分页文章信息失败:", articles[0])
	}
}
