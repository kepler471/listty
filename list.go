package main

type page item

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
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

func removeChild(child *item, parent *item) {
	for n := range parent.Tail {
		if child == &child.Parent.Tail[n] {
			parent.Tail = append(parent.Tail[:n], parent.Tail[n+1:]...)
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
