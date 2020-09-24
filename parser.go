package main

import (
	"io/ioutil"
	"log"
	"strings"
)

const (
	nextItem = "\n"
	indent   = "\t"
)

func parseTxt(filename string) {
	f, err := ioutil.ReadFile(filename + ".txt")
	if err != nil {
		log.Fatal("Did not load file: ", err)
	}
	s := string(f)
	root := item{
		Home:   true,
		Parent: nil,
		Head:   "root",
		Tail:   nil,
		depth:  -1,
	}
	lines := strings.Split(s, nextItem)

	splitAndSearch(&root, -1, lines)
}

func splitAndSearch(i *item, depth int, lines []string) {
	depth++
	var splits []int
	for n, l := range lines {
		d := strings.Count(l, indent)
		if d == depth {
			splits = append(splits, n)
			line := item{
				Home:   false,
				Parent: i,
				Head:   l,
				Tail:   nil,
				depth:  depth,
			}
			i.Tail = append(i.Tail, &line)
		}
	}
	for n, t := range i.Tail {
		if splits[n] == splits[len(splits)-1] {
			group := lines[splits[n]:]
			splitAndSearch(t, depth, group)
		} else {
			group := lines[splits[n]:splits[n+1]]
			splitAndSearch(t, depth, group)
		}
	}
}

// TODO: Add func to clean escape characters from line
