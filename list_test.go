package main

import (
	"testing"
)

func TestItem_InsertAt(t *testing.T) {
	notes := Item{Root: true, Text: "Homepage"}
	ch := []string{"£", "$", "%", "^", "&"}
	for _, c := range ch {
		notes.Children = append(notes.Children, &Item{Parent: &notes, Text: c})
	}
	index := 0
	t.Logf("Starting slice has %v items", len(notes.Children))
	if notes.Children[index].Locate() != index {
		t.Errorf(`notes.Tail[index].Locate() != index`)
	}
	t.Logf("Check `£`.Locate() = 0: %v", notes.Children[0].Locate())
	note := notes.Children[index]
	if len(append(note.Parent.Children, &Item{})) != len(note.Parent.Children)+1 {
		t.Errorf(`note.Parent.Tail, Item{})) != len(note.Parent.Tail) +1`)
	}
	asterix := Item{Text: "*", Parent: note.Parent}
	tail := append(note.Parent.Children, &asterix)
	index += 1
	t.Logf("Appended new Item to note.Parent.Tail, is now %v items long", len(tail))
	t.Logf("Post append: %v", tail)
	copy(tail[index+1:], tail[index:])
	t.Logf("Post copy: %v", tail)
	tail[index] = &asterix
	t.Logf("Post assign: %v", tail)

	t.Log(`// Now repeat by func newItem(i *Item) {...func (i *Item) AddSibling(j *Item, index int) {}}`)
	notes2 := Item{Root: true, Text: "Homepage"}
	for _, c := range ch {
		notes2.Children = append(notes2.Children, &Item{Parent: &notes2, Text: c})
	}
	t.Log(`Run newItem(&notes2.Tail[3])`)
	preLen := len(notes2.Children)
	newItem(notes2.Children[0])
	if len(notes2.Children) <= preLen {
		t.Errorf(`notes2.Tail is the same length`)
	}
	t.Logf("notes2 is now %v items long", len(notes2.Children))
}
