package utils

import (
	"database/sql"
	"encoding/csv"
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

type IdentityValue struct {
	ID        int
	URL       string
	LineageID string
}

func (cache IdentityCache) queryDB(sql_query string) (*sql.Rows, error) {
	// get a valid DB
	db, err := sql.Open("sqlite3", cache.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	init := `CREATE TABLE IF NOT EXISTS repo_identities (
          id INTEGER PRIMARY KEY,  
          source TEXT NOT NULL,
          lineage_id TEXT NOT NULL
       );`

	_, init_err := db.Exec(init)
	if init_err != nil {
		fmt.Println(init_err)
	}

	return db.Query(sql_query)
}

func (cache IdentityCache) GetAll() ([]IdentityValue, error) {
	results := []IdentityValue{}
	rows, query_err := cache.queryDB("SELECT * FROM repo_identities")
	if query_err != nil {
		fmt.Println(query_err)
	}
	defer rows.Close()

	for rows.Next() {

		var id int
		var source string
		var lineage_id string

		err := rows.Scan(&id, &source, &lineage_id)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, IdentityValue{
			id,
			source,
			lineage_id,
		})
	}
	return results, nil
}

func (cache IdentityCache) ExportAllToCSV(destination string) {
	_, err := cache.GetAll()
	if err != nil {
		fmt.Println(err)
	}

	csvFile, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	err = csvWriter.Write([]string{"678", "Jane", "jane@example.com", "$548,980"})
	if err != nil {
		// an error occurred during the flush
		fmt.Println(err)
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		// an error occurred during the flush
		fmt.Println(err)
	}
}
