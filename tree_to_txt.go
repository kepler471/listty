package main

import (
	"fmt"
	"io/ioutil"
)

const (
	INDENTATION = "\t"
	NEW_LINE    = "\n"
	PREFIX      = "- "
)

type ItemIteratee = func(child *item, depth int)

func treeToTxt(tree *item, fileName string) {
	err := ioutil.WriteFile(fileName+".txt", *treeToBytes(tree), 0644)

	if err != nil {
		fmt.Println(err)
	}
}

func treeToBytes(tree *item) *[]byte {
	var bytes []byte

	// purposely ignore root node, as it is readonly
	treeIterator(tree, 0, func(i *item, depth int) {
		bytes = append(bytes, []byte(formattedString(i.Head, depth))...)
	})

	return &bytes
}

func treeIterator(i *item, depth int, iteratee ItemIteratee) {
	i.ForEachChild(func(child *item, _ int) {
		iteratee(child, depth)
		if !child.IsLeaf() {
			treeIterator(child, depth+1, iteratee)
		}
	})
}

func formattedString(head string, depth int) string {
	return addIndentation(depth) + PREFIX + head + NEW_LINE
}

func addIndentation(depth int) string {
	indents := ""

	if depth == 0 {
		return indents
	}

	for i := depth; i > 0; i-- {
		indents += INDENTATION
	}

	return indents
}
