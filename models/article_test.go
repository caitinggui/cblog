package models

import (
	"testing"
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
		Category: &cate,
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
