package main

import (
	"flag"
	"log"

	"github.com/astrorick/seekret/internal/server"
)

func main() {
	// define flags to match command line arguments
	var (
		configFilePath string
		displayHelp    bool
	)

	// bind flags to variables
	flag.StringVar(&configFilePath, "config", "", "Config file path.")
	flag.BoolVar(&displayHelp, "help", false, "Display help.")

	// parse flags
	flag.Parse()

	// Display help and exit if help flag is set
	if displayHelp {
		flag.Usage()
		return
	}

	// init server with specified settings
	seekretServer, err := server.New(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// run server
	if err := seekretServer.Start(); err != nil {
		log.Fatal(err)
	}
}
