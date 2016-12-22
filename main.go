package main

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/api"
	"github.com/RulzUrLife/lasagna/config"
	"net/http"
	"time"
)

func main() {
	addr := fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        api.Mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	config.Info.Printf("Running on http://%s/\n", addr)
	config.Error.Fatal(s.ListenAndServe())
}
