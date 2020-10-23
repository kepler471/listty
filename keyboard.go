package main

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"strings"
)

// TODO: change inner conditional branches to switch cases

func handleEventKey(ev *tcell.EventKey, s tcell.Screen, c *Cursor) {
	if c.editMode {
		handleEdit(ev, s, c)
	} else {
		handleSelect(ev, s, c)
	}
}

// handleEdit controls the keyboard actions in EditMode.
// Text editing is handled in a naive manner, writing directly to the Item head,
// and moving the cursor as an index of the head string, c.x. c.x will point to
// the position
func handleEdit(ev *tcell.EventKey, s tcell.Screen, c *Cursor) {
	switch ev.Key() {
	// TODO: add some comments describing key actions
	case tcell.KeyEnter:
		c.ClearBuffer()

		if ev.Modifiers() == 1 {
			switch {
			case c.x == 0:
				c.i = c.i.AddSibling(&Item{Parent: c.i.Parent, Text: " "}, c.i.Locate())
			case len(c.i.Children) > 0:
				c.i = c.i.Children[0].AddSibling(&Item{Parent: c.i, Text: " "}, 0)
			case len(c.i.Children) == 0:
				c.i = c.i.AddSibling(&Item{Parent: c.i.Parent, Text: " "}, c.i.Locate()+1)
			}
			c.ResetX()
			return // to maintain editing state

		}
		_ = changeMode(c)

	case tcell.KeyEscape:
		c.ResetX()
		c.UnsetBuffer()
		_ = changeMode(c)

	case tcell.KeyUp:
		if ev.Modifiers() == 1 {
			c.i.MoveUp()
			return
		}

		c.ResetX()

		// If c.x = 0, then pressing up again could take you out of edit mode?
		// Maybe on a double press? Same with down? What if they
		// also take you to the next Item? This may start to blur
		// the line between edit and selection mode, which could
		// lead to removing the separate modes.

	case tcell.KeyDown:
		if ev.Modifiers() == 1 {
			c.i.MoveDown()
			return
		}

		c.x = len(c.i.Text) - 1

	case tcell.KeyLeft:
		if c.x > 0 {
			c.x--
		}

		if ev.Modifiers() == 1 {
			c.i.Unindent()
		}

	case tcell.KeyRight:
		if c.x < len(c.i.Text)-1 {
			c.x++
		}

		if ev.Modifiers() == 1 {
			c.i.Indent()
		}

	case tcell.KeyRune:
		c.i.Text = c.i.Text[:c.x] + string(ev.Rune()) + c.i.Text[c.x:]
		c.x++

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if c.x != 0 {
			c.i.Text = c.i.Text[:c.x-1] + c.i.Text[c.x:]
			c.x--
		}

	case tcell.KeyTab:
		c.i.Indent()

	case tcell.KeyBacktab:
		c.i.Unindent()

	case tcell.KeyCtrlS:
		treeToTxt(c.local, "text/example")

	case tcell.KeyCtrlQ:
		s.Fini()
		os.Exit(0)
	}
}

// handleSelect controls the keyboard actions when not in EditMode
func handleSelect(ev *tcell.EventKey, s tcell.Screen, c *Cursor) {
	switch ev.Key() {
	case tcell.KeyEnter:
		if ev.Modifiers() == 1 {
			if len(c.i.Children) > 0 {
				c.i = c.i.Children[0].AddSibling(&Item{Parent: c.i, Text: " "}, 0)
			} else {
				c.i = c.i.AddSibling(&Item{Parent: c.i.Parent, Text: " "}, c.i.Locate()+1)
			}
		}

		c.ResetX()
		c.SetBuffer()
		_ = changeMode(c)

	case tcell.KeyDelete:

	case tcell.KeyUp:
		if ev.Modifiers() == 1 {
			c.i.MoveUp()
			return
		}
		c.Up()

	case tcell.KeyDown:
		if ev.Modifiers() == 1 {
			c.i.MoveDown()
			return
		}
		c.Down()

	case tcell.KeyLeft:
		if ev.Modifiers() == 2 {
			// Increase scope to parent of current top level Item (dive out)
			if c.i.Parent != nil {
				c.local = c.i.Parent
			}
			return
		}

		if ev.Modifiers() == 1 {
			c.i.Unindent()
		}

		// Collapse selected item

	case tcell.KeyRight:
		if ev.Modifiers() == 2 {
			// Set selected Item as top level Item (dive in)
			// For now, limit to non-leaf items, as unsure how app handles for empty tails
			if len(c.i.Children) != 0 {
				c.local = c.i
			}
			return
		}

		if ev.Modifiers() == 1 {
			c.i.Indent()
		}

		// Expand selected item

	case tcell.KeyTab:
		c.i.Indent()

	case tcell.KeyBacktab:
		c.i.Unindent()

	case tcell.KeyCtrlS:
		treeToTxt(c.local, "text/example")

	case tcell.KeyCtrlQ:
		s.Fini()
		os.Exit(0)

	case tcell.KeyRune:
		switch strings.ToLower(string(ev.Rune())) {
		case "d":
			// Duplicate selected Item

		case ",":
			c.i.Unindent()

		case ".":
			c.i.Indent()

		// Would it be good to enter editing mode just by typing?
		// Try this out. If so, we would want the first letter
		// pressed to actually be entered.
		default:
			_ = changeMode(c)
		}
	}
}
