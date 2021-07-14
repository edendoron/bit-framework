package main

import (
	history ".."
	"../../../server"
	. "../../models"
	"context"
	"log"
)

func main() {

	history.LoadConfigs()

	RedirectLogger(history.Configs.BitHistoryCuratorPath)

	log.Printf("Service started - bit-history-curator")

	router := history.HistoryCuratorRoutes.NewRouter()

	srv := server.NewServer(router, history.Configs.BitHistoryCuratorPort)

	go func() {
		history.RemoveAgedData(history.GetCuratorTimeConfig())
		log.Printf("Service ended - bit-history-curator")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	//TODO: need to change to ListenAndServeTLS in order to support https
	//err := srv.ListenAndServeTLS(history.Configs.SSHCertPath, history.Configs.SSHKeyPath)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
