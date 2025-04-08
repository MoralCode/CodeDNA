package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/MoralCode/CodeDNA/utils"
)

type SimilarityTreeNode struct {
	Value string
	// mapping of a prefix
	children map[rune]*SimilarityTreeNode
	Parent   *SimilarityTreeNode

	// [16]*SimilarityTree
}

// Split a node's value into two nodes at the point specified by the given length
// This is done in a way that preserves the base node and returns the newly-split node as a value
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
	tail.children = (*tree).children

	// Step 3: Update HEAD (the original node)
	(*tree).children = newHeadChildren
	(*tree).Value = newHeadValue

	// Step 4: connect detached chain to TAIL
	for _, child := range tail.children {
		child.Parent = &tail
	}

	// Step 5: Cleanup
	// jk golang is garbage collected so this should just happen :tm:

	return &tail, nil
}

// internal function to add a null node to the given tree node
// this is only meant to be internal behavior, not something that general
// consumers of this tree structure should need to do
func (tree *SimilarityTreeNode) addNullNode() (*SimilarityTreeNode, error) {
	lookupVal, hasLookup := tree.children[rune(0)]
	if hasLookup {
		return lookupVal, nil
	} else {
		nullNode := SimilarityTreeNode{
			Parent:   tree,
			children: map[rune]*SimilarityTreeNode{},
			Value:    "",
		}
		tree.children[rune(0)] = &nullNode
		return &nullNode, nil
	}
}

// Add new nodes to the tree until the entire value has been added
// Returns:
//  1. the node that represents the value being added (either created or existing)
//  2. Any Auxiliary nodes that were created (such as the tail portion of a split node, or the null node for an add)
//  3. error (if any)
func (tree *SimilarityTreeNode) Add(value string) (*SimilarityTreeNode, *SimilarityTreeNode, error) {

	inValueLen := len(value)
	treeValueLen := len(tree.Value)
	sharedPrefixLen := len(utils.GetLongestPrefix(value, tree.Value))
	maxPossiblePrefixLen := min(inValueLen, treeValueLen)

	// if no value left, base case
	if inValueLen == 0 {
		// dont try and return nil if we called this on the root node (which has no parent)
		if tree.Parent != nil {
			return tree.Parent, nil, nil
		} else {
			return tree, nil, nil
		}
		//We do not create a new null node here since calling Add() with
		// a value of len 0 (empty string) is not considered valid.
		// null nodes should only be created when theres an exact value match with an existing node for some real value
	}

	if sharedPrefixLen > maxPossiblePrefixLen {
		return nil, nil, errors.New("shared prefix has exceeded its maximum possible length")
	} else if sharedPrefixLen < maxPossiblePrefixLen {
		//partial match but both the tree and incoming value have more to go
		// so split and create a new child
		newSplit, err := (*tree).Split(sharedPrefixLen)
		if err != nil {
			return nil, nil, err
		}

		newSubValue := value[sharedPrefixLen:]
		// create a new node representing the differing part of the value
		node := SimilarityTreeNode{
			Parent:   tree,
			children: map[rune]*SimilarityTreeNode{}, //empty map
			Value:    newSubValue,
		}
		// add it to the now-split root node
		(*tree).children[rune(newSubValue[0])] = &node
		return &node, newSplit, nil

	} else if sharedPrefixLen == maxPossiblePrefixLen {
		if inValueLen < treeValueLen {
			// perfect match for part of this node
			// split
			newSplit, err := (*tree).Split(sharedPrefixLen)
			if err != nil {
				return nil, nil, err
			}

			_, err = tree.addNullNode()
			if err != nil {
				return nil, nil, err
			}
			return tree, newSplit, nil

		} else if inValueLen == treeValueLen {
			// this node matches the value perfectly with no leftovers. search complete
			// Create null node, but still return the current node as primary
			nullNode, err := tree.addNullNode()
			if err != nil {
				return nil, nil, err
			}
			return tree, nullNode, nil
		} else if inValueLen > treeValueLen {
			// search limited by current tree value, traverse into children
			lookupRune := rune(value[sharedPrefixLen])
			lookupVal, hasLookup := tree.children[lookupRune]
			if hasLookup {
				return (*lookupVal).Add(value[sharedPrefixLen:])
			} else {
				// no sub value exists, create it
				node := SimilarityTreeNode{
					Parent:   tree,
					children: map[rune]*SimilarityTreeNode{}, //empty map
					Value:    value[sharedPrefixLen:],
				}
				tree.children[lookupRune] = &node
				return &node, nil, nil
			}
		}
	}
	return nil, nil, errors.New("add finished without result. This is not supposed to happen")
}

