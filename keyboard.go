package main

import (
	"github.com/gdamore/tcell/v2"
	"os"
)

// EditMode is a flag used to restrict keyboard actions. In edit mode, the user
//	can write text to an item. Outside of edit mode, the user can perform all
// 	the tree manipulation actions.
type EditMode bool

func handleEventKey(ev *tcell.EventKey, s tcell.Screen, c *Cursor, m *EditMode) {
	if *m {
		handleEdit(ev, s, c, nil)
	} else {
		handleManipulate(ev, s, c, nil)
	}
}

// handleEdit controls the keyboard actions in EditMode
func handleEdit(ev *tcell.EventKey, s tcell.Screen, c *Cursor, m *EditMode) {
	switch ev.Key() {
	case tcell.KeyEnter:
		// S-Enter creates a new item below, and saves any changes
		if ev.Modifiers() == 1 {
			// TODO: newItem needs fixing
			// It needs to be aware of the current item's tail, and should
			// 	add if there is a tail, that is expanded, it should add to
			// 	the top of that tail.
			// It seems like the pointer to the new item should be returned,
			// 	which could be a way to ensure cursor is directed correctly.
			newItem(c.i)
			c.Down()
			return // to maintain editing state
		}
		// C-S-Enter could create a new item above, and saves and changes
		// 	if in edit mode, however it seems difficult to access this keybind

		// Leave edit mode by saving changes to item head
		*m = false
	case tcell.KeyBackspace2:
	case tcell.KeyUp:
		c.Up()
	case tcell.KeyDown:
		c.Down()
	case tcell.KeyLeft:
	case tcell.KeyRight:
	case tcell.KeyRune:
	case tcell.KeyTab:
	case tcell.KeyBacktab:
	}
}

// handleManipulate controls the keyboard actions when not in EditMode
func handleManipulate(ev *tcell.EventKey, s tcell.Screen, c *Cursor, m *EditMode) {
	switch ev.Key() {
	case tcell.KeyEscape:
		//
		s.Fini()
		os.Exit(0)
	case tcell.KeyEnter:
		// S-Enter creates a new item below cursor
		if ev.Modifiers() == 1 {
			newItem(c.i)
			c.Down()
		}

		// Enter edit mode at cursor
		*m = true
	case tcell.KeyBackspace2:
	case tcell.KeyUp:
		c.Up()
	case tcell.KeyDown:
		c.Down()
	case tcell.KeyLeft:
	case tcell.KeyRight:
	case tcell.KeyRune:
	case tcell.KeyTab:
	case tcell.KeyBacktab:
	}
}
