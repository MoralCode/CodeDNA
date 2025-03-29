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
	Root *SimilarityTreeNode
	// map source to the leaf node
	Leaves map[string]*SimilarityTreeNode
}

// Split a node's value into two nodes at the point specified by the given length
func (tree *SimilarityTreeNode) Split(split_length int) (*SimilarityTreeNode, error) {
	// Step 0. Prerequisites
	if len((*tree).Value) < 2 {
		return nil, errors.New("not enough characters in value to successfully split")
	}

	if split_length > len(tree.Value) {
		return nil, errors.New("split length too long to successfully split")
	}

	if split_length <= 0 {
		return nil, errors.New("split length too short to successfully split")
	}

	// Step 1: Create
	tail := SimilarityTreeNode{
		Value:  (*tree).Value[split_length:],
		Parent: tree,
	}

	newHeadValue := (*tree).Value[:split_length]
	newHeadChildren := map[rune]*SimilarityTreeNode{
		rune(tail.Value[0]): &tail,
	}

	// Step 2: Transfer Children
	tail.Children = (*tree).Children

	// Step 3: Update HEAD (the original node)
	(*tree).Children = newHeadChildren
	(*tree).Value = newHeadValue

	// Step 4: connect detached chain to TAIL
	for _, child := range tail.Children {
		child.Parent = &tail
	}

	// Step 5: Cleanup
	// jk golang is garbage collected so this should just happen :tm:

	return &tail, nil
}

// Add new nodes to the tree until the entire value has been added
func (tree *SimilarityTreeNode) Add(value string) (*SimilarityTreeNode, error) {

	inValueLen := len(value)
	treeValueLen := len(tree.Value)
	sharedPrefixLen := len(utils.GetLongestPrefix(value, tree.Value))
	// if no value left, base case
	if len(value) == 0 {
		return tree.Parent, nil
	} else if sharedPrefixLen == treeValueLen {
		// if the value completely shares a prefix, traverse into child
		lookupRune := rune(value[sharedPrefixLen])
		lookupVal, hasLookup := tree.Children[lookupRune]
		if hasLookup {
			return (*lookupVal).Add(value[sharedPrefixLen:])
		} else {
			// no sub value exists, create it
			node := SimilarityTreeNode{
				Parent:   tree,
				Children: map[rune]*SimilarityTreeNode{}, //empty map
				Value:    value[sharedPrefixLen:],
			}
			tree.Children[lookupRune] = &node
			return &node, nil
		}
	} else if sharedPrefixLen == inValueLen {
		//if the incoming value ends before the end of the current value
		// split
		return (*tree).Split(sharedPrefixLen)
	} else {
		// if incoming value has a match ending in the middle of the current, length of tree value, we need to split it

		_, err := (*tree).Split(sharedPrefixLen)
		if err != nil {
			return nil, err
		}

		newSubValue := value[sharedPrefixLen:]
		// create a new node representing the differing part of the value
		node := SimilarityTreeNode{
			Parent:   tree,
			Children: map[rune]*SimilarityTreeNode{}, //empty map
			Value:    newSubValue,
		}
		// add it to the now-split root node
		(*tree).Children[rune(newSubValue[0])] = &node
		return &node, nil
	}
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

// Get the "full value" of this node (its value, prefixed with the value of all of its parents)
func (tree *SimilarityTreeNode) FullValueTo(node *SimilarityTreeNode) string {

	//base case: we found the target node
	if tree == node {
		return ""
	}

	//other base case: root node
	if tree.Parent == nil {
		return tree.Value
	}

	return tree.Parent.FullValueTo(node) + tree.Value
}

// Get the "distance" of this node to the root
func (tree *SimilarityTreeNode) Distance() int {
	// base case: root node
	if tree.Parent == nil {
		return 0
	}

	return tree.Parent.Distance() + 1
}

func (tree *SimilarityTreeNode) DistanceTo(node *SimilarityTreeNode) int {
	// base case: root node
	if tree.Parent == nil {
		return 0
	}

	// base case, target node found
	if tree == node {
		return 0
	}

	return tree.Parent.DistanceTo(node) + 1
}

func (graph *SimilarityTree) Add(source string, identifier string) error {

	if _, has := graph.Leaves[source]; !has {
		newNode, err := graph.Root.Add(identifier)
		if err != nil {
			return err
		}
		graph.Leaves[source] = newNode
	}
	return nil
}
