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
