package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestGetFullValue(t *testing.T) {
	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
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
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
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
	// 	children: map[rune]*SimilarityTreeNode{},
	// 	Parent:   nil,
	// }

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	// childNode.Parent = &rootNode

	if len(rootNode.Children()) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	targetIdVal := "abcdefgh"
	returnedChild, splitChild, err := rootNode.Add(targetIdVal)
	if err != nil {
		t.Errorf(`Error: %q`, err)
	}

	if len(rootNode.Children()) != 1 {
		t.Errorf(`rootNode failed ending conditions: children should be present, but were not`)
	}

	if splitChild != nil {
		t.Errorf(`simple append adds should not return a split child`)
	}

	fmt.Printf("Children after Add: %+v\n", rootNode.Children())

	newChild, exists := rootNode.Child(rune('e'))
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
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	originalValue := rootNode.FullValue()

	if len((*rootNode).Children()) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	targetIdVal := "abcdefgh"
	returnedChild, splitChild, err := (*rootNode).Add(targetIdVal)
	if err != nil {
		t.Errorf(`Error: %q`, err)
	}

	targetChildren := 2
	if l := len((*rootNode).Children()); l != targetChildren {
		t.Errorf(`rootNode failed ending conditions: %q should be present, but %q were actually`, targetChildren, l)
	}

	ogChild, exists := (*rootNode).Child(rune('f'))
	if !exists {
		t.Errorf("Expected child with key 'f', but it was not found")
	}
	expected := "fghi"
	if val := *ogChild; val.Value != expected {
		t.Errorf(`Value for original child = %q, was not %q`, val.Value, expected)
	}

	if originalValue != splitChild.FullValue() {
		t.Errorf("Split value returned from Add() doesnt match the origi")

	}

	newChild, exists := (*rootNode).Child(rune('e'))
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
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	if len((*rootNode).Children()) != 0 {
		t.Errorf(`rootNode failed initial conditions: children should not be present, but some were`)
	}

	(*rootNode).Add("abc")

	targetChildren := 1
	if l := len((*rootNode).Children()); l != targetChildren {
		t.Errorf(`rootNode failed ending conditions: %q should be present, but %q were actually`, targetChildren, l)
	}

	ogChild, exists := (*rootNode).Child('d')
	if !exists {
		t.Errorf("Expected child with key 'd', but it was not found")
	}
	expected := "dfghi"
	if val := *ogChild; val.Value != expected {
		t.Errorf(`Value for original child = %q, was not %q`, val.Value, expected)
	}
}

