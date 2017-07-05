// + bu ild OMIT

package main

import (
	"fmt"
	"log"
)

// startImport OMIT
import (
	"database/sql" // HLsql

	_ "github.com/go-sql-driver/mysql" // HLdrv
)

// end OMIT

func sqlEInsert(db *sql.DB) error {
	// startEInsert OMIT
	query := `insert into gophers set name=?`
	_, err := db.Exec(query, "Gordon") // HL
	if err != nil {
		return err
	}
	fmt.Println("<- inserted ->")
	// end OMIT
	return nil
}

func sqlExecer(query string, args ...interface{}) func(*sql.DB) error {
	return func(db *sql.DB) error {
		_, err := db.Exec(query, args...)
		if err != nil {
			return err
		}
		return nil
	}
}

type dummyCloser struct{}

func (c dummyCloser) Close() error {
	return nil
}

func main() {
	resource := dummyCloser{}
	var err error
	// startDCShort OMIT
	defer resource.Close()
	// end OMIT

	// startDCDetailed OMIT
	defer func() {
		err := resource.Close()
		if err != nil {
			log.Printf("Closing failed: %v\n", err)
		}
	}()
	// end OMIT

	// startOpen OMIT
	dsn := "gomeetup:20170706@tcp(127.0.0.1:3306)/gotest" // HL
	db, err := sql.Open("mysql", dsn)                     // HL
	if err != nil {
		log.Fatal("db is unavailable")
	}
	defer db.Close()
	// end OMIT

	createTable := "" +
		// startECreate OMIT
		`create table if not exists gophers(
			id int unsigned not null auto_increment,
			name varchar(255),
			primary key (id)
		)`
	// end OMIT

	funcs := []struct {
		desc string
		f    func(*sql.DB) error
	}{
		{"Exec create table", sqlExecer(createTable)},
		{"Truncate gophers", sqlExecer("truncate gophers")},
		{"Exec insert into table", sqlEInsert},
	}
	for _, entry := range funcs {
		err := entry.f(db)
		if err != nil {
			log.Fatalf("could not %s: %v\n", entry.desc, err)
		}
	}
}
