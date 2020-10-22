package main

type Cursor struct {
	x      int
	y      int
	i      *item
	buffer string
	m      EditMode
	lks    string
	mks    string
	// TODO: compose item with cursor?
	// 	makes changes globally and test.
	//*item
}

func NewCursor(i *item) Cursor {
	return Cursor{
		i: i.Children[0],
		m: EditMode(false),
	}
}

// TODO: Support Indent and Unindent with cursor movement
// There may be circumstances where the cursor will need to move after
//	these actions (thinking mainly about when collapsibility is added.

// Down moves cursor down a single row, and selects the correct item.
func (c *Cursor) Down() {
	if len(c.i.Children) > 0 {
		c.i = c.i.Children[0]
		c.ResetX()
		c.y++
		return
	}
	c.i = c.SearchDown(c.i)
}

// searchDown finds the next item in an ordered tree. If it cannot
// find a suitable next item, it returns the original item.
func (c *Cursor) SearchDown(i *item) *item {
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

func (c *Cursor) SearchUp(i *item) *item {
	if len(i.Children) == 0 {
		c.ResetX()
		c.y--
		return i
	}
	return c.SearchUp(i.Children[len(i.Children)-1]) // recurse on the last item in the tail
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

// TODO: add c.x inc/dec methods
