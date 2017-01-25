package api

import (
	"github.com/RulzUrLife/lasagna/common"
	"net/http"
	"strconv"
	"text/template"
)

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

type serveMux struct{ *http.ServeMux }

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	common.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(w, r)
}

func (mux *serveMux) NewEndpoint(name string, resource common.Endpoint) {
	var handler http.HandlerFunc
	var listTemplate = templates(name, "list.html")
	var getTemplate = templates(name, "get.html")
	var newTemplate = templates(name, "new.html")

	handler = func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path == "new" {
				defer r.Body.Close()
				if err := r.ParseForm(); err != nil {
					errTplt.Execute(w, common.New400Error(err.Error()))
				} else if e, err := resource.ValidateForm(r.Form); err != nil {
					errTplt.Execute(w, common.New400Error(err.Error()))
				} else if err := e.Save(); err != nil {
					common.Error.Printf("%s", err)
					errTplt.Execute(w, common.New500Error(err.Error()))
				} else {
					http.Redirect(w, r, common.Url(name, e.Hash()), http.StatusSeeOther)
				}
				break
			}
			fallthrough
		case http.MethodGet:
			switch url := r.URL.Path; url {
			case "":
				if data, err := resource.List(); err != nil {
					errTplt.Execute(w, err)
				} else if err := listTemplate.Execute(w, data); err != nil {
					common.Error.Printf("Rendering failed: %s", err)
				} else {
					common.Trace.Printf("%q", data)
				}
			case "new":
				newTemplate.Execute(w, struct{}{})
			default:
				if id, err := strconv.Atoi(url); err != nil {
					errTplt.Execute(w, common.New400Error("Invalid %s id '%s'", name, url))
				} else if data, err := resource.Get(id); err != nil {
					errTplt.Execute(w, err)
				} else if err := getTemplate.Execute(w, data); err != nil {
					common.Error.Printf("Rendering failed: %s", err)
				} else {
					common.Trace.Printf("%q", data)
				}
			}
		default:
			errTplt.Execute(w, common.New404Error("Resource does not exist"))
		}
	}

	mux.Handle(common.Url(name, ""), http.StripPrefix(common.Url(name, ""), handler))
}

var Mux = &serveMux{http.NewServeMux()}
