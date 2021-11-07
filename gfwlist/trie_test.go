package gfwlist

import "testing"

func TestDomainTrie_Search(t *testing.T) {
	trie := NewTrie()
	if trie.Search("www.google.com") {
		t.Fatal()
	}
	trie.Insert("google.com")
	if !trie.Search("www.google.com") {
		t.Fatal()
	}
	if trie.Search("ogle.com") {
		t.Fatal()
	}
	if trie.Search("com") {
		t.Fatal()
	}
	if trie.Search("auctions.yahoo.co.jp") {
		t.Fatal()
	}

	trie.Insert("auctions.yahoo.co.jp")

	if trie.Search("yahoo.co.jp") {
		t.Fatal()
	}
	if !trie.Search("auctions.yahoo.co.jp") {
		t.Fatal()
	}
	if trie.Search("xxauctions.yahoo.co.jp") {
		t.Fatal()
	}
	if !trie.Search("auctions.yahoo.co.jp") {
		t.Fatal()
	}

	trie.Insert("news.163.com")
	if trie.Search("163.com") {
		t.Fatal()
	}
	trie.Insert("ent.163.com")
	if trie.Search("163.com") {
		t.Fatal()
	}
	trie.Insert("163.com")
	if !trie.Search("163.com") {
		t.Fatal()
	}
	if trie.Search("com") {
		t.Fatal()
	}
}
