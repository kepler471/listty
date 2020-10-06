package main

import (
	"fmt"
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
		m:      EditMode(false),
	}

	// Open listty with edit mode off
	x, _ := s.Size()

	crsfmt := "Cursor X: %v Y: %v a: %v"
	keyfmt := "Keys: %s"
	modfmt := "Mods: %s"
	var lks, mks string

	// TODO: init c.x at end of line for first item
	// TODO: manage c.x movement in c.Down and c.Up
	for {
		// Pad item heads to give room for cursor
		if c.i.Head[len(c.i.Head)-1] != ' ' {
			c.i.Head += " "
		}
		//if c.i.Head == "" {
		//	c.i.Head += " "
		//} else {
		//	c.i.Head += " "
		//}
		s.Clear()
		emitStr(s, 4, 3+c.y, tcell.StyleDefault, ">")
		m := make(map[int]*item)
		f := make(map[int]int) // track indentations for cursor
		TreeMap(local, m)
		for row := 0; row < len(m); row++ {
			f[row] = 4 * (len(m[row].PathTo(local)) - 1)
			emitStr(s, 5, 3+row, black, strings.Repeat(" ", f[row])+m[row].Head)
		}
		emitStr(s, 5+f[c.y]+c.x, 3+c.y, white, string(c.i.Head[c.x]))

		// Info bar
		emitStr(s, 5, 0, black, "File: "+local.Head)
		emitStr(s, 5, 1, black, "Item path: "+strings.Join(c.i.Path(), " > "))

		// Info box
		drawBox(s, x-41, 0, x-1, 6, white, ' ')
		emitStr(s, x-40, 1, white, "Press Q to exit")
		emitStr(s, x-40, 3, white, fmt.Sprintf(crsfmt, c.x, c.y, string(c.i.Head[c.x])))
		emitStr(s, x-40, 4, white, fmt.Sprintf(keyfmt, lks))
		emitStr(s, x-40, 5, white, fmt.Sprintf(modfmt, mks))

		s.Show()

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			lks = ev.Name()
			mks = strconv.Itoa(int(ev.Modifiers()))
			handleEventKey(ev, s, &c, local)
		}
	}
}
