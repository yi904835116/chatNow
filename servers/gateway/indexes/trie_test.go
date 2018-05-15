package indexes

//TODO: implement automated tests for your trie data structure

import (
	"testing"
)

func TestSearch(t *testing.T) {

	cases := []struct {
		name            string
		keys            []string
		prefix          string
		resultsLimit    int
		expectedtLength int
	}{
		{
			"shared prefix",
			[]string{"do", "dog", "domain", "doout", "dig", "apple"},
			"do",
			20,
			4,
		},
		{
			"no shared prefix",
			[]string{"awesome", "big", "happy"},
			"b",
			20,
			1,
		},
		{
			"empty prefix",
			[]string{"awesome", "big", "happy"},
			"",
			20,
			0,
		},
		{
			"empty trie",
			[]string{},
			"",
			20,
			0,
		},
		{
			"exceeds list limitation",
			[]string{"do", "dog", "domain", "doout", "dig", "apple"},
			"d",
			3,
			3,
		},
		{
			"duplicated keys",
			[]string{"do", "do", "do", "door", "desk", "cat"},
			"do",
			4,
			4,
		},
		{
			"duplicated keys with results limit",
			[]string{"dog", "dog", "dog", "door", "desk", "cat"},
			"do",
			2,
			2,
		},
		{
			"different casting",
			[]string{"Dog", "DOG", "dog", "door", "deSk", "cat"},
			"d",
			20,
			5,
		},
	}

	for _, each := range cases {
		// For each case, construct a new trie.
		trie := NewTrie()

		// Build our test trie.
		id := int64(0)
		for _, key := range each.keys {
			trie.Insert(key, id)
			id++
		}

		result := trie.Search(each.resultsLimit, each.prefix)
		if len(result) != each.expectedtLength {
			t.Errorf("\ncase: %v\ngot: %v\nwant: %v", each.name, len(result), each.expectedtLength)
		}
	}

	specialCases := []struct {
		name                 string
		keys                 []string
		prefix               string
		expectedResultLength int
	}{
		{
			"different keys have same values",
			[]string{"dog", "do", "dope"},
			"do",
			1,
		},
	}

	for _, c := range specialCases {
		userID := int64(0)
		trie := NewTrie()
		for _, key := range c.keys {
			trie.Insert(key, userID)
			userID++
		}

		result := trie.Search(20, c.prefix)
		if len(result) != c.expectedResultLength {
			t.Errorf("\ncase: %v\ngot: %v\nwant: %v", c.name, len(result), c.expectedResultLength)
		}
	}
}

func TestRemove(t *testing.T) {

	// Fake values that represent user ID.
	values := []int64{
		int64(0),
		int64(1),
		int64(2),
		int64(3),
		int64(4),
	}

	cases := []struct {
		name           string
		keys           []string
		key            string
		value          int64
		expectedLength int
	}{
		{
			"target node has child nodes",
			[]string{"dog", "do", "dope", "cat"},
			"do",
			values[1],
			2,
		},
		{
			"target node has no child nodes",
			[]string{"dog", "do", "dope", "cat"},
			"dog",
			values[0],
			0,
		},
		{
			"target node has multiple values",
			[]string{"do", "do", "do", "dog", "dope"},
			"do",
			values[0],
			4,
		},
		{
			"case-insensitive remove",
			[]string{"do", "do", "do", "dog", "dope"},
			"DO",
			values[0],
			4,
		},
		{
			"remove empty key",
			[]string{"do", "dooog"},
			"",
			values[0],
			0,
		},
		{
			"empty trie",
			[]string{},
			"do",
			values[1],
			0,
		},
	}

	for _, each := range cases {
		trie := NewTrie()

		for i, key := range each.keys {
			trie.Insert(key, values[i])
		}

		trie.Remove(each.key, each.value)

		result := trie.Search(20, each.key)
		if len(result) != each.expectedLength {
			t.Errorf("\ncase: %v\ngot: %v\nwant: %v", each.name, len(result), each.expectedLength)
		}
	}

	// Test removing special cases.
	specialCases := []struct {
		name           string
		keys           []string
		key            string
		value          int64
		testKey        string
		expectedOutput int // Expected child nodes length.
	}{
		{
			"remove useless with its node",
			[]string{"do", "dog"},
			"dog",
			values[1],
			"do",
			0,
		},
		{
			"remove multiple useless nodes",
			[]string{"do", "doooont"},
			"doooont",
			values[1],
			"do",
			0,
		},
		{
			"remove multiple useless nodes when parent node has multiple child nodes",
			[]string{"do", "dooog", "dot", "dog"},
			"dooog",
			values[1],
			"do",
			2,
		},
	}

	for _, each := range specialCases {
		trie := NewTrie()

		for i, key := range each.keys {
			trie.Insert(key, values[i])
		}

		trie.Remove(each.key, each.value)

		// Find the node pointing to the last character in the test key.
		cur := trie.root
		for _, char := range each.testKey {
			_, hasKey := cur.children[char]
			if !hasKey {
				t.Error("error finding node")
			}
			cur = cur.children[char]
		}

		if len(cur.children) != each.expectedOutput {
			t.Errorf("\ncase: %v\ngot: %v\nwant: %v", each.name, len(cur.children), each.expectedOutput)
		}
	}
}