// Traverse down the tree to find the leaf node representing the given value
func (tree *SimilarityTreeNode) Find(value string) (*SimilarityTreeNode, error) {
	inValueLen := len(value)
	treeValueLen := len(tree.Value)
	// sharedPrefix :=
	sharedPrefixLen := len(utils.GetLongestPrefix(value, tree.Value))
	maxPossiblePrefixLen := min(inValueLen, treeValueLen)

	if inValueLen == 0 {
		// value found (prior node)
		// dont try and return nil if we called this on the root node (which has no parent)\
		if tree.Parent != nil {
			return tree.Parent, nil
		} else {
			return tree, nil
		}
	} else if sharedPrefixLen < maxPossiblePrefixLen {
		return nil, errors.New("node does not exist. matches stopped in the middle of a node")

	} else if sharedPrefixLen == maxPossiblePrefixLen {
		if inValueLen < treeValueLen {
			// perfect match for part of this node
			return nil, errors.New("node does not exist. search key was exhausted in the middle of a node")

		} else if inValueLen == treeValueLen {
			// this node matches the value perfectly with no leftovers. search complete
			return tree, nil

		} else if inValueLen > treeValueLen {
			// search limited by node value, traverse into children
			lookupRune := rune(value[sharedPrefixLen])
			lookupVal, hasLookup := tree.children[lookupRune]
			if hasLookup {
				return (*lookupVal).Find(value[sharedPrefixLen:])
			} else {
				// no sub value exists, error
				return nil, errors.New("node does not exist. child could not be found")
			}
		}
	}
	return nil, errors.New("search finished without result")
}

