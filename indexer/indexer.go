package indexer

// FileIndexer is a data structure to hold the content of the file
type FileIndexer struct {
	FileName    string
	FileContent string
}

// NewFileIndexer is a function to create a new File Indexer
func NewFileIndexer(fname, fcontent string) FileIndexer {
	temp := FileIndexer{
		FileName:    fname,
		FileContent: fcontent,
	}

	return temp
}
