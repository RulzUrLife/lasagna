package db

import (
	"github.com/RulzUrLife/lasagna/config"
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

func dedup(q string, params ...interface{}) (res []*Recipe, err error) {
	recipes := map[int]*Recipe{}
	config.Trace.Printf("%v%s", params, q)

	rows, err := DB.Query(q, params...)
	if err != nil {
		config.Error.Println(err)
		return
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
			config.Error.Println(err)
			return
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
	return
}

func ListRecipes() (*Recipes, error) {
	if recipes, err := dedup(query); err != nil {
		return nil, err
	} else {
		return &Recipes{recipes}, nil
	}
}

func GetRecipe(id int) (*Recipe, error) {
	recipes, err := dedup(query+"WHERE gordon.recipe.id = $1", id)
	if err != nil {
		return nil, err
	} else {
		return recipes[0], nil
	}
}
