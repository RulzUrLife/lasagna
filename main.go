package main

import (
	"fmt"
	"net/http"
	"time"
)

type Router struct{}

func (_ *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	addr := fmt.Sprintf("%s:%d", Config.Host, Config.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        &Router{},
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	Info.Printf("* Running on http://%s/\n", addr)
	Error.Fatal(s.ListenAndServe())
}
