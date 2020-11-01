package main

import (
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

func main() {

	s := startup()

	c := NewCursor(parseTxt("text/example"))

	for {
		// Pad Item heads to give room for cursor
		if c.i.Text == "" || string(c.i.Text[len(c.i.Text)-1]) != " " {
			c.i.Text += " "
		}

		s.Clear() // unsure how tcell Clear works
		// test further if any chars are deleted or left after deletion

		c.m = make(map[int]*Item)
		TreeMap(c.local, c.m)

		for row := 0; row < len(c.m); row++ {
			c.f[row] = tabx * (len(c.m[row].PathTo(c.local)) - 1)
			emitStr(s, lpad, tity+row, black, strings.Repeat(" ", c.f[row])+c.m[row].Text)

			if c.m[row] == c.i {
				c.y = row
			}
		}

		// Cursor
		if c.editMode {
			s.ShowCursor(lpad+c.f[c.y]+c.x, tity+c.y)
		} else {
			s.HideCursor()
			emitStr(s, lpad+c.f[c.y], tity+c.y, white, c.i.Text)
		}

		drawInfo(s, c, c.local)

		// TODO: update only the elements on the screen that would be changed.
		// may help to use views, or text cells.
		// see implementations in tcell and tview
		s.Show()

		// Events
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			c.lks = ev.Name()
			c.mks = strconv.Itoa(int(ev.Modifiers()))
			handleEventKey(ev, s, c)
		}
	}
}
