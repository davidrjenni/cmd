// Copyright (c) 2015 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type key int

const (
	up key = iota
	down
	left
	right
	quit
)

type point struct{ x, y int }

type snake []point

func newSnake(x, y, l int) snake {
	s := snake([]point{{x, y}})
	for i := 0; i < l; i++ {
		s.grow()
	}
	return s
}

func (s snake) ate(food point) bool {
	return s[0].x == food.x && s[0].y == food.y
}

func (s *snake) grow() {
	*s = append(*s, point{-1, -1})
}

func (s snake) move(k key, minx, miny, maxx, maxy int) bool {
	x, y := s[0].x, s[0].y
	switch k {
	case up:
		s[0].y = checkMin(s[0].y-1, miny, maxy)
	case down:
		s[0].y = checkMax(s[0].y+1, miny, maxy)
	case left:
		s[0].x = checkMin(s[0].x-1, minx, maxx)
	case right:
		s[0].x = checkMax(s[0].x+1, minx, maxx)
	}
	for i := 1; i < len(s); i++ {
		tx, ty := s[i].x, s[i].y
		s[i].x, s[i].y = x, y
		x, y = tx, ty
	}
	return s.collision()
}

func checkMin(next, min, max int) int {
	if next < min {
		return max
	}
	return next
}

func checkMax(next, min, max int) int {
	if next > max {
		return min
	}
	return next
}

func (s snake) collision() bool {
	for i := 1; i < len(s); i++ {
		if s[i].x == s[0].x && s[i].y == s[0].y {
			return true
		}
	}
	return false
}
