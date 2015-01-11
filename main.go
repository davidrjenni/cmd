// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
rproxy is a simple websocket-aware HTTP reverse proxy.

Installation:
	% go get github.com/davidrjenni/rproxy

Usage:
	% rproxy -target http[s]://... [-addr ...]

Example
	% rproxy -target "https://example.com:8000" -addr ":8080"
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	addr   = flag.String("addr", ":8080", "HTTP listen address")
	target = flag.String("target", "", "pass all requests to this address")
)

// reverseProxy represents a websocket-aware HTTP reverse proxy.
type reverseProxy struct {
	proxy  *httputil.ReverseProxy
	target *url.URL
}

func newReverseProxy(target string) (*reverseProxy, error) {
	p := &reverseProxy{}
	var err error
	p.target, err = url.Parse(target)
	if err != nil {
		return nil, err
	}
	p.proxy = httputil.NewSingleHostReverseProxy(p.target)
	return p, nil
}

func (p *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrade := r.Header.Get("Upgrade")
	if strings.ToLower(upgrade) == "websocket" {
		fmt.Println("is websocket")
		p.handleWebsocket(w, r)
	} else {
		p.proxy.ServeHTTP(w, r)
	}
}

func (p *reverseProxy) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	dst, err := net.Dial("tcp", p.target.Host)
	if err != nil {
		http.Error(w, "Error dialing target.", 500)
		log.Println("1")
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "ResponseWriter is not a hijacker?", 500)
		log.Println("2")
		return
	}
	src, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, "Hijack error: "+err.Error(), 500)
		log.Println("3")
		return
	}
	defer src.Close()
	defer dst.Close()

	err = r.Write(dst)
	if err != nil {
		http.Error(w, "Error copying request to target: "+err.Error(), 500)
		log.Println("4")
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(dst, src)
	go cp(src, dst)
	<-errc
}

func main() {
	flag.Parse()
	if *target == "" {
		usage()
	}
	proxy, err := newReverseProxy(*target)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func usage() {
	fmt.Println("Usage: rproxy -host http[s]://... [-addr ...]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}
