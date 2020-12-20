package models

import (
	"fmt"
	logger "github.com/caitinggui/seelog"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/yanyiwu/gojieba"
	"strings"
)

// searcher is coroutine safe
var Searcher *riot.Engine
var Tokenizer *gojieba.Jieba

type IndexSearcher struct {
	searcher  *riot.Engine
	tokenizer *gojieba.Jieba
}

func (self *IndexSearcher) Close() {
	self.searcher.Close()
	self.tokenizer.Free()
}

func InitSearcher() {
	if Searcher != nil {
		logger.Info("Searcher has been init")
		return
	}
	Searcher = &riot.Engine{}
	Tokenizer = gojieba.NewJieba()
	opts := types.EngineOpts{
		NotUseGse: true,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		UseStore: true,
		// StoreFolder: path,
		//StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
	}
	Searcher.Init(opts)
}

// index full article
func IndexDoc(doc Article) {
	if doc.ID == 0 || !doc.IsPublished() {
		return
	}
	words := Tokenizer.Tokenize(doc.Title+doc.Abstract+doc.Body+doc.Category.Name, gojieba.SearchMode, false)
	tokensMap := make(map[string][]int, len(words))
	tokens := make([]types.TokenData, 0, len(words))
	for i := 0; i < len(words); i++ {
		if _, ok := tokensMap[words[i].Str]; ok {
			tokensMap[words[i].Str] = append(tokensMap[words[i].Str], words[i].Start)
			continue
		}
		tokensMap[words[i].Str] = []int{words[i].Start}
	}
	for k, v := range tokensMap {
		tokens = append(tokens, types.TokenData{Text: k, Locations: v})
	}
	Searcher.Index(fmt.Sprint(doc.ID), types.DocData{
		Tokens: tokens,
		Attri:  doc,
	})
	return
}

// index full article by id
func IndexArticleById(id string) {
	arti, err := GetFullArticleById(id)
	if err != nil {
		logger.Infof("article(%v) doesn't exist: %v", id, err)
		return
	}

	IndexDoc(arti)
}

// index full article by ids
func IndexArticleByIds(ids []string) {
	if len(ids) == 0 {
		return
	}
	artis, err := GetFullArticleByIds(ids)
	if err != nil {
		logger.Warnf("get %v articles failed: ", ids)
	}
	for _, x := range artis {
		IndexDoc(x)
	}
}

// remove single doc index
func RemoveIndexById(id string) {
	Searcher.RemoveDoc(id, true)
}

func InitIndex() *IndexSearcher {
	// gob.Register(MyAttriStruct{})
	//gob.Register(Article{})
	// var path = "./riot-index"
	InitSearcher()
	// Searcher can't close here, otherwise we can't index doc to disk
	//defer Searcher.Close()
	// os.MkdirAll(path, 0777)

	go func() {
		arts, err := GetFullArticle()
		if err != nil {
			panic("load article index failed: " + err.Error())
		}
		for _, v := range arts {
			IndexDoc(v)
		}
		// Wait for the index to refresh
		Searcher.Flush()
		logger.Info("Created index number: ", Searcher.NumDocsIndexed())
	}()

	return &IndexSearcher{searcher: Searcher, tokenizer: Tokenizer}
}

// search full article
func SearchFullArticle(text string, page, pageSize int) ([]Article, int) {
	se := Searcher.SearchDoc(types.SearchReq{
		Tokens: Tokenizer.CutForSearch(text, true),
		RankOpts: &types.RankOpts{
			OutputOffset: (page - 1) * pageSize,
			MaxOutputs:   pageSize,
		}})
	res := []Article{}
	logger.Infof("total resul %v of keyword: %v from query: %v", se.NumDocs, se.Tokens, text)
	for _, doc := range se.Docs {
		d := doc.Attri.(Article)
		for _, t := range se.Tokens {
			d.Title = strings.Replace(d.Title, t, "<font color=red>"+t+"</font>", -1)
			d.Abstract = strings.Replace(d.Abstract, t, "<font color=red>"+t+"</font>", -1)
			d.Category.Name = strings.Replace(d.Category.Name, t, "<font color=red>"+t+"</font>", -1)
		}
		res = append(res, d)
	}
	return res, se.NumDocs
}
