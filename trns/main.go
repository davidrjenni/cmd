// Copyright (c) 2018 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
trns is a simple translation tool for the terminal using the DeepL
service.

Installation:
	% go get github.com/davidrjenni/cmd/trns

Usage:
	% trns

Enter words to translate or commands. Commands are prefixed with a
colon. The commands are:
	:c		# copy the last translation to the clipboard, using xclip
	:src <lang>	# set the language to translate from
	:dst <lang>	# set the language to translate to

where <lang> is one of
	en	for English
	de	for German
	fr	for French
	es	for Spanish
	it	for Italian
	nl	for Dutch
	pl	for Polish

Default is German for the source language and English for the
destination language.
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ravernkoh/deepl"
	"golang.org/x/crypto/ssh/terminal"
)

var langs = map[string]bool{
	deepl.English: true,
	deepl.German:  true,
	deepl.French:  true,
	deepl.Dutch:   true,
	deepl.Spanish: true,
	deepl.Italian: true,
	deepl.Polish:  true,
	deepl.Auto:    true,
}

type translator struct {
	cli *deepl.Client
	src string // source language
	dst string // destination language

	term         *terminal.Terminal
	translations []string // all the translations for the last text
	last         string   // last translation given
}

type shell struct {
	r io.Reader
	w io.Writer
}

func (sh *shell) Read(data []byte) (n int, err error)  { return sh.r.Read(data) }
func (sh *shell) Write(data []byte) (n int, err error) { return sh.w.Write(data) }

func newTranslator(src, dst string) (t *translator, cleanup func(), err error) {
	fd := int(os.Stdin.Fd())
	old, err := terminal.MakeRaw(fd)
	if err != nil {
		return nil, nil, err
	}
	term := terminal.NewTerminal(&shell{r: os.Stdin, w: os.Stdout}, "> ")
	if term == nil {
		return nil, nil, errors.New("could not create terminal")
	}
	return &translator{
		cli:  deepl.NewClient(),
		src:  src,
		dst:  dst,
		term: term,
	}, func() { terminal.Restore(fd, old) }, nil
}

func (t *translator) scan() (string, error) { return t.term.ReadLine() }

func (t *translator) write(s string) { t.term.Write([]byte(s)) }

func (t *translator) setDst(lang string) error {
	if lang = strings.ToUpper(lang); langs[lang] {
		t.dst = lang
		return nil
	}
	return fmt.Errorf("no such language: %q", strings.ToLower(lang))
}

func (t *translator) setSrc(lang string) error {
	if lang = strings.ToUpper(lang); langs[lang] {
		t.src = lang
		return nil
	}
	return fmt.Errorf("no such language: %q", strings.ToLower(lang))
}

func (t *translator) translate(text string) error {
	res, err := t.cli.Translate([]string{text}, t.src, t.dst)
	if err != nil {
		return err
	}
	if len(res) == 0 || len(res[0]) == 0 {
		return fmt.Errorf("no translation for %q", text)
	}
	t.translations = res[0]
	return nil
}

func (t *translator) next() (string, int) {
	if len(t.translations) == 0 {
		return "", 0
	}
	t.last = t.translations[0]
	if len(t.translations) == 1 {
		return t.last, 0
	}
	t.translations = t.translations[1:]
	return t.last, len(t.translations) - 1
}

func (t *translator) copy() error {
	cmd := exec.Command("xclip", "-i", "-sel", "clip")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(t.last)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

func main() {
	t, cleanup, err := newTranslator(deepl.German, deepl.English)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for {
		input, err := t.scan()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "! %s\r\n", err)
			continue
		}

	Interpret:
		switch {
		case input == "":
			continue
		case strings.HasPrefix(input, ":src "):
			if err := t.setSrc(input[5:]); err != nil {
				fmt.Fprintf(os.Stderr, "! %s\r\n", err)
			}
		case strings.HasPrefix(input, ":dst "):
			if err := t.setDst(input[5:]); err != nil {
				fmt.Fprintf(os.Stderr, "! %s\r\n", err)
			}
		case input == ":c":
			if err := t.copy(); err != nil {
				fmt.Fprintf(os.Stderr, "! %s\r\n", err)
			}
		case len(input) >= 1 && input[0] == ':':
			fmt.Fprintf(os.Stderr, "! no such command: %q\r\n", input)

		default:
			if err := t.translate(input); err != nil {
				fmt.Fprintf(os.Stderr, "! %s\r\n", err)
				continue
			}

			for {
				translation, n := t.next()
				t.write(translation)
				if n == 0 {
					t.write("\n")
					break
				}
				t.write(fmt.Sprintf(" (and %d more, hit 'enter' for more)\n", n))

			Scan:
				input, err = t.scan()
				if err != nil {
					if err == io.EOF {
						return
					}
					fmt.Fprintf(os.Stderr, "! %s\r\n", err)
					continue
				}

				if input == ":c" {
					if err := t.copy(); err != nil {
						fmt.Fprintf(os.Stderr, "! %s\r\n", err)
					}
					goto Scan
				}

				if input != "" {
					goto Interpret
				}
			}
		}
	}
}
