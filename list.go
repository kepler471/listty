package main

type page item

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
}

func createItem(parent *item) {
	child := item{
		Home:   false,
		Parent: parent,
		Head:   "",
		Tail:   nil,
	}
	parent.Tail = append(parent.Tail, child)
}

func moveUp(person *item) {
	for n, i := range person.Parent.Tail[1:] {
		if &i == person {
			swapItem(person.Parent.Tail, n, n-1)
			break
		}
	}
}

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
