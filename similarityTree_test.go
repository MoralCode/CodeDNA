package main

import (
	"testing"
)

func TestGetFullValue(t *testing.T) {
	childNode := SimilarityTreeNode{
		Value:    "efgh",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		Children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
		Parent:   nil,
	}

	childNode.Parent = &rootNode

	expected := "abcdefgh"
	if val := childNode.FullValue(); val != expected {
		t.Errorf(`FullValue() for child node = %q, was not %q`, val, expected)
	}
	expected = "abcd"
	if val := rootNode.FullValue(); val != expected {
		t.Errorf(`FullValue() for root node = %q, was not %q`, val, expected)
	}
}
