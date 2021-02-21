package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type req struct {
	Base64                 string `json:"Base64"`
	ContentManagementFlags struct {
		PdfContentManagement struct {
			Metadata int `json:"Metadata"`
		} `json:"PdfContentManagement"`
	} `json:"ContentManagementFlags"`
}

func rebuildzip(w http.ResponseWriter, r *http.Request) {

	//m max 5 MB file name we can change ut
	r.ParseMultipartForm(5 << 20)

	//handling json , not implemeted yet
	log.Println(r.PostFormValue("contentManagementFlagJson"))

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

	if handler.Header.Get("Content-Type") != "application/zip" || http.DetectContentType(buf) != "application/zip" {
		log.Println("mediatype is", handler.Header.Get("Content-Type"))

		http.Error(w, "uploaded file should be zip format", http.StatusUnsupportedMediaType)

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
