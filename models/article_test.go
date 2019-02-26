package models

import (
	"testing"
)

func TestGetArticleByCategory(t *testing.T) {
	var articles []Article
	cate := Category{Name: "TestGetArticleByCategory"}
	err := cate.Insert()
	if err != nil {
		t.Fatal("插入类型失败: ", err)
	}
	article := Article{
		Title:      "标题1",
		Body:       "正文1",
		CategoryId: cate.ID,
	}
	err = article.Insert()
	if err != nil {
		t.Fatal("创建文章失败: ", err)
	}
	articles, err = GetArticlesByCategory("TestGetArticleByCategory")
	if len(articles) == 0 {
		t.Fatal("根据类型查询文章失败: ", articles)
	}
	cate.Delete()
	article.Delete()
}
