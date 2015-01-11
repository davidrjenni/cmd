// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
rproxy is a simple HTTP reverse proxy.

Installation:
	% go get github.com/davidrjenni/rproxy

Usage:
	% rproxy -host http[s]://... [-addr ...]

Example
	% rproxy -host "https://example.com:8000" -addr ":8080"
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	addr = flag.String("addr", ":8080", "HTTP listen address")
	host = flag.String("host", "", "pass all requests to this host")
)

func main() {
	flag.Parse()
	if *host == "" {
		usage()
	}
	u, err := url.Parse(*host)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	s := &http.Server{
		Addr:    *addr,
		Handler: proxy,
	}
	log.Fatal(s.ListenAndServe())
}

func usage() {
	fmt.Println("Usage: rproxy -host http[s]://... [-addr ...]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}
