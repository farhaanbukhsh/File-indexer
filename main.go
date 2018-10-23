package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/farhaanbukhsh/file-indexer/conf"
	"github.com/farhaanbukhsh/file-indexer/indexer"
	"github.com/farhaanbukhsh/file-indexer/utility"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var configFlag = flag.String("config", "NULL", "To pass a different configuration file")

var config conf.Configuration

func fileIndexing(indexfilename string, fileIndexer []indexer.FileIndexer) error {
	err := utility.DeleteExistingIndex(config.IndexFilename)
	if err != nil {
		return err
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexfilename, mapping)
	if err != nil {
		return err
	}
	for _, fileIndex := range fileIndexer {
		index.Index(fileIndex.FileName, fileIndex.FileContent)
	}
	defer index.Close()
	return nil
}

func fileNameContentMap() []indexer.FileIndexer {
	var ROOTPATH = config.RootDirectory
	var files []string
	var fileIndexer []indexer.FileIndexer

	err := filepath.Walk(ROOTPATH, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	must(err)
	for _, filename := range files {
		content := utility.GetContent(filename)
		filesIndex := indexer.NewFileIndexer(filename, content)
		fileIndexer = append(fileIndexer, filesIndex)
	}
	return fileIndexer
}

func createIndex(c conf.Configuration) error {
	fileIndexer := fileNameContentMap()
	err := fileIndexing(c.IndexFilename, fileIndexer)
	if err != nil {
		return err
	}
	return nil
}

// IndexFile is the controller that helps with indexing the file
func IndexFile(w http.ResponseWriter, r *http.Request) {
	err := createIndex(config)
	json.NewEncoder(w).Encode(err)
	return
}

// SearchFile is the controller that helps with indexing the file
func SearchFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchResult := indexer.Search(config.IndexFilename, params["query"])
	json.NewEncoder(w).Encode(searchResult.Hits)
	return
}

func main() {
	flag.Parse()
	config = conf.NewConfig(*configFlag)
	fmt.Println("Refreshing the index")
	err := createIndex(config)
	must(err)
	fmt.Printf("Serving on %v \n", config.Port)
	router := mux.NewRouter()
	router.HandleFunc("/index", IndexFile).Methods("GET")
	router.HandleFunc("/search/{query}", SearchFile).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(config.Port, loggedRouter))
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
