package main

import "fmt"

type item struct {
	Home   bool
	Parent *item
	Head   string
	Tail   []item
	//open   bool // is this Tail collapsed
}

//type line []*line

func main() {

	house := item{
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
	for n, i := range person.Parent.Tail {
		if &i == person {
			x := person.Parent.Tail[n-1]
			person.Parent.Tail[n-1] = i
			person.Parent.Tail[n] = x
		}
	}

}
func moveUp() {
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
