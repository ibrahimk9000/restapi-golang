package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request) {

	//m max 5 MB file name we can change ut
	r.ParseMultipartForm(5 << 20)

	base64enc := r.Body

	defer base64enc.Close()

	cont, err := ioutil.ReadAll(base64enc)
	if err != nil {
		log.Println("ioutilReadAll", err)
		http.Error(w, "empty or malformed request body", http.StatusBadRequest)

		return
	}

	var mp map[string]json.RawMessage

	err = json.Unmarshal(cont, &mp)
	if err != nil {
		log.Println("unmarshal json", err)
		http.Error(w, "malformed json format", http.StatusBadRequest)

		return
	}

	var str string
	err = json.Unmarshal(mp["Base64"], &str)
	if err != nil {
		log.Println("unmarshal json base64", err)
		http.Error(w, "malformed json format ", http.StatusBadRequest)

		return
	}

	buf, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println("base64 decoding", err)
		http.Error(w, "malformed base64 encoding", http.StatusBadRequest)
		return

	}

	log.Printf("%v\n", http.DetectContentType(buf))

	//glasswall custom header
	addgwheader(w, temp)

	_, e := w.Write(buf)
	if e != nil {
		log.Println(e)

	}

}
