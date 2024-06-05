package main

import (
	"log"

	seekret "github.com/astrorick/seekret/internal/seekret"
)

func main() {
	// TODO: read command line flags

	if err := seekret.New("path/to/config/file").Start(); err != nil {
		log.Fatal(err)
	}
}
