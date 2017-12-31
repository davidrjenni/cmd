// Copyright (c) 2015 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/nsf/termbox-go"

type win struct {
	input chan key
	w, h  int
}

func newWin(w, h int) (win, error) {
	if err := termbox.Init(); err != nil {
		return win{}, err
	}
	win := win{make(chan key), w, h}
	go win.poll()
	return win, nil
}

func (w win) close() { termbox.Close() }

func (w win) draw(s snake, food point) {
	const (
		cherry = '◆'
		wall   = '▓'
		head   = '◎'
		tail   = '●'
	)
	for x := 0; x < w.w; x++ {
		termbox.SetCell(x, 0, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	for y := 1; y < w.h-1; y++ {
		termbox.SetCell(0, y, wall, termbox.ColorWhite, termbox.ColorBlack)
		for x := 1; x < w.w-1; x++ {
			if x == food.x && y == food.y {
				termbox.SetCell(x, y, cherry, termbox.ColorRed, termbox.ColorBlack)
			} else if x == s[0].x && y == s[0].y {
				termbox.SetCell(x, y, head, termbox.ColorGreen, termbox.ColorBlack)
			} else {
				found := false
				for i := 1; i < len(s); i++ {
					if x == s[i].x && y == s[i].y {
						termbox.SetCell(x, y, tail, termbox.ColorGreen, termbox.ColorBlack)
						found = true
					}
				}
				if !found {
					termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorBlack)
				}
			}
		}
		termbox.SetCell(w.w-1, y, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	for x := 0; x < w.w; x++ {
		termbox.SetCell(x, w.h-1, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}

func (w win) poll() {
	for {
		e := termbox.PollEvent()
		switch e.Ch {
		case 'k':
			w.input <- up
		case 'j':
			w.input <- down
		case 'h':
			w.input <- left
		case 'l':
			w.input <- right
		case 'q':
			w.input <- quit
		}
		switch e.Key {
		case termbox.KeyArrowUp:
			w.input <- up
		case termbox.KeyArrowDown:
			w.input <- down
		case termbox.KeyArrowLeft:
			w.input <- left
		case termbox.KeyArrowRight:
			w.input <- right
		case termbox.KeyEsc, termbox.KeyCtrlC:
			w.input <- quit
		}
	}
}
