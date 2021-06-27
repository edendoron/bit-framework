package main

import (
	handler "../internal/bitHandler"
	. "../internal/models"
	"../server"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	RedirectLogger("./internal/bitHandler")

	log.Printf("Server started - bit-handler")

	os.Setenv("HANDLERTRIGGERPERIOD", "1")
	os.Setenv("HANDLERTRIGGERBITTYPE", "CBIT")

	periodSec, _ := strconv.ParseFloat(os.Getenv("HANDLERTRIGGERPERIOD"), 64)
	bitType := os.Getenv("HANDLERTRIGGERBITTYPE")

	handler.CurrentTrigger.PeriodSec = periodSec
	handler.CurrentTrigger.BitType = bitType

	fmt.Println("default trigger period per second is:", handler.CurrentTrigger.PeriodSec, "Type is:", handler.CurrentTrigger.BitType)

	router := handler.HandlerRoutes.NewRouter()

	srv := server.NewServer(router, ":8085")

	go handler.StatusScheduler()

	//TODO: need to change to ListenAndServeTLS in order to support https
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
