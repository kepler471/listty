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
		s += i.Tail[index].Head + "\t"
	}
	return s
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

func (i *item) InsertAt(tail []item, index int) {
	tail = append(tail, *i)
	copy(tail[index+1:], tail[index:])
	tail[index] = *i
}

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
	i.InsertAt(i.Parent.Parent.Tail, index+1)
}

func newItem(i *item) {
	index := i.Locate()
	blank := item{Parent: i.Parent}
	blank.InsertAt(i.Parent.Tail, index+1)
}

func swapItem(tail []item, current, next int) {
	tail[current], tail[next] = tail[next], tail[current]
}
