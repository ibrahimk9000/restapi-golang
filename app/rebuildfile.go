package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Success struct {
	IsDisallowed bool   `json:"IsDisallowed"`
	ErrorMessage string `json:"ErrorMessage"`
}

type Badrequest struct {
	File []string `json:"file"`
}

type Forbidden struct {
	Message string `json:"message"`
}
type JsonMeta struct {
	Metadata           int `json:"Metadata"`
	InternalHyperlinks int `json:"InternalHyperlinks"`
	ExternalHyperlinks int `json:"ExternalHyperlinks"`
	EmbeddedFiles      int `json:"EmbeddedFiles"`
	EmbeddedImages     int `json:"EmbeddedImages"`
	Javascript         int `json:"Javascript"`
	Acroform           int `json:"Acroform"`
	ActionsAll         int `json:"ActionsAll"`
}
type ContentManagementFlagJson struct {
	PdfContentManagement        JsonMeta `json:"PdfContentManagement"`
	ExcelContentManagement      JsonMeta `json:"ExcelContentManagement"`
	PowerPointContentManagement JsonMeta `json:"PowerPointContentManagement"`
	WordContentManagement       JsonMeta `json:"WordContentManagement"`
}

type RebuildUrl struct {
	InputGetUrl               string                    `json:"InputGetUrl"`
	OutputPutUrl              string                    `json:"OutputPutUrl"`
	ContentManagementFlagJson ContentManagementFlagJson `json:"ContentManagementFlagJson"`
}

func rebuildfile(w http.ResponseWriter, r *http.Request) {

	//m max 5 MB file name we can change ut
	r.ParseMultipartForm(5 << 20)

	log.Printf("json payload : %v\n", r.PostFormValue("contentManagementFlagJson"))

	cont := r.PostFormValue("contentManagementFlagJson")

	var mp map[string]json.RawMessage

	err := json.Unmarshal([]byte(cont), &mp)
	if err != nil {
		log.Println("unmarshal json:", err)
		http.Error(w, "malformed json format", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")

	if err != nil {
		log.Println("formfile", err)
		http.Error(w, "file not found or wrong form field  name", http.StatusBadRequest)

		return
	}

	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("ioutilReadAll", err)
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}

	//uploaded file log info
	log.Printf("Filename: %v\n", handler.Filename)
	log.Printf("File size: %v\n", handler.Size)
	log.Printf("Content-Type: %v\n", handler.Header.Get("Content-Type"))
	log.Printf("Content-Type: %v\n", http.DetectContentType(buf))

	//glaswall custom header
	addgwheader(w, temp)

	_, e := w.Write(buf)
	if e != nil {
		log.Println(e)

	}

}
