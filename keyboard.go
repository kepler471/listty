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
		handleManipulate(ev, s, c, local)
	}
}

// handleEdit controls the keyboard actions in EditMode.
// Text editing is handled in a naive manner, writing directly to the item head,
// and moving the cursor as an index of the head string, c.x. c.x will point to
// the position
func handleEdit(ev *tcell.EventKey, s tcell.Screen, c *Cursor, local *item) {
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
		_ = changeMode(c)
	case tcell.KeyEscape:
		c.i.Head = c.buffer
		_ = changeMode(c)
	case tcell.KeyUp:
		if ev.Modifiers() == 1 {
			c.i.MoveUp()
			c.Up()
			return
		}
		c.x = 0
	case tcell.KeyDown:
		if ev.Modifiers() == 1 {
			c.i.MoveDown()
			c.Down()
			return
		}
		c.x = len(c.i.Head) - 1
	case tcell.KeyLeft:
		if c.x > 0 {
			c.x--
		}
	case tcell.KeyRight:
		if c.x < len(c.i.Head)-1 {
			c.x++
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
	}
}

// handleManipulate controls the keyboard actions when not in EditMode
func handleManipulate(ev *tcell.EventKey, s tcell.Screen, c *Cursor, local *item) {
	switch ev.Key() {
	case tcell.KeyEnter:
		// S-Enter creates a new item below cursor
		if ev.Modifiers() == 1 {
			newItem(c.i)
			c.Down()
		}
		// Enter edit mode at cursor
		c.buffer = c.i.Head
		_ = changeMode(c)
	case tcell.KeyDelete:
	case tcell.KeyUp:
		if ev.Modifiers() == 1 {
			c.i.MoveUp()
		}
		c.Up()
	case tcell.KeyDown:
		if ev.Modifiers() == 1 {
			c.i.MoveDown()
		}
		c.Down()
	case tcell.KeyLeft:
		if ev.Modifiers() == 2 {
			// Increase scope to parent of current top level item (dive out)
			if c.i.Parent != nil {
				local = c.i.Parent
			}
			return
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
		// Expand selected item
	case tcell.KeyTab:
		c.i.Indent()
	case tcell.KeyBacktab:
		c.i.Unindent()
	case tcell.KeyCtrlS:
		treeToTxt(local, "save")
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
		}
	}
}
