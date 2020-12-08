package models

import (
	"cblog/config"
	"cblog/utils"
	"encoding/gob"
	"fmt"
	logger "github.com/caitinggui/seelog"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"strings"
)

// searcher is coroutine safe
var Searcher *riot.Engine

func InitSearcher() {
	if Searcher != nil {
		logger.Info("Searcher has been init")
		return
	}
	Searcher = &riot.Engine{}
	opts := types.EngineOpts{
		Using: 1,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.FrequenciesIndex,
		},
		UseStore: !config.Config.Searcher.IsTest,
		// StoreFolder: path,
		StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
	}
	if utils.IfPathExist(config.Config.Searcher.DictoryPath) && utils.IfPathExist(config.Config.Searcher.StopWordPath) {
		logger.Info("stopworkpath exists")
		opts.StopTokenFile = config.Config.Searcher.StopWordPath
		opts.GseDict = config.Config.Searcher.DictoryPath
	}
	Searcher.Init(opts)
}

// index full article
func IndexDoc(doc Article) {
	if doc.ID == 0 {
		return
	}
	Searcher.Index(fmt.Sprint(doc.ID), types.DocData{
		Content: fmt.Sprintf("%s %s %s %s", doc.Title, doc.Abstract, doc.Body, doc.Category.Name),
		Attri:   doc,
	})
	return
}

// index full article by id
func IndexArticleById(id string) {
	arti, err := GetFullArticleById(id)
	if err != nil {
		logger.Info("article(%v) doesn't exist: %v", id, err)
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

func InitIndex() *riot.Engine {
	// gob.Register(MyAttriStruct{})
	gob.Register(Article{})

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

	return Searcher
}

// return *riot.Engine to close
// must use *riot.Engine not riot.Engine
func RestoreIndex() *riot.Engine {
	gob.Register(Article{})
	InitSearcher()
	//defer Searcher.Close()

	// Wait for the index to refresh
	Searcher.Flush()

	logger.Info("recover index number: ", Searcher.NumDocsIndexed())
	return Searcher
}

// search full article
func SearchFullArticle(text string, page, pageSize int) ([]Article, int) {
	se := Searcher.SearchDoc(types.SearchReq{
		Text: text,
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

func TestSearch() {
	InitIndex()
	//RestoreIndex()

	sea := Searcher.SearchDoc(types.SearchReq{
		Text: "tupian",
		RankOpts: &types.RankOpts{
			OutputOffset: 0,
			MaxOutputs:   100,
		}})
	res := []Article{}
	for _, doc := range sea.Docs {
		res = append(res, doc.Attri.(Article))
	}
	//_, data := Searcher.GetDBAllDocs()
	//logger.Info("index data:", data)
	logger.Info("search response: ", sea, "; docs = ", res)

	// os.RemoveAll("riot-index")
}
