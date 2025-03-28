package main

import (
	"errors"

	"github.com/MoralCode/CodeDNA/utils"
)

type SimilarityTreeNode struct {
	Value string
	// mapping of a prefix
	Children map[rune]*SimilarityTreeNode
	Parent   *SimilarityTreeNode

	// [16]*SimilarityTree
}

type SimilarityTree struct {
	Root   *SimilarityTreeNode
	Leaves []SimilarityTreeNode
}

// Split a node's value into two nodes at the point specified by the given length
func (tree *SimilarityTreeNode) Split(split_length int) error {
	// Step 0. Prerequisites
	if len(tree.Value) < 2 {
		return errors.New("not enough characters in value to successfully split")
	}

	if split_length > len(tree.Value) {
		return errors.New("split length too long to successfully split")
	}

	if split_length <= 0 {
		return errors.New("split length too short to successfully split")
	}

	// Step 1: Create

	tail := SimilarityTreeNode{
		Value: tree.Value[split_length+1:],
	}

	head := SimilarityTreeNode{
		Value: tree.Value[:split_length],
		Children: map[rune]*SimilarityTreeNode{
			rune(tail.Value[0]): &tail,
		},
	}
	tail.Parent = &head

	// Step 2: Transfer Children
	tail.Children = tree.Children

	// Step 3: Add HEAD (rootmost node in new pair) to chain
	tree.Parent.Children[rune(tree.Value[0])] = &head
	head.Parent = tree.Parent

	// Step 4: connect detached chain to TAIL
	for _, child := range tree.Parent.Children {
		child.Parent = &tail
	}

	// Step 5: Cleanup
	// jk golang is garbage collected so this should just happen :tm:

	return nil
}

func (tree *SimilarityTreeNode) Add(value string) {

	inValueLen := len(value)
	treeValueLen := len(tree.Value)
	sharedPrefixLen := len(utils.GetLongestPrefix(value, tree.Value))
	// if no value left, base case
	if len(value) == 0 {
		return
	} else if sharedPrefixLen == treeValueLen {
		// if the value completely matches, traverse into child
		lookupRune := rune(value[sharedPrefixLen])
		lookupVal, hasLookup := tree.Children[lookupRune]
		if hasLookup {
			(*lookupVal).Add(value[sharedPrefixLen:])
		} else {
			// no sub value exists, create it
			node := SimilarityTreeNode{
				Parent:   tree,
				Children: map[rune]*SimilarityTreeNode{}, //empty map
				Value:    value[sharedPrefixLen:],
			}
			tree.Children[lookupRune] = &node
			// TODO: log leaf
		}
	} else if sharedPrefixLen == inValueLen {
		//if the incoming value ends before the end of the current value
		// split

		// TODO: log leaf
	} else {
		// if incoming value has a match ending in the middle of the current, length of tree value, we need to split it
	}
	// lookupVal, hasLookup := tree.Children[lookupRune]
	// short circuit: simple just add case if there is no child matching the first rune of the value
	// if !hasLookup {
	// 	node := SimilarityTreeNode{
	// 		Parent:   tree,
	// 		Children: map[rune]*SimilarityTreeNode{},
	// 		Value:    value,
	// 	}
	// 	tree.Children[lookupRune] = &node
	// } else {
	// 	// if a key for the split already exists
	// 	// check prefix length of the value

	// }
	// compare the first <number of chars in the node> to each branch currently in the tree and see where it should go
}

func (tree *SimilarityTreeNode) IsLeaf() bool {
	return len(tree.Children) == 0
}

// Get the "full value" of this node (its value, prefixed with the value of all of its parents)
func (tree *SimilarityTreeNode) FullValue() string {
	// base case: root node
	if tree.Parent == nil {
		return tree.Value
	}

	return tree.Parent.FullValue() + tree.Value
}
