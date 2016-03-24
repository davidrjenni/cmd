// Copyright (c) 2015 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "math/rand"

type key int

const (
	up key = iota
	down
	left
	right
	space
	quit
	nop
)

type field [][]int

func newField(w, h int) field {
	f := make([][]int, w)
	for x := 0; x < w; x++ {
		f[x] = make([]int, h)
	}
	return field(f)
}

func (f *field) reset() {
	for x := 0; x < f.width(); x++ {
		for y := 0; y < f.height(); y++ {
			if (*f)[x][y] == 2 {
				(*f)[x][y] = 0
			}
		}
	}
}

func (f field) full() bool {
	for x := 0; x < f.width(); x++ {
		if f[x][1] == 1 {
			return true
		}
	}
	return false
}

func (f *field) delFullLines() (int, bool) {
	base := false
	var lines []int
	for y := 0; y < f.height(); y++ {
		full := true
		for x := 0; x < f.width(); x++ {
			if (*f)[x][y] != 1 {
				full = false
				break
			}
		}
		if full {
			lines = append(lines, y)
			if y == f.height()-1 {
				base = true
			}
		}
	}
	for _, l := range lines {
		for y := l; y > 1; y-- {
			for x := 0; x < f.width(); x++ {
				(*f)[x][y] = (*f)[x][y-1]
			}
		}
	}
	return len(lines), base
}

func (f field) possible(b block, x, y int) bool {
	for fx, bx := x, 0; bx < 5; fx, bx = fx+1, bx+1 {
		for fy, by := y, 0; by < 5; fy, by = fy+1, by+1 {
			if !b.at(bx, by) {
				continue
			}
			if fx < 0 || fy < 0 || fx >= f.width() || fy >= f.height() || f[fx][fy] == 1 {
				return false
			}
		}
	}
	return true
}

func (f field) width() int { return len(f) }

func (f field) height() int { return len(f[0]) }

func (f *field) put(b block) {
	for fx, bx := b.x, 0; bx < 5; bx, fx = bx+1, fx+1 {
		for fy, by := b.y, 0; by < 5; by, fy = by+1, fy+1 {
			if !b.at(bx, by) {
				continue
			}
			if fx >= 0 || fy >= 0 || fx < f.width() || fy < f.height() {
				if b.active {
					(*f)[fx][fy] = 2
				} else {
					(*f)[fx][fy] = 1
				}
			}
		}
	}
}

const (
	rotN  = 4
	typeN = 7
)

type block struct {
	active   bool
	rot, typ int
	x, y     int
}

func randBlock() block {
	r := rand.Intn(rotN)
	t := rand.Intn(typeN)
	return block{active: true, rot: r, typ: t, x: initPos[t][r][0], y: initPos[t][r][1]}
}

func (b *block) rotate() { b.rot = (b.rot + 1) % rotN }

func (b block) at(x, y int) bool {
	return allBlocks[b.typ][b.rot][y][x] == 1
}

var allBlocks = [][][][]int{
	{ // I.
		{
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{1, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 1, 1},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // T.
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // S.
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // Z.
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // L.
		{
			{0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 1, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // J.
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0},
		},
	},
	{ // O.
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		},
	},
}

var initPos = [][][]int{
	{ // I.
		{-2, 0},
		{-2, -2},
		{-2, -1},
		{-2, -2},
	},
	{ // T.
		{-2, -1},
		{-2, -1},
		{-2, -2},
		{-2, -1},
	},
	{ // S.
		{-2, -1},
		{-2, -1},
		{-2, -1},
		{-2, -2},
	},
	{ // Z.
		{-2, -1},
		{-2, -1},
		{-2, -1},
		{-2, -2},
	},
	{ // L.
		{-2, -1},
		{-2, -2},
		{-2, -1},
		{-2, -1},
	},
	{ // J.
		{-2, -1},
		{-2, -1},
		{-2, -1},
		{-2, -2},
	},
	{ // O.
		{-2, -2},
		{-2, -2},
		{-2, -2},
		{-2, -2},
	},
}
