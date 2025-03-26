package main

type LineageIdentifier string

type LineageID struct {
	idData string
	// the number of bytes used from the start of each commit
	prefixLength uint8
}

// https://stackoverflow.com/a/10030772/
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func LineageIDFromHashes(commit_hashes []string) *LineageID {
	lineageID := ""

	for _, element := range commit_hashes {
		lineageID += string(element[0])
	}
	return &LineageID{
		idData:       Reverse(lineageID),
		prefixLength: 4,
	}
}

func (lineageID *LineageID) String() string {
	return lineageID.idData
}
