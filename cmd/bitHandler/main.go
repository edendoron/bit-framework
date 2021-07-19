package main

import (
	handler "github.com/edendoron/bit-framework/internal/bitHandler"
	. "github.com/edendoron/bit-framework/internal/models"
	"github.com/edendoron/bit-framework/server"
	"log"
)

func main() {

	handler.LoadConfigs()

	RedirectLogger(handler.Configs.BitHandlerPath)

	log.Printf("Server started - bit-handler")

	log.Println("default trigger period per second is:", handler.CurrentTrigger.PeriodSec, "Type is:", handler.CurrentTrigger.BitType)

	router := handler.HandlerRoutes.NewRouter()

	srv := server.NewServer(router, handler.Configs.Host+handler.Configs.BitHandlerPort)

	go handler.StatusScheduler()

	if handler.Configs.UseHTTPS {
		err := srv.ListenAndServeTLS(handler.Configs.SSHCertPath, handler.Configs.SSHKeyPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
