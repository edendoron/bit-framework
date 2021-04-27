package main

import (
	. "../internal/bitHistoryCurator"
	. "../internal/bitIndexer"
	. "../internal/bitStorageAccess"
	. "../internal/bitTestResultsExporter"
	"log"
	"sync"
)

const NumOfServices = 4

func main() {
	log.Printf("Server started")

	// create a 'WaitGroup'
	wg := new(sync.WaitGroup)
	wg.Add(NumOfServices)

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
	go func() {
		BitHistoryCurator()
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}
