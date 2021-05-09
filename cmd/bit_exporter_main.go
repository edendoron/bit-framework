package main

import (
	exporter "../internal/bitTestResultsExporter"
	"../server"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	log.Printf("Server started - bit-test-result-exporter")

	//os.Setenv("EXPORTERBWSIZE", "20")
	//os.Setenv("EXPORTERBWUNITS", "KiB")

	size, _ := strconv.ParseFloat(os.Getenv("EXPORTERBWSIZE"), 32)
	units := os.Getenv("EXPORTERBWUNITS")

	exporter.CurrentBW.Size = float32(size)
	exporter.CurrentBW.UnitsPerSecond = units

	fmt.Println("default bandwidth size is:", exporter.CurrentBW.Size, "UPS:", exporter.CurrentBW.UnitsPerSecond)

	router := exporter.ExporterRoutes.NewRouter()

	srv := server.NewServer(router, ":8080")

	// NOTE: requests may be sent in 0.04 of a second deviation of the requested duration
	go exporter.ReportsScheduler(time.Second)

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS("../localhost.crt", "../localhost.key")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
