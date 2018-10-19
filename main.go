package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Configuration loads the config file
type Configuration struct {
	RootDirectory string
	IndexFilename string
	Port          string
}

// FileIndexer is a data structure to hold the content of the file
type FileIndexer struct {
	Filename    string
	FileContent string
}

var configFlag = flag.String("config", "NULL", "To pass a different configuration file")

func initializeConfig() Configuration {
	config := Configuration{}
	name := "config.json"
	if *configFlag != "NULL" {
		name = *configFlag
	}
	fmt.Println(*configFlag)
	file, err := os.Open(name)
	defer file.Close()
	checkerr(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	checkerr(err)
	return config
}

var config = initializeConfig()

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}

func getContent(name string) string {
	data, err := ioutil.ReadFile(name)
	checkerr(err)
	return string(data)
}

func fileIndexing(indexfilename string, fileIndexer []FileIndexer) error {
	err := deleteExistingIndex(config.IndexFilename)
	if err != nil {
		return err
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexfilename, mapping)
	if err != nil {
		return err
	}
	for _, fileIndex := range fileIndexer {
		index.Index(fileIndex.Filename, fileIndex.FileContent)
	}
	defer index.Close()
	return nil
}

func searchResults(indexFilename string, searchWord string) *bleve.SearchResult {
	index, _ := bleve.Open(indexFilename)
	defer index.Close()
	query := bleve.NewQueryStringQuery(searchWord)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	return searchResult
}

func fileNameContentMap() []FileIndexer {
	var ROOTPATH = config.RootDirectory
	var files []string
	var filesIndex FileIndexer
	var fileIndexer []FileIndexer

	err := filepath.Walk(ROOTPATH, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	checkerr(err)
	for _, filename := range files {
		content := getContent(filename)
		filesIndex = FileIndexer{Filename: filename, FileContent: content}
		fileIndexer = append(fileIndexer, filesIndex)
	}
	return fileIndexer
}

func creatIndex() error {
	var fileIndexer = fileNameContentMap()
	err := fileIndexing(config.IndexFilename, fileIndexer)
	if err != nil {
		return err
	}
	return nil
}

// IndexFile is the controller that helps with indexing the file
func IndexFile(w http.ResponseWriter, r *http.Request) {
	err := creatIndex()
	json.NewEncoder(w).Encode(err)
	return
}

// SearchFile is the controller that helps with indexing the file
func SearchFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchResult := searchResults(config.IndexFilename, params["query"])
	json.NewEncoder(w).Encode(searchResult.Hits)
	return
}

// Check if the index exist if it does, then flushes it off
func deleteExistingIndex(name string) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(name); err != nil {
			return fmt.Errorf("Can't Delete file: %v", err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	fmt.Println("Refreshing the index")
	err := creatIndex()
	checkerr(err)
	fmt.Printf("serving on %v \n", config.Port)
	router := mux.NewRouter()
	router.HandleFunc("/index", IndexFile).Methods("GET")
	router.HandleFunc("/search/{query}", SearchFile).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(config.Port, loggedRouter))
}
