package main

import (
	"fmt"
)

type Item struct {
	Root      bool
	Text      string
	Parent    *Item
	Children  []*Item
	Collapsed bool
}

type TreeIteratee func(i *Item)
type TailIteratee func(i *Item, idx int)

// TODO: Need to look at which methods should return pointers to items.
// 	Would make tracking Item movement much easier.

func (i *Item) Print() {
	fmt.Println("Notes:")
	for n := range i.Children {
		fmt.Printf("\t %v\t%p\n", i.Children[n].Text, &i.Children[n])
	}
}

func (i *Item) IsLeaf() bool {
	return i.Children == nil || len(i.Children) == 0
}

func (i *Item) StringChildren() (s string) {
	for index := range i.Children {
		s += i.Children[index].Text + ","
	}
	return
}

// Path returns a slice of items from root to i.
func (i *Item) Path() []string {
	var p []string
	return reverse(i.path(p))
}

func (i *Item) path(p []string) []string {
	p = append(p, i.Text)
	if i.Parent == nil {
		return p
	}
	return i.Parent.path(p)
}

// Path returns a slice of items from j to i.
func (i *Item) PathTo(j *Item) []string {
	var p []string
	return reverse(i.pathTo(j, p))
}

func (i *Item) pathTo(j *Item, p []string) []string {
	p = append(p, i.Text)
	if i.Parent == j {
		return p
	}
	return i.Parent.pathTo(j, p)
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Remove the Item from its parent's tail
func (i *Item) Remove() {
	index := i.Locate()
	i.Parent.Children = append(i.Parent.Children[:index], i.Parent.Children[index+1:]...)
}

func swapItem(tail []*Item, current, next int) {
	tail[current], tail[next] = tail[next], tail[current]
}

func (i *Item) MoveUp() {
	index := i.Locate()
	if index == 0 {
		return
	}
	swapItem(i.Parent.Children, index, index-1)
}

func (i *Item) MoveDown() {
	index := i.Locate()
	if index == len(i.Parent.Children)-1 {
		return
	}
	swapItem(i.Parent.Children, index, index+1)
}

// Locate returns the index value for the selected Item, within its parent's tail
func (i *Item) Locate() int {
	var loc int
	for index := range i.Parent.Children {
		if i == i.Parent.Children[index] {
			loc = index
		}
	}
	return loc
}

// Indent moves the Item to the end of the preceding Item's tail
func (i *Item) Indent() *Item {
	index := i.Locate()
	if index == 0 {
		return i
	}
	j := i.Parent.Children[index-1]
	i.Remove()
	j.Children = append(j.Children, i)
	i.Parent = j
	return i
}

// Unindent moves an Item after its Parent Item, in its Parent slice
func (i *Item) Unindent() *Item {
	if i.Parent.Root {
		return i
	}
	j := i.Parent.Parent
	index := i.Parent.Locate()
	i.Parent.AddSibling(i, index+1)
	i.Remove()
	i.Parent = j
	return i
}

// InsertAlongside places itself next to a target Item
func (i *Item) InsertAlongside(j *Item, index int) {
	j.Parent.Children = append(j.Parent.Children, i)
	copy(j.Parent.Children[index+1:], j.Parent.Children[index:])
	j.Parent.Children[index] = i
}

// AddSibling places the target Item j alongside i in i's
// parent's tail. Usuall, is called with index = i.Locate()+1.
func (i *Item) AddSibling(j *Item, index int) *Item {
	i.Parent.Children = append(i.Parent.Children, j)
	copy(i.Parent.Children[index+1:], i.Parent.Children[index:])
	i.Parent.Children[index] = j
	return j
}

func (i *Item) AddChild(j *Item) *Item {
	index := i.Locate()
	i.AddSibling(j, index+1).Indent()
	return j
}

func newItem(i *Item) *Item {
	index := i.Locate()
	j := Item{Parent: i.Parent, Text: "~"}
	i.AddSibling(&j, index+1)
	return &j
}

// invoke TreeIteratee on each Item in Children
func (i *Item) ForEachChild(iteratee TailIteratee) {
	if i.Children == nil {
		return
	}

	for j := 0; j < len(i.Children); j++ {
		curItem := &i.Children[j]

		if curItem != nil {
			iteratee(*curItem, j)
		}
	}
}

// Implementation always uses root: "From root, get to current Item using the PositionStack"
func getCurrentItem(root *Item, stack *PositionStack) *Item {
	currentItem := root

	currentItemIterator(root, stack, func(nextItem *Item) {
		currentItem = nextItem
	})

	return currentItem
}

// From root get to last Item invoking TreeIteratee on each Item
func currentItemIterator(root *Item, stack *PositionStack, iteratee TreeIteratee) {
	// could set count to a different depth to iterate from -> to other nodes in tree 🤔
	_toLastItemInStack(root, stack, iteratee, 0)
}

func _toLastItemInStack(root *Item, stack *PositionStack, iteratee TreeIteratee, count int) {
	if stack.GetLast().Depth == count {
		iteratee(root)
	} else {
		_toLastItemInStack(root.Children[stack.GetRow(count)], stack, iteratee, count+1)
	}
}

// TreeMap build a flat data structure for an Item tree, with row numbers for easy
// printing to screen.
func TreeMap(i *Item, m map[int]*Item) {
	if !i.IsLeaf() {
		for _, t := range i.Children {
			m[len(m)] = t
			TreeMap(t, m)
		}
	}
}
