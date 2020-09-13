package main

import "fmt"

type page item

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
}

func (notes *item) Print() {
	fmt.Println("Notes:")
	for i := range notes.Tail {
		fmt.Printf("\t %v\t%p\n", notes.Tail[i].Head, &notes.Tail[i])
	}
}

func createItem(selected *item) {
	for n := range selected.Parent.Tail {
		if selected == &selected.Parent.Tail[n] {
			blank := item{Parent: selected.Parent}
			selected.Parent.Tail = append(selected.Parent.Tail, blank)
			copy(selected.Parent.Tail[n+1:], selected.Parent.Tail[n:])
			selected.Parent.Tail[n] = blank
			break
		}
	}
}

func moveUp(selected *item) {
	for n := range selected.Parent.Tail[1:] {
		if selected == &selected.Parent.Tail[n] {
			swapItem(selected.Parent.Tail, n, n-1)
			break
		}
	}
}

func moveDown(selected *item) {
	for n := range selected.Parent.Tail[:len(selected.Parent.Tail)-1] {
		if selected == &selected.Parent.Tail[n] {
			swapItem(selected.Parent.Tail, n, n+1)
			break
		}
	}
}

func swapItem(items []item, current, next int) {
	items[current], items[next] = items[next], items[current]
}

func removeSelected(selected *item) {
	for n := range selected.Parent.Tail {
		if selected == &selected.Parent.Tail[n] {
			selected.Parent.Tail = append(
				selected.Parent.Tail[:n],
				selected.Parent.Tail[n+1:]...)
			break
		}
	}
}

// moves item into the Tail of the preceeding item in the slice depth it lived in
func indent(i *item) {
}

// places item after its Parent item, in its Parent items slice
func unindent(i *item) {
}
