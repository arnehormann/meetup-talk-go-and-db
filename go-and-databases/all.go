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

func sqlQRBasic(db *sql.DB) error {
	// startQR OMIT
	row := db.QueryRow(`select 1 num, "hello" txt`) // HL

	var num int
	var txt string
	err := row.Scan(&num, &txt) // HL
	if err != nil {
		return err
	}
	fmt.Printf("num=%d, txt=%q\n", num, txt)
	// end OMIT
	return nil
}

func sqlQBasic(db *sql.DB) error {
	// startQBasic OMIT
	rows, err := db.Query(`select 1 union select 2 union select 3`) // HL
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer rows.Close() // HL
	for rows.Next() {  // HL
		var n int
		err = rows.Scan(&n) // HL
		if err != nil {
			return fmt.Errorf("Scan failed: %v", err)
		}
		fmt.Println(n)
	}
	if err = rows.Err(); err != nil { // HL
		return fmt.Errorf("Rows error: %v", err)
	}
	// end OMIT
	return nil
}

func sqlQRInterface(db *sql.DB) error {
	// startQRInterface OMIT
	inputs := []interface{}{3, 3.0, '3', "3"}
	for _, in := range inputs {
		row := db.QueryRow(`select ?`, in) // HL

		var val interface{}
		err := row.Scan(&val) // HL
		if err != nil {
			return err
		}
		fmt.Printf("in=%#[1]v, out=%#[2]v of type %[2]T\n", in, val)
	}
	// end OMIT
	return nil
}

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

func sqlQColumns(db *sql.DB) error {
	rows, err := db.Query(`select 1 union select 2 union select 3`)
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer rows.Close()
	// startQColumnNames OMIT
	names, err := rows.Columns() // HL
	// end OMIT
	if err != nil {
		return fmt.Errorf("Columns failed: %v", err)
	}
	// startQColumns OMIT
	cols, err := rows.ColumnTypes() // HL
	col0 := cols[0]

	dbtype := col0.DatabaseTypeName()     // CHAR, INT, ...
	prec, scale, ok := col0.DecimalSize() // overall and fractional digits
	length, ok := col0.Length()           // for text or binaries
	colName := col0.Name()                // column name like before 1.8
	nullable, ok := col0.Nullable()       // nullable type, applies
	reflectType := col0.ScanType()        // the type in reflect
	// end OMIT

	if err != nil {
		return fmt.Errorf("ColumnTypes failed: %v", err)
	}
	// silence compiler unused warning
	var (
		_ = ok
		_ = names
		_ = dbtype
		_ = prec
		_ = scale
		_ = length
		_ = colName
		_ = nullable
		_ = reflectType
	)
	n := 0
	for rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return fmt.Errorf("Scan failed: %v", err)
		}
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Rows error: %v", err)
	}
	return nil
}

func sqlEBasic(db *sql.DB) error {
	// startEBasic OMIT
	query := `DO 1`
	_, err := db.Exec(query) // HL
	if err != nil {
		return err
	}
	// end OMIT
	return nil
}

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

func sqlEPrepInsert(db *sql.DB) error {
	// startEPrepInsert OMIT
	stmt, err := db.Prepare(`insert into gophers set name=?`) // HL
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, gopher := range []string{"Gwen", "Gladys", "George"} {
		res, err := stmt.Exec(gopher) // HL
		if err != nil {
			return err
		}
		num, _ := res.RowsAffected() // HL
		if num != 1 {
			log.Printf("could not add %s\n", gopher)
		}
	}
	fmt.Println("<- inserted ->")
	// end OMIT
	return nil
}

func sqlEPrepInsertSelect(db *sql.DB) error {
	// startEPrepInsertSelect OMIT
	rows, err := db.Query(`select name from gophers where name like "%o%"`) // HL
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return fmt.Errorf("Scan failed: %v", err)
		}
		fmt.Println(name)
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
		{"QueryRow", sqlQRBasic},
		{"Query", sqlQBasic},
		{"QueryRow with interface{}", sqlQRInterface},
		{"Prepared query", sqlQPrep},
		{"Exec", sqlEBasic},
		{"Exec create table", sqlExecer(createTable)},
		{"Truncate gophers", sqlExecer("truncate gophers")},
		{"select o-gophers", sqlEPrepInsertSelect},
		{"Exec insert into table", sqlEInsert},
		{"select o-gophers", sqlEPrepInsertSelect},
		{"Prepared Exec insert into table", sqlEPrepInsert},
		{"select o-gophers", sqlEPrepInsertSelect},
	}
	for _, entry := range funcs {
		err := entry.f(db)
		if err != nil {
			log.Fatalf("could not run %s: %v\n", entry.desc, err)
		}
	}
}
