package utility

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/farhaanbukhsh/file-indexer/logger"
)

// GetContent is a function for retrieving data from file
func GetContent(name string) string {
	data, err := ioutil.ReadFile(name)
	logger.Must(err)
	return string(data)
}

// DeleteExistingIndex checks if the index exist if it does, then flushes it off
func DeleteExistingIndex(name string) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(name); err != nil {
			return fmt.Errorf("Can't Delete file: %v", err)
		}
	}
	return nil
}
