package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func hashesFromStrings(hashes []string) []CommitHash {
	hashdata := []CommitHash{}

	for _, v := range hashes {
		h, err := hex.DecodeString(v)
		if err != nil {
			fmt.Errorf("failed to decode")
		}
		hashdata = append(hashdata, CommitHash(h))
	}
	return hashdata
}

func TestFromHashes(t *testing.T) {
	hashes := []string{
		"c157c5bb882fffe4932853ee413a36af63c337d9",
		"75288f635132b98b366e6993be945f3c9ddf8f05",
		"3cafb499963675d22f44007c91b906e77d45dfb5",
		"e48d65529880ebd2d061c8bfa13e78b74c411204",
		"e3a4055fb9d8afe217d73591bfb2724662fa86fc",
		"94e7ba5ba88de06ad0943bcb6facf12f9a9c2eee",
	}
	hashdata := hashesFromStrings(hashes)

	id := LineageIDFromHashes(hashdata, 4)
	if id.StringHex() != "9ee37c" {
		t.Errorf(`LineageIDFromHashes() = %q, was not %q`, id.StringHex(), "9ee37c")
	}
}

func TestAssembleBytes(t *testing.T) {
	nibbles := []uint8{0xA, 0xB, 0xC} // Example input
	result, err := AssembleBytesFromNibbles(nibbles)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Assembled bytes: %X\n", result)
	if r := hex.EncodeToString(result); r != "abc0" {
		t.Errorf(`AssembleBytesFromNibbles() = %q, was not %q`, r, "abc0")
	}
}

func TestReverseBytes(t *testing.T) {
	nibbles := []uint8{0xA, 0xB, 0xC} // Example input
	result := ReverseBytes(nibbles)

	if result[0] != 0xC || result[1] != 0xB || result[2] != 0xA {
		t.Errorf(`ReverseBytes() = %X, was not %X`, result, []uint8{0xC, 0xB, 0xA})
	}

}

func TestReverseNibbles(t *testing.T) {
	nibbles := []uint8{0x1, 0x2, 0x3} // Example input
	result := ReverseNibbles(nibbles)

	if result[0] != 0x30 || result[1] != 0x20 || result[2] != 0x10 {
		t.Errorf(`ReverseBytes() = %X, was not %X`, result, []uint8{0x30, 0x20, 0x10})
	}

}
