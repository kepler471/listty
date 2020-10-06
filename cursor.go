package main

type Cursor struct {
	x      int
	y      int
	i      *item
	buffer string
	m      EditMode
	lks    string
	mks    string
}

// TODO: Support Indent and Unindent with cursor movement
// There may be circumstances where the cursor will need to move after
//	these actions (thinking mainly about when collapsibility is added.

// Down moves cursor down a single row, and selects the correct item.
func (c *Cursor) Down() {
	if len(c.i.Tail) > 0 {
		c.i = c.i.Tail[0]
		c.x = 0
		c.y++
		return
	}
	c.i = c.SearchDown(c.i)
}

// searchDown finds the next item in an ordered tree. If it cannot
// find a suitable next item, it returns the original item.
func (c *Cursor) SearchDown(i *item) *item {
	index := i.Locate()
	if len(i.Parent.Tail) >= index+2 { // Can it move along parent's tail?
		i = i.Parent.Tail[index+1]
		c.x = 0
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
		c.x = 0
		c.y--
		return
	}
	c.i = c.SearchUp(c.i.Parent.Tail[index-1]) // search on preceding sibling
}

func (c *Cursor) SearchUp(i *item) *item {
	if len(i.Tail) == 0 {
		c.x = 0
		c.y--
		return i
	}
	return c.SearchUp(i.Tail[len(i.Tail)-1]) // recurse on the last item in the tail
}
