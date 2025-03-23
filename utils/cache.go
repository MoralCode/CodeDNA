
func check_cache(path string) string {
	pathparts := strings.Split(path, "/")
	reponame := pathparts[len(pathparts)-1]
	cacheFilename := "./" + reponame

	if !strings.HasSuffix(cacheFilename, ".txt") {
		cacheFilename += ".txt"
	}

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
	}

	log.Fatal(errors.New("Cachepath shouldnt be a directory"))
}

func write_cache(path string, lineageID string) {
	pathparts := strings.Split(path, "/")
	reponame := pathparts[len(pathparts)-1]
	cacheFilename := "./" + reponame

	if !strings.HasSuffix(cacheFilename, ".txt") {
		cacheFilename += ".txt"
	}

	d1 := []byte(lineageID)

	err := os.WriteFile(cacheFilename, d1, 0644)
	check(err)
}
