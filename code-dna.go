package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getLineageID(repo *git.Repository) string {
	// ... retrieving the HEAD reference
	ref, err := repo.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	// since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	// until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
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
	return lineageID
}

func main() {
	CheckArgs("<path>")
	path := os.Args[1]

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Length of the HEAD history
	// Info("git rev-list HEAD --count")
	fmt.Println(getLineageID(r))
}
