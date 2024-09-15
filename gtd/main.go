// Copyright (c) 2018 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
gtd is a task manager.

Installation:

	% go get github.com/davidrjenni/cmd/gtd

Usage:

	% gtd
*/
package main

import (
	"encoding/json"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

type priority uint8

const (
	no priority = iota
	low
	medium
	high
	critical
)

func (p priority) attr() termbox.Attribute {
	switch p {
	case no:
		return termbox.ColorWhite
	case low:
		return termbox.ColorBlue
	case medium:
		return termbox.ColorYellow
	case high:
		return termbox.ColorRed
	case critical:
		return termbox.ColorRed | termbox.AttrReverse
	}
	panic("unhandled priority")
}

func (p priority) rune() rune {
	switch p {
	case no:
		return ' '
	case low:
		return 'l'
	case medium:
		return 'm'
	case high:
		return 'h'
	case critical:
		return 'c'
	}
	panic("unhandled priority")
}

// Task represents a single entry in the task list.
type Task struct {
	Title    string   `json:"title"`
	Priority priority `json:"priority"`
	Done     bool     `json:"done"`
	Deleted  bool     `json:"deleted"`
}

type store struct {
	filename string
}

func (s store) open() ([]Task, error) {
	b, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	if err = json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s store) save(tasks []Task) error {
	b, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", b, 0644)
}

func main() {
	log.SetPrefix("gtd: ")
	log.SetFlags(0)

	s := store{filename: "tasks.json"}
	w, err := newWindow(s)
	if err != nil {
		log.Fatal(err)
	}
	defer w.close()

	if err := w.draw(); err != nil {
		log.Fatal(err)
	}
	if err := w.poll(); err != nil {
		log.Fatal(err)
	}
}
