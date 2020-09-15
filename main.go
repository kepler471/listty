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

	notes := item{Home: true, Head: "Homepage"}
	ch := []string{"£", "$", "%", "^", "&", "*"}
	for _, c := range ch {
		notes.Tail = append(notes.Tail, item{Parent: &notes, Head: c})
	}

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
	drawNotes(s, &notes, cy)
	cx += start

	for {
		// Block empty tail from existing
		if len(notes.Tail) == 0 {
			notes.Tail = append(notes.Tail, item{Parent: &notes})
			continue
		}
		drawNotes(s, &notes, cy)
		drawBox(s, X-42-1, Y-7-1, X-2, Y-2, white, ' ')
		emitStr(s, X-42, Y-7, white, "Press ESC twice to exit")
		emitStr(s, X-42, Y-6, white, fmt.Sprintf(posfmt, mx, my))
		emitStr(s, X-42, Y-5, white, fmt.Sprintf(btnfmt, bstr))
		emitStr(s, X-42, Y-4, white, fmt.Sprintf(keyfmt, lks))
		emitStr(s, X-42, Y-3, white, fmt.Sprintf(modfmt, mks))
		emitStr(s, 5, 1, tcell.StyleDefault, notes.Head)
		s.Show()

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
				newItem(&notes.Tail[cy])
				cy++
				cy %= len(notes.Tail)
			case tcell.KeyBackspace2:
				if ev.Modifiers() == 4 {
					notes.Tail[cy].Remove()
					if cy != 0 {
						cy--
					}
					continue
				}
				if len(notes.Tail[cy].Head) == 0 {
					continue
				}
				notes.Tail[cy].Head = notes.Tail[cy].Head[:len(notes.Tail[cy].Head)-1]
			case tcell.KeyUp:
				if ev.Modifiers() == 4 {
					if cy == 0 {
						continue
					}
					notes.Tail[cy].MoveUp()
				}
				if cy == 0 {
					cy = len(notes.Tail) - 1
					continue
				}
				cy--
			case tcell.KeyDown:
				if ev.Modifiers() == 4 {
					if cy == len(notes.Tail)-1 {
						continue
					}
					notes.Tail[cy].MoveDown()
				}
				cy++
				cy %= len(notes.Tail)
			case tcell.KeyRune:
				if ev.Modifiers() >= 1 {
					continue
				}
				input := string(ev.Rune())
				notes.Tail[cy].Head += input
			}
		}
	}
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

func drawNotes(s tcell.Screen, n *item, cursor int) {
	//w, h := s.Size()
	s.Clear()
	start := 5
	emitStr(s, start-1, start+cursor, tcell.StyleDefault, ">")
	emitStr(s, start, start-3, tcell.StyleDefault, "Data: "+n.StringChildren())
	emitStr(s, start, start-1, tcell.StyleDefault, "Cursor Y value: "+strconv.Itoa(cursor))
	for index := range n.Tail {
		emitStr(s, start, start+index, tcell.StyleDefault, n.Tail[index].Head)
	}
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
