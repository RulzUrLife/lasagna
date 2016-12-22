package common

import ()

type Ingredient struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Ingredients struct {
	Ingredients []*Ingredient `json:"ingredients"`
}
