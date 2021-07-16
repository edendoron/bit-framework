package bitQuery

import (
	"encoding/json"
	. "github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	. "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	. "github.com/edendoron/bit-framework/internal/models"
	"log"
	"net/http"
)

// bitStatusQueryHandler send @req to storage in order to read bitStatus reports, and writes it back to user through @w
func bitStatusQueryHandler(w http.ResponseWriter, req *http.Request, userGroup string) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer resp.Body.Close()
	var bitStatusList []BitStatus
	err = json.NewDecoder(resp.Body).Decode(&bitStatusList)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	maskedTestIds, e := getUserGroupsFiltering(userGroup)
	if e != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error extract user groups filtering", e)
	} else {
		FilterBitStatus(&bitStatusList, maskedTestIds)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(&bitStatusList)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// reportsQueryHandler send @req to storage in order to read reports, and writes it back to user through @w
func reportsQueryHandler(w http.ResponseWriter, req *http.Request, filter string, values []string) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer resp.Body.Close()
	var reports []TestReport
	err = json.NewDecoder(resp.Body).Decode(&reports)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	FilterReports(&reports, filter, values)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(&reports)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// userGroupsQueryHandler send @req to storage in order to read userGroups filtering rules, and writes it back to user through @w
func userGroupsQueryHandler(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer resp.Body.Close()
	var userGroups []string
	err = json.NewDecoder(resp.Body).Decode(&userGroups)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(&userGroups)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// add bitStatus query parameters according to params of request @r, return userGroup of request
func bitStatusRequestHandler(r *http.Request, req *http.Request) (bool, string) {

	params := req.URL.Query()

	start := r.URL.Query()["start"]
	end := r.URL.Query()["end"]
	userGroup := r.URL.Query()["user_group"]

	if len(start) == 0 || len(end) == 0 || len(userGroup) == 0 {
		return true, ""
	}
	params.Add("start", start[0])
	params.Add("end", end[0])
	params.Add("bit_status", "")

	req.URL.RawQuery = params.Encode()
	return false, userGroup[0]
}

// add report query parameters according to params of request @r, return filter type and values of request
func reportsRequestHandler(r *http.Request, req *http.Request) (bool, string, []string) {

	params := req.URL.Query()

	start := r.URL.Query()["start"]
	end := r.URL.Query()["end"]
	filter := r.URL.Query()["filter"]
	var values []string

	if len(start) == 0 || len(end) == 0 || len(filter) == 0 || (filter[0] != "time" && filter[0] != "tag" && filter[0] != "field") {
		return true, "", values
	}

	if filter[0] == "tag" {
		tagKey := r.URL.Query()["tag_key"]
		tagValue := r.URL.Query()["tag_value"]
		if len(tagKey) == 0 || len(tagValue) == 0 {
			return true, "", values
		}
		values = append(values, tagKey[0], tagValue[0])
	} else if filter[0] == "field" {
		field := r.URL.Query()["field"]
		if len(field) == 0 {
			return true, "", values
		}
		values = append(values, field[0])
	}

	params.Add("reports", "")
	params.Add("start", start[0])
	params.Add("end", end[0])
	params.Add("filter", "time")

	req.URL.RawQuery = params.Encode()
	return false, filter[0], values
}

/*
send request of user groups filtering rules to storage
return masked tests ID's according to @param userGroup, and error if occur
*/
func getUserGroupsFiltering(userGroup string) ([]uint64, error) {
	var res []uint64

	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		return res, err
	}

	params := req.URL.Query()
	params.Add("config_user_groups_filtering", "")
	params.Add("id", userGroup)

	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		return res, err
	}
	defer storageResponse.Body.Close()

	err = json.NewDecoder(storageResponse.Body).Decode(&res)
	if err != nil {
		log.Printf("error decode storage response body")
		return res, err
	}

	return res, nil
}

// FilterBitStatus filter failures from @param statusList according to @param maskedTestIds
func FilterBitStatus(statusList *[]BitStatus, maskedTestIds []uint64) {
	idx := 0
	for _, status := range *statusList {
		n := 0
		found := false
		for _, failure := range status.Failures {
			for i := range maskedTestIds {
				if failure.FailureData.TestId == maskedTestIds[i] {
					found = true
					break
				}
			}
			if !found {
				status.Failures[n] = failure
				n++
			}
			found = false
		}
		if n == 0 {
			*statusList = append((*statusList)[:idx], (*statusList)[idx+1:]...)
		} else {
			(*statusList)[idx].Failures = status.Failures[:n]
			idx++
		}
	}
}

// FilterReports filter reports from @param reportList according to @params filter, values
func FilterReports(reportList *[]TestReport, filter string, values []string) {
	if filter == "time" {
		return
	}
	n := 0
	for _, report := range *reportList {
		switch filter {
		case "tag":
			for _, tagSet := range report.TagSet {
				if values[0] == tagSet.Key && values[1] == tagSet.Value {
					(*reportList)[n] = report
					n++
					break
				}
			}
		case "field":
			for _, fieldSet := range report.FieldSet {
				if values[0] == fieldSet.Key {
					(*reportList)[n] = report
					n++
					break
				}
			}
		}
	}
	*reportList = (*reportList)[:n]
}
