package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func list_ingredients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List ingredients")
}

func get_ingredients(w http.ResponseWriter, r *http.Request) {
	if url := r.URL.Path[13:]; url == "" {
		http.Redirect(w, r, "/ingredients", http.StatusSeeOther)
	} else if i, err := strconv.Atoi(url); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	} else {
		fmt.Fprintf(w, "Get ingredient %d", i)
	}
}

func init() {
	Mux.HandleFunc("/ingredients", list_ingredients)
	Mux.HandleFunc("/ingredients/", get_ingredients)
}
