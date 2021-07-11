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
		//TODO: handle error
		srv.Shutdown(context.Background())
	}()

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}