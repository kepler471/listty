package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"log"
	"strings"
)

var (
	white = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorRed)

	black = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)
)

type Cursor struct {
	x      int
	y      int
	i      *item
	buffer string
	m      EditMode
}

// TODO: Support Indent and Unindent with cursor movement
// There may be circumstances where the cursor will need to move after
//	these actions (thinking mainly about when collapsibility is added.

// Down moves cursor down a single row, and selects the correct item.
func (c *Cursor) Down() {
	if len(c.i.Tail) > 0 {
		c.i = c.i.Tail[0]
		c.y++
		return
	}
	c.i = c.searchDown(c.i)
}

// searchDown finds the next item in an ordered tree. If it cannot
// find a suitable next item, it returns the original item.
func (c *Cursor) searchDown(i *item) *item {
	index := i.Locate()
	if len(i.Parent.Tail) >= index+2 { // Can it move along parent's tail?
		i = i.Parent.Tail[index+1]
		c.y++
		return i
	}
	if i.Parent.Parent == nil { // Protection searching above root
		return c.i
	}
	return c.searchDown(i.Parent)
}

func (c *Cursor) Up() {
	index := c.i.Locate()
	if index == 0 {
		if c.i.Parent.Parent == nil { // Protection searching above root
			return
		}
		c.i = c.i.Parent
		c.y--
		return
	}
	c.i = c.searchUp(c.i.Parent.Tail[index-1]) // search on preceding sibling
}

func (c *Cursor) searchUp(i *item) *item {
	if len(i.Tail) == 0 {
		c.y--
		return i
	}
	return c.searchUp(i.Tail[len(i.Tail)-1]) // recurse on the last item in the tail
}

func startup() tcell.Screen {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal(e)
	}
	if e := s.Init(); e != nil {
		log.Fatal(e)
	}

	s.SetStyle(tcell.StyleDefault)
	return s
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
		s.Show()
	}
}

func (i *item) Plot(s tcell.Screen, m map[*item]int, style tcell.Style) {
	if len(i.Tail) > 0 {
		for _, t := range i.Tail {
			depth := len(t.Path()) - 1
			m[t]++
			emitStr(s, 5, 5+m[t], style, strings.Repeat("\t", depth)+t.Head)
			t.Plot(s, m, style)
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}
