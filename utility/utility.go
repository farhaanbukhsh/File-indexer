package utility

import "io/ioutil"

// GetContent is a function for retrieving data from file
func GetContent(name string) string {
	data, err := ioutil.ReadFile(name)
	must(err)
	return string(data)
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
