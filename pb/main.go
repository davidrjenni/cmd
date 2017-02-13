// Copyright (c) 2017 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
pb prints the bits of an integer.

Installation:
	% go get github.com/davidrjenni/cmd/pb

Usage:
	% pb integer
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pb: ")

	if len(os.Args) < 2 {
		log.Fatal("Usage: pb integer")
	}

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%b\n", i)
}
