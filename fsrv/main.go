// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
fsrv is a file server.

Installation:
	% go get github.com/davidrjenni/cmd/fsrv

Usage:
	% fsrv [options]

Options:
	-addr	HTTP listen address (default: :8080)
	-auth	credentials for basic auth (default: none)
		example: fsrv -auth="user:pw"
	-dir	directory (default: .)
*/
package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"strings"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("fsrv: ")

	addr := flag.String("http", ":8080", "HTTP listen address")
	auth := flag.String("auth", "", "colon separated credentials for basic auth")
	dir := flag.String("dir", ".", "directory")
	flag.Parse()

	h := http.FileServer(http.Dir(*dir))
	if pair := strings.Split(*auth, ":"); len(pair) == 2 {
		h = basicAuth(pair[0], pair[1], h)
	}
	http.Handle("/", h)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func basicAuth(user, pw string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header["Authorization"]
		if header == nil {
			unauthorized(w)
			return
		}
		auth := strings.SplitN(header[0], " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			unauthorized(w)
			return
		}
		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			unauthorized(w)
			return
		}
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != user || pair[1] != pw {
			unauthorized(w)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"user\"")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
