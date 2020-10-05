package main

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"strings"
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

// handleEdit controls the keyboard actions in EditMode.
// Text editing is handled in a naive manner, writing directly to the item head,
// and moving the cursor as an index of the head string, c.x. c.x will point to
// the position
func handleEdit(ev *tcell.EventKey, s tcell.Screen, c *Cursor, m *EditMode) {
	switch ev.Key() {
	case tcell.KeyEnter:
		c.buffer = ""

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
		// 	if in edit mode, however it seems difficult to access this keybind.

		*m = false
	case tcell.KeyEscape:
		c.i.Head = c.buffer
		*m = false
	case tcell.KeyUp:
		if ev.Modifiers() == 2 {
			c.i.MoveUp()
			c.Up()
			return
		}
		c.x = 0
	case tcell.KeyDown:
		if ev.Modifiers() == 2 {
			c.i.MoveDown()
			c.Down()
			return
		}
		c.x = len(c.i.Head)
	case tcell.KeyLeft:
		if c.x > 0 {
			c.x--
		}
	case tcell.KeyRight:
		if c.x < len(c.i.Head) {
			c.x++
		}
	case tcell.KeyRune:
		c.i.Head = c.i.Head[:c.x] + string(ev.Rune()) + c.i.Head[c.x:]
		c.x++
	case tcell.KeyBackspace2:
		c.i.Head = c.i.Head[:c.x-1] + c.i.Head[c.x:]
		c.x--
	case tcell.KeyTab:
		c.i.Indent()
	case tcell.KeyBacktab:
		c.i.Unindent()
	}
}

// handleManipulate controls the keyboard actions when not in EditMode
func handleManipulate(ev *tcell.EventKey, s tcell.Screen, c *Cursor, m *EditMode) {
	switch ev.Key() {
	case tcell.KeyEnter:
		// S-Enter creates a new item below cursor
		if ev.Modifiers() == 1 {
			newItem(c.i)
			c.Down()
		}
		// Enter edit mode at cursor
		c.buffer = c.i.Head
		*m = true
	case tcell.KeyBackspace2:
		// Delete an item, maybe with a double press, or a delete mode?
	case tcell.KeyUp:
		// S-up moves item up a tail
		if ev.Modifiers() == 2 {
			c.i.MoveUp()
		}
		c.Up()
	case tcell.KeyDown:
		// S-down moves item down a tail
		if ev.Modifiers() == 2 {
			c.i.MoveDown()
		}
		c.Down()
	case tcell.KeyLeft:
		if ev.Modifiers() == 2 {
			// View parent of local root (dive out)
			// return
		}
		// Collapse
	case tcell.KeyRight:
		if ev.Modifiers() == 2 {
			// Set current item as local root (dive in)
			// return
		}
		// Expand
	case tcell.KeyTab:
		c.i.Indent()
	case tcell.KeyBacktab:
		c.i.Unindent()
	case tcell.KeyRune:
		switch strings.ToLower(string(ev.Rune())) {
		case "d":
			// Duplicate selected item
		case "q":
			s.Fini()
			os.Exit(0)
		}
	}
}
