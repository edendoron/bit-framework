package main

import (
	. "../internal/bitIndexer"
	. "../internal/bitStorageAccess"
	. "../internal/bitTestResultsExporter"
	"log"
	"sync"
)

func main() {
	log.Printf("Server started")

	// create a 'WaitGroup'
	wg := new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		BitExporter()
		wg.Done()
	}()
	go func() {
		BitStorageAccess()
		wg.Done()
	}()
	go func() {
		BitIndexer()
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}

