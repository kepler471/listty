// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"os"
	"strconv"
	"strings"
)

func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	root := item{Home: true, Head: "Homepage"}

	ch := []string{"£", "$", "%", "^", "&", "*"}
	for _, c := range ch {
		root.Tail = append(root.Tail, item{Parent: &root, Head: c})
	}

	root.Tail[0].Tail = append(root.Tail[0].Tail, item{Parent: &root.Tail[0], Head: "hello", depth: 1})

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	modfmt := "Mods: %s"
	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	mx, my := -1, -1
	var bstr, lks, mks string
	X, Y := s.Size()
	//s.EnableMouse()
	ecnt := 0
	start := 5
	var cx, cy int
	cx += start
	c := Cursor{
		x: cx,
		y: cy,
		i: &root,
	}
	depth:=0
	currentItem := &root
	drawInfo(s, currentItem, c.y, depth)
	row := 0
	currentItem.Plot(s, &row)


	for {
		currentItem := getCurrentItem(&root, &stack)

		// Block empty tail from existing
		if len(currentItem.Tail) == 0 {
			currentItem.Tail = append(currentItem.Tail, item{
				Parent: currentItem,
				Head:   "Parent: " + currentItem.Head + ", " + " Depth: " + strconv.Itoa(depth) + " Row: " + strconv.Itoa(cy),
			})
		}
		drawInfo(s, currentItem, c.y, depth)
		row = 0
		currentItem.Plot(s, &row)
		drawBox(s, X-42-1, Y-7-1, X-2, Y-2, white, ' ')
		emitStr(s, X-42, Y-7, white, "Press ESC twice to exit")
		emitStr(s, X-42, Y-6, white, fmt.Sprintf(posfmt, mx, my))
		emitStr(s, X-42, Y-5, white, fmt.Sprintf(btnfmt, bstr))
		emitStr(s, X-42, Y-4, white, fmt.Sprintf(keyfmt, lks))
		emitStr(s, X-42, Y-3, white, fmt.Sprintf(modfmt, mks))
		emitStr(s, 5, 1, tcell.StyleDefault, currentItem.Head)
		s.Show()

		// TODO: Replace currentItem with cursor.item
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			lks = ev.Name()
			mks = strconv.Itoa(int(ev.Modifiers()))
			switch ev.Key() {
			case tcell.KeyEscape:
				ecnt++
				if ecnt > 1 {
					s.Fini()
					os.Exit(0)
				}
			case tcell.KeyEnter:
				newItem(&currentItem.Tail[c.y])
				c.y++
				c.y %= row // len(currentItem.Tail)
			case tcell.KeyBackspace2:
				if ev.Modifiers() == 4 {
					currentItem.Tail[c.y].Remove()
					if c.y != 0 {
						c.y--
					}
					continue
				}
				if len(currentItem.Tail[c.y].Head) == 0 {
					continue
				}
				currentItem.Tail[c.y].Head = currentItem.Tail[c.y].Head[:len(currentItem.Tail[c.y].Head)-1]
			case tcell.KeyUp:
				if ev.Modifiers() == 4 {
					if c.y == 0 {
						continue
					}
					currentItem.Tail[c.y].MoveUp()
				}
				if c.y == 0 {
					c.y = len(currentItem.Tail) - 1
					continue
				}
				c.y--
			case tcell.KeyDown:
				if ev.Modifiers() == 4 {
					if c.y == len(currentItem.Tail)-1 {
						continue
					}
					currentItem.Tail[c.y].MoveDown()
				}
				c.y++
				c.y %= row // len(currentItem.Tail)
			case tcell.KeyRune:
				if ev.Modifiers() >= 1 {
					continue
				}
				currentItem.Tail[c.y].Head += string(ev.Rune())
			case tcell.KeyRight:
				if currentItem == nil || len(currentItem.Tail) == 0 {
					currentItem.Tail[c.y].AddChild(&item{Head: "newborn baby", Parent: &currentItem.Tail[c.y]})
					c.y = 0

					continue
				}
				depth++
				c.y = 0

			case tcell.KeyLeft:
				if !currentItem.Home {
					stack.Pop()
					depth--
					cy = stack.GetRow(depth)
				}
				stack.Pop()
				depth--
			case tcell.KeyTab:
				currentItem.Tail[c.y].Head = "\t" + currentItem.Tail[c.y].Head
				currentItem.Tail[c.y].Indent()
				c.y--
				//case tcell.KeyBacktab:
				//notes.Tail[c.y].Head = notes.Tail[c.y].Head[1:]
				//notes.Tail[c.y].Unindent()
			}
		}
	}
}

func (i *item) Plot(s tcell.Screen, row *int) {
	if len(i.Tail) > 0 {
		for index, t := range i.Tail {
			emitStr(s, 5, 5+*row, tcell.StyleDefault, fmt.Sprintf(
				"Tail Id: %v, Depth: %v, Head: %v", strconv.Itoa(index), t.depth, strings.Repeat("\t", t.depth)+t.Head),
			)
			*row++
			t.Plot(s, row)
		}
	} else {
		return
	}
}

type Cursor struct {
	x int
	y int
	i *item
}

// Down moves cursor down a single row, and selects the correct item.
func (c *Cursor) Down() {
	if len(c.i.Tail) > 0 {
		c.i = &c.i.Tail[0]
		c.y++
		return
	}
	c.i = c.searchDown(c.i)
}

// searchDown finds the next item in an ordered tree. If it cannot
// find a suitable next item, it returns the original item.
func (c *Cursor) searchDown(i *item) *item {
	index := i.Locate()
	if len(i.Parent.Tail) >= index+2 { // Can it move along parent's tail?
		i = &i.Parent.Tail[index+1]
		c.y++
		return i
	}
	if i.Parent.Parent == nil { // Protection searching above root
		return c.i
	}
	return c.searchDown(i.Parent)
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawFromCursor() {}

func drawInfo(s tcell.Screen, currentItem *item, cursor int, depth int) {
	//w, h := s.Size()
	s.Clear()
	start := 5
	path := currentItem.Path()
	emitStr(s, start-1, start+cursor, tcell.StyleDefault, ">")
	emitStr(s, start, start-3, tcell.StyleDefault, "Data: "+currentItem.StringChildren())
	emitStr(s, start, start-1, tcell.StyleDefault, "Cursor Y value: "+strconv.Itoa(cursor))
	emitStr(s, start, start+20, tcell.StyleDefault, "Item path: "+strings.Join(path, " > "))

	//for index := range currentItem.Tail {
	//	emitStr(s, start, start+index, tcell.StyleDefault, currentItem.Tail[index].Head)
	//}

	s.Show()
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}
