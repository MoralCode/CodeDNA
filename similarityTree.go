package main

import "github.com/MoralCode/CodeDNA/utils"

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

func (tree *SimilarityTreeNode) Add(value string) {

	inValueLen := len(value)
	treeValueLen := len(tree.Value)
	sharedPrefixLen := len(utils.GetLongestPrefix(value, tree.Value))
	// if no value left, base case
	if len(value) == 0 {
		return
	} else if sharedPrefixLen == treeValueLen {
		// if the value completely matches, traverse into child
		lookupRune := rune(value[sharedPrefixLen+1])
		lookupVal, hasLookup := tree.Children[lookupRune]
		if hasLookup {
			lookupVal.Add(value[sharedPrefixLen:])
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
