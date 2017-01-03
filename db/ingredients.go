package db

import (
	"github.com/RulzUrLife/lasagna/config"
)

type Ingredient struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Ingredients struct {
	Ingredients []*Ingredient `json:"ingredients"`
}

func ListIngredients() (*Ingredients, error) {
	q := "SELECT gordon.ingredient.id, gordon.ingredient.name FROM gordon.ingredient"
	config.Trace.Println(q)

	rows, err := DB.Query(q)
	if err != nil {
		config.Error.Println(err)
		return nil, err
	}
	defer rows.Close()

	ingredients := []*Ingredient{}

	for rows.Next() {
		ingredient := &Ingredient{}
		err = rows.Scan(&ingredient.Id, &ingredient.Name)
		if err != nil {
			config.Error.Println(err)
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return &Ingredients{ingredients}, nil
}

func GetIngredient(id int) (*Ingredient, error) {
	ingredient := &Ingredient{}
	q := "SELECT gordon.ingredient.id, gordon.ingredient.name "
	q = q + "FROM gordon.ingredient WHERE gordon.ingredient.id = $1"

	config.Trace.Printf("%s, id=%d", q, id)

	err := DB.QueryRow(q, id).Scan(&ingredient.Id, &ingredient.Name)
	if err != nil {
		config.Error.Println(err)
	}
	return ingredient, err
}
