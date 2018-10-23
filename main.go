package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/farhaanbukhsh/file-indexer/conf"
	"github.com/farhaanbukhsh/file-indexer/indexer"
	"github.com/farhaanbukhsh/file-indexer/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var configFlag = flag.String("config", "NULL", "To pass a different configuration file")

var config conf.Configuration

// IndexFile is the controller that helps with indexing the file
func IndexFile(w http.ResponseWriter, r *http.Request) {
	err := indexer.NewIndex(config)
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
	err := indexer.NewIndex(config)
	logger.Must(err)
	fmt.Printf("Serving on %v \n", config.Port)
	router := mux.NewRouter()
	router.HandleFunc("/index", IndexFile).Methods("GET")
	router.HandleFunc("/search/{query}", SearchFile).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(config.Port, loggedRouter))
}
