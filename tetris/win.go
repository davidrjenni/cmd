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

func (w win) draw(f field) {
	const wall = '▓'
	for x := 0; x < w.w+1; x++ {
		termbox.SetCell(x, 0, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	for y := 0; y < f.height()+1; y++ {
		termbox.SetCell(0, y, wall, termbox.ColorWhite, termbox.ColorBlack)
		if y > 0 && y < f.height()+1 {
			for x := 1; x < f.width()+1; x++ {
				switch f[x-1][y-1] {
				case 1:
					termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorWhite)
				case 2:
					termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorGreen)
				default:
					termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorBlack)
				}
			}
		}
		termbox.SetCell(w.w, y, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	for x := 0; x < w.w+1; x++ {
		termbox.SetCell(x, w.h, wall, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}

func (w win) poll() {
	for {
		switch termbox.PollEvent().Key {
		case termbox.KeyArrowUp:
			w.input <- up
		case termbox.KeyArrowDown:
			w.input <- down
		case termbox.KeyArrowLeft:
			w.input <- left
		case termbox.KeyArrowRight:
			w.input <- right
		case termbox.KeySpace:
			w.input <- space
		case termbox.KeyEsc, termbox.KeyCtrlC:
			w.input <- quit
		}
	}
}
