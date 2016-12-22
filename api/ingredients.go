package api

import (
	"github.com/RulzUrLife/lasagna/db"
	"net/http"
	"strconv"
)

func list_ingredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := db.ListIngredients()
	if err == nil {
		w.(*ResponseWriter).WriteResp(ingredients)
	} else {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func get_ingredients(w http.ResponseWriter, r *http.Request) {
	if url := r.URL.Path[13:]; url == "" {
		http.Redirect(w, r, "/ingredients", http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	} else if ingredient, err := db.GetIngredient(i); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		w.(*ResponseWriter).WriteResp(ingredient)
	}
}

func init() {
	Mux.HandleFunc("/ingredients", list_ingredients)
	Mux.HandleFunc("/ingredients/", get_ingredients)
}
