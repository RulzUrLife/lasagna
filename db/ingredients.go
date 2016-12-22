package db

import (
	"github.com/RulzUrLife/lasagna/config"
)

type Ingredient struct {
	Id   int
	Name string
}

type Ingredients []*Ingredient

func ListIngredients() (ingredients Ingredients, err error) {
	q := "SELECT gordon.ingredient.id, gordon.ingredient.name FROM gordon.ingredient"

	config.Trace.Println(q)

	rows, err := DB.Query(q)
	if err != nil {
		config.Error.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		ingredient := &Ingredient{}
		err = rows.Scan(&ingredient.Id, &ingredient.Name)
		if err != nil {
			config.Error.Println(err)
			return
		}
		ingredients = append(ingredients, ingredient)
	}
	return
}
