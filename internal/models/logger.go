package models

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// RedirectLogger redirects logs to each service unique path
func RedirectLogger(servicePath string) {
	currTime := time.Now().Format("20060102150405")

	fileName := servicePath + "/log/" + currTime + "_logs.txt"
	file, e := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		fmt.Println(e)
	}
	log.SetOutput(file)
}
