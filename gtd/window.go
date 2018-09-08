// Copyright (c) 2018 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

type window struct {
	line   int
	height int
	width  int
	tasks  []Task
	store  store
}

func newWindow(s store) (*window, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	tasks, err := s.open()
	if err != nil {
		return nil, err
	}
	w := &window{store: s, tasks: tasks}
	w.width, w.height = termbox.Size()
	return w, nil
}

func (w window) close() { termbox.Close() }

func (w window) save() error { return w.store.save(w.tasks) }

func (w window) draw() error {
	for y := 0; y < w.height; y++ {
		// Empty line
		if len(w.tasks) <= y {
			for x := 0; x < w.width; x++ {
				termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlack)
			}
			continue
		}

		task := w.tasks[y]
		termbox.SetCell(0, y, '[', termbox.ColorWhite, termbox.ColorBlack)
		if task.Done {
			termbox.SetCell(1, y, '✔', termbox.ColorWhite, termbox.ColorBlack)
		} else {
			termbox.SetCell(1, y, ' ', termbox.ColorWhite, termbox.ColorBlack)
		}
		termbox.SetCell(2, y, ']', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(3, y, ' ', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(4, y, task.Priority.rune(), task.Priority.attr(), termbox.ColorBlack)
		termbox.SetCell(5, y, ' ', termbox.ColorWhite, termbox.ColorBlack)

		title := []rune(task.Title)
		for x := 6; x < w.width; x++ {
			attr := termbox.ColorWhite
			if y == w.line {
				attr |= termbox.AttrBold
			}
			if x-6 < len(title) {
				termbox.SetCell(x, y, title[x-6], attr, termbox.ColorBlack)
			} else {
				termbox.SetCell(x, y, ' ', attr, termbox.ColorBlack)
			}
		}
	}
	return termbox.Flush()
}

func (w window) poll() error {
	for {
		e := termbox.PollEvent()
	Decide:
		switch {
		case unicode.IsDigit(e.Ch):
			n := int(e.Ch - '0')
			switch e = termbox.PollEvent(); e.Ch {
			case 'g':
				if w.line = n; w.line >= len(w.tasks) || w.line >= w.height {
					w.line = min(len(w.tasks), w.height) - 1
				}
			case 'j':
				if w.line += n; w.line >= len(w.tasks) || w.line >= w.height {
					w.line = min(len(w.tasks), w.height) - 1
				}
			case 'k':
				if w.line -= n; w.line < 0 {
					w.line = 0
				}
			default:
				goto Decide
			}
		case e.Ch == 'g':
			w.line = 0
		case e.Ch == 'G':
			if w.line = w.height - 1; w.line >= len(w.tasks) {
				w.line = len(w.tasks) - 1
			}
		case e.Ch == 'k' || e.Key == termbox.KeyArrowUp || e.Key == termbox.KeyBackspace2:
			if w.line--; w.line < 0 {
				if w.line = w.height - 1; w.line >= len(w.tasks) {
					w.line = len(w.tasks) - 1
				}
			}
		case e.Ch == 'j' || e.Key == termbox.KeyArrowDown:
			if w.line++; w.line >= w.height || w.line >= len(w.tasks) {
				w.line = 0
			}
		case e.Ch == '+':
			task := w.tasks[w.line]
			if task.Priority < Critical {
				w.tasks[w.line].Priority++
				if err := w.save(); err != nil {
					return err
				}
			}
		case e.Ch == '-':
			task := w.tasks[w.line]
			if task.Priority > No {
				w.tasks[w.line].Priority--
				if err := w.save(); err != nil {
					return err
				}
			}
		case e.Ch == 'a':
			termbox.Close()
			title, err := editor("")
			if err != nil {
				return nil
			}
			w.tasks = append(w.tasks, Task{Title: title})
			if w.line = len(w.tasks) - 1; w.line >= w.height {
				w.line = w.height - 1
			}
			if err := w.save(); err != nil {
				return err
			}
			if err := termbox.Init(); err != nil {
				return err
			}
		case e.Ch == 'e':
			termbox.Close()
			title, err := editor(w.tasks[w.line].Title)
			if err != nil {
				return nil
			}
			w.tasks[w.line].Title = title
			if err := w.save(); err != nil {
				return err
			}
			if err := termbox.Init(); err != nil {
				return err
			}
		case e.Key == termbox.KeySpace:
			w.tasks[w.line].Done = !w.tasks[w.line].Done
			if err := w.save(); err != nil {
				return err
			}
		case e.Ch == 'd' || e.Key == termbox.KeyDelete:
			if len(w.tasks) >= 1 {
				w.tasks = append(w.tasks[:w.line], w.tasks[w.line+1:]...)
				if w.line >= len(w.tasks) {
					w.line = len(w.tasks) - 1
				}
				if err := w.save(); err != nil {
					return err
				}
			}
		case e.Ch == 's':
			switch e = termbox.PollEvent(); e.Ch {
			case 'p':
				sort.SliceStable(w.tasks, func(i, j int) bool {
					return w.tasks[i].Priority > w.tasks[j].Priority
				})
				w.save()
			case 'c':
				sort.SliceStable(w.tasks, func(i, j int) bool {
					return !w.tasks[i].Done && w.tasks[j].Done
				})
				w.save()
			case 't':
				sort.SliceStable(w.tasks, func(i, j int) bool {
					return w.tasks[i].Title < w.tasks[j].Title
				})
				w.save()
			default:
				goto Decide
			}
		case e.Key == termbox.MouseLeft:
			if w.line = e.MouseY; w.line >= len(w.tasks) {
				w.line = len(w.tasks) - 1
			}
			if 0 <= e.MouseX && e.MouseX <= 2 {
				w.tasks[w.line].Done = !w.tasks[w.line].Done
				if err := w.save(); err != nil {
					return err
				}
			} else if 3 <= e.MouseX && e.MouseX <= 5 {
				if w.tasks[w.line].Priority == Critical {
					w.tasks[w.line].Priority = No
				} else {
					w.tasks[w.line].Priority++
				}
				if err := w.save(); err != nil {
					return err
				}
			}
		case e.Type == termbox.EventResize:
			w.width, w.height = e.Width, e.Height
			if w.line >= w.height {
				if w.line = w.height - 1; w.line >= len(w.tasks) {
					w.line = len(w.tasks) - 1
				}
			}
		case e.Ch == 'q' || e.Key == termbox.KeyEsc || e.Key == termbox.KeyCtrlC:
			termbox.Close()
			os.Exit(0)
		}
		if err := w.draw(); err != nil {
			return err
		}
	}
}

func editor(text string) (string, error) {
	file, err := ioutil.TempFile("", "tasks-")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	if err := ioutil.WriteFile(file.Name(), []byte(text), 0644); err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return "", err
	}
	if err = cmd.Wait(); err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(file.Name())
	return string(b), err
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
