package main

import (
	"log"

	"github.com/astrorick/seekret/internal/seekret"
)

func main() {
	if err := seekret.New().Start(); err != nil {
		log.Fatal(err)
	}
}
