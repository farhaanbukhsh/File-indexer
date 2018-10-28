package indexer

import (
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/farhaanbukhsh/file-indexer/conf"
	"github.com/farhaanbukhsh/file-indexer/logger"
	"github.com/farhaanbukhsh/file-indexer/utility"
	"github.com/radovskyb/watcher"
)

// FileIndexer is a data structure to hold the content of the file
type FileIndexer struct {
	FileName    string
	FileContent string
}

type FileIndexerArray struct {
	IndexerArray    []FileIndexer
	FileIndexLogger *logger.Logger
}

func Search(indexFilename string, searchWord string) *bleve.SearchResult {
	index, _ := bleve.Open(indexFilename)
	defer index.Close()
	query := bleve.NewQueryStringQuery(searchWord)
	request := bleve.NewSearchRequest(query)
	result, _ := index.Search(request)
	return result
}

func fileIndexing(fileIndexer FileIndexerArray, c conf.Configuration) {
	err := utility.DeleteExistingIndex(c.IndexFilename)
	fileIndexer.FileIndexLogger.Must(err, "Successfully deleted previous index")
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(c.IndexFilename, mapping)
	fileIndexer.FileIndexLogger.Must(err, "Successfully ran bleve for indexing")
	for _, fileIndex := range fileIndexer.IndexerArray {
		index.Index(fileIndex.FileName, fileIndex.FileContent)
	}
	defer index.Close()
}

func fileNameContentMap(c conf.Configuration, lg *logger.Logger) FileIndexerArray {
	var root = c.RootDirectory
	var files []string
	fileIndexer := FileIndexerArray{
		FileIndexLogger: lg,
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	fileIndexer.FileIndexLogger.Must(err, "Successfully traversed "+root)
	for _, filename := range files {
		content, err := utility.GetContent(filename)
		fileIndexer.FileIndexLogger.Must(err, "Successfully obtained content from "+filename)
		filesIndex := NewFileIndexer(filename, content)
		fileIndexer.IndexerArray = append(fileIndexer.IndexerArray, filesIndex)
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
func NewIndex(c conf.Configuration, lg *logger.Logger) {
	lg.Must(nil, "Refreshing the index")
	w := watcher.New()
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Remove, watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				lg.Must(nil, event.Name())
				fileIndexer := fileNameContentMap(c, lg)
				fileIndexing(fileIndexer, c)
			case err := <-w.Error:
				lg.Must(err, "")
			case <-w.Closed:
				return
			}
		}
	}()

	err := w.AddRecursive(c.RootDirectory)
	lg.Must(err, "Successfully added "+c.RootDirectory+" to the watcher")

	go func() {
		w.Wait()
	}()

	err = w.Start(time.Millisecond * 100)
	lg.Must(err, "Successfully started the watcher")
}
