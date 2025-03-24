package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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
	db       *gorm.DB
}

type IdentityValue struct {
	ID        uint   `gorm:"primaryKey"`
	Nickname  string `gorm:"unique"`
	Timestamp time.Time
	URL       string `gorm:"unique"`
	LineageID string
}

func (cache *IdentityCache) connect(automigrate bool) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cache.Filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if automigrate {
		// Perform database migration
		err = db.AutoMigrate(&IdentityValue{})
		if err != nil {
			log.Fatal(err)
		}
	}
	cache.db = db
	return db, nil
}

func (cache *IdentityCache) GetAll() ([]IdentityValue, error) {
	if cache.db == nil {
		cache.connect(true)
		// return nil, fmt.Errorf("database connection is nil")
	}
	if cache.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var identities []IdentityValue
	result := cache.db.Find(&identities)
	if result.Error != nil {
		return nil, result.Error
	}
	return identities, nil
}

func (cache *IdentityCache) Add(identity IdentityValue) error {
	if cache.db == nil {
		cache.connect(true)
		// return nil, fmt.Errorf("database connection is nil")
	}
	result := cache.db.Create(&identity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cache *IdentityCache) ExportAllToCSV(destination string) {
	data, err := cache.GetAll()
	if err != nil {
		fmt.Println(err)
	}

	csvFile, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write([]string{"id", "source", "lineage_id"})
	if err != nil {
		// an error occurred during the flush
		fmt.Println(err)
	}
	for _, v := range data {
		err = csvWriter.Write([]string{strconv.FormatUint(uint64(v.ID), 10), v.URL, v.LineageID})
		if err != nil {
			// an error occurred during the flush
			fmt.Println(err)
		}
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		// an error occurred during the flush
		fmt.Println(err)
	}
}
