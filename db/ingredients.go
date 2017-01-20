package db

import (
	"database/sql"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"strconv"
)

type Ingredient struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

const ingredient_query = `
SELECT gordon.ingredient.id, gordon.ingredient.name
FROM gordon.ingredient
`

func (i *Ingredient) Hash() string {
	return strconv.Itoa(i.Id)
}

func (_ *Ingredient) List() (interface{}, *common.HTTPError) {
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
	return struct {
		Ingredients []*Ingredient `json:"ingredients"`
	}{ingredients}, nil
}

func (_ *Ingredient) Get(id int) (interface{}, *common.HTTPError) {
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
