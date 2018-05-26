package indexes

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64

import (
	"sort"
	"strings"
	"sync"
)

// A Trie represents a trie data structure.
type Trie struct {
	root *node
	mx   sync.RWMutex
}

// NewTrie constructs a new empty trie with a root node.
func NewTrie() *Trie {
	return &Trie{
		root: newNode(0),
	}
}

// Insert inserts a new key/value pair entry into the trie,
// where the key is the key and value is user ID.
func (trie *Trie) Insert(key string, userID int64) {
	// Make all keys lowercase, so our search is case-insensitive.
	key = strings.ToLower(key)

	arr := strings.Split(key, " ")

	trie.mx.Lock()
	defer trie.mx.Unlock()

	// handle field value contains multiple words separated by spaces
	for _, key := range arr {
		trie.root.insert(key, userID)
	}
}

// Search retrieves the first n values that match a given prefix string from the trie.
func (trie *Trie) Search(numberToReturn int, prefix string) map[int64]bool {

	trie.mx.RLock()
	defer trie.mx.RUnlock()

	// results is a set
	results := make(map[int64]bool)

	prefix = strings.ToLower(prefix)

	if len(prefix) == 0 {
		return results
	}

	curNode := trie.root

	for _, char := range prefix {
		// Find the child node of current node associated with that character.
		_, exist := curNode.children[char]

		if !exist {
			return results
		}
		curNode = curNode.children[char]
	}

	return curNode.search(numberToReturn, results, 0)
}

// Remove removes a key/value pair entry from the trie,
// where key is a word and value is user ID.
func (trie *Trie) Remove(key string, value int64) {
	key = strings.ToLower(key)
	trie.mx.Lock()
	trie.root.remove(key, value)
	trie.mx.Unlock()
}

// node represents a single node in the trie.
type node struct {
	char     rune
	values   map[int64]bool
	children map[rune]*node
	parent   *node
}

// newNode constrcuts a new node with a given character.
func newNode(char rune) *node {
	return &node{
		char:     char,
		values:   make(map[int64]bool),
		children: make(map[rune]*node),
		parent:   nil,
	}
}

func (root *node) insert(key string, userID int64) {
	curNode := root
	// Loop through each character in the key.
	for _, char := range key {
		// For each character
		// find the child node of current node associated with that character.
		_, hasChild := curNode.children[char]

		// If there is no child node associated with that character,
		// create a new node
		// and add it to current node as a child associated with the character.
		if !hasChild {
			node := newNode(char)
			curNode.children[char] = node
			curNode.children[char].parent = curNode
		}
		// Update current node.
		curNode = curNode.children[char]
	}
	// Add value to current node, which represents
	// the last character in the key.
	_, hasUserID := curNode.values[userID]
	// Ensure only unique user ID can be added.
	if !hasUserID {
		curNode.values[userID] = true
	}
}

// root here is not the root of the trie.
// It represents a node whose char is the last character of the prefix.
func (root *node) search(numberToReturn int, results map[int64]bool, totalResults int) map[int64]bool {
	// Store all user IDs of the current node
	// if there are any.
	if len(root.values) != 0 {
		for userID := range root.values {
			// when our total results length reaches the limit,
			// return the results.
			if totalResults == numberToReturn {
				return results
			}
			// Populate our set.
			_, hasUserID := results[userID]
			if !hasUserID {
				results[userID] = true
			}
			totalResults++
		}
	}

	// Explore all child nodes.
	if len(root.children) > 0 {
		branchResults := make(map[int64]bool)
		sortedChars := []rune{}
		for char := range root.children {
			sortedChars = append(sortedChars, char)
		}
		sort.Slice(sortedChars, func(i, j int) bool {
			return sortedChars[i] < sortedChars[j]
		})

		for _, char := range sortedChars {
			branchResults := root.children[char].search(numberToReturn, branchResults, totalResults)
			for userID := range branchResults {
				// Before add each user ID to results,
				// make sure the limit is not reached yet.
				// If it is already reached, return it.
				if len(results) == numberToReturn {
					return results
				}
				results[userID] = true
			}
			totalResults = len(results)
			// Stop exploring branches if the limit is already reached.
			if totalResults == numberToReturn {
				return results
			}
		}
		return results
	}

	// If the current node has no more children,
	// return the results and trace it back.
	return results
}

func (root *node) remove(key string, value int64) {
	// Find the node whose value we want to remove for a given key.
	curNode := root
	for _, char := range key {
		_, hasChild := curNode.children[char]
		if !hasChild {
			return
		}
		curNode = curNode.children[char]
	}
	// Now our current node is pointing at the node want to remove.
	// Remove the value.
	delete(curNode.values, value)
	curNode.removeUselessNodes()
}

// Trace up and remove dangling nodes.
func (root *node) removeUselessNodes() {

	parentNode := root.parent

	// Remove the node if no other values found in the same node
	// and no child nodes are attached.
	if len(root.values) == 0 && len(root.children) == 0 {
		delete(parentNode.children, root.char)
		parentNode.removeUselessNodes()
	}
	return
}
