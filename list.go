package main

import (
	"fmt"
)

type item struct {
	Root     bool
	Text     string
	Parent   *item
	Children []*item
}

type TreeIteratee func(i *item)
type TailIteratee func(i *item, idx int)

// TODO: Need to look at which methods should return pointers to items.
// 	Would make tracking item movement much easier.

func (i *item) Print() {
	fmt.Println("Notes:")
	for n := range i.Children {
		fmt.Printf("\t %v\t%p\n", i.Children[n].Text, &i.Children[n])
	}
}

func (i *item) IsLeaf() bool {
	return i.Children == nil || len(i.Children) == 0
}

func (i *item) StringChildren() (s string) {
	for index := range i.Children {
		s += i.Children[index].Text + ","
	}
	return
}

// Path returns a slice of items from root to i.
func (i *item) Path() []string {
	var p []string
	return reverse(i.path(p))
}

func (i *item) path(p []string) []string {
	p = append(p, i.Text)
	if i.Parent == nil {
		return p
	}
	return i.Parent.path(p)
}

// Path returns a slice of items from j to i.
func (i *item) PathTo(j *item) []string {
	var p []string
	return reverse(i.pathTo(j, p))
}

func (i *item) pathTo(j *item, p []string) []string {
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

// Remove the item from its parent's tail
func (i *item) Remove() {
	index := i.Locate()
	i.Parent.Children = append(i.Parent.Children[:index], i.Parent.Children[index+1:]...)
}

func swapItem(tail []*item, current, next int) {
	tail[current], tail[next] = tail[next], tail[current]
}

func (i *item) MoveUp() {
	index := i.Locate()
	if index == 0 {
		return
	}
	swapItem(i.Parent.Children, index, index-1)
}

func (i *item) MoveDown() {
	index := i.Locate()
	if index == len(i.Parent.Children)-1 {
		return
	}
	swapItem(i.Parent.Children, index, index+1)
}

// Locate returns the index value for the selected item, within its parent's tail
func (i *item) Locate() int {
	var loc int
	for index := range i.Parent.Children {
		if i == i.Parent.Children[index] {
			loc = index
		}
	}
	return loc
}

// Indent moves the item to the end of the preceding item's tail
func (i *item) Indent() *item {
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

// Unindent moves an item after its Parent item, in its Parent slice
func (i *item) Unindent() *item {
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

// InsertAlongside places itself next to a target item
func (i *item) InsertAlongside(j *item, index int) {
	j.Parent.Children = append(j.Parent.Children, i)
	copy(j.Parent.Children[index+1:], j.Parent.Children[index:])
	j.Parent.Children[index] = i
}

// AddSibling places the target item j alongside i in i's
// parent's tail. Usuall, is called with index = i.Locate()+1.
func (i *item) AddSibling(j *item, index int) *item {
	i.Parent.Children = append(i.Parent.Children, j)
	copy(i.Parent.Children[index+1:], i.Parent.Children[index:])
	i.Parent.Children[index] = j
	return j
}

func (i *item) AddChild(j *item) *item {
	index := i.Locate()
	i.AddSibling(j, index+1).Indent()
	return j
}

func newItem(i *item) *item {
	index := i.Locate()
	j := item{Parent: i.Parent, Text: "~"}
	i.AddSibling(&j, index+1)
	return &j
}

// invoke TreeIteratee on each item in Children
func (i *item) ForEachChild(iteratee TailIteratee) {
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

// Implementation always uses root: "From root, get to current item using the PositionStack"
func getCurrentItem(root *item, stack *PositionStack) *item {
	currentItem := root

	currentItemIterator(root, stack, func(nextItem *item) {
		currentItem = nextItem
	})

	return currentItem
}

// From root get to last item invoking TreeIteratee on each item
func currentItemIterator(root *item, stack *PositionStack, iteratee TreeIteratee) {
	// could set count to a different depth to iterate from -> to other nodes in tree 🤔
	_toLastItemInStack(root, stack, iteratee, 0)
}

func _toLastItemInStack(root *item, stack *PositionStack, iteratee TreeIteratee, count int) {
	if stack.GetLast().Depth == count {
		iteratee(root)
	} else {
		_toLastItemInStack(root.Children[stack.GetRow(count)], stack, iteratee, count+1)
	}
}

// TreeMap build a flat data structure for an item tree, with row numbers for easy
// printing to screen.
func TreeMap(i *item, m map[int]*item) {
	if !i.IsLeaf() {
		for _, t := range i.Children {
			m[len(m)] = t
			TreeMap(t, m)
		}
	}
}
