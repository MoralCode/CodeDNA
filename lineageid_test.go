package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

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
