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
	"io/ioutil"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

type priority uint8

const (
	No priority = iota
	Low
	Medium
	High
	Critical
)

func (p priority) attr() termbox.Attribute {
	switch p {
	case No:
		return termbox.ColorWhite
	case Low:
		return termbox.ColorBlue
	case Medium:
		return termbox.ColorYellow
	case High:
		return termbox.ColorRed
	case Critical:
		return termbox.ColorRed | termbox.AttrReverse
	}
	panic("unhandled priority")
}

func (p priority) rune() rune {
	switch p {
	case No:
		return ' '
	case Low:
		return 'l'
	case Medium:
		return 'm'
	case High:
		return 'h'
	case Critical:
		return 'c'
	}
	panic("unhandled priority")
}

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
	b, err := ioutil.ReadFile("tasks.json")
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
	return ioutil.WriteFile("tasks.json", b, 0644)
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
