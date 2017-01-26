package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"strconv"
)

type Utensil struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

const utensil_query = `
SELECT gordon.utensil.id, gordon.utensil.name
FROM gordon.utensil
`

const utensil_save_query = `
INSERT INTO gordon.utensil (name)
VALUES ($1)
RETURNING gordon.utensil.id, gordon.utensil.name
`

func (u *Utensil) Hash() string {
	return strconv.Itoa(u.Id)
}

func (_ *Utensil) New() common.Endpoint {
	return &Utensil{}
}

func (_ *Utensil) List() (interface{}, *common.HTTPError) {
	common.Trace.Println(utensil_query)

	rows, err := DB.Query(utensil_query)
	if err != nil {
		common.Error.Println(err)
		return nil, common.New500Error("Retrieving of utensils failed")
	}
	defer rows.Close()

	utensils := []*Utensil{}

	for rows.Next() {
		utensil := &Utensil{}
		err = rows.Scan(&utensil.Id, &utensil.Name)
		if err != nil {
			common.Error.Println(err)
			return nil, common.New500Error("Retrieving of utensils failed")
		}
		utensils = append(utensils, utensil)
	}
	return struct {
		Utensils []*Utensil `json:"utensils"`
	}{utensils}, nil
}

func (_ *Utensil) Get(id int) (interface{}, *common.HTTPError) {
	utensil := &Utensil{}
	q := utensil_query + "WHERE gordon.utensil.id = $1"
	common.Trace.Printf("[%d]%s", id, q)

	err := DB.QueryRow(q, id).Scan(&utensil.Id, &utensil.Name)
	if err != nil {
		common.Error.Println(err)
		if err == sql.ErrNoRows {
			return nil, common.New404Error(fmt.Sprintf("Utensil %d not found", id))
		} else {
			return nil, common.New500Error("Retrieving of utensil failed")
		}
	}
	return utensil, nil
}

func (_ *Utensil) Validate(values map[string][]string) (common.Endpoint, error) {
	if names, ok := values["name"]; !ok {
		return nil, errors.New(fmt.Sprintf(missing, "name"))
	} else if len(names) != 1 || names[0] == "" {
		return nil, errors.New(fmt.Sprintf(invalid, "name"))
	} else {
		return &Utensil{Name: names[0]}, nil
	}
}

func (u *Utensil) Save() error {
	return DB.QueryRow(utensil_save_query, u.Name).Scan(&u.Id, &u.Name)
}
