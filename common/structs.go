package common

import (
	"encoding/json"
)

type Deduplier interface {
	Hash() string
}

type OrderedMap struct {
	elements map[string]Deduplier
	order    []string
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{map[string]Deduplier{}, []string{}}
}

func (om *OrderedMap) Values() []Deduplier {
	res := make([]Deduplier, 0)
	for _, index := range om.order {
		res = append(res, om.elements[index])
	}
	return res
}

func (om *OrderedMap) Get(index int) Deduplier {
	if index < len(om.order) {
		return om.elements[om.order[index]]
	} else {
		return nil
	}
}

func (om *OrderedMap) Append(d Deduplier) Deduplier {
	if id := d.Hash(); id == "" {
		d = nil
	} else if v, ok := om.elements[id]; ok {
		d = v
	} else {
		om.elements[id] = d
		om.order = append(om.order, id)
	}
	return d
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(om.Values())
}
