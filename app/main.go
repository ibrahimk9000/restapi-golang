package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {
	mx := mux.NewRouter()
	mx.Use(Logmiddleware, AuthMiddleware)

	mx.HandleFunc("/api", Get).Methods("GET")
	mx.HandleFunc("/api", Post).Methods("POST")
	mx.HandleFunc("/api", Put).Methods("PUT")
	//mx.HandleFunc("/api/rebuild/s3tos3", rebuilds3tos3).Methods("POST")

	fmt.Println("Server is ready to handle requests at port 8100")
	log.Fatal(http.ListenAndServe(":8100", mx))
}
