package main

import "testing"

func TestParseTxt3Short(t *testing.T) {
	parseTxt("textfiles/all_short")
}
func TestParseTxt3Flat(t *testing.T) {
	parseTxt("textfiles/flat")
}
func TestParseTxt3Long(t *testing.T) {
	parseTxt("textfiles/example2")
}
func Test_TreeToBytes(t *testing.T) {
	ExpectedStr := "- £\n\t- a\n\t\t- d\n\t\t- e\n\t\t- f\n\t- b\n\t- c\n- $\n- %\n"
	Root := Item{
		Root:     true,
		Parent:   nil,
		Text:     "ROOT_NODE",
		Children: nil,
	}

	depth0 := []string{"£", "$", "%"}
	depth1 := []string{"a", "b", "c"}
	depth2 := []string{"d", "e", "f"}

	// create a nested structure
	for j := 0; j < len(depth0); j++ {
		Root.Children = append(Root.Children, &Item{Parent: &Root, Text: depth0[j], Root: false})
		Root.Children[0].Children = append(Root.Children[0].Children, &Item{Parent: Root.Children[0], Text: depth1[j], Root: false})
		Root.Children[0].Children[0].Children = append(Root.Children[0].Children[0].Children, &Item{Parent: Root.Children[0].Children[0], Text: depth2[j], Root: false})
	}

	treeByteArray := *treeToBytes(&Root)
	treeBytesAsString := string(treeByteArray)

	if ExpectedStr != treeBytesAsString {
		t.Errorf("❌ Converted string is incorrect: %s \nexpected: %s", treeBytesAsString, ExpectedStr)
	}
}
