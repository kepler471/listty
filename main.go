package main

import (
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

func main() {
	s := startup()

	root := parseTxt("text/example")
	// local will be the top level item on view in the application.
	local := root

	c := Cursor{
		x:      0,
		y:      0,
		i:      root.Tail[0],
		buffer: "",
		m:      EditMode(false), // Open listty with edit mode off
	}

	for {
		// Pad item heads to give room for cursor
		if c.i.Head == "" || string(c.i.Head[len(c.i.Head)-1]) != " " {
			c.i.Head += " "
		}

		s.Clear() // unsure how tcell Clear works
		// test further if any chars are deleted or left after deletion

		m := make(map[int]*item)
		f := make(map[int]int) // track indentations for cursor
		TreeMap(local, m)

		for row := 0; row < len(m); row++ {
			f[row] = tabx * (len(m[row].PathTo(local)) - 1)
			emitStr(s, lpad, tity+row, black, strings.Repeat(" ", f[row])+m[row].Head)

			// As long as the cursor know which item to look at, it is found easily.
			// This may lead to not having to use any movement methods on cursor.
			if m[row] == c.i {
				c.y = row
			}
		}

		// Cursor
		s.ShowCursor(5+f[c.y]+c.x, 3+c.y)

		drawInfo(s, &c, local)

		s.Show()

		// Events
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			c.lks = ev.Name()
			c.mks = strconv.Itoa(int(ev.Modifiers()))
			handleEventKey(ev, s, &c, local)
		}
	}
}
