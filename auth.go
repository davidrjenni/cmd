// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

// fileTokenStore persists a token.
type fileTokenStore struct {
	filename string // filename of the cache file
}

// ReadToken reads a token from the cache file.
func (s *fileTokenStore) ReadToken() (tok *oauth2.Token, err error) {
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	err = gob.NewDecoder(f).Decode(&tok)
	return
}

// WriteToken writes a token to the cache file.
func (s *fileTokenStore) WriteToken(token *oauth2.Token) {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

// oauthClient returns an authorized HTTP client.
func oauthClient(conf *oauth2.Config) (*http.Client, error) {
	store := &fileTokenStore{filename: cacheFilename(conf)}
	tok, err := store.ReadToken()
	if err != nil || (tok != nil && tok.Expired()) {
		code := codeFromWeb(conf)
		tok, err = conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return nil, err
		}
		store.WriteToken(tok)
	}
	return conf.Client(oauth2.NoContext, tok), nil
}

// codeFromWeb returns an authorization code from the web.
func codeFromWeb(conf *oauth2.Config) string {
	ch := make(chan string)
	state := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		if r.FormValue("state") != state {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if code := r.FormValue("code"); code != "" {
			fmt.Fprintf(w, "<h1>Authorized.</h1>")
			ch <- code
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
	}))
	defer ts.Close()

	conf.RedirectURL = ts.URL
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)

	fmt.Printf("Visit the URL for the auth dialog: [%v]", url)
	return <-ch
}

// cacheFilename generates a filename for a given OAuth2 configuration.
func cacheFilename(conf *oauth2.Config) string {
	h := fnv.New32a()
	h.Write([]byte(conf.ClientID))
	h.Write([]byte(conf.ClientSecret))
	for _, s := range conf.Scopes {
		h.Write([]byte(s))
	}
	name := fmt.Sprintf("taskfs-token-%v", h.Sum32())
	return filepath.Join(osUserCacheDir(), url.QueryEscape(name))
}

// osUserCacheDir returns the cache directory for the current user.
func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	log.Printf("TODO: osUserCacheDir on GOOS %q", runtime.GOOS)
	return "."
}