func TestFind(t *testing.T) {
	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	childNodeA := SimilarityTreeNode{
		Value:    "wxyz",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootValueNode := SimilarityTreeNode{
		Value: "abcd",
		children: map[rune]*SimilarityTreeNode{
			rune('e'): &childNode,
			rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	rootNode := SimilarityTreeNode{
		Value: "",
		children: map[rune]*SimilarityTreeNode{
			rune('a'): &rootValueNode,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	childNodeA.Parent = &rootValueNode
	childNode2.Parent = &childNode
	rootValueNode.Parent = &rootNode

	if d, err := rootNode.Find(""); d != &rootNode || err != nil {
		t.Errorf(`find exact value of root node (empty) returned node with value %q, but should have been node with value %q`, d.Value, rootNode.Value)
	}

	if d, err := rootNode.Find("abcd"); d != &rootValueNode || err != nil {
		t.Errorf(`find exact value of root node returned node with value %q, but should have been node with value %q`, d.Value, rootValueNode.Value)
	}

	if d, err := rootNode.Find("abcde"); d != nil || err == nil {
		t.Errorf(`find value after root node returned node with value %q, but should have been an error`, d.Value)
	}

	if d, err := rootNode.Find("abcdehgf"); d != nil || err == nil {
		t.Errorf(`find value when matches stop partway returned value %q, but should have been an error`, d.Value)
	}

	if d, err := rootNode.Find("abcdihgf"); d != nil || err == nil {
		t.Errorf(`find value when matches stop at node boundary returned value %q, but should have been an error`, d.Value)
	}

	if d, err := rootNode.Find("abcdwxyz"); d != &childNodeA || err != nil {
		t.Errorf(`find value of child returned node with value %q, but should have been node with value %q`, d.Value, childNodeA.Value)
	}

}

func TestDistance(t *testing.T) {
	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
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
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value:    "abcd",
		children: map[rune]*SimilarityTreeNode{rune('e'): &childNode},
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
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	childNodeA := SimilarityTreeNode{
		Value:    "wxyz",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value: "abcd",
		children: map[rune]*SimilarityTreeNode{
			rune('e'): &childNode,
			rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	childNodeA.Parent = &rootNode
	childNode2.Parent = &childNode

	if a, err := childNode.CommonAncestorWith(&childNode2); a != &childNode && err != nil {
		t.Errorf(`common ancestor between childNode and childnode2 was node with value: %q, but should have been node with value %q`, a.Value, childNode.Value)
	}

	if a, err := childNodeA.CommonAncestorWith(&childNode2); a != &rootNode && err != nil {
		t.Errorf(`common ancestor between childNode and childnode2 was node with value: %q, but should have been node with value %q`, a.Value, rootNode.Value)
	}

}

func TestLeafDetection(t *testing.T) {
	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	childNodeA := SimilarityTreeNode{
		Value:    "wxyz",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	nullNode := SimilarityTreeNode{
		Value:    "",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value: "abcd",
		children: map[rune]*SimilarityTreeNode{
			rune(0):   &nullNode,
			rune('e'): &childNode,
			rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	childNodeA.Parent = &rootNode
	childNode2.Parent = &childNode
	nullNode.Parent = &rootNode

	leaves := rootNode.Leaves()

	if l := len(leaves); l != 3 {
		t.Errorf(`Leaves() returned incorrect number of leaves, expected %q, got %q`, 3, l)
	}

	if !slices.Contains(leaves, &childNode2) {
		t.Errorf(`Leaves() should contain childNode2, but didnt`)
	}

	if !slices.Contains(leaves, &childNodeA) {
		t.Errorf(`Leaves() should contain childNodeA, but didnt`)
	}

	if !slices.Contains(leaves, &nullNode) {
		t.Errorf(`Leaves() should contain nullNode, but didnt`)
	}

}

func TestSiblingDetection(t *testing.T) {
	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	childNodeA := SimilarityTreeNode{
		Value:    "wxyz",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	nullNode := SimilarityTreeNode{
		Value:    "",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	rootNode := SimilarityTreeNode{
		Value: "abcd",
		children: map[rune]*SimilarityTreeNode{
			rune(0):   &nullNode,
			rune('e'): &childNode,
			rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	childNodeA.Parent = &rootNode
	childNode2.Parent = &childNode
	nullNode.Parent = &rootNode

	siblings := rootNode.Siblings()

	if l := len(siblings); l != 0 {
		t.Errorf(`Siblings() returned incorrect number of siblings for root node, expected %d, got %d`, 0, l)
	}

	siblings = childNode.Siblings()

	if l := len(siblings); l != 2 {
		t.Errorf(`Siblings() returned incorrect number of siblings for childNode, expected %d, got %d`, 2, l)
	}

	if slices.Contains(siblings, &childNode2) {
		t.Errorf(`Siblings() of childNode should contain childNode2, but didnt`)
	}

	if !slices.Contains(siblings, &childNodeA) {
		t.Errorf(`Siblings() of childNode should contain childNodeA, but didnt`)
	}

	if slices.Contains(siblings, &nullNode) {
		t.Errorf(`Siblings() should not contain nullNode, but didnt`)
	}

}

func TestLeafMaintainance(t *testing.T) {

	childNode2 := SimilarityTreeNode{
		Value:    "ijkl",
		children: map[rune]*SimilarityTreeNode{},
		Parent:   nil,
	}

	childNode := SimilarityTreeNode{
		Value:    "efgh",
		children: map[rune]*SimilarityTreeNode{rune('i'): &childNode2},
		Parent:   nil,
	}

	// childNodeA := SimilarityTreeNode{
	// 	Value:    "wxyz",
	// 	children: map[rune]*SimilarityTreeNode{},
	// 	Parent:   nil,
	// }

	rootNode := SimilarityTreeNode{
		Value: "abcd",
		children: map[rune]*SimilarityTreeNode{
			rune('e'): &childNode,
			// rune('w'): &childNodeA,
		},
		Parent: nil,
	}

	childNode.Parent = &rootNode
	// childNodeA.Parent = &rootNode
	childNode2.Parent = &childNode

	test := SimilarityTree{
		Root:   &rootNode,
		Leaves: map[string]*SimilarityTreeNode{},
	}

	// simple add
	if v := len(test.Leaves); v != 0 {
		t.Errorf(`similarity tree was expected to have no leaves but instead had: %d`, v)
	}

	test.Add("w-x-y-z", "abcdwxyz")

	if _, has := test.Leaves["w-x-y-z"]; has != true {
		t.Errorf(`similarity tree was expected to have a leaf at key %q: with value %q, but it doesnt exist`, "w-x-y-z", "wxyz")
	}

	if v, has := test.Leaves["w-x-y-z"]; has == true && v.Value != "wxyz" {
		t.Errorf(`similarity tree was expected to have a leaf at key %q: with value %q, but instead had node with value %q`, "w-x-y-z", "wxyz", v.Value)
	}

	// add a value that will split the just-added value
	test.Add("w-x", "abcdwx")

	if _, has := test.Leaves["w-x"]; has != true {
		t.Errorf(`similarity tree was expected to have a leaf at key %q: with value %q, but it doesnt exist`, "w-x", "wx")
	}

	if v, has := test.Leaves["w-x"]; has == true && v.Value != "wx" {
		t.Errorf(`similarity tree was expected to have a leaf at key %q: with value %q, but instead had node with value %q`, "w-x", "wx", v.Value)
	}

}
