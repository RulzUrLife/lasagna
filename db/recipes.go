package db

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
)

const query = `
SELECT gordon.recipe.id, gordon.recipe.name,
	gordon.ingredient.id, gordon.ingredient.name,
	gordon.recipe_ingredients.quantity, gordon.recipe_ingredients.measurement
FROM gordon.recipe
LEFT OUTER JOIN gordon.recipe_ingredients
ON (gordon.recipe.id = gordon.recipe_ingredients.fk_recipe)
LEFT OUTER JOIN gordon.ingredient
ON (gordon.recipe_ingredients.fk_ingredient = gordon.ingredient.id)
`

type RecipeIngredient struct {
	Measurement string `json:"measurement"`
	Quantity    int    `json:"quantity"`
	Ingredient
}

type Recipe struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	Ingredients []*RecipeIngredient `json:"ingredients"`
}

type Recipes struct {
	Recipes []*Recipe `json:"recipes"`
}

func dedup(q string, params ...interface{}) ([]*Recipe, *common.HTTPError) {
	res := make([]*Recipe, 0)
	recipes := map[int]*Recipe{}
	common.Trace.Printf("%v%s", params, q)

	rows, err := DB.Query(q, params...)
	if err != nil {
		common.Error.Println(err)
		return nil, common.New500Error(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		recipe := &Recipe{}
		ingredient := &RecipeIngredient{}

		err = rows.Scan(
			&recipe.Id, &recipe.Name, &ingredient.Id, &ingredient.Name,
			&ingredient.Quantity, &ingredient.Measurement,
		)
		if err != nil {
			common.Error.Println(err)
			return nil, common.New500Error(err.Error())
		}
		if v, ok := recipes[recipe.Id]; ok {
			recipe = v
		} else {
			recipes[recipe.Id] = recipe
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}
	for _, recipe := range recipes {
		res = append(res, recipe)
	}
	return res, nil
}

func ListRecipes() (*Recipes, *common.HTTPError) {
	if recipes, err := dedup(query); err != nil {
		return nil, err
	} else {
		return &Recipes{recipes}, nil
	}
}

func GetRecipe(id int) (*Recipe, *common.HTTPError) {
	recipes, err := dedup(query+"WHERE gordon.recipe.id = $1", id)
	if err != nil {
		return nil, err
	} else if len(recipes) < 1 {
		return nil, common.New404Error(fmt.Sprintf("Recipe %d not Found", id))
	} else {
		return recipes[0], nil
	}
}
