package indexer

import (
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/farhaanbukhsh/file-indexer/conf"
	"github.com/farhaanbukhsh/file-indexer/utility"
)

// FileIndexer is a data structure to hold the content of the file
type FileIndexer struct {
	FileName    string
	FileContent string
}

func Search(indexFilename string, searchWord string) *bleve.SearchResult {
	index, _ := bleve.Open(indexFilename)
	defer index.Close()
	query := bleve.NewQueryStringQuery(searchWord)
	request := bleve.NewSearchRequest(query)
	result, _ := index.Search(request)
	return result
}

func fileIndexing(fileIndexer []FileIndexer, c conf.Configuration) error {
	err := utility.DeleteExistingIndex(c.IndexFilename)
	if err != nil {
		return err
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(c.IndexFilename, mapping)
	if err != nil {
		return err
	}
	for _, fileIndex := range fileIndexer {
		index.Index(fileIndex.FileName, fileIndex.FileContent)
	}
	defer index.Close()
	return nil
}

func fileNameContentMap(c conf.Configuration) []FileIndexer {
	var ROOTPATH = c.RootDirectory
	var files []string
	var fileIndexer []FileIndexer

	err := filepath.Walk(ROOTPATH, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	must(err)
	for _, filename := range files {
		content := utility.GetContent(filename)
		filesIndex := NewFileIndexer(filename, content)
		fileIndexer = append(fileIndexer, filesIndex)
	}
	return fileIndexer
}

// NewFileIndexer is a function to create a new File Indexer
func NewFileIndexer(fname, fcontent string) FileIndexer {
	temp := FileIndexer{
		FileName:    fname,
		FileContent: fcontent,
	}

	return temp
}

// NewIndex is a function to create new indexes
func NewIndex(c conf.Configuration) error {
	fileIndexer := fileNameContentMap(c)
	err := fileIndexing(fileIndexer, c)
	if err != nil {
		return err
	}
	return nil
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
