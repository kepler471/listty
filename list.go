package main

import (
	"fmt"
)

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
}

func (i *item) Print() {
	fmt.Println("Notes:")
	for n := range i.Tail {
		fmt.Printf("\t %v\t%p\n", i.Tail[n].Head, &i.Tail[n])
	}
}

func (i *item) StringChildren() (s string) {
	for index := range i.Tail {
		s += i.Tail[index].Head + ","
	}
	return
}

// Locate returns the index value for the selected item, within its parent's tail
func (i *item) Locate() (index int) {
	for index = range i.Parent.Tail {
		if i == &i.Parent.Tail[index] {
			return
		}
	}
	return
}

// InsertAlongside places itself next to a target item
func (i *item) InsertAlongside(j *item, index int) {
	j.Parent.Tail = append(j.Parent.Tail, *i)
	copy(j.Parent.Tail[index+1:], j.Parent.Tail[index:])
	j.Parent.Tail[index] = *i
}

// AddSibling places the target item alongside itself in its
// parent's tail
func (i *item) AddSibling(j *item, index int) {
	i.Parent.Tail = append(i.Parent.Tail, *j)
	copy(i.Parent.Tail[index+1:], i.Parent.Tail[index:])
	i.Parent.Tail[index] = *j
}

// Remove the item from its parent's tail
func (i *item) Remove() {
	index := i.Locate()
	i.Parent.Tail = append(i.Parent.Tail[:index], i.Parent.Tail[index+1:]...)
}

func (i *item) MoveUp() {
	index := i.Locate()
	if index == 0 {
		return
	}
	swapItem(i.Parent.Tail, index, index-1)
}

func (i *item) MoveDown() {
	index := i.Locate()
	if index == len(i.Parent.Tail)-1 {
		return
	}
	swapItem(i.Parent.Tail, index, index+1)
}

// Indent moves the item to the end of the preceding item's tail
func (i *item) Indent() {
	index := i.Locate()
	i.Remove()
	i.Parent.Tail[index-1].Tail = append(i.Parent.Tail[index-1].Tail, *i)
}

// Unindent moves an item after its Parent item, in its Parent slice
func (i *item) Unindent() {
	i.Remove()
	index := i.Parent.Locate()
	i.Parent.AddSibling(i, index+1)
}

func newItem(i *item) {
	index := i.Locate()
	blank := item{Parent: i.Parent}
	i.AddSibling(&blank, index+1)
}

func swapItem(tail []item, current, next int) {
	tail[current], tail[next] = tail[next], tail[current]
}

func unPack(itemToUnPack *item, cursor int, depth int, currentDepth int) *item {
	if itemToUnPack == nil {
		panic("😱")
	}

	if currentDepth == depth {
		return itemToUnPack
	}

	return unPack(&itemToUnPack.Tail[cursor], cursor, depth, currentDepth + 1)
}