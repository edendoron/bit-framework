package bitconfig

import (
	"bytes"
	"encoding/json"
	"github.com/edendoron/bit-framework/configs/rafael.com/bina/bit"
	"github.com/edendoron/bit-framework/internal/models"
	"io/ioutil"
	"log"
	"net/http"
)

// PostFailuresData writes the configuration failures found in Configs.BitConfigFailuresPath to storage
func PostFailuresData() {
	failure := bit.Failure{}
	files, err := ioutil.ReadDir(Configs.BitConfigFailuresPath)
	if err != nil {
		log.Fatal("error reading config_failures dir error: ", err)
	}
	for i, f := range files {
		content, e := ioutil.ReadFile(Configs.BitConfigFailuresPath + f.Name())
		if e != nil {
			log.Printf("error reading config_failures file number %v, error: %v", i+1, e)
			continue
		}
		e = json.Unmarshal(content, &failure)
		if e != nil {
			log.Printf("error unmarshal config_failures file number %v, error: %v", i+1, e)
			continue
		}
		e = models.ValidateType(failure)
		if e != nil {
			log.Printf("error validate config_failures file number %v, error: %v", i+1, e)
			continue
		}
		message := models.KeyValue{Key: "config_failure", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Printf("error marshal config_failures file number %v, error: %v", i+1, e)
			continue
		}

		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Printf("error post to storage config_failures file number %v, error: %v", i+1, e)
		}
		storageResp.Body.Close()
	}
}

// PostGroupFilterData writes the user group filtering rules found in Configs.BitConfigUserGroupPath to storage
func PostGroupFilterData() {
	groupFilter := bit.UserGroupsFiltering_FilteredFailures{}
	files, err := ioutil.ReadDir(Configs.BitConfigUserGroupPath)
	if err != nil {
		log.Printf("error reading config_user_groups_filtering dir. error: %v", err)
		return
	}
	for i, f := range files {
		content, e := ioutil.ReadFile(Configs.BitConfigUserGroupPath + f.Name())
		if e != nil {
			log.Printf("error reading config_user_groups_filtering file number %v, error: %v", i+1, e)
			continue
		}
		e = json.Unmarshal(content, &groupFilter)
		if e != nil {
			log.Printf("error unmarshal config_user_groups_filtering file number %v, error: %v", i+1, e)
			continue
		}
		e = models.ValidateType(groupFilter)
		if e != nil {
			log.Printf("error validate config_user_groups_filtering file number %v, error: %v", i+1, e)
			continue
		}
		message := models.KeyValue{Key: "config_user_group_filtering", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Printf("error marshal config_user_groups_filtering file number %v, error: %v", i+1, e)
			continue
		}
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Printf("error post to storage config_user_groups_filtering file number %v, error: %v", i+1, e)
		}
		storageResp.Body.Close()
	}
}
