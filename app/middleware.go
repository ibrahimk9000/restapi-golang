package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

//there will middleware chain here

func Logmiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", "my-service").
			Str("host", "host").
			Logger()

		c := alice.New()

		// Install the logger handler with default output on the console
		c = c.Append(hlog.NewHandler(log))

		// Install some provided extra handler to set some request's context fields.
		// Thanks to that handler, all our logs will come with some prepopulated fields.
		c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))
		c = c.Append(hlog.RemoteAddrHandler("ip"))
		c = c.Append(hlog.UserAgentHandler("user_agent"))
		c = c.Append(hlog.RefererHandler("referer"))
		c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

		// Here is your final handler
		h := c.Then(next)
		h.ServeHTTP(w, r)
		// Output: {"level":"info","time":"2001-02-03T04:05:06Z","role":"my-service","host":"local-hostname","req_id":"b4g0l5t6tfid6dtrapu0","user":"current user","status":"ok","message":"Something happened"
	})
}

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		errauth := "you don't have  valid authoriaztion token"
		erremptyauth := "you didn't provide authoriaztion token"

		//log about request
		//there will be logging middleware soon
		//log.Printf("method: %v\n", r.Method)
		//log.Printf("URL: %v\n", r.URL)
		//log.Printf("RemoteAddr: %v\n", r.RemoteAddr)
		//log.Printf("Host: %v\n", r.Host)
		//log.Printf("Content-Type: %v\n", r.Header.Get("Content-Type"))
		//log.Printf("RequestURI: %v\n", r.RequestURI)

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
