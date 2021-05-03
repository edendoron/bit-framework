package bitConfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const storageURL = "http://localhost:8082/failure/raw"

func PostFailuresData() {
	//failure := Failure{}
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
		//err = json.Unmarshal(content, &failure)

		//TODO: handle error
		//err = ValidateType(failure)

		postBody := bytes.NewReader(content)
		storageResp, e := http.Post(storageURL, "application/json; charset=UTF-8", postBody)
		if e != nil || storageResp.StatusCode != http.StatusOK {
			//TODO: handle this error
			return
		}
		//TODO: handle this error
		storageResp.Body.Close()
	}
}
