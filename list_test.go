package main

import (
	"testing"
)

func TestItem_InsertAt(t *testing.T) {
	notes := item{Home: true, Head: "Homepage"}
	ch := []string{"£", "$", "%", "^", "&"}
	for _, c := range ch {
		notes.Tail = append(notes.Tail, &item{Parent: &notes, Head: c})
	}
	index := 0
	t.Logf("Starting slice has %v items", len(notes.Tail))
	if notes.Tail[index].Locate() != index {
		t.Errorf(`notes.Tail[index].Locate() != index`)
	}
	t.Logf("Check `£`.Locate() = 0: %v", notes.Tail[0].Locate())
	note := notes.Tail[index]
	if len(append(note.Parent.Tail, &item{})) != len(note.Parent.Tail)+1 {
		t.Errorf(`note.Parent.Tail, item{})) != len(note.Parent.Tail) +1`)
	}
	asterix := item{Head: "*", Parent: note.Parent}
	tail := append(note.Parent.Tail, &asterix)
	index += 1
	t.Logf("Appended new item to note.Parent.Tail, is now %v items long", len(tail))
	t.Logf("Post append: %v", tail)
	copy(tail[index+1:], tail[index:])
	t.Logf("Post copy: %v", tail)
	tail[index] = &asterix
	t.Logf("Post assign: %v", tail)

	t.Log(`// Now repeat by func newItem(i *item) {...func (i *item) AddSibling(j *item, index int) {}}`)
	notes2 := item{Home: true, Head: "Homepage"}
	for _, c := range ch {
		notes2.Tail = append(notes2.Tail, &item{Parent: &notes2, Head: c})
	}
	t.Log(`Run newItem(&notes2.Tail[3])`)
	preLen := len(notes2.Tail)
	newItem(notes2.Tail[0])
	if len(notes2.Tail) <= preLen {
		t.Errorf(`notes2.Tail is the same length`)
	}
	t.Logf("notes2 is now %v items long", len(notes2.Tail))
}
