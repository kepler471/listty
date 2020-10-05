package main

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

func main() {
	s := startup()

	root := parseTxt("text/example")

	c := Cursor{
		x: 0,
		y: 0,
		i: root.Tail[0],
	}

	// Open listty with edit mode off
	m := EditMode(false)

	for {
		s.Clear()
		emitStr(s, 4, 3+c.y, tcell.StyleDefault, ">")
		lines := strings.Split(string(*treeToBytes(root)), nextItem)
		for index, line := range lines {
			emitStr(s, 5, 3+index, black, line)
		}
		emitStr(s, 5, 0, black, "File: "+root.Head)
		emitStr(s, 5, 1, black, "Item path: "+strings.Join(c.i.Path(), " > "))
		s.Show()

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			handleEventKey(ev, s, &c, &m)
		}
	}
}
