package tests

import (
	"bytes"
	"encoding/json"
	. "github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	. "github.com/edendoron/bit-framework/internal/bitHandler"
	. "github.com/edendoron/bit-framework/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

var configPath = "./configs/prog_configs/configs.yml"

//test the normal flow of the framework
func TestSystemFlow(t *testing.T) {
	killServicesCmd := exec.Command("Taskkill", "/IM", "main.exe", "/F")

	cmdRunStorage := exec.Command("go", "run", "../cmd/bitStorageAccess/main.go", "-config-file", configPath)

	cmdRunConfig := exec.Command("go", "run", "../cmd/bitConfig/main.go", "-config-file", configPath)

	cmdRunExporter := exec.Command("go", "run", "../cmd/bitTestResultsExporter/main.go", "-config-file", configPath)

	cmdRunIndexer := exec.Command("go", "run", "../cmd/bitIndexer/main.go", "-config-file", configPath)

	cleanTestStorageConfigDirs() //clean test environment storage before test start

	err := cmdRunStorage.Start()
	if err != nil {
		log.Fatalln("Couldn't start bitStorageAccess service", err)
	}

	time.Sleep(time.Second)

	err = cmdRunConfig.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitConfig service", err)
	}

	time.Sleep(3 * time.Second)

	configFiles, err := ioutil.ReadDir("./configs/config_failures")
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Can't read filtering rules directory", err)
	}

	configFailuresStorage := fetchConfigFailuresFromStorage()
	if len(configFiles) != len(configFailuresStorage) {
		cleanTest(killServicesCmd)
		t.Errorf("At least one config failure is missing from storage %v", err)
		return
	}

	userGroups := fetchUserGroupsFromStorage()
	if len(userGroups) == 0 {
		cleanTest(killServicesCmd)
		t.Errorf("No user groups were generated in storage %v", err)
		return
	}

	sentReports := ReportBody{Reports: []TestReport{report0, report1, report2, report3, report4, report5,
		report6, report7, report8, report9, report10}}

	reportsMarshaled, err := json.MarshalIndent(sentReports, "", " ")
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Error marshaling sent reports", err)
	}
	body := bytes.NewBuffer(reportsMarshaled)

	err = cmdRunIndexer.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitIndexer service", err)
	}

	err = cmdRunExporter.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitTestResultsExporter service", err)
	}

	time.Sleep(time.Second)

	exporterRes, err := http.Post("http://localhost:8087/report/raw", "application/json; charset=UTF-8", body)
	if err != nil || exporterRes.StatusCode != http.StatusOK {
		cleanTest(killServicesCmd)
		t.Errorf("Reports couldn't reach exporter  %v", err)
		return
	}

	time.Sleep(3 * time.Second)

	reports := fetchReportsFromStorage()

	if len(reports) != len(sentReports.Reports) {
		cleanTest(killServicesCmd)
		t.Errorf("At least one report is missing from the storage.")
		return
	}

	cleanTest(killServicesCmd)

	err = cmdRunIndexer.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunIndexer:\n", err)
	}

	err = cmdRunExporter.Wait()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatalln("Error waiting on cmdRunExporter:\n", err)
	}

	err = cmdRunConfig.Wait()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatalln("Error waiting on cmdRunConfig:\n", err)
	}

	err = cmdRunStorage.Wait()
	if err != nil && err.Error() != "exit status 1" {
		log.Fatalln("Error waiting on cmdRunStorage:\n", err)
	}
}

// test the normal flow of the framework with bitHandler
func TestSystemFlowWithHandler(t *testing.T) {
	killServicesCmd := exec.Command("Taskkill", "/IM", "main.exe", "/F")

	cmdRunStorage := exec.Command("go", "run", "../cmd/bitStorageAccess/main.go", "-config-file", configPath)

	cmdRunConfig := exec.Command("go", "run", "../cmd/bitConfig/main.go", "-config-file", configPath)

	cmdRunExporter := exec.Command("go", "run", "../cmd/bitTestResultsExporter/main.go", "-config-file", configPath)

	cmdRunIndexer := exec.Command("go", "run", "../cmd/bitIndexer/main.go", "-config-file", configPath)

	cmdRunHandler := exec.Command("go", "run", "../cmd/bitHandler/main.go", "-config-file", configPath)

	cmdRunQuery := exec.Command("go", "run", "../cmd/bitQuery/main.go", "-config-file", configPath)

	cmdRunCurator := exec.Command("go", "run", "../cmd/bitHistoryCurator/main.go", "-config-file", configPath)

	cleanTestStorageConfigDirs() // clean test environment storage before test start

	err := cmdRunStorage.Start()
	if err != nil {
		log.Fatalln("Couldn't start bitStorageAccess service", err)
	}

	time.Sleep(time.Second)

	err = cmdRunConfig.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitConfig service", err)
	}

	time.Sleep(3 * time.Second)

	err = cmdRunHandler.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitHandler service\n", err)
	}


	err = cmdRunIndexer.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitIndexer service", err)
	}

	err = cmdRunExporter.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitTestResultsExporter service", err)
	}

	time.Sleep(time.Second)

	go func() {
		for i := 1; i < 10; i++ {
			report := TestReport{
				TestId: 		float64(i),
				ReportPriority: 12,
				Timestamp:      time.Now().Add(5*time.Second),
				TagSet: []KeyValue{
					{Key: "hostname", Value: "server02"},
				},
				FieldSet: []KeyValue{
					{Key: "volts", Value: "50"},
				},
			}
			sentReports := ReportBody{Reports: []TestReport{report}}

			reportsMarshaled, err := json.MarshalIndent(sentReports, "", " ")
			if err != nil {
				cleanTest(killServicesCmd)
				log.Fatalln("Error marshaling sent reports:\n", err)
			}

			body := bytes.NewBuffer(reportsMarshaled)
			exporterRes, err := http.Post("http://localhost:8087/report/raw", "application/json; charset=UTF-8", body)
			if err != nil || exporterRes.StatusCode != http.StatusOK {
				cleanTest(killServicesCmd)
				log.Fatalln("Reports couldn't reach exporter\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	err = cmdRunQuery.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitQuery service:\n", err)
	}

	time.Sleep(10 * time.Second)

	bitStatuses := fetchStatusFromQuery()

	if len(bitStatuses) > 0 && len(expectedResult.Failures) != len(bitStatuses[0].Failures) {
		cleanTest(killServicesCmd)
		t.Errorf("Handler didn't catch high voltage failure:\n %v", err)
	}

	cleanTestStorageConfigDirs()

	err = cmdRunCurator.Start()
	if err != nil {
		cleanTest(killServicesCmd)
		log.Fatalln("Couldn't start bitQuery service:\n", err)
	}

	err = killServicesCmd.Run()
	if err != nil {
		log.Fatalln("Error killing services", err)
	}

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

	err = cmdRunCurator.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunStorage:\n", err)
	}

	err = cmdRunStorage.Wait()
	if err != nil && err.Error() != "exit status 1"{
		log.Fatalln("Error waiting on cmdRunStorage:\n", err)
	}
}

