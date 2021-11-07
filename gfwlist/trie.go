package gfwlist

import "strings"

type node struct {
	domain  string
	nodes   map[string]*node
	matches bool
}

type DomainTrie struct {
	root *node
}

func NewTrie() *DomainTrie {
	return &DomainTrie{
		root: &node{
			domain: ".",
			nodes:  map[string]*node{},
		},
	}
}

func (g *DomainTrie) Insert(domain string) {
	tokens := strings.Split(domain, ".")
	curr := g.root
	for i := len(tokens) - 1; i >= 0; i-- {
		sub := tokens[i]
		if n, ok := curr.nodes[sub]; ok {
			if n.matches {
				return
			}
		} else {
			curr.nodes[sub] = &node{
				domain: sub,
				nodes:  map[string]*node{},
			}
		}
		curr = curr.nodes[sub]
	}
	curr.matches = true
}

func (g *DomainTrie) Search(domain string) bool {
	tokens := strings.Split(domain, ".")
	curr := g.root
	for i := len(tokens) - 1; i >= 0; i-- {
		sub := tokens[i]
		if n, ok := curr.nodes[sub]; ok {
			if n.matches {
				return true
			}
			curr = n
		} else {
			break
		}
	}
	return false
}
