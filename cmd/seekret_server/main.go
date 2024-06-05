package main

import (
	"flag"
	"log"

	"github.com/astrorick/seekret/internal/seekret"
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

	// run server with specified settings
	if err := seekret.New(configFilePath).Start(); err != nil {
		log.Fatal(err)
	}
}
