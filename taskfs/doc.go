// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Taskfs is a 9P file server which presents your Google tasks in the local namespace.

Installation:
	% go get github.com/davidrjenni/cmd/taskfs
plan9port needs to be installed: http://swtch.com/plan9port

Usage:
	% taskfs &				# You may need to open the authorization dialog in your browser.
	% 9 srv 'tcp!localhost!8080' taskfs	# start network file service

	# You can manipulate your tasks with the 9p command.
	% 9p ls taskfs					# explore your task lists
	% echo -n Foo | 9p write taskfs/clone		# create a new task list
	% echo -n Foo | 9p write taskfs/<list>/clone	# create a new task
	% 9p ls taskfs/<list>				# explore your tasks
	% 9p read taskfs/<list>/<task>/notes		# read notes of a task
	% echo -n 'notes...' | 9p write taskfs/<list>/<task>/notes	# write notes of a task

	# Or you can mount the filesystem.
	% 9 mount 'unix!/tmp/<namespace>/taskfs' <dir>
	% cd <dir> && ls		# explore your task lists
	% echo -n <name> >> clone	# create a new task list
	% cd <list> && ls		# explore your tasks
	% echo -n <name> >> clone	# create a new task
	% cat <task>/notes		# read notes of a task
	% echo -n <notes> >> notes	# write notes of a task
	% 9 unmount <dir>
The changes are visible in the Google calendar or Gmail web app (after a refresh).
*/
package main
