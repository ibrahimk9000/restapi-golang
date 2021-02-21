package main

import (
	"log"
	"net/http"
	"strings"
)

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
