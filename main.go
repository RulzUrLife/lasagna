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
	"/static/", http.FileServer(http.Dir(common.Config.Assets.Static)),
)

func main() {
	common.Info.Println("Register URL patterns")
	mux := &api.ServeMux{http.NewServeMux()}

	mux.HandleFunc("/", index)
	mux.Handle("/static/", static)

	mux.NewEndpoint("ingredients", db.ListIngredients, db.GetIngredient)
	mux.NewEndpoint("utensils", db.ListUtensils, db.GetUtensil)
	mux.NewEndpoint("recipes", db.ListRecipes, db.GetRecipe)

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
