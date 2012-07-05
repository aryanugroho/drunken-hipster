package main

import (
	"net/http"
	"log"
	"io"
	"os"
	"fmt"
)

type RequestHandler struct {
	Transport *http.Transport
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("incoming request: %#v", *r)

	r.RequestURI = ""
	r.URL.Scheme = "http"
	r.URL.Host = "127.0.0.1:8000"

	resp, err := h.Transport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func main() {
	listen_addr := ":9000"

	mux := http.NewServeMux()
	mux.Handle("/", &RequestHandler{Transport:&http.Transport{DisableKeepAlives:false,DisableCompression:false}})

	srv := &http.Server{Handler: mux, Addr: listen_addr}

	err := srv.ListenAndServe()

	if err != nil {
		log.Printf("ListenAndServe: %v", err)
		os.Exit(1)
	}
}
