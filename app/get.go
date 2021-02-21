package main

import (
	"log"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {

	_, e := w.Write([]byte("get requset success"))
	if e != nil {
		log.Println(e)

	}

}
