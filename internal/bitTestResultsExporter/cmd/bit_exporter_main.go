package main

import (
	exporter ".."
	"../../../server"
	. "../../models"
	"log"
	"time"
)

func main() {

	exporter.LoadConfigs()

	RedirectLogger(exporter.Configs.BitExporterPath)

	log.Printf("Server started - bit-test-result-exporter")

	router := exporter.ExporterRoutes.NewRouter()

	srv := server.NewServer(router, exporter.Configs.BitExporterPort)

	// NOTE: requests may be sent in 0.1 second deviation of the requested duration
	go exporter.ReportsScheduler(time.Second)

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS(exporter.Configs.SSHCertPath, exporter.Configs.SSHKeyPath)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
