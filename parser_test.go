package main

import "testing"

//func TestParseTxt(t *testing.T) {
//	parseTxt()
//}
//
//func TestParseTxt2Flat(t *testing.T) {
//	parseTxt2("flat")
//}
//func TestParseTxt2Indents(t *testing.T) {
//	parseTxt2("indents")
//}
//func TestParseTxt2VariousShort(t *testing.T) {
//	parseTxt2("all_short")
//}
//
//func TestParseTxt2VariousLong(t *testing.T) {
//	parseTxt2("example2")
//}

func TestParseTxt3Short(t *testing.T) {
	parseTxt3("textfiles/all_short")
}
func TestParseTxt3Flat(t *testing.T) {
	parseTxt3("textfiles/flat")
}
func TestParseTxt3Long(t *testing.T) {
	parseTxt3("textfiles/example2")
}
