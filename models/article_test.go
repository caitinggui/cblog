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
