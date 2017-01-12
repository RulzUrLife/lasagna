package main

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/api"
	"github.com/RulzUrLife/lasagna/common"
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

var static = http.StripPrefix(
	"/static/", http.FileServer(http.Dir(common.Config.StaticDir)),
)

func main() {

	common.Info.Println("Register URL patterns")
	mux := &api.ServeMux{http.NewServeMux()}
	mux.HandleFunc("/", index)

	mux.Handle("/static/", static)
	mux.NewEndpoint("/ingredients",
		func() (interface{}, *common.HTTPError) { return db.ListIngredients() },
		func(id int) (interface{}, *common.HTTPError) { return db.GetIngredient(id) },
	)
	mux.NewEndpoint("/utensils",
		func() (interface{}, *common.HTTPError) { return db.ListUtensils() },
		func(id int) (interface{}, *common.HTTPError) { return db.GetUtensil(id) },
	)
	mux.NewEndpoint("/recipes",
		func() (interface{}, *common.HTTPError) { return db.ListRecipes() },
		func(id int) (interface{}, *common.HTTPError) { return db.GetRecipe(id) },
	)

	addr := fmt.Sprintf("%s:%d", common.Config.Host, common.Config.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	common.Info.Printf("Running on http://%s/\n", addr)
	common.Error.Fatal(s.ListenAndServe())
}
