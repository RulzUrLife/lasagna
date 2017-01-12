package api

import (
	"encoding/json"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

type ListMethod func() (interface{}, *common.HTTPError)
type GetMethod func(int) (interface{}, *common.HTTPError)

type Templates struct {
	HTML *template.Template
	XML  *template.Template
}

type Endpoint struct {
	Name string
	List struct {
		Method ListMethod
		*Templates
	}
	Get struct {
		Method GetMethod
		*Templates
	}
}

var base = path.Join("templates", "base.html")
var errTplt = &Templates{
	template.Must(template.ParseFiles(base, path.Join("templates", "error.html"))),
	nil,
}

type ResponseWriter interface {
	Render(w http.ResponseWriter, data interface{})
	Error(w http.ResponseWriter, err *common.HTTPError)
}

type HTMLResponseWriter struct {
	*Endpoint
	*template.Template
	ErrorTemplate *template.Template
}

func (rw *HTMLResponseWriter) Render(w http.ResponseWriter, data interface{}) {
	rw.Template.Execute(w, data)
}

func (rw *HTMLResponseWriter) Error(w http.ResponseWriter, err *common.HTTPError) {
	w.WriteHeader(err.Code)
	rw.ErrorTemplate.Execute(w, err)
}

type JSONResponseWriter struct {
	*Endpoint
}

func (rw *JSONResponseWriter) Render(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func (rw *JSONResponseWriter) Error(w http.ResponseWriter, err *common.HTTPError) {
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}

func (e *Endpoint) parseHeaders(w http.ResponseWriter, r *http.Request, t *Templates) ResponseWriter {
	accepts := r.Header["Accept"]
	for _, media := range strings.Split(accepts[0], ",") {
		media = strings.Split(media, ";")[0]
		switch media {
		case "text/html":
			common.Trace.Printf("HTML rendering")
			w.Header().Set("Content-Type", "text/html")
			return &HTMLResponseWriter{e, t.HTML, errTplt.HTML}
		case "application/xhtml+xml":
			common.Trace.Printf("XML rendering")
			w.Header().Set("Content-Type", "application/xhtml+xml")
			return nil
		}
	}
	common.Trace.Printf("JSON rendering")
	w.Header().Set("Content-Type", "application/json")
	return &JSONResponseWriter{e}
}

func (e *Endpoint) get(w http.ResponseWriter, r *http.Request) {
	rw := e.parseHeaders(w, r, e.Get.Templates)
	if url := r.URL.Path[len(e.Name)+1:]; url == "" {
		http.Redirect(w, r, e.Name, http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		rw.Error(w, common.NewHTTPError("400 Bad Request", http.StatusBadRequest))
	} else if data, err := e.Get.Method(i); err != nil {
		rw.Error(w, err)
	} else {
		rw.Render(w, data)
	}
}

func (e *Endpoint) list(w http.ResponseWriter, r *http.Request) {
	rw := e.parseHeaders(w, r, e.List.Templates)
	if data, err := e.List.Method(); err != nil {
		rw.Error(w, err)
	} else {
		common.Trace.Printf("%q", data)
		rw.Render(w, data)
	}
}

type ServeMux struct{ *http.ServeMux }

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	common.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *ServeMux) NewEndpoint(name string, list ListMethod, get GetMethod) {
	endpoint := &Endpoint{Name: name}
	dir := path.Join("templates", name)

	endpoint.List.Method = list
	endpoint.List.Templates = &Templates{
		template.Must(template.ParseFiles(base, path.Join(dir, "list.html"))),
		nil,
	}

	endpoint.Get.Method = get
	endpoint.Get.Templates = &Templates{
		template.Must(template.ParseFiles(base, path.Join(dir, "get.html"))),
		nil,
	}

	mux.HandleFunc(name, endpoint.list)
	mux.HandleFunc(fmt.Sprintf("%s/", name), endpoint.get)
}
