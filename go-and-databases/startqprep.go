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

func sqlQPrep(db *sql.DB) error {
	// startQPrep OMIT
	stmt, err := db.Prepare(`select ?`) // HL
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer stmt.Close() // HL

	in := "value"
	rows, err := stmt.Query(in) // HL
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var val interface{}
		err = rows.Scan(&val)
		if err != nil {
			return fmt.Errorf("Scan failed: %v", err)
		}
		fmt.Printf("in=%#[1]v, out=%#[2]v of type %[2]T\n", in, val)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Rows error: %v", err)
	}
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

	funcs := []struct {
		desc string
		f    func(*sql.DB) error
	}{
		{"Prepared Query", sqlQPrep},
	}
	for _, entry := range funcs {
		err := entry.f(db)
		if err != nil {
			log.Fatalf("could not run %s: %v\n", entry.desc, err)
		}
	}
}
