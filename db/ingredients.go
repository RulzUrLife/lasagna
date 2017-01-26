package db

import (
	"database/sql"
	"errors"
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
const ingredient_save_query = `
INSERT INTO gordon.ingredient (name)
VALUES ($1)
RETURNING gordon.ingredient.id, gordon.ingredient.name
`

func (i *Ingredient) Hash() string {
	return strconv.Itoa(i.Id)
}

func (_ *Ingredient) New() common.Endpoint {
	return &Ingredient{}
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

func (_ *Ingredient) Validate(values map[string][]string) (common.Endpoint, error) {
	if names, ok := values["name"]; !ok {
		return nil, errors.New(fmt.Sprintf(missing, "name"))
	} else if len(names) != 1 || names[0] == "" {
		return nil, errors.New(fmt.Sprintf(invalid, "name"))
	} else {
		return &Ingredient{Name: names[0]}, nil
	}
}

func (i *Ingredient) Save() error {
	return DB.QueryRow(ingredient_save_query, i.Name).Scan(&i.Id, &i.Name)
}
