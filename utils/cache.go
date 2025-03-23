package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type CSVCache interface {
	Has()
	Get()
	Put()
	Delete()
}

type IdentityCache struct {
	Filename string
}

func (cache IdentityCache) queryDB(sql_query string) (*sql.Rows, error) {
	// get a valid DB
	db, err := sql.Open("sqlite3", cache.Filename)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	init := `CREATE TABLE IF NOT EXISTS repo_identities (
          id INTEGER PRIMARY KEY,  
          source TEXT NOT NULL,
          lineage_id TEXT NOT NULL
       );`

	_, init_err := db.Query(init)
	if init_err != nil {
		fmt.Println(init_err)
	}

	return db.Query(sql_query)
}

func (cache IdentityCache) GetAll() (*sql.Rows, error) {
	return cache.queryDB("SELECT * FROM repo_identities")
}

func (cache IdentityCache) ExportAllToCSV(destination string) {
	rows, query_err := cache.GetAll()
	if query_err != nil {
		fmt.Println(query_err)
	}
	defer rows.Close()

	// create file
	f, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()
	f.WriteString("id, source, lineage_id\n")

	for rows.Next() {

		var id int
		var source string
		var lineage_id string

		err := rows.Scan(&id, &source, &lineage_id)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(fmt.Sprintf("%d, %s, %s\n", id, source, lineage_id))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cache IdentityCache) check_cache(path string) string {
	inputfile, err := os.Stat(path)
	if err != nil {

		// Checking if the given file exists or not
		// Using IsNotExist() function
		if os.IsNotExist(err) {
			log.Fatal("File not Found !!")
		}
	}
	mode := inputfile.Mode()

	if mode.IsRegular() {

		// do file stuff
		data, err := os.ReadFile(cache.Filename)
		if err != nil {
			fmt.Println(err)
		}
		return string(data)
	}

	panic(errors.New("Cachepath shouldnt be a directory"))
}

func (cache IdentityCache) write_cache(path string, lineageID string) {

	d1 := []byte(lineageID)

	err := os.WriteFile(cache.Filename, d1, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
