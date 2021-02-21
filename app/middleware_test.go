package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func ExampleResponseRecorder() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Output:
	// 200
	// text/html; charset=utf-8
	// <html><body>Hello World!</body></html>
}

func ExampleServer() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
	// Output: Hello, client
}

func TestAuthMiddleware(t *testing.T) {

	var tests = []struct {
		hk     string
		hv     string
		status int
	}{
		{"Authorization", "Bearer mysecrettoken", http.StatusOK},
		{"malformed", "Bearer mysecrettoken", http.StatusUnauthorized},
		{"Authorization", "earer mysecrettoken", http.StatusUnauthorized},
		{"Authorization", "Bearer notmysecrettoken", http.StatusUnauthorized},
		{"", "", http.StatusUnauthorized},
	}

	for _, test := range tests {
		if output, _ := middlewareconnect(test.hk, test.hv); output != test.status {
			t.Errorf("Test Failed: {%s} header, {%s} header value, status: {%d},expected: {%d}", test.hk, test.hv, test.status, output)
		}
	}
}

func middlewareconnect(hk, hv string) (int, string) {

	endpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "autorization succeful")
	})

	ts := httptest.NewServer(customMiddleware(endpoint))

	client := &http.Client{}

	req, err := http.NewRequest("POST", ts.URL, strings.NewReader("empty message"))
	if hk != "" || hv != "" {
		req.Header.Set(hk, hv)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	status := resp.StatusCode
	statusm := resp.Status

	return status, statusm
}

/*
   var tests = []struct {
       input    int
       expected int
   }{
       {2, 4},
       {-1, 1},
       {0, 2},
       {-5, -3},
       {99999, 100001},
   }

   for _, test := range tests {
       if output := Calculate(test.input); output != test.expected {
           t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.input, test.expected, output)
       }
   }
*/
