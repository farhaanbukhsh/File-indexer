package conf

import (
	"encoding/json"
	"os"
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
	must(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	must(err)
	return config
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
