// Package trie provides a simple trie for string prefix matching
package trie

import (
    "strings"
)

type node struct {
    key      string
    children []*node
    isRoot   bool
}

func newNode(s string) (n *node) {
    return &node{key: s, children: make([]*node, 0), isRoot: false}
}

func (n *node) findChild(s string) *node {
    // improve later with some sort of sorted way

    for _, v := range n.children {
        if strings.Compare(v.key, s) == 0 {
            return v
        }
    }
    return nil
}

func (n *node) addChild(s string) *node {
    // improve later with some sort of sorted way
    node := n.findChild(s)
    if node == nil {
        node = newNode(s)
        n.children = append(n.children, node)
    }
    return node
}

// This is a naive implementation of a trie and is not thread safe, any modifications of this trie
// should be mutex locked.
// Trie shape(elements with a ! is a match) when provided with {我,不} {我,是} {w,o,r,d} {w,a,s} {w,o}
//
//         /- 是!
//    /- 我  - 不!
//  *
//    \- w - o! - r - d!
//         \- a - s!
//
// Only those marked with a "!" is considered a match when finding
//
type Trie struct {
    root *node
}

// Creates a new Trie with a "*" root
//
func NewTrie() *Trie {
    return &Trie{root: newNode("*")}
}

// Inserts a new sub tree based on the string slice provided
//
func (t *Trie) Insert(pieces []string) (n *node) {
    if len(pieces) == 0 {
        return nil
    }
    node := t.root
    for _, v := range pieces {
        node = node.addChild(v)
    }
    node.isRoot = true
    return node
}

// This function find matches based on string slice
// returns a channel, True if found, False if not found
//
func (t *Trie) FindMatch(set []string) (result chan bool) {

    result = make(chan bool)
    if len(set) == 0 {
        result <- false
        return
    }
    go func() {
        node := t.root

        for _, v := range set {
            node = node.findChild(v)
            if node == nil {
                result <- false
                return
            }
            if node.isRoot {
                result <- true
                return
            }
        }
        result <- node.isRoot
        return
    }()
    return
}

func (t *Trie) FindMatchLen(set []string) (result chan int) {

    result = make(chan int)
    if len(set) == 0 {
        result <- 0
        return
    }
    go func() {
        node := t.root
        count := 0

        for _, v := range set {
            node = node.findChild(v)
            count = count + 1
            if node == nil {
                result <- 0
                return
            }
            if node.isRoot {
                result <- count
                return
            }
        }
        result <- count
        return
    }()
    return
}
