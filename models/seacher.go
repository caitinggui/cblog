package models

import (
	"encoding/gob"
	"fmt"
	logger "github.com/caitinggui/seelog"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	// searcher is coroutine safe
	Searcher = riot.Engine{}

	text  = "Google Is Experimenting With Virtual Reality Advertising"
	text1 = `Google accidentally pushed Bluetooth update for Home
	speaker early`
	text2 = `Google is testing another Search results layout with 
	rounded cards, new colors, and the 4 mysterious colored dots again`

	opts = types.EngineOpts{
		Using:         1,
		GseDict:       "static/dict/dictionary.txt",
		StopTokenFile: "static/dict/stop_tokens.txt",
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.FrequenciesIndex,
		},
		UseStore: true,
		// StoreFolder: path,
		StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
	}
)

func InitIndex() {
	// gob.Register(MyAttriStruct{})
	gob.Register(Article{})

	// var path = "./riot-index"
	Searcher.Init(opts)
	defer Searcher.Close()
	// os.MkdirAll(path, 0777)

	arts, err := GetFullArticle()
	if err != nil {
		panic("load article index failed: " + err.Error())
	}
	for _, v := range arts {
		Searcher.Index(fmt.Sprint(v.ID), types.DocData{
			Content: fmt.Sprintf("%s %s %s %s", v.Title, v.Abstract, v.Body, v.Category.Name),
			Attri:   v,
		})
	}

	// Wait for the index to refresh
	Searcher.Flush()

	logger.Info("Created index number: ", Searcher.NumDocsIndexed())
}

func RestoreIndex() {
	// var path = "./riot-index"
	gob.Register(Article{})
	Searcher.Init(opts)
	defer Searcher.Close()
	// os.MkdirAll(path, 0777)

	// Wait for the index to refresh
	Searcher.Flush()

	logger.Info("recover index number: ", Searcher.NumDocsIndexed())
}

func Search() {
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
