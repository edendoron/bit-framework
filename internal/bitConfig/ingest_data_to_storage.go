package bitConfig

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func PostFailuresData() {
	failure := Failure{}
	files, err := ioutil.ReadDir(Configs.BitConfigFailuresPath)
	if err != nil {
		log.Println("error reading config_failures dir")
		log.Println(err)
		return
	}
	for i, f := range files {
		content, e := ioutil.ReadFile(Configs.BitConfigFailuresPath + f.Name())
		if e != nil {
			log.Println("error reading config_failures file number ", i+1)
			log.Println(e)
			return
		}
		e = json.Unmarshal(content, &failure)
		//TODO: check 64bit encoding problem with protobuff
		if e != nil {
			log.Println("error unmarshal config_failures file number ", i+1)
			log.Println(e)
			return
		}
		e = ValidateType(failure)
		if e != nil {
			log.Println("error validate config_failures file number ", i+1)
			log.Println(e)
			return
		}
		message := KeyValue{Key: "config_failure", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Println("error marshal config_failures file number ", i+1)
			log.Println(e)
			return
		}
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Println("error post to storage config_failures file number ", i+1)
			log.Println(e)
			return
		}
		e = storageResp.Body.Close()
		if e != nil {
			log.Println("error close request body ", i+1)
			log.Println(e)
			return
		}
	}
}

func PostGroupFilterData() {
	groupFilter := UserGroupsFiltering_FilteredFailures{}
	files, err := ioutil.ReadDir(Configs.BitConfigUserGroupPath)
	if err != nil {
		log.Println("error reading config_user_groups_filtering dir")
		log.Println(err)
		return
	}
	for i, f := range files {
		content, e := ioutil.ReadFile(Configs.BitConfigUserGroupPath + f.Name())
		if e != nil {
			log.Println("error reading config_user_groups_filtering file number ", i+1)
			log.Println(e)
			return
		}
		e = json.Unmarshal(content, &groupFilter)
		if e != nil {
			log.Println("error unmarshal config_user_groups_filtering file number ", i+1)
			log.Println(e)
			return
		}
		e = ValidateType(groupFilter)
		if e != nil {
			log.Println("error validate config_user_groups_filtering file number ", i+1)
			log.Println(e)
			return
		}
		message := KeyValue{Key: "config_user_group_filtering", Value: string(content)}
		postBody, e := json.MarshalIndent(message, "", " ")
		if e != nil {
			log.Println("error marshal config_user_groups_filtering file number ", i+1)
			log.Println(e)
			return
		}
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(Configs.StorageWriteURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			log.Println("error post to storage config_user_groups_filtering file number ", i+1)
			log.Println(e)
			return
		}
		e = storageResp.Body.Close()
		if e != nil {
			log.Println("error close request body ", i+1)
			log.Println(e)
			return
		}
	}
}
