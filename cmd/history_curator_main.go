package main

import (
	history "../internal/bitHistoryCurator"
	. "../internal/models"
	"../server"
	"context"
	"log"
)

func main() {

	RedirectLogger("./internal/bitHistoryCurator")

	log.Printf("Service started - bit-history-curator")

	router := history.HistoryCuratorRoutes.NewRouter()

	srv := server.NewServer(router, ":8083")

	go func() {
		history.RemoveAgedData(history.GetConfig().AgedTimeDuration)
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
