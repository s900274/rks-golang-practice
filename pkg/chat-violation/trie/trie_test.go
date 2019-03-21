package trie

import (
    . "./"
    "testing"
)


func Example(t *testing.T) {
    trie := NewTrie()
    trie.Insert([]string{"t","e","s","t"})
    trie.Insert([]string{"t","e","s"})
    result := <- trie.FindMatch([]string{"t","e"})
    if result {
        t.Errorf("te should not match test");
    }
    result = <- trie.FindMatch([]string{"t","e", "s", "t"})
    if !result {
        t.Errorf("test should match test");
    }
    result = <- trie.FindMatch([]string{"t","e", "s"})
    if !result {
        t.Errorf("tes should match tes");
    }

}

