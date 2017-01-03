package main

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/api"
	"github.com/RulzUrLife/lasagna/config"
	"github.com/RulzUrLife/lasagna/db"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello World!")
}

func main() {

	config.Info.Println("Register URL patterns")
	mux := &api.ServeMux{http.NewServeMux()}
	mux.HandleFunc("/", index)
	mux.NewEndpoint("/ingredients",
		func() (interface{}, error) { return db.ListIngredients() },
		func(id int) (interface{}, error) { return db.GetIngredient(id) },
	)

	addr := fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	config.Info.Printf("Running on http://%s/\n", addr)
	config.Error.Fatal(s.ListenAndServe())
}
