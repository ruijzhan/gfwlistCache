package gfwlist

type Interface interface {
	Insert(list, domain string)
	ListsContain(domain string) []string
}

func New() Interface {
	return &listSet{
		domainLists: make(map[string]*DomainTrie),
	}
}

type listSet struct {
	domainLists map[string]*DomainTrie
}

func (ls *listSet) Insert(list, domain string) {
	trie, ok := ls.domainLists[list]
	if !ok {
		trie = NewTrie()
		ls.domainLists[list] = trie
	}
	trie.Insert(domain)
}

func (ls *listSet) ListsContain(domain string) []string {
	res := make([]string, 0, len(ls.domainLists))

	for list, trie := range ls.domainLists {
		if trie.Search(domain) {
			res = append(res, list)
		}
	}
	return res
}
