package main

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"strconv"
	"strings"
)

// EditMode is a flag used to restrict keyboard actions. In edit mode, the user
//	can write text to an item. Outside of edit mode, the user can perform all
// 	the tree manipulation actions.
type EditMode bool

// ModeError is used to catch an error during mode switch
type ModeError struct {
	EditMode EditMode
}

func (m *ModeError) Error() string {
	return "EditMode was " + strconv.FormatBool(bool(m.EditMode)) + " but was not flipped to " + strconv.FormatBool(bool(!m.EditMode))
}

func changeMode(c *Cursor) error {
	clone := c
	c.m = !c.m
	if clone != c {
		return &ModeError{EditMode: c.m}
	}
	return nil
}

func handleEventKey(ev *tcell.EventKey, s tcell.Screen, c *Cursor, local *item) {
	if c.m {
		handleEdit(ev, s, c, local)
	} else {
		handleSelect(ev, s, c, local)
	}
}

// handleEdit controls the keyboard actions in EditMode.
// Text editing is handled in a naive manner, writing directly to the item head,
// and moving the cursor as an index of the head string, c.x. c.x will point to
// the position
func handleEdit(ev *tcell.EventKey, s tcell.Screen, c *Cursor, local *item) {
	switch ev.Key() {
	// TODO: add some comments describing key actions
	case tcell.KeyEnter:
		c.buffer = ""

		if ev.Modifiers() == 1 {
			switch {
			case c.x == 0:
				c.i = c.i.AddSibling(&item{Parent: c.i.Parent, Head: " "}, c.i.Locate())
			case len(c.i.Tail) > 0:
				c.i = c.i.Tail[0].AddSibling(&item{Parent: c.i, Head: " "}, 0)
			case len(c.i.Tail) == 0:
				c.i = c.i.AddSibling(&item{Parent: c.i.Parent, Head: " "}, c.i.Locate()+1)
			}
			c.x = 0
			return // to maintain editing state

		}
		_ = changeMode(c)

	case tcell.KeyEscape:
		c.x = 0
		c.i.Head = c.buffer
		c.buffer = ""
		_ = changeMode(c)

	case tcell.KeyUp:
		if ev.Modifiers() == 1 {
			c.i.MoveUp()
			return
		}

		c.x = 0

		// If c.x = 0, then pressing up again could take you out of edit mode?
		// Maybe on a double press? Same with down? What if they
		// also take you to the next item? This may start to blur
		// the line between edit and selection mode, which could
		// lead to removing the separate modes.

	case tcell.KeyDown:
		if ev.Modifiers() == 1 {
			c.i.MoveDown()
			return
		}

		c.x = len(c.i.Head) - 1

	case tcell.KeyLeft:
		if c.x > 0 {
			c.x--
		}

		if ev.Modifiers() == 1 {
			c.i.Unindent()
		}

	case tcell.KeyRight:
		if c.x < len(c.i.Head)-1 {
			c.x++
		}

		if ev.Modifiers() == 1 {
			c.i.Indent()
		}

	case tcell.KeyRune:
		c.i.Head = c.i.Head[:c.x] + string(ev.Rune()) + c.i.Head[c.x:]
		c.x++

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if c.x != 0 {
			c.i.Head = c.i.Head[:c.x-1] + c.i.Head[c.x:]
			c.x--
		}

	case tcell.KeyTab:
		c.i.Indent()

	case tcell.KeyBacktab:
		c.i.Unindent()

	case tcell.KeyCtrlS:
		treeToTxt(local, "text/example")

	case tcell.KeyCtrlQ:
		s.Fini()
		os.Exit(0)
	}
}

// handleSelect controls the keyboard actions when not in EditMode
func handleSelect(ev *tcell.EventKey, s tcell.Screen, c *Cursor, local *item) {
	switch ev.Key() {
	case tcell.KeyEnter:
		if ev.Modifiers() == 1 {
			if len(c.i.Tail) > 0 {
				c.i = c.i.Tail[0].AddSibling(&item{Parent: c.i, Head: " "}, 0)
			} else {
				c.i = c.i.AddSibling(&item{Parent: c.i.Parent, Head: " "}, c.i.Locate()+1)
			}
		}

		c.x = 0
		c.buffer = c.i.Head
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

	//Shift - Left/Right used for indentation? Try it out
	case tcell.KeyLeft:
		if ev.Modifiers() == 2 {
			// Increase scope to parent of current top level item (dive out)
			if c.i.Parent != nil {
				local = c.i.Parent
			}
			return
		}

		if ev.Modifiers() == 1 {
			c.i.Unindent()
		}

		// Collapse selected item

	case tcell.KeyRight:
		if ev.Modifiers() == 2 {
			// Set selected item as top level item (dive in)
			// For now, limit to non-leaf items, as unsure how app handles for empty tails
			if len(c.i.Tail) != 0 {
				local = c.i
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
		treeToTxt(local, "text/example")

	case tcell.KeyCtrlQ:
		s.Fini()
		os.Exit(0)

	case tcell.KeyRune:
		switch strings.ToLower(string(ev.Rune())) {
		case "d":
			// Duplicate selected item

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
