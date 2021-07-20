package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRunSimulation(t *testing.T) {

	cmdRunStorage := exec.Command("go", "run", "../cmd/bitStorageAccess/main.go", "-config-file", ConfigPath)

	cmdRunConfig := exec.Command("go", "run", "../cmd/bitConfig/main.go", "-config-file", ConfigPath)

	cmdRunExporter := exec.Command("go", "run", "../cmd/bitTestResultsExporter/main.go", "-config-file", ConfigPath)

	cmdRunIndexer := exec.Command("go", "run", "../cmd/bitIndexer/main.go", "-config-file", ConfigPath)

	cmdRunHandler := exec.Command("go", "run", "../cmd/bitHandler/main.go", "-config-file", ConfigPath)

	cmdRunQuery := exec.Command("go", "run", "../cmd/bitQuery/main.go", "-config-file", ConfigPath)

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

	time.Sleep(3 * time.Second)

	ticker := time.NewTicker(200 * time.Millisecond)
	stopTicker := make(chan struct{})
	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case <- ticker.C:
				report := TestReport{
					TestId: 		float64(rand.Intn(math.MaxInt64)),
					ReportPriority: float64(rand.Intn(math.MaxInt64)),
					Timestamp:      time.Now(),
					TagSet: []KeyValue{
						{Key: "hostname", Value: "server02"},
					},
					FieldSet: []KeyValue{
						{Key: "volts", Value: fmt.Sprint(rand.Intn(math.MaxInt64))},
					},
				}
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

	time.Sleep(10 * time.Second)

	fmt.Println("Press any key to stop simulation:")
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if len(input.Text()) > 0 {
			close(stopTicker)
			break
		}
	}


	cleanTest()

	err = cmdRunHandler.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunHandler:\n", err)
	}

	err = cmdRunIndexer.Wait()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatalln("Error waiting on cmdRunIndexer:\n", err)
	}

	err = cmdRunExporter.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunExporter:\n", err)
	}

	err = cmdRunConfig.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunConfig:\n", err)
	}

	err = cmdRunStorage.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunStorage:\n", err)
	}
}
