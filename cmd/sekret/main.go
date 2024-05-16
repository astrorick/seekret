package main

import (
	"log"

	"github.com/astrorick/sekret/internal/sekret"
)

func main() {
	if err := sekret.New().Start(); err != nil {
		log.Fatal(err)
	}
}
