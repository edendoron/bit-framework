package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestRunSimulation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping simulation")
	}

	fmt.Println("---------Starting simulation---------")

	cmdRunStorage := exec.Command("go", "run", "../cmd/bitStorageAccess/main.go", "-config-file", ConfigPath)

	cmdRunConfig := exec.Command("go", "run", "../cmd/bitConfig/main.go", "-config-file", ConfigPath)

	cmdRunExporter := exec.Command("go", "run", "../cmd/bitTestResultsExporter/main.go", "-config-file", ConfigPath)

	cmdRunIndexer := exec.Command("go", "run", "../cmd/bitIndexer/main.go", "-config-file", ConfigPath)

	cmdRunHandler := exec.Command("go", "run", "../cmd/bitHandler/main.go", "-config-file", ConfigPath)

	cmdRunQuery := exec.Command("go", "run", "../cmd/bitQuery/main.go", "-config-file", ConfigPath)

	services = append(services, cmdRunStorage, cmdRunExporter, cmdRunIndexer, cmdRunHandler, cmdRunQuery)

	cleanTestStorageConfigDirs() // clean test environment storage before test start

	err := cmdRunStorage.Start()
	if err != nil {
		log.Fatalln("Couldn't start bitStorageAccess service", err)
	}

	time.Sleep(3 * time.Second)

	err = cmdRunConfig.Start()
	if err != nil {
		cleanTest()
		log.Fatalln("Couldn't start bitConfig service", err)
	}

	time.Sleep(3 * time.Second)

	err = cmdRunHandler.Start()
	if err != nil {
		cleanTest()
		log.Fatalln("Couldn't start bitHandler service\n", err)
	}

	err = cmdRunIndexer.Start()
	if err != nil {
		cleanTest()
		log.Fatalln("Couldn't start bitIndexer service", err)
	}

	err = cmdRunExporter.Start()
	if err != nil {
		cleanTest()
		log.Fatalln("Couldn't start bitTestResultsExporter service", err)
	}

	time.Sleep(5 * time.Second)

	ticker := time.NewTicker(200 * time.Millisecond)
	stopTicker := make(chan struct{})

	go func() {
		for {
			select {
			case <- ticker.C:
				report := generateReport()
				sentReports := ReportBody{Reports: []TestReport{report}}

				reportsMarshaled, err := json.MarshalIndent(sentReports, "", " ")
				if err != nil {
					cleanTest()
					log.Fatalln("Error marshaling sent reports:\n", err)
				}

				body := bytes.NewBuffer(reportsMarshaled)
				exporterRes, err := http.Post("http://localhost:8087/report/raw", "application/json; charset=UTF-8", body)
				if err != nil || exporterRes.StatusCode != http.StatusOK {
					cleanTest()
					log.Fatalln("Reports couldn't reach exporter\n", err)
				}
				time.Sleep(time.Second)
			case <-stopTicker:
				ticker.Stop()
				return
			}
		}
	}()

	err = cmdRunQuery.Start()
	if err != nil {
		cleanTest()
		log.Fatalln("Couldn't start bitQuery service:\n", err)
	}

	fmt.Println("Press ctrl+c to stop simulation:")
	startTime := time.Now()
	for range time.Tick(30 * time.Second){
		fmt.Println("time elapsed: " + time.Since(startTime).Round(time.Second).String())
	}

	cleanTest()
}
