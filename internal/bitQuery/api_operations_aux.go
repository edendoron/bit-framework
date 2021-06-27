package bitQuery

import (
	. "../../configs/rafael.com/bina/bit"
	. "../apiResponseHandlers"
	. "../models"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const storageDataReadURL = "http://localhost:8082/data/read"

func QueryHandler(w http.ResponseWriter, req *http.Request, requestedData string, userGroup string) {
	client := &http.Client{}
	//fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	switch requestedData {
	case "bit_status":
		var bitStatusList []BitStatus
		err = json.NewDecoder(resp.Body).Decode(&bitStatusList)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		e := filterBitStatus(&bitStatusList, userGroup)
		if e != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(&bitStatusList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
	case "reports":
		var reports []TestReport
		err = json.NewDecoder(resp.Body).Decode(&reports)
		if err != nil {
			ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(&reports)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
	}

	err = resp.Body.Close()
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
	}
}

func paramsHandler(r *http.Request, params url.Values, filter string) {
	switch filter {
	// cases are almost identical, query keys and var names are different for readability - need support on client side
	case "time":
		start := r.URL.Query()["start"][0]
		end := r.URL.Query()["end"][0]
		params.Add("start", start)
		params.Add("end", end)
	case "tag":
		tag := r.URL.Query()["tag"]
		params.Add("tag_key", tag[0])
		params.Add("tag_value", tag[1])
	case "field":
		field := r.URL.Query()["field"]
		params.Add("field", field[0])
	}
}

func getUserGroupsFiltering(userGroup string) ([]uint64, error) {
	var res []uint64

	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		return res, err
	}
	//TODO: defer req.Body.Close()

	params := req.URL.Query()
	params.Add("config_user_groups_filtering", "")
	params.Add("id", userGroup)

	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	storageResponse, err := client.Do(req)
	if err != nil || storageResponse.StatusCode != http.StatusOK {
		return res, err
	}

	err = json.NewDecoder(storageResponse.Body).Decode(&res)
	if err != nil {
		log.Printf("error decode storage response body")
		return res, err
	}
	err = storageResponse.Body.Close()
	if err != nil {
		log.Printf("error close storage response body")
		return res, err
	}
	return res, nil
}

func filterBitStatus(statusList *[]BitStatus, userGroup string) error {
	maskedTestIds, err := getUserGroupsFiltering(userGroup)
	if err != nil {
		return err
	}

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
		}
		status.Failures = status.Failures[:n]
	}
	return nil
}
