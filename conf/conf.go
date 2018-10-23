package conf

import (
	"encoding/json"
	"os"

	"github.com/farhaanbukhsh/file-indexer/logger"
)

// Configuration loads the config file
type Configuration struct {
	RootDirectory string
	IndexFilename string
	Port          string
}

// NewConfig is a function for creating a new configuration
func NewConfig(filename string) Configuration {
	config := Configuration{}
	name := "config.json"
	if filename != "NULL" {
		name = filename
	}
	file, err := os.Open(name)
	defer file.Close()
	logger.Must(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	logger.Must(err)
	return config
}
