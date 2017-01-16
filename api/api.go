package api

import (
	"encoding/json"
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

type List struct {
	Method ListMethod
	*Templates
}

type Get struct {
	Name   string
	Method GetMethod
	*Templates
}

var errTplt = templates("error.html")

func templates(paths ...string) *Templates {
	base := path.Join(common.Config.Assets.Templates, "base.html")
	p := path.Join(append([]string{common.Config.Assets.Templates}, paths...)...)
	funcs := template.FuncMap{
		"div": div, "slice": slice,
	}
	common.Info.Printf("Register template %s", p)
	return &Templates{
		template.Must(template.New("base.html").Funcs(funcs).ParseFiles(base, p)),
		nil,
	}
}

type ResponseWriter interface {
	Render(w http.ResponseWriter, data interface{}) error
	Error(w http.ResponseWriter, err *common.HTTPError) error
}

type HTMLResponseWriter struct {
	*template.Template
}

func (rw *HTMLResponseWriter) Render(w http.ResponseWriter, data interface{}) error {
	return rw.Template.Execute(w, data)
}

func (rw *HTMLResponseWriter) Error(w http.ResponseWriter, err *common.HTTPError) error {
	w.WriteHeader(err.Code)
	return errTplt.HTML.Execute(w, err)
}

type JSONResponseWriter struct{}

func (rw *JSONResponseWriter) Render(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func (rw *JSONResponseWriter) Error(w http.ResponseWriter, err *common.HTTPError) error {
	w.WriteHeader(err.Code)
	return json.NewEncoder(w).Encode(err)
}

func parseHeaders(w http.ResponseWriter, r *http.Request, t *Templates) ResponseWriter {
	accepts := r.Header["Accept"]
	for _, media := range strings.Split(accepts[0], ",") {
		media = strings.Split(media, ";")[0]
		switch media {
		case "text/html":
			common.Trace.Printf("HTML rendering")
			w.Header().Set("Content-Type", "text/html")
			return &HTMLResponseWriter{t.HTML}
		case "application/xhtml+xml":
			common.Trace.Printf("XML rendering")
			w.Header().Set("Content-Type", "application/xhtml+xml")
			return nil
		}
	}
	common.Trace.Printf("JSON rendering")
	w.Header().Set("Content-Type", "application/json")
	return &JSONResponseWriter{}
}

func (g *Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := parseHeaders(w, r, g.Templates)
	if url := r.URL.Path[len(g.Name)+1:]; url == "" {
		http.Redirect(w, r, g.Name, http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		rw.Error(w, common.NewHTTPError("400 Bad Request", http.StatusBadRequest))
	} else if data, err := g.Method(i); err != nil {
		rw.Error(w, err)
	} else if err := rw.Render(w, data); err != nil {
		common.Error.Printf("Rendering failed: %s", err)
	} else {
		common.Trace.Printf("%q", data)
	}
}

func (l *List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := parseHeaders(w, r, l.Templates)
	if data, err := l.Method(); err != nil {
		rw.Error(w, err)
	} else if err := rw.Render(w, data); err != nil {
		common.Error.Printf("Rendering failed: %s", err)
	} else {
		common.Trace.Printf("%q", data)
	}
}

type ServeMux struct{ *http.ServeMux }

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	common.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *ServeMux) NewEndpoint(name string, list ListMethod, get GetMethod) {
	mux.Handle(name, &List{list, templates(name, "list.html")})
	mux.Handle(
		strings.Join([]string{name, "/"}, ""),
		&Get{name, get, templates(name, "get.html")},
	)
}
