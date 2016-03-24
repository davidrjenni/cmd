// Copyright (c) 2015 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestNewSnake(t *testing.T) {
	s := newSnake(1, 2, 0)

	if s[0].x != 1 || s[0].y != 2 {
		t.Errorf("Expected (1, 2), got (%d, %d)", s[0].x, s[0].y)
	}
	if len(s) != 1 {
		t.Errorf("Expected 3, got %d", len(s))
	}
}

func TestNewSnakeWithLength(t *testing.T) {
	s := newSnake(1, 1, 2)

	if len(s) != 3 {
		t.Errorf("Expected 3, got %d", len(s))
	}
}

func TestAte(t *testing.T) {
	s := newSnake(1, 1, 3)

	if !s.ate(point{1, 1}) {
		t.Errorf("Expected that the snake ate the food")
	}
	if s.ate(point{1, 0}) {
		t.Errorf("Expected that the snake did not ate the food")
	}
	if s.ate(point{0, 1}) {
		t.Errorf("Expected that the snake did not ate the food")
	}
}

func TestGrow(t *testing.T) {
	s := newSnake(1, 1, 0)

	s.grow()

	if len(s) != 2 {
		t.Fatalf("Expected 2, got %d", len(s))
	}
	if s[1].x >= 0 {
		t.Errorf("Expected x < 0, got %d", s[1].x)
	}
	if s[1].y >= 0 {
		t.Errorf("Expected y < 0, got %d", s[1].y)
	}
}

func TestMoveUp(t *testing.T) {
	s := newSnake(2, 2, 1)

	coll := s.move(up, 0, 0, 10, 10)

	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 2 || s[0].y != 1 {
		t.Errorf("Expected (2, 1), got (%d, %d)", s[0].x, s[0].y)
	}
	if s[1].x != 2 || s[1].y != 2 {
		t.Errorf("Expected (2, 2), got (%d, %d)", s[1].x, s[1].y)
	}
}

func TestMoveUpperWall(t *testing.T) {
	s := newSnake(5, 0, 0)

	coll := s.move(up, 0, 0, 10, 10)
	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 5 || s[0].y != 10 {
		t.Errorf("Expected (5, 10), got (%d, %d)", s[0].x, s[0].y)
	}
}

func TestMoveDown(t *testing.T) {
	s := newSnake(2, 2, 1)

	coll := s.move(down, 0, 0, 10, 10)

	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 2 || s[0].y != 3 {
		t.Errorf("Expected (2, 1), got (%d, %d)", s[0].x, s[0].y)
	}
	if s[1].x != 2 || s[1].y != 2 {
		t.Errorf("Expected (2, 2), got (%d, %d)", s[1].x, s[1].y)
	}
}

func TestMoveLowerWall(t *testing.T) {
	s := newSnake(5, 10, 0)

	coll := s.move(down, 0, 0, 10, 10)
	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 5 || s[0].y != 0 {
		t.Errorf("Expected (5, 0), got (%d, %d)", s[0].x, s[0].y)
	}
}

func TestMoveLeft(t *testing.T) {
	s := newSnake(2, 2, 1)

	coll := s.move(left, 0, 0, 10, 10)

	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 1 || s[0].y != 2 {
		t.Errorf("Expected (2, 1), got (%d, %d)", s[0].x, s[0].y)
	}
	if s[1].x != 2 || s[1].y != 2 {
		t.Errorf("Expected (2, 2), got (%d, %d)", s[1].x, s[1].y)
	}
}

func TestMoveLeftWall(t *testing.T) {
	s := newSnake(0, 5, 0)

	coll := s.move(left, 0, 0, 10, 10)
	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 10 || s[0].y != 5 {
		t.Errorf("Expected (10, 5), got (%d, %d)", s[0].x, s[0].y)
	}
}

func TestMoveRight(t *testing.T) {
	s := newSnake(2, 2, 1)

	coll := s.move(right, 0, 0, 10, 10)

	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 3 || s[0].y != 2 {
		t.Errorf("Expected (2, 1), got (%d, %d)", s[0].x, s[0].y)
	}
	if s[1].x != 2 || s[1].y != 2 {
		t.Errorf("Expected (2, 2), got (%d, %d)", s[1].x, s[1].y)
	}
}

func TestMoveRightWall(t *testing.T) {
	s := newSnake(10, 5, 0)

	coll := s.move(right, 0, 0, 10, 10)
	if coll {
		t.Errorf("Expected no collision")
	}
	if s[0].x != 0 || s[0].y != 5 {
		t.Errorf("Expected (0, 5), got (%d, %d)", s[0].x, s[0].y)
	}
}

func TestCollide(t *testing.T) {
	s := newSnake(2, 2, 4)

	s.move(up, 0, 0, 10, 10)
	s.move(left, 0, 0, 10, 10)
	s.move(down, 0, 0, 10, 10)

	coll := s.move(right, 0, 0, 10, 10)
	if !coll {
		t.Errorf("Expected collision")
	}
}
