package common

import (
	"encoding/json"
	"io"
)

func Dump(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

type Ingredient struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Ingredients struct {
	Ingredients []*Ingredient `json:"ingredients"`
}

func (i *Ingredient) Dump(w io.Writer) error  { return Dump(i, w) }
func (i *Ingredients) Dump(w io.Writer) error { return Dump(i, w) }
