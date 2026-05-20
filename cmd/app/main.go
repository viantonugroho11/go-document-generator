package main

import (
	"log"
	"os"

	"go-document-generator/internal/bootstrap"
)

func main() {
	if err := bootstrap.RunApp(); err != nil {
		log.Fatalf("app: %v", err)
	}
	os.Exit(0)
}
