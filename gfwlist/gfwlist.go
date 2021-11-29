package gfwlist

import (
	"bytes"
	"io"
)

type gfwList struct {
	trie *DomainTrie
}

var (
	g *gfwList
)

func init() {
	g = new()
}

// Has determins if a domain is in gfwlist
func Has(domain string) bool {
	return g.has(domain)
}

func (g *gfwList) has(domain string) bool {
	return g.trie.Search(domain)
}

func new() (gfw *gfwList) {
	gfw = &gfwList{
		trie: NewTrie(),
	}
	rules, err := Asset("gfwlist.txt")
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(rules)
	buff := &bytes.Buffer{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if b == '\n' {
			gfw.trie.Insert(buff.String())
			buff.Reset()
		} else {
			buff.WriteByte(b)
		}
	}
	return
}
