package blocklist

import (
	"strings"

	radix "github.com/armon/go-radix"
)

type Blocklist interface {
	Create()
	Add(domain string)
	Contains(domain string) bool
	Len() int
}

type BasicBlocklist struct {
	blockList map[string]struct{}
}

func NewBasicBlocklist() Blocklist {
	bbl := &BasicBlocklist{}
	bbl.Create()
	return bbl
}

func (bbl *BasicBlocklist) Create() {
	bbl.blockList = make(map[string]struct{})
}

func (bbl *BasicBlocklist) Add(name string) {
	bbl.blockList[name] = struct{}{}
}

func (bbl *BasicBlocklist) Len() int {
	return len(bbl.blockList)
}

func (bbl *BasicBlocklist) Contains(name string) bool {
	_, ok := bbl.blockList[name]
	return ok
}

type RadixBlocklist struct {
	blockTree *radix.Tree
}

func NewRadixBlocklist() Blocklist {
	rbl := &RadixBlocklist{}
	rbl.Create()
	return rbl
}

func (rbl *RadixBlocklist) Create() {
	rbl.blockTree = radix.New()
}

func (rbl *RadixBlocklist) Add(name string) {
	rbl.blockTree.Insert(reverseString(name), 1)
}

func (rbl *RadixBlocklist) Len() int {
	return rbl.blockTree.Len()
}

func (rbl *RadixBlocklist) Contains(name string) bool {
	name = reverseString(name)
	match, _, ok := rbl.blockTree.LongestPrefix(name)
	if !ok {
		return false
	}

	return strings.HasPrefix(name, match)
}

// Reverse a string
// Taken from: https://stackoverflow.com/a/10030772
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
