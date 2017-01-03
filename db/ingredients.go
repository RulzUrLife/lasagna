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

const ingredient_query = `
SELECT gordon.ingredient.id, gordon.ingredient.name
FROM gordon.ingredient
`

func ListIngredients() (*Ingredients, error) {
	config.Trace.Println(ingredient_query)

	rows, err := DB.Query(ingredient_query)
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
	q := ingredient_query + "WHERE gordon.ingredient.id = $1"
	config.Trace.Printf("[%d]%s", id, q)

	err := DB.QueryRow(q, id).Scan(&ingredient.Id, &ingredient.Name)
	if err != nil {
		config.Error.Println(err)
	}
	return ingredient, err
}
