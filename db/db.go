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

const (
	invalid = "Invalid attribute '%s'"
	missing = "Missing attribut '%s'"
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

func (ni NullInt64) Hash() (id string) {
	if ni.Valid {
		id = strconv.FormatInt(ni.Int64, 10)
	}
	return
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return []byte(strconv.FormatInt(ni.Int64, 10)), nil
	} else {
		return []byte("null"), nil
	}
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
