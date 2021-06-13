package bitConfig

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const storageURL = "http://localhost:8082/data/write"

func PostFailuresData() {
	failure := Failure{}
	files, err := ioutil.ReadDir("./configs/config_failures")
	if err != nil {
		log.Println("error reading config_failures dir")
		return
	}
	for i, f := range files {
		fmt.Println(f.Name())
		content, e := ioutil.ReadFile("./configs/config_failures/" + f.Name())
		if e != nil {
			log.Println("error reading config_failures file number ", i)
			return
		}
		e = json.Unmarshal(content, &failure)
		if e != nil {
			log.Println("error unmarshal config_failures file number ", i)
			return
		}
		e = ValidateType(failure)
		if e != nil {
			log.Println("error validate config_failures file number ", i)
			return
		}
		message := KeyValue{Key: "config_failure", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Println("error marshal config_failures file number ", i)
			return
		}
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(storageURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Println("error post to storage config_failures file number ", i)
			return
		}
		e = storageResp.Body.Close()
		if e != nil {
			log.Println("error close request body ", i)
			return
		}
	}
}

func PostGroupFilterData() {
	groupFilter := UserGroupsFiltering_FilteredFailures{}
	files, err := ioutil.ReadDir("./configs/config_user_groups_filtering")
	if err != nil {
		log.Println("error reading config_user_groups_filtering dir")
		return
	}
	for i, f := range files {
		fmt.Println(f.Name())
		content, e := ioutil.ReadFile("./configs/config_user_groups_filtering/" + f.Name())
		if e != nil {
			log.Println("error reading config_user_groups_filtering file number ", i)
			return
		}
		e = json.Unmarshal(content, &groupFilter)
		if e != nil {
			log.Println("error unmarshal config_user_groups_filtering file number ", i)
			return
		}
		e = ValidateType(groupFilter)
		if e != nil {
			log.Println("error validate config_user_groups_filtering file number ", i)
			return
		}
		message := KeyValue{Key: "config_user_group_filtering", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Println("error marshal config_user_groups_filtering file number ", i)
			return
		}
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(storageURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Println("error post to storage config_user_groups_filtering file number ", i)
			return
		}
		e = storageResp.Body.Close()
		if e != nil {
			log.Println("error close request body ", i)
			return
		}
	}
}
