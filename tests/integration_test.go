package tests

import (
	"bytes"
	"encoding/json"
	. "github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	. "github.com/edendoron/bit-framework/internal/bitHandler"
	. "github.com/edendoron/bit-framework/internal/models"
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

	cmdRunHandler := exec.Command("go", "run", "../cmd/bitHandler/main.go", "-config-file", configPath)

	//cmdRunQuery := exec.Command("go", "run", "../cmd/bitQuery/main.go", "-config-file", configPath)

	//cmdRunCurator := exec.Command("go", "run", "../cmd/bitCurator/main.go", "-config-file", configPath)

	err := cmdRunStorage.Start()
	if err != nil {
		log.Fatalln("Couldn't start bitStorageAccess service", err)
	}

	time.Sleep(time.Second)

	cleanStorage() //clean test environment storage before test start

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
	}

	userGroups := fetchUserGroupsFromStorage()
	if len(userGroups) == 0 {
		cleanTest(killServicesCmd)
		t.Errorf("No user groups were generated in storage %v", err)
	}

	//cmdRunHandler.Start()

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
	}

	time.Sleep(3 * time.Second)

	reports := fetchReportsFromStorage()

	if len(reports) != len(sentReports.Reports) {
		cleanTest(killServicesCmd)
		t.Errorf("At least one report is missing from the storage.")
	}

	cleanStorage()
	err = killServicesCmd.Run()
	if err != nil {
		log.Fatalln("Error killing services", err)
	}


	cmdRunHandler.Wait()
	cmdRunIndexer.Wait()
	cmdRunExporter.Wait()
	cmdRunConfig.Wait()
	cmdRunStorage.Wait()
}





// temp variables for tests

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
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
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

var failure1 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "temperature test",
			TestId:         2,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds temperature problem",
			Purpose:        "check temperature is within 60-80 range, with a deviation of 8",
			Severity:       2,
			OperatorFailure: []string{
				"can't ignite",
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
			MatchingField: "TemperatureCelsius",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname"),
				Value: []byte("server02"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       60,
					Miximum:       80,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_VALUE,
						Value: 8,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_NO_WINDOW,
					WindowSize:     5,
					FailuresCCount: 2,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_LATCH_UNTIL_RESET,
			IndicationSeconds: 0,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"TemperatureCelsius group",
				"groupRafael",
				"group general",
			},
			MasksOtherGroup: []string{
				"group1",
			},
		},
	},
}

var failure2 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "pressure test",
			TestId:         2,
			BitType:        []string{"CBIT, PBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds air pressure problem",
			Purpose:        "check pressure is out of 0-20 range, with a deviation of 16 percent",
			Severity:       24,
			OperatorFailure: []string{
				"OperatorFailure1",
				"OperatorFailure2",
				"OperatorFailure3",
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
			MatchingField: "AirPressure",
			MatchingTag: &KeyValuePair{
				Key:   []byte("hostname123"),
				Value: []byte("north"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       0,
					Miximum:       20,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_OUTOF,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_PERCENT,
						Value: 16,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     3,
					FailuresCCount: 2,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_LATCH_FOREVER,
			IndicationSeconds: 0,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group temp",
				"group1",
			},
			MasksOtherGroup: []string{
				"group2",
			},
		},
	},
	Time: time.Now(),
}

var failure3 = ExtendedFailure{
	Failure: Failure{
		Description: &FailureDescription{
			UnitName:       "system test check",
			TestName:       "oil test",
			TestId:         1,
			BitType:        []string{"CBIT"},
			Description:    "this is a mock failure to test services",
			AdditionalInfo: "the failure finds oil problem",
			Purpose:        "check oil is within 0-9 range, with a deviation of 1",
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
			MatchingField: "oil",
			MatchingTag: &KeyValuePair{
				Key:   []byte("ss"),
				Value: []byte("longstanding"),
			},
			FailureCriteria: &FailureExaminationRule_FailureCriteria{
				ValueCriteria: &FailureExaminationRule_FailureCriteria_FailureValueCriteria{
					Minimum:       0,
					Miximum:       9,
					ThresholdMode: FailureExaminationRule_FailureCriteria_FailureValueCriteria_WITHIN,
					Exceeding: &FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding{
						Type:  FailureExaminationRule_FailureCriteria_FailureValueCriteria_Exceeding_VALUE,
						Value: 1,
					},
				},
				TimeCriteria: &FailureExaminationRule_FailureCriteria_FailureTimeCriteria{
					WindowType:     FailureExaminationRule_FailureCriteria_FailureTimeCriteria_SLIDING,
					WindowSize:     1,
					FailuresCCount: 1,
				},
			},
		},
		ReportDuration: &FailureReportDuration{
			Indication:        FailureReportDuration_NUM_OF_SECONDS,
			IndicationSeconds: 3,
		},
		Dependencies: &Failure_FailureDependencies{
			BelongsToGroup: []string{
				"group1",
			},
			MasksOtherGroup: []string{
				"group4",
			},
		},
	},
}

var report0 = TestReport{
	TestId:         123,
	ReportPriority: 12,
	Timestamp:      time.Now(),
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
