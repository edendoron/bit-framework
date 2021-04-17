package main

import (
	. "../internal/bit-indexer"
	. "../internal/bit-test-results-exporter"
	"log"
)

func main() {
	log.Printf("Server started")
	go BitExporter()
	BitIndexer()
}
