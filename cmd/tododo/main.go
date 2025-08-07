package main

import (
	"log"

	"github.com/bmarse/tododo/pkg/ui"
)

func main() {
	err := ui.Run()
	if err != nil {
		log.Fatalf("Error running tododo: %v", err)
	}
}
