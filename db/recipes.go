package db

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"github.com/lib/pq"
	"strconv"
)

const query = `
SELECT gordon.recipe.id, gordon.recipe.name, gordon.recipe.directions,
	gordon.recipe.description, gordon.recipe.difficulty, gordon.recipe.duration,
	gordon.recipe.people, gordon.recipe.category,
	gordon.utensil.id, gordon.utensil.name,
	gordon.ingredient.id, gordon.ingredient.name,
	gordon.recipe_ingredients.quantity, gordon.recipe_ingredients.measurement
FROM gordon.recipe
LEFT OUTER JOIN gordon.recipe_utensils
ON (gordon.recipe.id = gordon.recipe_utensils.fk_recipe)
LEFT OUTER JOIN gordon.utensil
ON (gordon.recipe_utensils.fk_utensil = gordon.utensil.id)
LEFT OUTER JOIN gordon.recipe_ingredients
ON (gordon.recipe.id = gordon.recipe_ingredients.fk_recipe)
LEFT OUTER JOIN gordon.ingredient
ON (gordon.recipe_ingredients.fk_ingredient = gordon.ingredient.id)
`

type Recipe struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Ingredients *common.OrderedMap `json:"ingredients"`
	Utensils    *common.OrderedMap `json:"utensils"`
	Directions  []Direction        `json:"directions"`
	Description NullString         `json:"description"`
	Difficulty  NullInt64          `json:"difficulty"`
	Duration    NullString         `json:"duration"`
	People      NullInt64          `json:"people"`
	Category    NullString         `json:"category"`
}

type RecipeIngredient struct {
	Measurement NullString `json:"measurement"`
	Quantity    NullInt64  `json:"quantity"`
	Id          NullInt64  `json:"id"`
	Name        NullString `json:"name"`
}

func (ri RecipeIngredient) Hash() (id string) {
	return ri.Id.Hash()
}

type RecipeUtensil struct {
	Id   NullInt64  `json:"id"`
	Name NullString `json:"name"`
}

func (ru RecipeUtensil) Hash() (id string) {
	return ru.Id.Hash()
}

type Direction struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (r *Recipe) Hash() string {
	return strconv.Itoa(r.Id)
}

func (_ *Recipe) New() interface{} {
	return &Recipe{}
}

func (d *Direction) Scan(src interface{}) (err error) {
	res := scan(src)
	d.Title, d.Text = string(res[0]), string(res[1])
	return nil
}

func dedup(q string, params ...interface{}) (*common.OrderedMap, *common.HTTPError) {
	recipes := common.NewOrderedMap()
	common.Trace.Printf("%v%s", params, q)

	rows, err := DB.Query(q, params...)
	if err != nil {
		common.Error.Println(err)
		return nil, common.New500Error(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		recipe := &Recipe{
			Ingredients: common.NewOrderedMap(), Utensils: common.NewOrderedMap(),
		}
		ingredient := RecipeIngredient{}
		utensil := RecipeUtensil{}

		err = rows.Scan(
			&recipe.Id, &recipe.Name, &pq.GenericArray{&recipe.Directions},
			&recipe.Description, &recipe.Difficulty, &recipe.Duration,
			&recipe.People, &recipe.Category,
			&utensil.Id, &utensil.Name,
			&ingredient.Id, &ingredient.Name,
			&ingredient.Quantity, &ingredient.Measurement,
		)
		if err != nil {
			common.Error.Println(err)
			return nil, common.New500Error(err.Error())
		}
		if len(recipe.Directions) == 0 {
			// initialise empty arrays allow to marshal them correctly
			// see https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
			recipe.Directions = make([]Direction, 0)
		}
		recipe, _ = recipes.Append(recipe).(*Recipe)
		recipe.Ingredients.Append(ingredient)
		recipe.Utensils.Append(utensil)
	}
	return recipes, nil
}

func (_ *Recipe) List() (interface{}, *common.HTTPError) {
	if recipes, err := dedup(query); err != nil {
		return nil, err
	} else {
		return struct {
			Recipes interface{} `json:"recipes"` // envelope wrapper
		}{recipes}, nil
	}
}

func (_ *Recipe) Get(id int) (interface{}, *common.HTTPError) {
	if recipes, err := dedup(query+"WHERE gordon.recipe.id = $1", id); err != nil {
		return nil, err
	} else if recipe := recipes.Get(0); recipe == nil {
		return nil, common.New404Error(fmt.Sprintf("Recipe %d not Found", id))
	} else {
		return recipe, nil
	}
}
