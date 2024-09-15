// Copyright (c) 2015 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

func main() {
	const (
		height = 20
		width  = 14
	)
	w, err := newWin(width, height)
	if err != nil {
		log.Fatal(err)
	}
	var score, level, lines int
	f := newField(width-1, height-1)
	d := nop

	for !f.full() {
		b := randBlock()
		b.x += width/2 - 2
		b.y++
		for b.active {
			timeout := time.Duration(500 * math.Pow(0.75, float64(level-1)))
			for {
				f.reset()
				switch d {
				case up:
					tmp := b
					tmp.rotate()
					if f.possible(tmp, tmp.x, tmp.y) {
						b.rotate()
					}
				case down:
					timeout = time.Duration(40)
				case left:
					if f.possible(b, b.x-1, b.y) {
						b.x--
					}
				case right:
					if f.possible(b, b.x+1, b.y) {
						b.x++
					}
				case space:
					for f.possible(b, b.x, b.y+1) {
						b.y++
					}
				}
				f.put(b)
				w.draw(f)
				select {
				case d = <-w.input:
					if d == quit {
						goto Exit
					}
				case <-time.After(timeout * time.Millisecond):
					goto Break
				}
			}
		Break:
			d = nop
			if f.possible(b, b.x, b.y+1) {
				b.y++
			} else {
				b.active = false
				f.put(b)
			}
		}
		n, baseLine := f.delFullLines()
		lines += n
		level = 1 + lines/10
		switch n {
		case 1:
			score += 4 * level
		case 2:
			score += 10 * level
		case 3:
			score += 30 * level
		case 4:
			score += 120 * level
		}
		if baseLine {
			score += 1000 * level
		}
		w.draw(f)
	}
Exit:
	w.close()
	fmt.Printf("Score: %d\n", score)
}
