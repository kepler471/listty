package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const (
	nextItem = "\n"
	indent   = "\t"
)

// parseTxt creates a item tree structure from a txt, and returns the root item
func parseTxt(filename string) *item {
	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}
	// TODO: use a Reader, so the whole file is not loaded at once
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Did not load file: ", err)
	}
	s := string(f)
	lines := strings.Split(s, nextItem)
	root := item{
		Home: true,
		Head: strings.TrimSuffix(path.Dir(filename), ".txt"),
	}
	splitAndSearch(&root, -1, lines)
	return &root
}

// splitAndSearch creates a tree structure on a given item i, by searching for
// indentations in lines of text. It will recursively search on all direct children
// of i, with subgroups of the original set of lines.
func splitAndSearch(i *item, depth int, lines []string) {
	depth++
	var splits []int
	for n, l := range lines {
		d := strings.Count(l, indent)
		if d == depth {
			splits = append(splits, n)
			line := item{
				Parent: i,
				Head:   strings.TrimLeft(l, "\t"),
			}
			i.Tail = append(i.Tail, &line)
		}
	}
	// splits recorded the indices at which to divide lines, and the function
	// is then called on each of these subgroups, to search the next depth level.
	for n, t := range i.Tail {
		if splits[n] == splits[len(splits)-1] {
			splitAndSearch(t, depth, lines[splits[n]:])
		} else {
			splitAndSearch(t, depth, lines[splits[n]:splits[n+1]])
		}
	}
}
