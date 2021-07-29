package tests

import (
	"encoding/json"
	"fmt"
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	"github.com/edendoron/bit-framework/internal/models"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const layout = "2006-January-02 15:4:5"
const ConfigPath = "./configs/prog_configs/configs.yml"

var services []*exec.Cmd

func fetchConfigFailuresFromStorage() []bit.Failure {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
		killServices()
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("config_failures", "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		killServices()
		log.Fatalln("Couldn't get config failures from storage:\n", err)
	}

	var configFailures []bit.Failure
	err = json.NewDecoder(resp.Body).Decode(&configFailures)
	if err != nil {
		killServices()
		log.Fatalln("Can't decode response from storage", err)
	}

	return configFailures
}

func fetchUserGroupsFromStorage() []string {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
		killServices()
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("user_groups", "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		killServices()
		log.Fatalln("Couldn't get user groups from storage:\n", err)
	}

	var userGroups []string
	err = json.NewDecoder(resp.Body).Decode(&userGroups)
	if err != nil {
		killServices()
		log.Fatalln("Can't decode response from storage", err)
	}

	return userGroups
}

func fetchReportsFromStorage() []models.TestReport {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
		killServices()
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("reports", "")
	params.Add("start", time.Now().Add(-5*time.Minute).Format(layout))
	params.Add("end", time.Now().Add(5*time.Minute).Format(layout))
	params.Add("filter", "time")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		killServices()
		log.Fatalln("Couldn't get reports from storage:\n", err)
	}

	var reports []models.TestReport
	err = json.NewDecoder(resp.Body).Decode(&reports)
	if err != nil {
		killServices()
		log.Fatalln("Can't decode response from storage", err)
	}

	return reports
}

func fetchStatusFromQuery() []bit.BitStatus {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8085/status", nil)
	if err != nil {
		killServices()
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("user_group", "engine")
	params.Add("start", time.Now().Add(-5*time.Minute).Format(layout))
	params.Add("end", time.Now().Add(5*time.Minute).Format(layout))
	params.Add("filter", "time")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		killServices()
		log.Fatalln("Couldn't get status from query:\n", err)
	}

	var bitStatus []bit.BitStatus
	err = json.NewDecoder(resp.Body).Decode(&bitStatus)
	if err != nil {
		killServices()
		log.Fatalln("Can't decode response from query:\n", err)
	}

	return bitStatus
}

func generateReport() models.TestReport {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(3)
	report := models.TestReport{
		TestId:         float64(rand.Intn(math.MaxInt64)),
		ReportPriority: float64(rand.Intn(256)),
		Timestamp:      time.Now(),
		TagSet: []models.KeyValue{
			{Key: "hostname", Value: "server02"},
		},
		FieldSet: []models.KeyValue{
			{Key: "", Value: ""},
		},
	}

	switch randomInt {
	case 0:
		report.FieldSet[0].Key = "volts"
		report.FieldSet[0].Value = fmt.Sprint(rand.Intn(20))
	case 1:
		report.FieldSet[0].Key = "AirPressure"
		report.FieldSet[0].Value = fmt.Sprint(rand.Intn(100))
	case 2:
		report.FieldSet[0].Key = "TemperatureCelsius"
		report.FieldSet[0].Value = fmt.Sprint(rand.Intn(100))
	}

	return report
}

func cleanTest() {
	deleteTestFilesFromStorage()
	killServices()
	cleanTestStorageConfigDirs()
}

func cleanTestStorageConfigDirs() {
	configDirs, err := os.ReadDir("./storage/config")
	if err != nil {
		log.Fatalln("Can't read config directories", err)
	}
	for _, dir := range configDirs {
		innerFiles, err := os.ReadDir("./storage/config/" + dir.Name())
		if err != nil {
			log.Fatalln("Can't read "+dir.Name(), err)
		}
		for _, file := range innerFiles {
			if file.Name() != ".gitignore" {
				err = os.Remove("./storage/config/" + dir.Name() + "/" + file.Name())
				if err != nil {
					log.Fatalln("Can't remove file "+file.Name()+"\n", err)
				}
			}
		}
	}
}

func deleteTestFilesFromStorage() {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8082/data/delete", nil)
	if err != nil {
		killServices()
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("timestamp", time.Now().Add(time.Minute).Format(layout))
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		killServices()
		log.Fatalln("Delete call to storage failed:\n", err)
	}
}

func killServices() {
	if runtime.GOOS == "windows" {
		if err := exec.Command("Taskkill", "/IM", "main.exe", "/F").Run(); err != nil {
			log.Fatalln("Error killing service: ", err)
		}
	} else {
		for _, service := range services {
			if err := service.Process.Signal(os.Kill); err != nil {
				log.Fatalln("Error killing service: ", err)
			}
		}
	}
}
