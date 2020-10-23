package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const (
	nextItem = "\n"
	indent   = "\t"
)

// parseTxt creates a Item tree structure from a txt, and returns the root Item
func parseTxt(filename string) *Item {
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
	root := Item{
		Root: true,
		Text: strings.TrimSpace(strings.TrimSuffix(path.Base(filename), ".txt")),
	}
	splitAndSearch(&root, -1, lines)
	return &root
}

// splitAndSearch creates a tree structure on a given Item i, by searching for
// indentations in lines of text. It will recursively search on all direct children
// of i, with subgroups of the original set of lines.
func splitAndSearch(i *Item, depth int, lines []string) {
	depth++
	var splits []int
	for n, l := range lines {
		d := strings.Count(l, indent)
		if d == depth {
			splits = append(splits, n)
			line := Item{
				Parent: i,
				// TODO: add single whitespace at line end
				Text: strings.TrimSpace(l),
			}
			i.Children = append(i.Children, &line)
		}
	}
	// splits recorded the indices at which to divide lines, and the function
	// is then called on each of these subgroups, to search the next depth level.
	for n, t := range i.Children {
		if splits[n] == splits[len(splits)-1] {
			splitAndSearch(t, depth, lines[splits[n]:])
		} else {
			splitAndSearch(t, depth, lines[splits[n]:splits[n+1]])
		}
	}
}

const (
	INDENTATION = "\t"
	NEW_LINE    = "\n"
	PREFIX      = ""
)

type ItemIteratee = func(child *Item, depth int)

func treeToTxt(tree *Item, fileName string) {
	err := ioutil.WriteFile(fileName+".txt", *treeToBytes(tree), 0644)

	if err != nil {
		fmt.Println(err)
	}
}

func treeToBytes(tree *Item) *[]byte {
	var bytes []byte

	// purposely ignore root node, as it is readonly
	treeIterator(tree, 0, func(i *Item, depth int) {
		bytes = append(bytes, []byte(formattedString(i.Text, depth))...)
	})

	return &bytes
}

func treeIterator(i *Item, depth int, iteratee ItemIteratee) {
	i.ForEachChild(func(child *Item, _ int) {
		iteratee(child, depth)
		if !child.IsLeaf() {
			treeIterator(child, depth+1, iteratee)
		}
	})
}

func formattedString(head string, depth int) string {
	return strings.Repeat(INDENTATION, depth) + PREFIX + head + NEW_LINE
}
