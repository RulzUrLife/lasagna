package db

import (
	"database/sql"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	_ "github.com/lib/pq"
	"regexp"
	"strconv"
)

var (
	DB = connect()
)

type NullString struct {
	sql.NullString
}

type NullInt64 struct {
	sql.NullInt64
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return []byte(fmt.Sprintf("\"%s\"", ns.String)), nil
	} else {
		return []byte("null"), nil
	}

}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return []byte(strconv.FormatInt(ni.Int64, 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func connect() *sql.DB {
	dbConfig := common.Config.Database
	re := regexp.MustCompile("password=.* ")

	connString := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		dbConfig.Name, dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port,
	)

	common.Info.Printf(
		"Connecting to database with following options: %s",
		re.ReplaceAllString(connString, "password=******* "),
	)
	db, err := sql.Open("postgres", connString)

	if err != nil {
		common.Error.Fatalf(
			"Something bad happened during database connection: %s", err,
		)
	}
	return db
}
