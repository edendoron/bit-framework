package bitQuery

import (
	. "../apiResponseHandlers"
	"net/http"
)

func BitStatusQuery(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	filter := r.URL.Query()["filter"][0]
	userGroup := r.URL.Query()["user_group"][0]

	params := req.URL.Query()
	params.Add("bit_status", "")
	params.Add("filter", filter)

	paramsHandler(r, params, filter)

	req.URL.RawQuery = params.Encode()

	bitStatusQueryHandler(w, req, "bit_status", userGroup)
}

func ReportQuery(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, storageDataReadURL, nil)
	if err != nil {
		ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	filter := r.URL.Query()["filter"][0]

	params := req.URL.Query()
	params.Add("report", "")
	params.Add("filter", filter)

	paramsHandler(r, params, filter)

	req.URL.RawQuery = params.Encode()

	bitStatusQueryHandler(w, req, "report", "")
}
