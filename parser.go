package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	nextItem = "\n"
	indent   = "\t"
)

func parseTxt3(filename string) {
	f, err := ioutil.ReadFile(filename + ".txt")
	if err != nil {
		log.Fatal("Did not load file: ", err)
	}
	s := string(f)
	root := item{
		Home:   true,
		Parent: nil,
		Head:   "root",
		Tail:   nil,
		depth:  -1,
		row:    -1,
	}
	lines := strings.Split(s, nextItem)

	splitAndSearch(&root, -1, lines)
	fmt.Println(lines)
}

//func tempCreate()

func splitAndSearch(i *item, depth int, lines []string) {
	depth++
	var splits []int
	for n, l := range lines {
		d := strings.Count(l, indent)
		if d == depth {
			splits = append(splits, n)
			line := item{
				Home:   false,
				Parent: i,
				Head:   l,
				Tail:   nil,
				depth:  depth,
				row:    n,
			}
			i.Tail = append(i.Tail, &line)
		}
	}
	fmt.Println("\n", i.Head)
	i.Print()
	for n, t := range i.Tail {
		if splits[n] == splits[len(splits)-1] {
			group := lines[splits[n]:]
			splitAndSearch(t, depth, group)
		} else {
			group := lines[splits[n]:splits[n+1]]
			splitAndSearch(t, depth, group)
		}
	}
	fmt.Println("Return")
}

