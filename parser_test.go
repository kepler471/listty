package main

import "testing"

//func TestParseTxt(t *testing.T) {
//	parseTxt()
//}

func TestParseTxt2Flat(t *testing.T) {
	parseTxt2("flat")
}
func TestParseTxt2Indents(t *testing.T) {
	parseTxt2("indents")
}
func TestParseTxt2VariousShort(t *testing.T) {
	parseTxt2("all_short")
}

func TestParseTxt2VariousLong(t *testing.T) {
	parseTxt2("example2")
}
