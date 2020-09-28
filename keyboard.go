package main

import (
	"github.com/gdamore/tcell/v2"
	"os"
)

func handleEventKey(ev *tcell.EventKey, s tcell.Screen, c *Cursor) {
	switch ev.Key() {
	case tcell.KeyEscape:
		s.Fini()
		os.Exit(0)
	case tcell.KeyEnter:
		newItem(c.i)
	case tcell.KeyBackspace2:
	case tcell.KeyUp:
		c.Up()
	case tcell.KeyDown:
		c.Down()
	case tcell.KeyRune:
	case tcell.KeyRight:
	case tcell.KeyLeft:
	case tcell.KeyTab:
	case tcell.KeyBacktab:
	}
}
