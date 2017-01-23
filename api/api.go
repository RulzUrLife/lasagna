package api

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"net/http"
	"strconv"
	"text/template"
)

type Create struct {
	*BaseHandler
}

type List struct {
	*BaseHandler
}

type Get struct {
	*BaseHandler
}

type BaseHandler struct {
	Name string
	*template.Template
	common.Endpoint
}

func NewHandler(name string, html string, ressource common.Endpoint) *BaseHandler {
	return &BaseHandler{common.Url(name), templates(name, html), ressource}
}

var errTplt = templates("error.html")

func templates(paths ...string) *template.Template {
	base := common.Path(common.Config.Assets.Templates, "base.html")
	p := common.Path(append([]string{common.Config.Assets.Templates}, paths...)...)
	funcs := template.FuncMap{
		"div": div, "slice": slice, "url": url,
	}
	common.Info.Printf("Register template %s", p)
	return template.Must(template.New("base.html").Funcs(funcs).ParseFiles(base, p))
}

func (g *Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	del := func(id int) {
		if err := g.Delete(id); err != nil {
			errTplt.Execute(w, err)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
	get := func(id int) {
		if data, err := g.Get(id); err != nil {
			errTplt.Execute(w, err)
		} else if err := g.Template.Execute(w, data); err != nil {
			common.Error.Printf("Rendering failed: %s", err)
		} else {
			common.Trace.Printf("%q", data)
		}
	}

	if url := r.URL.Path[len(g.Name)+1:]; url == "" {
		// trailing / redirect to base url
		http.Redirect(w, r, g.Name, http.StatusSeeOther)
	} else if id, err := strconv.Atoi(url); err != nil {
		// non parsable parameter
		errTplt.Execute(w, common.New400Error(
			fmt.Sprintf("Invalid %s id '%s'", g.Name, url),
		))
	} else {
		// switch through http methods, invoke the correct one
		switch r.Method {
		case http.MethodDelete:
			del(id)
		case http.MethodGet:
			get(id)
		default:
			errTplt.Execute(w, common.New404Error("Method does not exist on this endpoint"))
		}
	}
}

func (l *List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			errTplt.Execute(w, common.New400Error(err.Error()))
		} else if e, err := l.ValidateForm(r.Form); err != nil {
			errTplt.Execute(w, common.New400Error(err.Error()))
		} else if err := e.Save(); err != nil {
			common.Error.Printf("%s", err)
			errTplt.Execute(w, common.New500Error(err.Error()))
		} else {
			http.Redirect(w, r, common.Url(l.Name, e.Hash()), http.StatusSeeOther)
		}
	} else {
		if data, err := l.List(); err != nil {
			errTplt.Execute(w, err)
		} else if err := l.Template.Execute(w, data); err != nil {
			common.Error.Printf("Rendering failed: %s", err)
		} else {
			common.Trace.Printf("%q", data)
		}
	}
}

func (c *Create) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Template.Execute(w, struct{}{})
}

type serveMux struct{ *http.ServeMux }

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	common.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *serveMux) NewEndpoint(name string, ressource common.Endpoint) {
	mux.Handle(common.Url(name), &List{NewHandler(name, "list.html", ressource)})
	mux.Handle(common.Url(name, ""), &Get{NewHandler(name, "get.html", ressource)})
	mux.Handle(common.Url(name, "new"), &Create{NewHandler(name, "create.html", ressource)})
}

var Mux = &serveMux{http.NewServeMux()}
