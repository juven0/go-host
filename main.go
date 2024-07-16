package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	customTransport = http.DefaultTransport
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		log.Printf("Error creating proxy request: %v", err)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		log.Printf("Error forwarding request: %v", err)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

func main() {
	remote, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			w.Header().Set("X-Ben", "Rad")
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
