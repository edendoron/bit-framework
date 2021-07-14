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

	if exporter.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(exporter.Configs.SSHCertPath, exporter.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
