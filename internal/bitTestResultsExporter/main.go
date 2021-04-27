package bitExporter

import (
	"../../server"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func BitExporter() {

	log.Printf("Server started - bit-test-result-exporter")

	os.Setenv("EXPORTERBWSIZE", "2")
	os.Setenv("EXPORTERBWUNITS", "KiB")

	size, _ := strconv.ParseFloat(os.Getenv("EXPORTERBWSIZE"), 32)
	units := os.Getenv("EXPORTERBWUNITS")

	currentBW.Size = float32(size)
	currentBW.UnitsPerSecond = units

	fmt.Println("default bandwidth size is:", currentBW.Size, "UPS:", currentBW.UnitsPerSecond)

	router := routes.NewRouter()

	//heap.Init(&reportsQueue)
	//reportsQueue.Init(reportsQueue.maxSize)

	srv := server.NewServer(router, ":8079")

	go reportsScheduler(5 * time.Second)

	go postIndexer()

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("localhost.crt", "localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
