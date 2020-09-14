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
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
	"os"
	"strconv"
	//"github.com/rivo/tview"
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

	//defStyle := tcell.StyleDefault.
	//	Background(tcell.ColorBlack).
	//	Foreground(tcell.ColorWhite)
	//s.SetStyle(defStyle)

	notes := item{Home: true, Head: "Homepage"}
	ch := []string{"£", "$", "%", "^", "&", "*"}
	for _, c := range ch {
		notes.Tail = append(notes.Tail, item{Parent: &notes, Head: c})
	}
	cursor := 0

	displayHelloWorld(s, &notes, cursor)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHelloWorld(s, &notes, cursor)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			case tcell.KeyRight:
				newItem(&notes.Tail[cursor])
				cursor++
				cursorLim := len(notes.Tail)
				cursor %= cursorLim
				displayHelloWorld(s, &notes, cursor)
			case tcell.KeyLeft:
				notes.Tail[cursor].Remove()
				if cursor != 0 {
					cursor--
				}
				displayHelloWorld(s, &notes, cursor)
			//case tcell.KeyUp:
			//	if cursor == 0 {
			//		continue
			//	}
			//	notes.Tail[cursor].MoveUp()
			//	cursor--
			//	displayHelloWorld(s, &notes, cursor)
			//case tcell.KeyDown:
			//	notes.Tail[cursor].MoveDown()
			//	cursor++
			//	cursor %= len(notes.Tail)
			//	displayHelloWorld(s, &notes, cursor)
			case tcell.KeyUp:
				if cursor == 0 {
					continue
				}
				cursor--
				displayHelloWorld(s, &notes, cursor)
			case tcell.KeyDown:
				cursor++
				cursorLim := len(notes.Tail)
				cursor %= cursorLim
				displayHelloWorld(s, &notes, cursor)
			//case tcell.KeyEnter:

			case tcell.KeyRune:
				input := string(ev.Rune())
				notes.Tail[cursor].Head += input
				displayHelloWorld(s, &notes, cursor)
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

func displayHelloWorld(s tcell.Screen, n *item, cursor int) {
	//w, h := s.Size()
	s.Clear()
	start := 10
	emitStr(s, start-1, start+cursor, tcell.StyleDefault, ">")
	emitStr(s, start, start-5, tcell.StyleDefault, "Press ESC to exit.")
	emitStr(s, start, start-4, tcell.StyleDefault, n.Head)
	emitStr(s, start, start-3, tcell.StyleDefault, "Should see all these: "+n.StringChildren())
	emitStr(s, start, start-1, tcell.StyleDefault, "Cursor ID: "+strconv.Itoa(cursor))
	for index := range n.Tail {
		emitStr(s, start, start+index, tcell.StyleDefault, n.Tail[index].Head)
	}
	s.Sync()
}