//
//func parseTxt2(filename string) {
//	f, err := ioutil.ReadFile(filename + ".txt")
//	if err != nil {
//		log.Fatal("Did not load file: ", err)
//	}
//	s := string(f)
//	root := item{
//		Home:   true,
//		Parent: nil,
//		Head:   "root",
//		Tail:   nil,
//		depth:  -1,
//		row:    -1,
//	}
//	lines := strings.Split(s, nextItem)
//	first := item{
//		Home:   false,
//		Parent: &root,
//		Head:   lines[0],
//		Tail:   nil,
//		depth:  0,
//		row:    0,
//	}
//	root.Tail = append(root.Tail, first)
//	_fillTail(&root, &first, 0, 0, lines)
//	fmt.Println("asdsdas")
//}
//
//// Recursively iterates through the list of notes, determining indentation on each line and
//// 	handles notes on equal depth, indented by 1, or unindented by n.
//// ISSUES:
////	The calls back up the stack mean that some tails have items inserted in the wrong
//// 		order. This may be easily fixed by walking through the tree and sorting every tail by row.
////	depth == depth0 + 1 branch returns on an item, which is not added to root. It needs to
////		be combined somehow, but at least the structure is being built somewhere
//func _fillTail(root *item, i *item, row0 int, depth0 int, lines []string) *item {
//	row := row0 + 1
//	if row >= len(lines) {
//		return i
//	}
//	depth := strings.Count(lines[row], indent)
//
//	if depth == depth0 { // !hasTail || isLeaf
//		j := item{
//			Parent: i.Parent, // this was incorrectly assigned to i
//			Head:   lines[row],
//			depth:  depth,
//			row:    row,
//		}
//		fmt.Printf("50 :: row: %v, depth %v, head: %v\n", row, depth, lines[row])
//		i.Parent.Tail = append(i.Parent.Tail, j) // Add as sibling
//		return _fillTail(root, &j, j.row, i.depth, lines)
//	}
//	if depth == depth0+1 { // hasTail || !isLeaf
//		j := item{
//			Parent: i,
//			Head:   lines[row],
//			depth:  depth,
//			row:    row,
//		}
//		fmt.Printf("62 :: row: %v, depth %v, head: %v\n", row, depth, lines[row])
//		_fillTail(root, &j, j.row, j.depth, lines) // _fillTail on new item j ...
//		// The above line will keep parsing text and adding to j's tree. It needs to stop if
//		// it ever comes back to a depth of i.depth, as that item would be in i.Tail, a
//		// sibling of j.
//		// One way of doing this could be walking through the branch and creating a copy, up
//		// to this known limit, i.depth.
//		//branch := item{}
//		//trimmedBranch := _newBranch(&branch, &j, &j, j.row, j.depth, lines)
//		i.Tail = append(i.Tail, j) // ... then append j on i after it has recursed
//		// We also need the last item added to j, call k. This will then be passed on for the next iteration
//		//return _fillTail(root, &j, j.row, j.depth, lines)
//	}
//	if depth < depth0 {
//		j := item{
//			Head:  lines[row],
//			depth: depth,
//			row:   row,
//		}
//		fmt.Printf("77 :: row: %v, depth %v, head: %v\n", row, depth, lines[row])
//		_lastAtDepth(root, i, &j)
//		return _fillTail(root, &j, j.row, j.depth, lines)
//	}
//	return i
//}
//
////// Start at top of j, and build a new branch up to the row and depth limits
////func _newBranch(branch *item, i *item, j *item, rowLimit int, depthLimit int, lines []string) *item {
////	if i.depth == depthLimit && i.row > rowLimit {
////		return nil // here we need to break the recursion and return the new branch built up to this point
////	}
////	if len(i.Tail) > 0 {
////		for index, t := range i.Tail {
////			// Only need to append to tails, as we are still working with pointers, so
////			// parents should still be assigned correctly
////			branch.Tail = append(branch.Tail, t)
////			branch.Tail[index]
////		}
////	}
////	branch.Tail = i.Tail
////
////}
//
//// find the previous item at j.depth. Use i as a working variable.
//// Assign parent to j and place in the matching tail.
//func _lastAtDepth(root *item, i *item, j *item) {
//	if j.depth == i.Parent.depth {
//		j.Parent = i.Parent.Parent
//		i.Parent.Parent.Tail = append(i.Parent.Parent.Tail, *j)
//		return
//	}
//	// We do not need to catch for exceeding root as a valid tree
//	// will always give a match
//	//if i.Parent.Parent == nil {
//	//	return
//	//}
//	_lastAtDepth(root, i.Parent, j)
//}
//
////func _lastStringMatch(root *item, depth int, find string, targetDepth int, ans []*item) []*item {
////	for _, i := range root.Tail {
////		if i.Head == find && depth == targetDepth {
////			ans = append(ans, &i)
////		}
////		if !i.IsLeaf() {
////			_lastStringMatch(&i, depth+1, find, targetDepth)
////		}
////	}
////	return ans
////}
//
//func parseNext(root item, items []item, n int, lines []string) []item {
//	i := items[n-1]
//	dep := strings.Count(lines[n], indent)
//	j := item{
//		Head:  lines[n],
//		depth: dep,
//	}
//	switch {
//	case j.depth == i.depth:
//		index := i.Locate()
//		j.Parent = i.Parent
//		i.AddSibling(&j, index+1)
//	case j.depth == i.depth+1:
//		j.Parent = &i
//		i.Tail = append(i.Tail, j)
//		//case j.depth < i.depth:
//
//	}
//	items[n-1] = i
//	items = append(items, j)
//	return items
//}
//
//// parseTxt builds an item tree structure from a text document
//func parseTxt() {
//	f, err := ioutil.ReadFile("example2.txt")
//	if err != nil {
//		log.Fatal("Did not load file: ", err)
//	}
//	s := string(f)
//	lines := strings.Split(s, nextItem)
//	//root := item{
//	//	Home:   true,
//	//	Parent: nil,
//	//	Head:   "root",
//	//	Tail:   nil,
//	//	depth:  -1,
//	//}
//	//tsil := make(map[loc]*item)
//	list := make(map[int]*item)
//	for n, l := range lines {
//		depth := strings.Count(l, indent)
//		switch {
//		case n == 0:
//			root := item{
//				Home:   true,
//				Parent: nil,
//				Head:   "root",
//				Tail:   nil,
//				depth:  -1,
//			}
//			i := item{
//				Home:   false,
//				Parent: &root,
//				Head:   l,
//				Tail:   nil,
//				depth:  depth,
//			}
//			//root.Tail = append(root.Tail, i)
//			list[-1] = &root
//			list[-1].Tail = append(list[-1].Tail, i)
//			list[n] = &i
//		case depth == list[n-1].depth:
//			i := item{
//				Home:   false,
//				Parent: list[n-1].Parent,
//				Head:   l,
//				Tail:   nil,
//				depth:  depth,
//			}
//			list[n-1].Parent.Tail = append(list[n-1].Parent.Tail, i)
//			//list[n-1].AddSibling(&i, list[n-1].Locate()+1)
//			list[n] = &i
//		case depth == list[n-1].depth+1:
//			i := item{
//				Home:   false,
//				Parent: list[n-1],
//				Head:   l,
//				Tail:   nil,
//				depth:  depth,
//			}
//			list[n-1].Tail = append(list[n-1].Tail, i)
//			//list[n-1].AddChild(&i)
//			list[n] = &i
//		case depth < list[n-1].depth:
//			for dy := 1; true; dy++ {
//				if list[n-dy].depth == depth {
//					i := item{
//						Home:   false,
//						Parent: list[n-dy].Parent,
//						Head:   l,
//						Tail:   nil,
//						depth:  depth,
//					}
//					list[n-dy].Parent.Tail = append(list[n-dy].Parent.Tail, i)
//					//list[n-dy].AddSibling(&i, list[n-dy].Locate()+1)
//					list[n] = &i
//					break
//				}
//			}
//		}
//		//fmt.Printf("Depth: %v, %v\n", list[n].depth, list[n].Head)
//	}
//	for n := 0; n < len(list); n++ {
//		fmt.Printf("Depth: %v, %v\n", list[n].depth, list[n].Head)
//	}
//}
