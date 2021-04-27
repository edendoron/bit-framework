package main

import (
	. "../internal/bitIndexer"
	. "../internal/bitStorageAccess"
	. "../internal/bitTestResultsExporter"
	"log"
)

func main() {
	log.Printf("Server started")
	go BitExporter()
	go BitStorageAccess()
	BitIndexer()
}
