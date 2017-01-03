package api

import (
	"encoding/json"
	"fmt"
	"github.com/RulzUrLife/lasagna/config"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

type ListMethod func() (interface{}, error)
type GetMethod func(int) (interface{}, error)

type Endpoint struct {
	Name     string
	List     ListMethod
	Get      GetMethod
	Template *template.Template
}

func (e *Endpoint) get(w http.ResponseWriter, r *http.Request) {
	if url := r.URL.Path[len(e.Name)+1:]; url == "" {
		http.Redirect(w, r, "/ingredients", http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	} else if data, err := e.Get(i); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		e.dump(w, r, e.Template.Lookup("get.html"), data)
	}
}

func (e *Endpoint) list(w http.ResponseWriter, r *http.Request) {
	if data, err := e.List(); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		e.dump(w, r, e.Template.Lookup("list.html"), data)
	}
}

func (e *Endpoint) dump(
	w http.ResponseWriter, r *http.Request, t *template.Template, data interface{},
) {
	accepts := r.Header["Accept"]
	for _, media := range strings.Split(accepts[0], ",") {
		media = strings.Split(media, ";")[0]
		switch media {
		case "text/html":
			config.Trace.Printf("Rendering %s", path.Join(e.Name[1:], t.Name()))
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, data)
			return
		case "application/xhtml+xml":
			w.Header().Set("Content-Type", "application/xhtml+xml")
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

func (mux *ServeMux) NewEndpoint(name string, list ListMethod, get GetMethod) {
	tplt := template.New(name)
	base := path.Join("templates", name[1:])

	tplt, err := tplt.ParseGlob(path.Join(base, "*"))
	if err != nil {
		config.Error.Fatalf("Parsing %s templates failed", name)
	} else {
		for _, t := range tplt.Templates() {
			config.Info.Printf("Parsing %s", path.Join(base, t.Name()))
		}
	}
	endpoint := &Endpoint{name, list, get, tplt}
	mux.HandleFunc(name, endpoint.list)
	mux.HandleFunc(fmt.Sprintf("%s/", name), endpoint.get)
}
