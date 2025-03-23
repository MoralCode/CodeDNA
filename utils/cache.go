package utils

import (
	"errors"
	"log"
	"os"

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
		check(err)
		return string(data)
	}

	log.Fatal(errors.New("Cachepath shouldnt be a directory"))
}

func (cache IdentityCache) write_cache(path string, lineageID string) {

	d1 := []byte(lineageID)

	err := os.WriteFile(cache.Filename, d1, 0644)
	check(err)
}
