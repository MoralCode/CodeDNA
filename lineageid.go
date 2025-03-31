package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/bits"
)

type LineageIdentifier string

// 160 bit or 20 byte hash (40 hex digits)
type CommitHash [20]byte

type LineageID struct {
	idData []byte
	// the number of bits used from the start of each commit
	prefixLength uint8
}

// https://stackoverflow.com/a/10030772/
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func ReverseBits(s []byte) []byte {
	data := []byte(s)
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = bits.Reverse8(data[j]), bits.Reverse8(data[i])
	}
	// handle odd length inputs
	if len(s)%2 != 0 {
		middle := (len(s) - 1) / 2
		data[middle] = bits.Reverse8(data[middle])
	}
	return data
}

func ReverseNibbles(s []byte) []byte {
	data := []byte(s)
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = ReverseNibble(data[j]), ReverseNibble(data[i])
	}
	// handle odd length inputs
	if len(s)%2 != 0 {
		middle := (len(s) - 1) / 2
		data[middle] = ReverseNibble(data[middle])
	}
	return data
}

func ReverseNibble(b byte) byte {
	end := b >> 4
	start := b << 4

	return start | end

}

func ReverseBytes(s []byte) []byte {
	data := []byte(s)
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
func LineageIDFromHashes(commit_hashes []CommitHash, prefixLength uint8) *LineageID {
	lineageID := []byte{}

	if prefixLength > 8 {
		fmt.Errorf("prefix lengths > 8 are not yet supported")
	}

	for _, commithash := range commit_hashes {
		firstbyte := commithash[0]
		prefix := firstbyte >> (8 - prefixLength)
		lineageID = append(lineageID, prefix)
	}
	return &LineageID{
		idData:       ReverseBytes(lineageID),
		prefixLength: prefixLength,
	}
}

func (lineageID *LineageID) String() string {
	return lineageID.StringHex()
}

func (lineageID *LineageID) StringHex() string {
	return hex.EncodeToString(lineageID.idData)
}

func (lineageID *LineageID) StringB64() string {
	return base64.StdEncoding.EncodeToString(lineageID.idData)
}

func (lineageID *LineageID) Bytes() []byte {
	return lineageID.idData
}