func (tree *SimilarityTreeNode) IsLeaf() bool {
	return len(tree.children) == 0 || tree.children[rune(0)] != nil
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
func (tree *SimilarityTreeNode) Print(level int) {

	indents := strings.Repeat("\t", level)
	valLen := len(tree.Value)

	val := ""
	if valLen == 0 {
		val = "<empty value>"
	} else if valLen < 10 {
		val = tree.Value
	} else {
		val = tree.Value[0:5] + "..." + tree.Value[valLen-5:] + " (" + strconv.Itoa(valLen) + ")"
	}
	if tree.IsLeaf() {
		val += " [LEAF]"
	}

	fmt.Println(indents+"Value:", val)
	for k, v := range tree.children {
		if k != rune(0) {
			fmt.Println(indents + "Child " + string(k) + ":")
			v.Print(level + 1)
		}

	}
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

func (tree *SimilarityTreeNode) Siblings() []*SimilarityTreeNode {
	if tree.Parent == nil {
		return []*SimilarityTreeNode{}
	}

	siblingMap := tree.Parent.children

	v := make([]*SimilarityTreeNode, 0, len(siblingMap))

	for key, value := range siblingMap {
		// exclude null nodes because those only serve as pointers from the leaf detection parts of the graph
		if key != rune(0) {
			v = append(v, value)
		}
	}
	return v
}

func (tree *SimilarityTreeNode) Family() string {
	parents := tree.parentChain()

	for _, p := range parents {
		if p.Parent != nil && p.Parent.Parent != nil && len(p.Children()) > 0 {
			return p.Parent.FullValueTo(p.Parent.Parent)
		} else if p.Parent != nil && p.Parent.Parent != nil && len(p.Children()) > 0 {
			return p.Parent.FullValue()
		} else if p.Parent == nil {
			return "orphan"
		}
	}
	return "orphan"
}

// Allow callers to query the children of a node in the tree
// the purpose of this function is to both be an abstraction,
// and to filter out null nodes as they are an internal construct
func (tree *SimilarityTreeNode) Children() []*SimilarityTreeNode {
	childNodes := []*SimilarityTreeNode{}
	for k, v := range tree.children {
		if k != rune(0) {
			childNodes = append(childNodes, v)
		}
	}
	return childNodes
}

// Allow callers to query the presence of children in the tree
// the purpose of this function is to pass through the "has" capability
// of the golang map underlying this structure since it is private
func (tree *SimilarityTreeNode) Child(value rune) (*SimilarityTreeNode, bool) {
	if value == rune(0) {
		return nil, false
	}

	child, has := tree.children[value]
	return child, has
}

// Return all leaf nodes in a particular part of the tree.
// Leaf nodes represent/point to nodes whose FullValue() represents the LineageID for a particular repository
// Since these identifiers can end in the middle of other identifiers
// (such as an old abandoned repo that was later picked up by a new maintainer but the original remains as is)
// We expand the traditional "computer science" definition of leaf nodes (i.e. nodes that have no children)
// to also include nodes that are parents of a null node (child with a key of rune(0)), thus allowing "leaf nodes"
// to exist mid-tree (making them more similar to git branches than traditional leaf nodes)
func (tree *SimilarityTreeNode) Leaves() []*SimilarityTreeNode {
	leaves := make([]*SimilarityTreeNode, 0, 5)
	// base case: we are a leaf
	if len(tree.children) == 0 {
		leaves = append(leaves, tree)

	} else {

		for key, child := range tree.children {
			if key == rune(0) {
				// if we encounter a null node, that means the parent (i.e. the current tree) is also a leaf node
				leaves = append(leaves, tree)
			} else {
				leaves = append(leaves, child.Leaves()...)
			}
		}
	}
	// TODO: handle leaf values that end in the middle of a tree, perfectly on a node but where that node also has children
	// maybe this could use the null rune(0) value as the key?
	return leaves
}

func (tree *SimilarityTreeNode) TreePath() string {
	// base case: root node
	if tree.Parent == nil {
		if len(tree.Value) == 0 {
			return ""
		} else {
			return string(tree.Value[0])
		}
	}

	return tree.Parent.TreePath() + string(tree.Value[0])
}

func (tree *SimilarityTreeNode) parentChain() []*SimilarityTreeNode {
	// base case: root node
	if tree.Parent == nil {
		return []*SimilarityTreeNode{tree}
	}

	return append([]*SimilarityTreeNode{tree}, tree.Parent.parentChain()...)
}

// Find the closest common ancestor
func (a *SimilarityTreeNode) CommonAncestorWith(b *SimilarityTreeNode) (*SimilarityTreeNode, error) {

	chain := a.parentChain()

	parent := b

	for parent != nil {
		if idx := slices.Index(chain, parent); idx > -1 {
			return parent, nil
		}
		parent = parent.Parent
	}
	return nil, errors.New("no shared parentage between the nodes")
}

func (root *SimilarityTreeNode) SimilarityScore(source1Node *SimilarityTreeNode, source2Node *SimilarityTreeNode) (int, error) {

	commonAncestor, err := source1Node.CommonAncestorWith(source2Node)
	if err != nil {
		return -1, errors.Join(errors.New("failed to calculate common ancestor"), err)
	}
	// source1IndependentDistance := source1Node.DistanceTo(commonAncestor)
	// source2IndependentDistance := source2Node.DistanceTo(commonAncestor)

	source1IndependentDistance := len(source1Node.FullValueTo(commonAncestor))
	source2IndependentDistance := len(source2Node.FullValueTo(commonAncestor))

	fmt.Println(len(commonAncestor.FullValueTo(root)))
	return source1IndependentDistance + source2IndependentDistance, nil

}

// lol maybe this should be called the family tree instead to keep with the CodeDNA naming theme

// SimilarityTree is a higher level structure that exists to keep track of
// labelled leaves in the tree so that the nicknames or source URLs for each repo identified by it can be used to look up the node in the tree less-expensively than the lower level node.Find() function
type SimilarityTree struct {
	Root *SimilarityTreeNode
	// map source to the leaf node
	Leaves map[string]*SimilarityTreeNode
}

func NewSimilarityTree() SimilarityTree {
	return SimilarityTree{
		Root: &SimilarityTreeNode{
			Value:    "",
			children: map[rune]*SimilarityTreeNode{},
			Parent:   nil,
		},
		Leaves: map[string]*SimilarityTreeNode{},
	}
}

func (graph *SimilarityTree) Add(source string, identifier string) error {
	existingLeaf, has := graph.Leaves[source]
	var newNode *SimilarityTreeNode
	newNode, _, err := graph.Root.Add(identifier)
	if err != nil {
		return err
	}
	if !has || existingLeaf != newNode {
		graph.Leaves[source] = newNode
	}
	return nil
}
