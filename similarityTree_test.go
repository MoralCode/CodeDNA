package main

import (
	"fmt"
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

func TestGetFullValueTo(t *testing.T) {
	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		Children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		Children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
		Parent:   nil,
	}

	childNode.Parent = &rootNode
	childNode2.Parent = &childNode

	if d := childNode.FullValueTo(&childNode); d != "" {
		t.Errorf(`childNode distance to self was: %q, but should have been %q`, d, "")
	}

	if d := childNode2.FullValueTo(&childNode); d != "ijkl" {
		t.Errorf(`childNode2 distance to childnode was: %q, but should have been %q`, d, "ijkl")
	}

	if d := childNode2.FullValueTo(&rootNode); d != "efghijkl" {
		t.Errorf(`childNode2 distance to root was: %q, but should have been %q`, d, "efghijkl")
	}
}

func TestAddAppendCase(t *testing.T) {
	// childNode := SimilarityTreeNode{
	// 	Value:    "efgh",
	// 	Children: map[rune]*SimilarityTreeNode{},
	// 	Parent:   nil,
	// }

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	// childNode.Parent = &rootNode

	if len(rootNode.Children) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	targetIdVal := "abcdefgh"
	returnedChild, err := rootNode.Add(targetIdVal)
	if err != nil {
		t.Errorf(`Error: %q`, err)
	}

	if len(rootNode.Children) != 1 {
		t.Errorf(`rootNode failed ending conditions: children should be present, but were not`)
	}

	fmt.Printf("Children after Add: %+v\n", rootNode.Children)

	newChild, exists := rootNode.Children[rune('e')]
	if !exists {
		t.Errorf("Expected child with key 'e', but it was not found")
	}
	expected := "efgh"
	if val := newChild; val.Value != expected {
		t.Errorf(`Value for new child = %q, was not %q`, val.Value, expected)
	}

	if newChild.Parent != &rootNode {
		t.Errorf(`Parent incorrectly set for new child`)
	}

	if fullVal := (*newChild).FullValue(); fullVal != targetIdVal {
		t.Errorf(`full value %q did not match %q`, fullVal, targetIdVal)
	}

	if returnedChild != newChild {
		t.Errorf(`Returned Child is not the same`)
	}
}

func TestAddSplitCase(t *testing.T) {

	rootNode := &SimilarityTreeNode{
		Value:    "abcdfghi",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	if len((*rootNode).Children) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	targetIdVal := "abcdefgh"
	returnedChild, err := (*rootNode).Add(targetIdVal)
	if err != nil {
		t.Errorf(`Error: %q`, err)
	}

	targetChildren := 2
	if l := len((*rootNode).Children); l != targetChildren {
		t.Errorf(`rootNode failed ending conditions: %q should be present, but %q were actually`, targetChildren, l)
	}

	ogChild, exists := (*rootNode).Children[rune('f')]
	if !exists {
		t.Errorf("Expected child with key 'f', but it was not found")
	}
	expected := "fghi"
	if val := *ogChild; val.Value != expected {
		t.Errorf(`Value for original child = %q, was not %q`, val.Value, expected)
	}

	newChild, exists := (*rootNode).Children[rune('e')]
	if !exists {
		t.Errorf("Expected child with key 'e', but it was not found")
	}
	expected = "efgh"
	if val := newChild; val.Value != expected {
		t.Errorf(`Value for new child = %q, was not %q`, val.Value, expected)
	}

	if fullVal := (*newChild).FullValue(); fullVal != targetIdVal {
		t.Errorf(`full value %q did not match %q`, fullVal, targetIdVal)
	}

	if newChild.Parent != rootNode {
		t.Errorf(`Parent incorrectly set for new child`)
	}

	if returnedChild != newChild {
		t.Errorf(`Returned Child is not the same`)
	}
}

func TestAddShorterCase(t *testing.T) {

	rootNode := &SimilarityTreeNode{
		Value:    "abcdfghi",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	if len((*rootNode).Children) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	(*rootNode).Add("abc")

	targetChildren := 1
	if l := len((*rootNode).Children); l != targetChildren {
		t.Errorf(`rootNode failed ending conditions: %q should be present, but %q were actually`, targetChildren, l)
	}

	ogChild, exists := (*rootNode).Children[rune('d')]
	if !exists {
		t.Errorf("Expected child with key 'd', but it was not found")
	}
	expected := "dfghi"
	if val := *ogChild; val.Value != expected {
		t.Errorf(`Value for original child = %q, was not %q`, val.Value, expected)
	}
}

func TestDistance(t *testing.T) {
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

	if d := rootNode.Distance(); d != 0 {
		t.Errorf(`rootNode distance was: %q, but should have been %q`, d, 0)
	}

	if d := childNode.Distance(); d != 1 {
		t.Errorf(`childNode distance was: %q, but should have been %q`, d, 1)
	}

}

func TestDistanceTo(t *testing.T) {

	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		Children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		Children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
		Parent:   nil,
	}

	childNode.Parent = &rootNode
	childNode2.Parent = &childNode

	if d := childNode.DistanceTo(&childNode); d != 0 {
		t.Errorf(`childNode distance to self was: %q, but should have been %q`, d, 0)
	}

	if d := childNode2.DistanceTo(&childNode); d != 1 {
		t.Errorf(`childNode2 distance to childnode was: %q, but should have been %q`, d, 1)
	}

	if d := childNode2.DistanceTo(&rootNode); d != 2 {
		t.Errorf(`childNode2 distance to root was: %q, but should have been %q`, d, 2)
	}

}

func TestCommonAncestor(t *testing.T) {

	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		Children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	childNodeA := SimilarityTreeNode{
		Value:    "wxyz",
		Children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value: "abcd",
		Children: map[rune]*SimilarityTreeNode{
			rune('e'): &childNode,
			rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	childNodeA.Parent = &rootNode
	childNode2.Parent = &childNode

	test := SimilarityTree{
		Root:   &rootNode,
		Leaves: map[string]*SimilarityTreeNode{},
	}

	if a, err := test.CommonAncestor(&childNode, &childNode2); a != &childNode && err != nil {
		t.Errorf(`common ancestor between childNode and childnode2 was node with value: %q, but should have been node with value %q`, a.Value, childNode.Value)
	}

	if a, err := test.CommonAncestor(&childNodeA, &childNode2); a != &rootNode && err != nil {
		t.Errorf(`common ancestor between childNode and childnode2 was node with value: %q, but should have been node with value %q`, a.Value, rootNode.Value)
	}

}