// temp variables for tests

var expectedResult = BitStatus{
	Failures: []*BitStatus_RportedFailure{
		{
			FailureData: failure0.Failure.Description,
			Timestamp:   timestamppb.New(time.Now()),
			Count:       1,
		},
	},
}

var failure0 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "volts test",
			TestId:         1,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds voltage problem",
			Purpose:        "check voltage is within 1-7 range, with a deviation of 10%",
			Severity:       1,
			OperatorFailure: []string{
				"unable to start",
				"normal functionality is damaged",
			},
			LineReplacentUnits: []string{
				"line1",
				"line2",
			},
			FieldReplacemntUnits: []string{
				"field1",
				"field2",
				"field3",
			},
		},
		ExaminationRule: &FailureExaminationRule{
			MatchingField: "volts",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname"),
				Value: []byte("server02"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       2,
					Miximum:       7,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT,
						Value: 10,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     5,
					FailuresCCount: 0,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_NO_LATCH,
			IndicationSeconds: 600,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group1",
				"groupRafael",
				"groupField",
			},
			MasksOtherGroup: []string{
				"group3",
				"group4",
			},
		},
	},
}

var report0 = TestReport{
	TestId:         123,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(5*time.Second),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "volts", Value: "6.5"},
	},
}

var report1 = TestReport{
	TestId:         124,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname", Value: "server01"},
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "volts", Value: "10"},
		{Key: "oil", Value: "4"},
	},
}

var report2 = TestReport{
	TestId:         125,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "-1"},
	},
}

var report3 = TestReport{
	TestId:         126,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "-3.3"},
		{Key: "TemperatureCelsius", Value: "68"},
	},
}

var report4 = TestReport{
	TestId:         127,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "50"},
		{Key: "TemperatureCelsius", Value: "60"},
	},
}

var report5 = TestReport{
	TestId:         129,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "10"},
		{Key: "TemperatureCelsius", Value: "72"},
	},
}

var report6 = TestReport{
	TestId:         128,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(time.Second),
	TagSet: []KeyValue{
		{Key: "hostname123", Value: "north"},
		{Key: "host", Value: "north east"},
		{Key: "hostname", Value: "server02"},
	},
	FieldSet: []KeyValue{
		{Key: "AirPressure", Value: "23.25"},
		{Key: "TemperatureCelsius", Value: "70"},
	},
}

var report7 = TestReport{
	TestId:         130,
	ReportPriority: 12,
	Timestamp:      time.Now(),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "8"},
	},
}

var report8 = TestReport{
	TestId:         131,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(2 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "2"},
	},
}

var report9 = TestReport{
	TestId:         132,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(2 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "4.4"},
		{Key: "AirPressure", Value: "40.4"},
	},
}

var report10 = TestReport{
	TestId:         133,
	ReportPriority: 12,
	Timestamp:      time.Now().Add(4 * time.Second),
	TagSet: []KeyValue{
		{Key: "ss", Value: "longstanding"},
		{Key: "hostname123", Value: "north"},
	},
	FieldSet: []KeyValue{
		{Key: "oil", Value: "4.4"},
		{Key: "AirPressure", Value: "35"},
	},
}

//
//var report11 = TestReport{
//	TestId:         134,
//	ReportPriority: 12,
//	Timestamp:      time.Now().Add(4 * time.Second),
//}
//
//var report12 = TestReport{
//	TestId:         135,
//	ReportPriority: 12,
//	Timestamp:      time.Now().Add(time.Hour),
//}
