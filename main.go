package main

import (
	"flag"
	"os"

	"github.com/farhaanbukhsh/file-indexer/conf"
	"github.com/farhaanbukhsh/file-indexer/indexer"
	"github.com/farhaanbukhsh/file-indexer/logger"
	"github.com/farhaanbukhsh/file-indexer/server"
)

var configFlag = flag.String("c", "", "Use alternative config file.")
var verboseFlag = flag.Bool("v", false, "Verbosely print the log output to Standard Output.")
var helpFlag = flag.Bool("h", false, "Print this help information.")

func main() {
	flag.Parse()

	if *helpFlag == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	lg := logger.NewLogger("logfile", *verboseFlag)
	config := conf.NewConfig(*configFlag, lg)
	go indexer.NewIndex(config, lg)
	server.NewServer(config, lg).Start()
}
