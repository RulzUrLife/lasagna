package db

import (
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"github.com/lib/pq"
)

const query = `
SELECT gordon.recipe.id, gordon.recipe.name, gordon.recipe.directions,
	gordon.recipe.description, gordon.recipe.difficulty, gordon.recipe.duration,
	gordon.recipe.people, gordon.recipe.category,
	gordon.ingredient.id, gordon.ingredient.name,
	gordon.recipe_ingredients.quantity, gordon.recipe_ingredients.measurement
FROM gordon.recipe
LEFT OUTER JOIN gordon.recipe_ingredients
ON (gordon.recipe.id = gordon.recipe_ingredients.fk_recipe)
LEFT OUTER JOIN gordon.ingredient
ON (gordon.recipe_ingredients.fk_ingredient = gordon.ingredient.id)
`

type Recipe struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	Directions  []Direction        `json:"directions"`
	Description NullString         `json:"description"`
	Difficulty  NullInt64          `json:"difficulty"`
	Duration    NullString         `json:"duration"`
	People      NullInt64          `json:"people"`
	Category    NullString         `json:"category"`
}

type Recipes struct {
	Recipes []*Recipe `json:"recipes"`
}
type RecipeIngredient struct {
	Measurement NullString `json:"measurement"`
	Quantity    NullInt64  `json:"quantity"`
	Ingredient
}

type Direction struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (d *Direction) Scan(src interface{}) (err error) {
	res := scan(src)
	d.Title, d.Text = string(res[0]), string(res[1])
	return nil
}

func scan(src interface{}) (elems [][]byte) {
	var elem []byte

	bytes := src.([]byte)
	// remove enclosing parenthesis
	bytes = bytes[1 : len(bytes)-1]
	for i := 0; i < len(bytes); {
		elem, i = scanBytes(bytes, i)
		elems = append(elems, elem)
	}
	return
}

func scanBytes(bytes []byte, i int) (elem []byte, _ int) {
	var escape bool

	switch bytes[i] {
	case '"':
		for i++; i < len(bytes); i++ {
			if escape {
				if bytes[i] == ',' {
					break
				}
				elem = append(elem, bytes[i])
				escape = false
			} else {
				switch bytes[i] {
				default:
					elem = append(elem, bytes[i])
				case '\\', '"':
					escape = true
				}
			}
		}
	default:
		for ; i < len(bytes); i++ {
			if bytes[i] == ',' {
				break
			}
			elem = append(elem, bytes[i])
		}
	}
	return elem, i + 1
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
		// initialise empty arrays allow to marshal them correctly
		// see https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
		recipe.Ingredients = make([]RecipeIngredient, 0)
		ingredient := RecipeIngredient{}

		err = rows.Scan(
			&recipe.Id, &recipe.Name, &pq.GenericArray{&recipe.Directions},
			&recipe.Description, &recipe.Difficulty, &recipe.Duration,
			&recipe.People, &recipe.Category,
			&ingredient.Id, &ingredient.Name,
			&ingredient.Quantity, &ingredient.Measurement,
		)
		if err != nil {
			common.Error.Println(err)
			return nil, common.New500Error(err.Error())
		}
		if len(recipe.Directions) == 0 {
			recipe.Directions = make([]Direction, 0) // ensure empty array initialisation
		}
		if v, ok := recipes[recipe.Id]; ok {
			recipe = v
		} else {
			recipes[recipe.Id] = recipe
		}
		if ingredient.Id.Valid {
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
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
