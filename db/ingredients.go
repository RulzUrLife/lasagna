package db

import (
	"database/sql"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
)

type Ingredient struct {
	Id   NullInt64  `json:"id"`
	Name NullString `json:"name"`
}

type Ingredients struct {
	Ingredients []*Ingredient `json:"ingredients"`
}

const ingredient_query = `
SELECT gordon.ingredient.id, gordon.ingredient.name
FROM gordon.ingredient
`

func ListIngredients() (*Ingredients, *common.HTTPError) {
	common.Trace.Println(ingredient_query)

	rows, err := DB.Query(ingredient_query)
	if err != nil {
		common.Error.Println(err)
		return nil, common.New500Error("Retrieving of ingredients failed")
	}
	defer rows.Close()

	ingredients := []*Ingredient{}

	for rows.Next() {
		ingredient := &Ingredient{}
		err = rows.Scan(&ingredient.Id, &ingredient.Name)
		if err != nil {
			common.Error.Println(err)
			return nil, common.New500Error("Retrieving of ingredients failed")
		}
		ingredients = append(ingredients, ingredient)
	}
	return &Ingredients{ingredients}, nil
}

func GetIngredient(id int) (*Ingredient, *common.HTTPError) {
	ingredient := &Ingredient{}
	q := ingredient_query + "WHERE gordon.ingredient.id = $1"
	common.Trace.Printf("[%d]%s", id, q)

	err := DB.QueryRow(q, id).Scan(&ingredient.Id, &ingredient.Name)
	if err != nil {
		common.Error.Println(err)
		if err == sql.ErrNoRows {
			return nil, common.New404Error(fmt.Sprintf("Ingredient %d not found", id))
		} else {
			return nil, common.New500Error("Retrieving of ingredient failed")
		}
	}
	return ingredient, nil
}

func PostIngredient(name string) (*Ingredient, error) {
	return nil, nil
}
