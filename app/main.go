package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	GwMetricFormFileRead = "gw-metric-formfileread"
	GwMetricFileSize     = "gw-metric-filesize"
	GwVersion            = "gw-version"
	GwMetricDetect       = "gw-metric-detect"
	GwMetricRebuild      = "gw-metric-rebuild"
)

var (
	temp = gwcustomheader{"0.01", "5 mb", "1.39", "0.02", "0.03"}
)

type gwcustomheader struct {
	metricFormFileread string
	metricFileSize     string
	version            string
	metricDetect       string
	metricRebuild      string
}

func addgwheader(w http.ResponseWriter, v gwcustomheader) {
	w.Header().Set(GwMetricFormFileRead, v.metricFormFileread)
	w.Header().Set(GwMetricFileSize, v.metricFileSize)
	w.Header().Set(GwVersion, v.version)
	w.Header().Set(GwMetricDetect, v.metricDetect)
	w.Header().Set(GwMetricRebuild, v.metricRebuild)

}

//there will middleware chain here
func customMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		errauth := "you don't have  valid authoriaztion token"
		erremptyauth := "you didn't provide authoriaztion token"

		//log about request
		//there will be logging middleware soon
		log.Printf("method: %v\n", r.Method)
		log.Printf("URL: %v\n", r.URL)
		log.Printf("RemoteAddr: %v\n", r.RemoteAddr)
		log.Printf("Host: %v\n", r.Host)
		log.Printf("Content-Type: %v\n", r.Header.Get("Content-Type"))
		log.Printf("RequestURI: %v\n", r.RequestURI)

		//Authorization: Bearer
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, erremptyauth, http.StatusUnauthorized)

			return
		}

		authHeaderParts := strings.Fields(authHeader)

		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, errauth, http.StatusUnauthorized)

			return
		}

		if authHeaderParts[1] != "mysecrettoken" {

			http.Error(w, errauth, http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mx := mux.NewRouter()
	mx.Use(customMiddleware)

	mx.HandleFunc("/api/rebuild/file", rebuildfile).Methods("POST")
	mx.HandleFunc("/api/rebuild/zip", rebuildzip).Methods("POST")
	mx.HandleFunc("/api/rebuild/base64", rebuildbase64).Methods("POST")
	//mx.HandleFunc("/api/rebuild/s3tos3", rebuilds3tos3).Methods("POST")

	fmt.Println("Server is ready to handle requests at port 8100")
	log.Fatal(http.ListenAndServe(":8100", mx))
}
