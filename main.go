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
)

// This program just prints "Hello, World!".  Press ESC to exit.
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

	notes := setupExamples()
	cursor := 0
	cursorLim := getCursorLim(notes)

	displayHelloWorld(s, notes, cursor, cursorLim)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHelloWorld(s, notes, cursor, cursorLim)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			case tcell.KeyUp:
				newItem(&notes.Tail[cursor])
				displayHelloWorld(s, notes, cursor, cursorLim)
				cursor++
				s.Sync()
			case tcell.KeyDown:
				notes.Tail[cursor].Remove()
				displayHelloWorld(s, notes, cursor, cursorLim)
				s.Sync()
			case tcell.KeyLeft:
				notes.Tail[cursor].MoveUp()
				//cursor %= cursorLim
				displayHelloWorld(s, notes, cursor, cursorLim)
				cursor--
				s.Sync()
			case tcell.KeyRight:
				notes.Tail[cursor].MoveDown()
				displayHelloWorld(s, notes, cursor, cursorLim)
				cursor++
				//cursor %= cursorLim
				s.Sync()
			}
		}
	}
}

func setupExamples() *item {
	notes := item{Home: true, Head: "notes"}
	ch := []string{"£", "$", "%", "^", "&", "*"}
	for _, c := range ch {
		notes.Tail = append(notes.Tail, item{Parent: &notes, Head: c})
	}
	//fmt.Println(len(notes.Tail))
	newItem(&notes.Tail[3])
	//fmt.Println(len(notes.Tail))

	return &notes
}

func getCursorLim(selected *item) int {
	return len(selected.Tail)
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

func displayHelloWorld(s tcell.Screen, n *item, cursor int, cursorLim int) {
	w, h := s.Size()
	s.Clear()
	emitStr(s, w/2-7, h/2-2, tcell.StyleDefault, string(cursor))
	emitStr(s, w/2-7, h/2-1, tcell.StyleDefault, n.Head)
	emitStr(s, w/2-7, h/2, tcell.StyleDefault, n.StringChildren())
	emitStr(s, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
	s.Show()
}
