package common

import (
	"encoding/json"
)

type Resource interface {
	Hash() string
}

type Endpoint interface {
	Resource
	List() (interface{}, *HTTPError)
	Get(int) (interface{}, *HTTPError)
	New() Endpoint
	Validate() error
	ValidateForm(values map[string][]string) error
	Save() error
	Delete(int) *HTTPError
}

type OrderedMap struct {
	elements map[string]Resource
	order    []string
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{map[string]Resource{}, []string{}}
}

func (om *OrderedMap) Values() []Resource {
	res := make([]Resource, 0)
	for _, index := range om.order {
		res = append(res, om.elements[index])
	}
	return res
}

func (om *OrderedMap) Get(index int) Resource {
	if index < len(om.order) {
		return om.elements[om.order[index]]
	} else {
		return nil
	}
}

func (om *OrderedMap) Len() int {
	return len(om.order)
}

func (om *OrderedMap) Append(r Resource) Resource {
	if id := r.Hash(); id == "" {
		r = nil
	} else if v, ok := om.elements[id]; ok {
		r = v
	} else {
		om.elements[id] = r
		om.order = append(om.order, id)
	}
	return r
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(om.Values())
}
