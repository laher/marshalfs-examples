package examples

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/laher/marshalfs"
	_ "github.com/lib/pq"
)

func Example_sqlx() {

	type Person struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Email     string
	}

	var schema = `
	CREATE TABLE person (
	    first_name text,
	    last_name text,
	    email text
	);
	`
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatalln("POSTGRES_DSN not set. End")
	}
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	personGenerator := func(filename string) (interface{}, error) {
		if !strings.HasPrefix(filename, "person") {
			return nil, os.ErrNotExist
		}
		parts := filepath.SplitList(filename)
		if len(parts) != 3 {
			// should it be a different error?
			return nil, os.ErrNotExist
		}
		lastName := parts[2]
		firstName := parts[3]
		// Query the database, storing result in a Person (wrapped in []interface{})
		rows, err := db.Queryx("SELECT * FROM person where first_name = ? AND last_name = ?", firstName, lastName)
		if err != nil {
			return nil, err
		}
		person := &Person{}
		for rows.Next() {
			err := rows.StructScan(person)
			if err != nil {
				return nil, err
			}
			return person, nil
		}
		return nil, os.ErrNotExist
	}

	marshalfs.New(json.Marshal, marshalfs.FileMap{"person/*": marshalfs.NewFileGenerator(personGenerator)})

	// Output:
}
