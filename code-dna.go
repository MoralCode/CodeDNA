package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"

	"github.com/MoralCode/CodeDNA/utils"
	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v69/github"
)

func getOriginUrlFromRepo(repo *git.Repository) (string, error) {
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}

	return remote.Config().URLs[0], nil
}

func getLineageIDFromRepo(repo *git.Repository, prefixLength uint8) (string, error) {
	// ... retrieving the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	// ... retrieves the commit history
	// since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	// until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderDFSPostNoMerge})
	// , Since: &since, Until: &until
	if err != nil {
		return "", err
	}

	var commit_hashes []CommitHash

	err = cIter.ForEach(func(c *object.Commit) error {
		// here we convert the type so we arent passing around a plumbing.Hash everywhere
		cache := []byte{}
		// fmt.Println(c.Hash.String())
		for _, b := range c.Hash {
			cache = append(cache, b)
		}
		// fmt.Println(hex.EncodeToString(cache))
		commit_hashes = append(commit_hashes, CommitHash(cache))
		return nil
	})
	if err != nil {
		return "", err
	}
	lineageID := LineageIDFromHashes(commit_hashes, prefixLength)
	return lineageID.String(), nil
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func repoOwnerAndNameFromURL(repourl string) (string, string) {
	parsedurl, err := url.Parse(repourl)
	if err != nil {
		log.Fatal(err)
	}

	pathparts := strings.Split(parsedurl.Path, "/")
	reponame := pathparts[len(pathparts)-1]
	owner := pathparts[len(pathparts)-2]
	return owner, reponame
}

func lineageIDFromGitHub(repourl string, prefixLength uint8) string {

	// TODO: maybe use  https://github.com/shurcooL/githubv4
	if !isValidUrl(repourl) {
		log.Fatal(errors.New("url is not valid"))
	}

	client := github.NewClient(nil)
	ctx := context.Background()

	owner, reponame := repoOwnerAndNameFromURL(repourl)

	var allCommits []*github.RepositoryCommit

	var opt = &github.CommitsListOptions{
		SHA:         "HEAD",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {

		commits, resp, err := client.Repositories.ListCommits(ctx, owner, reponame, opt)
		check(err)
		// if err != nil {
		// 	return err
		// }
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Println("checking page", resp.NextPage, "from Github REST API")
		opt.Page = resp.NextPage
	}

	var commit_hashes []CommitHash

	for _, commit := range allCommits {
		hashbytes, err := hex.DecodeString(*commit.SHA)
		CheckIfError(err)
		commit_hashes = append(commit_hashes, CommitHash(hashbytes))

	}

	lineageID := LineageIDFromHashes(commit_hashes, prefixLength)

	// err = os.WriteFile(cacheFilename, d1, 0644)
	// check(err)
	return lineageID.String()
}

func cloneRepo(repourl string, into string) error {
	if !strings.HasPrefix(repourl, "http") {
		repourl = "https://" + repourl
	}
	_, err := git.PlainClone(into, true, &git.CloneOptions{
		URL:               repourl,
		RecurseSubmodules: 0,
		// Differently than the git CLI, by default go-git downloads
		// all tags and its related objects. To avoid unnecessary
		// data transmission and processing, opt-out tags.
		Tags:         git.NoTags,
		SingleBranch: true,
		// Not a net positive change for performance, this was added
		// to better align the output when compared with the git CLI.
		Progress: os.Stdout,
	})
	return err
}

func lineageIDFromGitClone(repourl string, tempdir string, prefixLength uint8) string {
	err := cloneRepo(repourl, tempdir)
	repo, err := git.PlainOpen(tempdir)
	CheckIfError(err)

	lineageId, err := getLineageIDFromRepo(repo, prefixLength)
	CheckIfError(err)
	return lineageId

}

type RepoImport struct {
	RepoSource string
	Nickname   string
}

func importManyRepos(filename string) ([]RepoImport, error) {
	// open file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var repos []RepoImport

	for _, element := range data {

		new := RepoImport{
			RepoSource: element[0],
		}

		if len(element) > 1 {
			new.Nickname = strings.TrimSpace(element[1])
		} else {
			new.Nickname = element[0]
		}

		repos = append(repos, new)
	}
	return repos, nil
}

func analyzeRepo(analysisPath string, prefixLength uint8) (string, string, error) {
	fmt.Println("Starting analysis for", analysisPath)
	var lineageID string
	var source string

	// classify path type
	if isValidUrl(analysisPath) {
		fmt.Println("Querying from github...")
		lineageID = lineageIDFromGitHub(analysisPath, prefixLength)
		source = analysisPath
	} else if _, err := os.Stat(analysisPath); errors.Is(err, os.ErrNotExist) {
		return "", "", err
	} else {
		fmt.Println("Reading from disk...")
		var repo *git.Repository

		// We instantiate a new repository object from the given path (the .git folder)
		repo, err := git.PlainOpen(analysisPath)
		CheckIfError(err)

		lineageID, err = getLineageIDFromRepo(repo, prefixLength)
		CheckIfError(err)
		source, err = getOriginUrlFromRepo(repo)
		CheckIfError(err)

	}
	return source, lineageID, nil
}

func bulkCloneTask(id int, cache *utils.IdentityCache, tempdir string, data chan RepoImport) {

	for repo := range data {
		owner, repoName := repoOwnerAndNameFromURL(repo.RepoSource)
		fmt.Println("Importing", repoName, "from", owner, "as \""+repo.Nickname+"\"")
		cloneDir := tempdir + "/" + owner + "_" + repoName

		err := os.MkdirAll(cloneDir, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = cloneRepo(repo.RepoSource, cloneDir)
		if err != nil {
			fmt.Println(err)
		}

		gitrepo, err := git.PlainOpen(cloneDir)
		if err != nil {
			fmt.Println("open")

			fmt.Println(err)
			continue
		}

		lineageID, err := getLineageIDFromRepo(gitrepo, 4)
		if err != nil {
			fmt.Println("get id")

			fmt.Println(err)
			continue
		}

		if !cache.Has(repo.RepoSource) {
			newValue := utils.IdentityValue{
				URL:       repo.RepoSource,
				LineageID: lineageID,
			}
			if repo.Nickname != "" {
				newValue.Nickname = repo.Nickname
			}
			err := cache.Add(newValue)
			if err != nil {
				fmt.Println("add to cache")
				fmt.Println(err)
				// errors <- err
				continue
			} else {
				err = os.RemoveAll(cloneDir)
				if err != nil {
					fmt.Println("cleanup")
					fmt.Println(err)
					// errors <- err
					continue
				}
			}
		}
	}
}

// https://github.com/jessevdk/go-flags/issues/405
// https://github.com/jessevdk/go-flags/issues/387
// I think this arg parsing lib is abandoned.....
type Analyze struct {
	Enabled bool `hidden:"true" no-ini:"true"`

	Args struct {
		Repository string `description:"The repository to analyze" required:"true"`
		Nickname   string `description:"A nickname to assign to the new record"`
	} ` positional-args:"yes"`
}

type Export struct {
	Enabled bool   `hidden:"true" no-ini:"true"`
	Path    string `long:"path" description:"The path to export to" default:"database.csv"`
}

type ImportCommand struct {
	Enabled       bool   `hidden:"true" no-ini:"true"`
	Path          string `long:"path" description:"The path to import from" required:"true"`
	CloneExisting bool   `long:"clone-existing" description:"whether or not to clone a repository if it exists in the cache"`
}

type SimilarityCommand struct {
	Enabled bool `hidden:"true" no-ini:"true"`
}

type MainCmd struct {
	Verbosity  []bool            `short:"v" long:"verbose" description:"Show verbose debug information"`
	CachePath  string            `long:"cachepath" default:"cache.sqlite" description:"The path to the cache database to use"`
	Analyze    Analyze           `command:"analyze" description:"Analyze a repository"`
	Export     Export            `command:"export" description:"export the database to CSV"`
	Import     ImportCommand     `command:"import" description:"import from CSV"`
	Similarity SimilarityCommand `command:"similarity" description:"run repo similarity report"`
}

// Detect when the subcommand is used.
func (c *Analyze) Execute(args []string) error {
	c.Enabled = true
	return nil
}
func (c *Export) Execute(args []string) error {
	c.Enabled = true
	return nil
}
func (c *ImportCommand) Execute(args []string) error {
	c.Enabled = true
	return nil
}
func (c *SimilarityCommand) Execute(args []string) error {
	c.Enabled = true
	return nil
}

func main() {
	var opts MainCmd

	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
	}

	cache := utils.IdentityCache{
		Filename: opts.CachePath,
	}

	if len(opts.Verbosity) >= 1 {
		fmt.Printf("%+v\n", opts)
	}

	if opts.Analyze.Enabled {
		analysisPath := opts.Analyze.Args.Repository
		source, lineageID, err := analyzeRepo(analysisPath, 4)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Could not Analyze. Attempting fetch from cache...")
			// assume its a name and fetch from cache
			cached, err := cache.GetByNickname(analysisPath)
			if err != nil {
				panic(err)
			}
			lineageID = cached.LineageID
			source = cached.URL
		}

		if !cache.Has(source) {
			newValue := utils.IdentityValue{
				URL:       source,
				LineageID: lineageID,
			}
			if opts.Analyze.Args.Nickname != "" {
				newValue.Nickname = opts.Analyze.Args.Nickname
			}
			cache.Add(newValue)
		}

		fmt.Println(lineageID)
		fmt.Println(source)

	}

	if opts.Import.Enabled {
		fmt.Println("Importing from", opts.Import.Path)
		repos, err := importManyRepos(opts.Import.Path)
		CheckIfError(err)

		tempdir := "./repositories"

		totalRepos := len(repos)
		fmt.Println("Beginning Cloning of", totalRepos, "repositories")

		// Creating a channel
		channel := make(chan RepoImport)
		var wg sync.WaitGroup

		// Creating workers to execute the task
		for i := 0; i < 8; i++ {
			fmt.Println("Main: Starting worker", i)
			wg.Add(1)
			go bulkCloneTask(i, &cache, tempdir, channel)
		}

		batch := repos[750:900]

		for _, repo := range batch {
			if !opts.Import.CloneExisting && cache.Has(repo.RepoSource) {
				fmt.Println("\t Source exists in cache, skipping")
				// processedRepos += 1
				continue
			}

			channel <- repo
			// processedRepos += 1

		}

	}

	if opts.Export.Enabled {
		fmt.Println("Exporting db to", opts.Export.Path)
		cache.ExportAllToCSV(opts.Export.Path)
	}

	if opts.Similarity.Enabled {
		tree := NewSimilarityTree()

		cache, err := cache.GetAll()
		CheckIfError(err)
		// Add all repos to tree
		for _, v := range cache {
			// TODO: use url if no nickname available
			err = tree.Add(v.Nickname, v.LineageID)
			CheckIfError(err)
		}

		tree.Root.Print(0)
		fmt.Println("")
		fmt.Println("===========")
		fmt.Println("")

		leafFamilies := []string{}

		for id, leaf := range tree.Leaves {
			idString := leaf.TreePath()
			idString += " \t( "
			idString += id
			idString += " ):\t "
			idString += leaf.Family()
			leafFamilies = append(leafFamilies, idString)
		}

		sort.Strings(leafFamilies)

		for _, str := range leafFamilies {
			fmt.Println(str)
		}

		// sanity check with prefix lengths

	}

}
