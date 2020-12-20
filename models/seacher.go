package models

import (
	"cblog/config"
	"cblog/utils"
	"encoding/gob"
	"fmt"
	logger "github.com/caitinggui/seelog"
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"strings"
)

// searcher is coroutine safe
var Searcher *engine.Engine

func InitSearcher() {
	if Searcher != nil {
		logger.Info("Searcher has been init")
		return
	}
	Searcher = &engine.Engine{}
	opts := types.EngineInitOptions{
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.DocIdsIndex,
		},
		UsePersistentStorage:    !config.Config.Searcher.IsTest,
		PersistentStorageFolder: "indexs",
		// StoreFolder: path,
		//StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
	}
	if utils.IfPathExist(config.Config.Searcher.DictoryPath) && utils.IfPathExist(config.Config.Searcher.StopWordPath) {
		logger.Info("stopworkpath exists")
		//opts.StopTokenFile = config.Config.Searcher.StopWordPath
		//opts.SegmenterDictionaries = config.Config.Searcher.DictoryPath
	}
	Searcher.Init(opts)
}

// index full article
func IndexDoc(doc Article) {
	if doc.ID == 0 || !doc.IsPublished() {
		return
	}
	Searcher.IndexDocument(doc.ID, types.DocumentIndexData{
		Content: fmt.Sprintf("%s %s %s %s", doc.Title, doc.Abstract, doc.Body, doc.Category.Name),
		Labels:  []string{doc.Category.Name},
	}, true)
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
func RemoveIndexById(id uint64) {
	Searcher.RemoveDocument(id, true)
}

func InitIndex() *engine.Engine {
	// gob.Register(MyAttriStruct{})
	//gob.Register(Article{})
	// var path = "./engine-index"
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
		Searcher.FlushIndex()
		logger.Info("Created index number: ", Searcher.NumDocumentsIndexed())
	}()

	return Searcher
}

// return *engine.Engine to close
// must use *engine.Engine not engine.Engine
func RestoreIndex() *engine.Engine {
	gob.Register(Article{})
	InitSearcher()
	//defer Searcher.Close()

	// Wait for the index to refresh
	Searcher.FlushIndex()

	logger.Info("recover index number: ", Searcher.NumDocumentsIndexed())
	return Searcher
}

// search full article
func SearchFullArticle(text string, page, pageSize int) ([]Article, int) {
	se := Searcher.Search(types.SearchRequest{
		Text: text,
		RankOptions: &types.RankOptions{
			OutputOffset: (page - 1) * pageSize,
			MaxOutputs:   pageSize,
		}})
	ids := make([]string, 0, se.NumDocs)
	for _, doc := range se.Docs {
		ids = append(ids, utils.ToStr(doc.DocId))
	}
	res, _ := GetFullArticleByIds(ids)
	logger.Infof("total resul %v of keyword: %v from query: %v", se.NumDocs, se.Tokens, text)
	for k, _ := range res {
		for _, t := range se.Tokens {
			res[k].Title = strings.Replace(res[k].Title, t, "<font color=red>"+t+"</font>", -1)
			res[k].Abstract = strings.Replace(res[k].Abstract, t, "<font color=red>"+t+"</font>", -1)
			res[k].Category.Name = strings.Replace(res[k].Category.Name, t, "<font color=red>"+t+"</font>", -1)
		}
	}
	return res, se.NumDocs
}

func TestSearch() {
	InitIndex()
	//RestoreIndex()

	sea := Searcher.Search(types.SearchRequest{
		Text: "tupian",
		RankOptions: &types.RankOptions{
			OutputOffset: 0,
			MaxOutputs:   100,
		}})
	res := []uint64{}
	for _, doc := range sea.Docs {
		res = append(res, doc.DocId)
	}
	//_, data := Searcher.GetDBAllDocs()
	//logger.Info("index data:", data)
	logger.Info("search response: ", sea, "; docs = ", res)

	// os.RemoveAll("engine-index")
}
