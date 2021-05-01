package bitStorageAccess

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
)

func convertToKeyValuePair(arr []KeyValue) []*KeyValuePair {
	copyArr := make([]*KeyValuePair, len(arr))
	for k, v := range arr {
		pair := KeyValuePair{Key: []byte(v.Key), Value: []byte(v.Value)}
		copyArr[k] = &pair
	}
	return copyArr
}

func convertToKeyValue(arr []*KeyValuePair) []KeyValue {
	copyArr := make([]KeyValue, len(arr))
	for k, v := range arr {
		pair := KeyValue{Key: string(v.Key), Value: string(v.Value)}
		copyArr[k] = pair
	}
	return copyArr
}

//func testResultToTestReport(tr TestResult) TestReport{
//
//}
//
//func testResultToTestReport(tr TestResult) TestReport{
//
//}
