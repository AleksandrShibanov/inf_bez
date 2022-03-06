package main

import (
	"io"
	"log"
	"net/http"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Proxy(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	proxyReq, _ := http.NewRequest(r.Method, r.URL.Path, nil)
	proxyReq.Host = r.URL.Host

	copyHeader(proxyReq.Header, r.Header)
	proxyReq.Header.Del("Proxy-Connection")

	return http.DefaultTransport.RoundTrip(r)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := Proxy(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
		}),
	}

	log.Fatal(server.ListenAndServe())
}
