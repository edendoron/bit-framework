package bitConfig

import (
	. "../../configs/rafael.com/bina/bit"
	. "../models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const storageURL = "http://localhost:8082/data/write"

func PostFailuresData() {
	failure := Failure{}
	files, err := ioutil.ReadDir("/failures")
	if err != nil {
		//TODO: handle error
		return
	}
	//TODO: handle error
	for _, f := range files {
		fmt.Println(f.Name())
		//TODO: handle error
		content, _ := ioutil.ReadFile(f.Name())

		//TODO: handle error
		err = json.Unmarshal(content, &failure)
		//TODO: handle error
		err = ValidateType(failure)

		message := KeyValue{Key: "failures", Value: string(content)}
		//TODO: handle error
		postBody, _ := json.MarshalIndent(message, "", " ")
		postBodyBuf := bytes.NewReader(postBody)
		storageResp, e := http.Post(storageURL, "application/json; charset=UTF-8", postBodyBuf)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			//TODO: handle this error
			return
		}
		//TODO: handle error
		storageResp.Body.Close()
	}
}
