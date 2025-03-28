package utils

import (
	"log"
	"strings"
)

// Takes two strings and trims both to be of equal length
// Returns: trimmed string 1, trimmed string 2, excess from string 1 (or empty), excess from string 2 (or empty)
func TrimStringsToEqualLength(str1 string, str2 string) (string, string, string, string) {
	len1 := len(str1)
	len2 := len(str2)

	if len1 == len2 {
		return str1, str2, "", ""
	} else if len1 > len2 {
		return str1[:len2], str2, str1[len2:], ""
	} else {
		// string2 must be longer
		return str1, str2[:len1], "", str2[len1:]
	}
}

func GetLongestPrefix(str1 string, str2 string) string {

	length := len(str1)

	if length != len(str2) {
		log.Fatal("must check prefix on strings of equal length")
	}

	if length == 0 {
		return ""
	} else if length == 1 {
		if str1 == str2 {
			return str1
		} else {
			return ""
		}
	}

	// integer division
	splitpoint := length / 2

	left1 := str1[:splitpoint]
	left2 := str2[:splitpoint]

	right1 := str1[splitpoint:]
	right2 := str2[splitpoint:]

	// if the left half is a prefix for the whole thing
	if strings.HasPrefix(str1, left2) {
		// we have half of the prefix, traverse right
		return left2 + GetLongestPrefix(right1, right2)
	} else {
		// prefix is shorter than half, go left, discard right
		return GetLongestPrefix(left1, left2)
	}
}
