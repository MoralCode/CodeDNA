package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// https://stackoverflow.com/a/10030772/
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getLineageID(repo *git.Repository) string {
	// ... retrieving the HEAD reference
	ref, err := repo.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	// since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	// until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderDFSPost})
	// , Since: &since, Until: &until
	CheckIfError(err)

	lineageID := ""

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		// fmt.Printf("%v (%T)", string(c.Hash.String()[0]), string(c.Hash.String()[0]))
		lineageID += string(c.Hash.String()[0])
		return nil
	})
	CheckIfError(err)
	return Reverse(lineageID)
}

// isValidUrl tests a string to determine if it is a well-structured url or not.
// from https://www.golangcode.com/how-to-check-if-a-string-is-a-url/
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func lineageIDForPath(path string) string {
	pathparts := strings.Split(path, "/")
	reponame := pathparts[len(pathparts)-1]
	cacheFilename := "./" + reponame + ".txt"

	var repo *git.Repository
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
		data, err := os.ReadFile(cacheFilename)
		check(err)
		return string(data)
	} else {
		// do directory stuff

		// We instantiate a new repository object from the given path (the .git folder)
		repo, err = git.PlainOpen(path)
		CheckIfError(err)

		// Length of the HEAD history
		// Info("git rev-list HEAD --count")
		lineageID := getLineageID(repo)

		d1 := []byte(lineageID)

		err = os.WriteFile(cacheFilename, d1, 0644)
		check(err)

		return lineageID
	}
}

func main() {
	CheckArgs("<path> <path2>")
	path := os.Args[1]

	lineageIDForPath(path)

}
