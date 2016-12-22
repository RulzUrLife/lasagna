package api

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/config"
	"net/http"
)

var (
	Mux = &ServeMux{http.NewServeMux()}
)

type ResponseWriter struct {
	http.ResponseWriter
	request *http.Request
}

func (rw *ResponseWriter) WriteResp(i interface{}) (int, error) {
	accepts := rw.request.Header["Accept"]
	for _, media := range accepts {
		config.Trace.Printf("%s", media)
	}
	config.Trace.Printf("Adapt response to headers")
	rw.Header().Set("Content-Type", "application/json")
	return rw.ResponseWriter.Write([]byte{})
}

type ServeMux struct {
	*http.ServeMux
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config.Info.Printf("%s %s", r.Method, r.URL.Path)
	mux.ServeMux.ServeHTTP(&ResponseWriter{w, r}, r)
}

func index(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello World!")
}

func init() {
	config.Info.Println("Register URL patterns")
	Mux.HandleFunc("/", index)
}
