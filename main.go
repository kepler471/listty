package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

func main() {
	s := startup()

	var bstr, lks, mks string
	X, Y := s.Size()

	root := parseTxt("text/example")

	c := Cursor{
		x: 5,
		y: 0,
		i: root.Tail[0],
	}

	for {
		s.Clear()
		emitStr(s, 4, 3+c.y, tcell.StyleDefault, ">")
		lines := strings.Split(string(*treeToBytes(root)), nextItem)
		for index, line := range lines {
			emitStr(s, 5, 3+index, black, line)
		}
		drawBox(s, X-42-1, Y-7, X-2, Y-2, white, ' ')
		emitStr(s, X-42, Y-6, white, "Press ESC to exit")
		emitStr(s, X-42, Y-5, white, fmt.Sprintf("Buttons: %s", bstr))
		emitStr(s, X-42, Y-4, white, fmt.Sprintf("Keys: %s", lks))
		emitStr(s, X-42, Y-3, white, fmt.Sprintf("Mods: %s", mks))
		emitStr(s, 5, 0, black, "File: "+root.Head)
		emitStr(s, 5, 1, black, "Item path: "+strings.Join(c.i.Path(), " > "))
		s.Show()

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			lks = ev.Name()
			mks = strconv.Itoa(int(ev.Modifiers()))
			handleEventKey(ev, s, &c)
		}
	}
}
