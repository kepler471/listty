package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"log"
	"strings"
)

// Draw formatting
const (
	exit = 1 + iota
	_
	crsr
	keys
	mods
	lpad   = 5
	box0   = 0
	boxy   = 6
	boxl   = 41
	boxr   = 1
	tabx   = 4
	tity   = 3
	fily   = 0
	pthy   = 1
	crsfmt = "Cursor X: %v Y: %v a: %v"
	keyfmt = "Keys: %s"
	modfmt = "Mods: %s"
)

var (
	white = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorRed)

	black = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)
)

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

func drawInfo(s tcell.Screen, c *Cursor, local *item) {
	width, _ := s.Size()
	// Info bar
	emitStr(s, lpad, fily, black, "File: "+local.Head)
	emitStr(s, lpad, pthy, black, "Item path: "+strings.Join(c.i.Path(), " > "))

	// Info box
	drawBox(s, width-boxl, box0, width-boxr, boxy, white, ' ')
	emitStr(s, width-(boxl-boxr), exit, white, "Press Ctrl-Q to exit")
	emitStr(s, width-(boxl-boxr), crsr, white, fmt.Sprintf(crsfmt, c.x, c.y, string(c.i.Head[c.x])))
	emitStr(s, width-(boxl-boxr), keys, white, fmt.Sprintf(keyfmt, c.lks))
	emitStr(s, width-(boxl-boxr), mods, white, fmt.Sprintf(modfmt, c.mks))

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
