package main

import (
	"testing"
)

func Test_TreeToBytes(t *testing.T) {
	ExpectedStr := "- £\n\t- a\n\t\t- d\n\t\t- e\n\t\t- f\n\t- b\n\t- c\n- $\n- %\n"
	Root := item{
		Home:   true,
		Parent: nil,
		Head:   "ROOT_NODE",
		Tail:   nil,
	}

	depth0 := []string{"£", "$", "%"}
	depth1 := []string{"a", "b", "c"}
	depth2 := []string{"d", "e", "f"}

	// create a nested structure
	for j := 0; j < len(depth0); j++ {
		Root.Tail = append(Root.Tail, &item{Parent: &Root, Head: depth0[j], Home: false})
		Root.Tail[0].Tail = append(Root.Tail[0].Tail, &item{Parent: Root.Tail[0], Head: depth1[j], Home: false})
		Root.Tail[0].Tail[0].Tail = append(Root.Tail[0].Tail[0].Tail, &item{Parent: Root.Tail[0].Tail[0], Head: depth2[j], Home: false})
	}

	treeByteArray := *treeToBytes(&Root)
	treeBytesAsString := string(treeByteArray)

	if ExpectedStr != treeBytesAsString {
		t.Errorf("❌ Converted string is incorrect: %s \nexpected: %s", treeBytesAsString, ExpectedStr)
	}
}
