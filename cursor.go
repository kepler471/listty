package main

// Cursor maintains the state of the view of the Item tree structure
type Cursor struct {
	x        int
	y        int
	i        *Item
	root     *Item
	local    *Item
	buffer   string
	editMode bool
	lks      string
	mks      string
	m        map[int]*Item
	f        map[int]int
	// TODO: compose Item with cursor?
	// 	makes changes globally and test.
	//*Item
}

func NewCursor(i *Item) *Cursor {
	m := make(map[int]*Item)
	f := make(map[int]int)
	return &Cursor{
		i:     i.Children[0],
		root:  i,
		local: i,
		m:     m,
		f:     f,
	}
}

// TODO: Support Indent and Unindent with cursor movement
// There may be circumstances where the cursor will need to move after
//	these actions (thinking mainly about when collapsibility is added.

// Down moves cursor down a single row, and selects the correct Item.
func (c *Cursor) Down() {
	if len(c.i.Children) > 0 {
		c.i = c.i.Children[0]
		c.ResetX()
		c.y++
		return
	}
	c.i = c.SearchDown(c.i)
}

// searchDown finds the next Item in an ordered tree. If it cannot
// find a suitable next Item, it returns the original Item.
func (c *Cursor) SearchDown(i *Item) *Item {
	index := i.Locate()
	if len(i.Parent.Children) >= index+2 { // Can it move along parent's tail?
		i = i.Parent.Children[index+1]
		c.ResetX()
		c.y++
		return i
	}
	if i.Parent.Parent == nil { // Protection searching above root
		return c.i
	}
	return c.SearchDown(i.Parent)
}

func (c *Cursor) Up() {
	index := c.i.Locate()
	if index == 0 {
		if c.i.Parent.Parent == nil { // Protection searching above root
			return
		}
		c.i = c.i.Parent
		c.ResetX()
		c.y--
		return
	}
	c.i = c.SearchUp(c.i.Parent.Children[index-1]) // search on preceding sibling
}

func (c *Cursor) SearchUp(i *Item) *Item {
	if len(i.Children) == 0 {
		c.ResetX()
		c.y--
		return i
	}
	return c.SearchUp(i.Children[len(i.Children)-1]) // recurse on the last Item in the tail
}

func (c *Cursor) ResetX() {
	c.x = 0
}

func (c *Cursor) ClearBuffer() {
	c.buffer = ""
}

func (c *Cursor) SetBuffer() {
	c.buffer = c.i.Text
}

func (c *Cursor) UnsetBuffer() {
	c.i.Text = c.buffer
	c.ClearBuffer()
}

func changeMode(c *Cursor) error {
	c.editMode = !c.editMode
	return nil
}

// TODO: add c.x inc/dec methods
