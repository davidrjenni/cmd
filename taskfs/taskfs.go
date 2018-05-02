// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"strings"

	p "github.com/mortdeus/go9p"
	"github.com/mortdeus/go9p/srv"
	"google.golang.org/api/tasks/v1"
)

// notes is a file which contains the notes of a task.
type notes struct {
	listID  string
	task    *tasks.Task
	service *tasks.TasksService
	srv.File
}

// addNotes adds the notes of a task to a given directory.
func addNotes(dir *srv.File, user p.User, task *tasks.Task, listID string, service *tasks.Service) error {
	n := new(notes)
	n.listID = listID
	n.task = task
	n.service = tasks.NewTasksService(service)
	return n.Add(dir, "notes", user, nil, 0666, n)
}

// Write handles a write to the notes.
func (n *notes) Write(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	n.Lock()
	defer n.Unlock()
	prevNotes := n.task.Notes
	l := offset + uint64(len(buf))
	b := make([]byte, l)
	c := copy(b[offset:], buf)
	n.task.Notes = string(b)
	_, err := n.service.Update(n.listID, n.task.Id, n.task).Do()
	if err != nil {
		n.task.Notes = prevNotes
		return 0, err
	}
	return c, nil
}

// Read handles a read of the notes.
func (n *notes) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	n.Lock()
	defer n.Unlock()
	var err error
	n.task, err = n.service.Get(n.listID, n.task.Id).Do()
	if err != nil {
		return 0, err
	}
	b := []byte(n.task.Notes)
	off := int(offset)
	if off >= len(b) {
		return 0, nil
	}
	return copy(buf, b[off:]), nil
}

// task is a directory which contains its properties as files.
type task struct {
	listID       string
	task         *tasks.Task
	tasksService *tasks.TasksService
	service      *tasks.Service
	user         p.User
	relevant     bool
	srv.File
}

// addTask adds a task to the given directory.
func addTask(dir *srv.File, user p.User, tsk *tasks.Task, listID string, service *tasks.Service) (*task, error) {
	t := new(task)
	t.listID = listID
	t.task = tsk
	t.tasksService = tasks.NewTasksService(service)
	t.service = service
	t.user = user
	if err := t.addNotes(); err != nil {
		return nil, err
	}

	// the title is used as directory name, so spaces are annoying.
	title := strings.Replace(tsk.Title, " ", "_", -1)
	return t, t.Add(dir, title, user, nil, p.DMDIR|0755, t)
}

// addNotes adds the notes file to the directory.
func (t *task) addNotes() error {
	return addNotes(&t.File, t.user, t.task, t.listID, t.service)
}

// cloneTask is a file which is used to create a new task.
type cloneTask struct {
	list    *list
	service *tasks.Service
	srv.File
}

// addCloneTask adds a cloneTask to a given directory.
func addCloneTask(dir *srv.File, user p.User, list *list, service *tasks.Service) error {
	f := new(cloneTask)
	f.list = list
	f.service = service
	return f.Add(dir, "clone", user, nil, 0666, f)
}

// Write handles a write to the cloneTask by adding a new task to the lists directory.
func (c *cloneTask) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	c.Lock()
	defer c.Unlock()
	task := &tasks.Task{Title: string(data)}
	task, err := tasks.NewTasksService(c.service).Insert(c.list.list.Id, task).Do()
	if err != nil {
		return 0, err
	}
	t, err := c.list.addTask(task)
	if err != nil {
		return 0, err
	}
	c.list.tasks[task.Id] = t
	return len(data), nil
}

// list is a directory which contains tasks.
type list struct {
	list             *tasks.TaskList
	relevant         bool
	service          *tasks.Service
	user             p.User
	tasks            map[string]*task
	tasklistsService *tasks.TasklistsService
	tasksService     *tasks.TasksService
	srv.File
}

// addList populates a given directory with the tasks of a given list.
func addList(dir *srv.File, user p.User, lst *tasks.TaskList, service *tasks.Service) (*list, error) {
	l := new(list)
	l.list = lst
	l.user = user
	l.service = service
	l.tasks = make(map[string]*task)
	l.tasklistsService = tasks.NewTasklistsService(service)
	l.tasksService = tasks.NewTasksService(service)
	if err := addCloneTask(&l.File, user, l, service); err != nil {
		return nil, err
	}
	return l, l.Add(dir, lst.Title, user, nil, p.DMDIR|0755, l)
}

