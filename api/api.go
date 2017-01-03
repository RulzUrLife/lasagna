package api

import (
	"encoding/json"
	"fmt"
	"github.com/RulzUrLife/lasagna/config"
	"net/http"
	"strconv"
	"strings"
)

type Endpoint struct {
	Name         string
	ListTemplate string
	GetTemplate  string
	List         func() (interface{}, error)
	Get          func(int) (interface{}, error)
}

func (e *Endpoint) get(w http.ResponseWriter, r *http.Request) {
	if url := r.URL.Path[len(e.Name)+1:]; url == "" {
		http.Redirect(w, r, "/ingredients", http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	} else if data, err := e.Get(i); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		e.dump(w, r, e.GetTemplate, data)
	}
}

func (e *Endpoint) list(w http.ResponseWriter, r *http.Request) {
	if data, err := e.List(); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		e.dump(w, r, e.ListTemplate, data)
	}
}

func (e *Endpoint) dump(w http.ResponseWriter, r *http.Request, tplt string, data interface{}) {
	accepts := r.Header["Accept"]
	for _, media := range strings.Split(accepts[0], ",") {
		switch media {
		case "text/html":
			w.Header().Set("Content-Type", "text/html")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type ServeMux struct{ *http.ServeMux }

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *ServeMux) HandleEndpoint(endpoint *Endpoint) {
	mux.HandleFunc(endpoint.Name, endpoint.list)
	mux.HandleFunc(fmt.Sprintf("%s/", endpoint.Name), endpoint.get)
}
