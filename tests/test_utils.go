package tests

import (
	"encoding/json"
	. "github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const layout = "2006-January-02 15:4:5"

func fetchConfigFailuresFromStorage() []Failure {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("config_failures", "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln("Couldn't get config failures from storage:\n", err)
	}

	var configFailures []Failure
	err = json.NewDecoder(resp.Body).Decode(&configFailures)
	if err != nil {
		log.Fatalln("Can't decode response from storage", err)
	}

	return configFailures
}

func fetchUserGroupsFromStorage() []string {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("user_groups", "")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln("Couldn't get user groups from storage:\n", err)
	}

	var userGroups []string
	err = json.NewDecoder(resp.Body).Decode(&userGroups)
	if err != nil {
		log.Fatalln("Can't decode response from storage", err)
	}

	return userGroups
}

func fetchReportsFromStorage() []TestReport {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082/data/read", nil)
	if err != nil {
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
		log.Fatalln("Couldn't get reports from storage:\n", err)
	}

	var reports []TestReport
	err = json.NewDecoder(resp.Body).Decode(&reports)
	if err != nil {
		log.Fatalln("Can't decode response from storage", err)
	}

	return reports
}

func cleanTest(killServicesCmd *exec.Cmd) {
	cleanStorage()
	err := killServicesCmd.Run()
	if err != nil {
		log.Fatalln("Error killing services", err)
	}
}

func cleanStorage() {
	deleteTestFilesFromStorage()
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
					log.Fatalln("Can't remove file "+file.Name(), err)
				}
			}
		}
	}
}

func deleteTestFilesFromStorage() {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8082/data/delete", nil)
	if err != nil {
		log.Fatalln("Can't make new http request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	params := req.URL.Query()
	params.Add("timestamp", time.Now().Format(layout))
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalln("Delete call to storage failed:\n", err)
	}
}
