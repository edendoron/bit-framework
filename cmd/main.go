package main

import (
	. "../internal/bitHistoryCurator"
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

	go func() {
		wg.Add(1)
		BitExporter()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		BitStorageAccess()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		BitIndexer()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		BitHistoryCurator()
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}
