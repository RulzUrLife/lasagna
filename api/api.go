package api

import (
	"encoding/json"
	"github.com/RulzUrLife/lasagna/common"
	"net/http"
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

type Create struct {
	*BaseHandler
}

type List struct {
	Method ListMethod
	*BaseHandler
}

type Get struct {
	Method GetMethod
	*BaseHandler
}

type BaseHandler struct {
	Name string
	*Templates
}

func (bh *BaseHandler) GetTemplates() *Templates {
	return bh.Templates
}

func NewHandler(name string, html string) *BaseHandler {
	return &BaseHandler{common.Url(name), templates(name, html)}
}

type Handler interface {
	ServeHTTP(w ResponseWriter, r *http.Request)
	GetTemplates() *Templates
}

func Handle(h Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := parseHeaders(w, r, h.GetTemplates())
		h.ServeHTTP(rw, r)
	}
}

var errTplt = templates("error.html")

func templates(paths ...string) *Templates {
	base := common.Path(common.Config.Assets.Templates, "base.html")
	p := common.Path(append([]string{common.Config.Assets.Templates}, paths...)...)
	funcs := template.FuncMap{
		"div": div, "slice": slice, "url": url,
	}
	common.Info.Printf("Register template %s", p)
	return &Templates{
		template.Must(template.New("base.html").Funcs(funcs).ParseFiles(base, p)),
		nil,
	}
}

type ResponseWriter interface {
	Write(data interface{}) error
	Writer() http.ResponseWriter
}

type HTMLResponseWriter struct {
	w http.ResponseWriter
	*template.Template
}

func (rw *HTMLResponseWriter) Write(data interface{}) error {
	if err, ok := data.(*common.HTTPError); ok {
		rw.w.WriteHeader(err.Code)
		return errTplt.HTML.Execute(rw.w, data)
	} else {
		return rw.Template.Execute(rw.w, data)
	}
}

func (rw *HTMLResponseWriter) Writer() http.ResponseWriter {
	return rw.w
}

type JSONResponseWriter struct {
	w http.ResponseWriter
}

func (rw *JSONResponseWriter) Write(data interface{}) error {
	if err, ok := data.(*common.HTTPError); ok {
		rw.w.WriteHeader(err.Code)
	}
	return json.NewEncoder(rw.w).Encode(data)
}

func (rw *JSONResponseWriter) Writer() http.ResponseWriter {
	return rw.w
}

func parseHeaders(w http.ResponseWriter, r *http.Request, t *Templates) ResponseWriter {
	accepts := r.Header["Accept"]
	for _, media := range strings.Split(accepts[0], ",") {
		media = strings.Split(media, ";")[0]
		switch media {
		case "text/html":
			common.Trace.Printf("HTML rendering")
			w.Header().Set("Content-Type", "text/html")
			return &HTMLResponseWriter{w, t.HTML}
		case "application/xhtml+xml":
			common.Trace.Printf("XML rendering")
			w.Header().Set("Content-Type", "application/xhtml+xml")
			return nil
		}
	}
	common.Trace.Printf("JSON rendering")
	w.Header().Set("Content-Type", "application/json")
	return &JSONResponseWriter{w}
}

func (g *Get) ServeHTTP(rw ResponseWriter, r *http.Request) {
	if url := r.URL.Path[len(g.Name)+1:]; url == "" {
		http.Redirect(rw.Writer(), r, g.Name, http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		rw.Write(common.NewHTTPError("400 Bad Request", http.StatusBadRequest))
	} else if data, err := g.Method(i); err != nil {
		rw.Write(err)
	} else if err := rw.Write(data); err != nil {
		common.Error.Printf("Rendering failed: %s", err)
	} else {
		common.Trace.Printf("%q", data)
	}
}

func (l *List) ServeHTTP(rw ResponseWriter, r *http.Request) {
	if data, err := l.Method(); err != nil {
		rw.Write(err)
	} else if err := rw.Write(data); err != nil {
		common.Error.Printf("Rendering failed: %s", err)
	} else {
		common.Trace.Printf("%q", data)
	}
}

func (c *Create) ServeHTTP(rw ResponseWriter, r *http.Request) {
	switch rw.(type) {
	case *HTMLResponseWriter:
		rw.Write(struct{}{})
	default:
		// endpoint not available for something else than html
		rw.Write(common.New404Error("Endpoint only exists for mimetype: text/html"))
	}
}

type serveMux struct{ *http.ServeMux }

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	common.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *serveMux) NewEndpoint(name string, list ListMethod, get GetMethod) {
	mux.HandleFunc(common.Url(name), Handle(&List{list, NewHandler(name, "list.html")}))
	mux.HandleFunc(common.Url(name, ""), Handle(&Get{get, NewHandler(name, "get.html")}))
	mux.HandleFunc(common.Url(name, "new"), Handle(&Create{NewHandler(name, "create.html")}))
}

var Mux = &serveMux{http.NewServeMux()}
