package bitquery

import (
	rh "github.com/edendoron/bit-framework/internal/apiResponseHandlers"
	"net/http"
)

func BitStatusQuery(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	failed, userGroup := bitStatusRequestHandler(r, req)
	if failed {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad Request, check query params", err)
		return
	}

	bitStatusQueryHandler(w, req, userGroup)
}

func ReportQuery(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	failed, filter, values := reportsRequestHandler(r, req)
	if failed {
		rh.ApiResponseHandler(w, http.StatusBadRequest, "Bad Request, check query params", err)
		return
	}

	reportsQueryHandler(w, req, filter, values)
}

func UserGroupQuery(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, Configs.StorageReadURL, nil)
	if err != nil {
		rh.ApiResponseHandler(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := req.URL.Query()
	params.Add("user_groups", "")

	req.URL.RawQuery = params.Encode()

	userGroupsQueryHandler(w, req)
}
