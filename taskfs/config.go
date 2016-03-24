// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/oauth2"
	"google.golang.org/api/tasks/v1"
)

// OAuth2 configuration.
var conf = &oauth2.Config{
	ClientID:     "YOUR CLIENT ID",
	ClientSecret: "YOUR CLIENT SECRET",
	Scopes:       []string{tasks.TasksScope},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	},
}
