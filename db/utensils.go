package db

import (
	"database/sql"
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

func (u *Utensil) Hash() string {
	return strconv.Itoa(u.Id)
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