// addTask adds a given task to the directory.
func (l *list) addTask(t *tasks.Task) (*task, error) {
	return addTask(&l.File, l.user, t, l.list.Id, l.service)
}

// Open populates the directory with the relevant tasks of the list.
func (l *list) Open(fid *srv.FFid, mode uint8) error {
	tasks, err := l.tasksService.List(l.list.Id).Do()
	if err != nil {
		return err
	}
	// Mark all as outdated.
	for id := range l.tasks {
		l.tasks[id].relevant = false
	}
	for _, task := range tasks.Items {
		if _, ok := l.tasks[task.Id]; !ok {
			t, err := l.addTask(task)
			if err != nil {
				return err
			}
			l.tasks[task.Id] = t
		}
		// Mark as relevant and update.
		t := l.tasks[task.Id]
		t.relevant = true
		t.task = task
		t.File.Rename(t.task.Title)
	}
	// Remove outdated.
	for id, task := range l.tasks {
		if !task.relevant {
			task.File.Remove()
			delete(l.tasks, id)
		}
	}
	return nil
}

// lists is a directory which contains task lists.
type lists struct {
	lists            map[string]*list
	service          *tasks.Service
	tasklistsService *tasks.TasklistsService
	user             p.User
	dir              *srv.File
}

// addLists populates a given directory with the lists provided by the service.
func addLists(dir *srv.File, user p.User, service *tasks.Service) *lists {
	l := new(lists)
	l.dir = dir
	l.lists = make(map[string]*list)
	l.service = service
	l.tasklistsService = tasks.NewTasklistsService(service)
	l.user = user
	return l
}

// addList adds the given list to the directory.
func (l *lists) addList(list *tasks.TaskList) (*list, error) {
	return addList(l.dir, l.user, list, l.service)
}

// Open populates the directory with all the relevant lists.
func (l *lists) Open(fid *srv.FFid, mode uint8) error {
	lists, err := l.tasklistsService.List().Do()
	if err != nil {
		return err
	}
	// Mark all as outdated.
	for id := range l.lists {
		l.lists[id].relevant = false
	}
	for _, list := range lists.Items {
		if _, ok := l.lists[list.Id]; !ok {
			lst, err := l.addList(list)
			if err != nil {
				return err
			}
			l.lists[list.Id] = lst
		}
		// Mark as relevant and update.
		lst := l.lists[list.Id]
		lst.relevant = true
		lst.list = list
		lst.File.Rename(lst.list.Title)
	}
	// Remove outdated lists.
	for id, list := range l.lists {
		if !list.relevant {
			list.File.Remove()
			delete(l.lists, id)
		}
	}
	return nil
}

// cloneList is a file which is used to create a new list.
type cloneList struct {
	lists   *lists
	service *tasks.Service
	srv.File
}

// addCloneList adds a cloneList to a given directory.
func addCloneList(dir *srv.File, user p.User, lists *lists, service *tasks.Service) error {
	f := new(cloneList)
	f.lists = lists
	f.service = service
	return f.Add(dir, "clone", user, nil, 0666, f)
}

// Write handles a write to the cloneList by adding a new list to the lists directory.
func (c *cloneList) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	c.Lock()
	defer c.Unlock()
	list := &tasks.TaskList{Title: string(data)}
	list, err := tasks.NewTasklistsService(c.service).Insert(list).Do()
	if err != nil {
		return 0, err
	}
	l, err := c.lists.addList(list)
	if err != nil {
		return 0, err
	}
	c.lists.lists[l.list.Id] = l
	return len(data), nil
}

// newTasklistsFsrv creates a file server which serves the task lists.
func newTasklistsFsrv(service *tasks.Service) (*srv.Fsrv, error) {
	user := p.OsUsers.Uid2User(os.Geteuid())
	root := new(srv.File)
	lists := addLists(root, user, service)
	err := root.Add(nil, "/", user, nil, p.DMDIR|0555, lists)
	if err != nil {
		return nil, err
	}
	if err := addCloneList(root, user, lists, service); err != nil {
		return nil, err
	}
	srv := srv.NewFileSrv(root)
	srv.Dotu = true
	srv.Debuglevel = 0
	return srv, nil
}

// TCP listen address for the fileserver.
var addr = flag.String("tcp", ":8080", "tcp listen address")

func main() {
	flag.Parse()

	client, err := oauthClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	service, err := tasks.New(client)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := newTasklistsFsrv(service)
	if err != nil {
		log.Fatal(err)
	}
	srv.Start(srv)
	err = srv.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
}
