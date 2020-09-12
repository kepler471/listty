package main

type page item

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
}

// Adds a new item at the same depth as the cursor.
// Given the cursor's selected item, createItem looks at the parent,
// and inserts a new item into the slice `selected.Parent.Tail`, directly
// after the current item.
func createItem(selected *item) {
	for n, i := range selected.Parent.Tail {
		if &i == selected.Parent {
			blank := item{Parent: selected.Parent}
			selected.Parent.Tail = append(selected.Parent.Tail, blank)
			copy(selected.Parent.Tail[n+1:], selected.Parent.Tail[n:])
			selected.Parent.Tail[n] = blank
			break
		}
	}
}

// move an item up its parent's slice
func moveUp(person *item) {
	for n, i := range person.Parent.Tail[1:] {
		if &i == person {
			swapItem(person.Parent.Tail, n, n-1)
			break
		}
	}
}

// move an item down its parent's slice
func moveDown(person *item) {
	for n, i := range person.Parent.Tail[:len(person.Parent.Tail)-1] {
		if &i == person {
			swapItem(person.Parent.Tail, n, n+1)
			break
		}
	}
}

func swapItem(items []item, current, next int) {
	items[current], items[next] = items[next], items[current]
}

func removeChild(child *item, parent *item) {
	for n, i := range parent.Tail {
		if &i == child {
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
