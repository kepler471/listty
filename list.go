package main

import (
	"fmt"
)

type page item

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
	//open   bool // is this Tail collapsed
}

//type line []*line

func main() {

	house := page{
		Home:   true,
		Parent: nil,
		Head:   "papa",
		Tail:   []item{},
	}
	fmt.Println(house.Head)
}

//mama := item{
//Parent: Home,
//Head:   "mama",
//Tail:   nil,
//}
//papa := item{
//Parent: nil,
//Head:   "papa",
//Tail:   nil,
//}

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
	if &person.Parent.Tail[0] == person {
		return
	}

	for n, i := range person.Parent.Tail {
		if &i == person {
			swapItem(person.Parent.Tail, &i, n, n-1)
		}
	}
}

func moveDown(person *item) {
	if &person.Parent.Tail[len(person.Parent.Tail)-1] == person {
		return
	}

	for n, i := range person.Parent.Tail {
		if &i == person {
			swapItem(person.Parent.Tail, &i, n, n+1)
		}
	}
}

func swapItem(items []item, i *item, currentPosition int, newPosition int) {
	x := items[newPosition]
	items[newPosition] = *i
	items[currentPosition] = x
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
func indent(i *item) item {

}

// places item after its Parent item, in its Parent items slice
func unindent(i *item) item {

}

/*
a
	b
	*
		c
		d

*/
