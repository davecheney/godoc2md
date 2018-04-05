package main

import (
	"testing"
)

func TestUrlFromPackage(t *testing.T) {
	testData := []struct {
		pkg      string
		expected string
	}{
		{"github.com/davecheney/godoc2md", "https://github.com/davecheney/godoc2md/tree/master"},
		{"github.com/davecheney/godoc2md/examples", "https://github.com/davecheney/godoc2md/tree/master/examples"},
		{"github.com/davecheney/godoc2md/examples/martini", "https://github.com/davecheney/godoc2md/tree/master/examples/martini"},
		{"bitbucket.org/atlassianlabs/bitbucket-golang-base", "https://bitbucket.org/atlassianlabs/bitbucket-golang-base/src/master"},
		{"bitbucket.org/atlassianlabs/bitbucket-golang-base/util", "https://bitbucket.org/atlassianlabs/bitbucket-golang-base/src/master/util"},
		{"time", "https://golang.org/src/time"},
		{"go/build", "https://golang.org/src/go/build"},
		{"go/build/build.go", "https://golang.org/src/go/build/build.go"},
		{"encoding/json", "https://golang.org/src/encoding/json"},
		{"golang.org/x/tools/godoc", "https://github.com/golang/tools/tree/master/godoc"},
		{"example.com/myuser/myrepo", "https://example.com/myuser/myrepo/src"},
	}
	for n, tt := range testData {
		got := urlFromPackage(tt.pkg)
		if got != tt.expected {
			t.Errorf("urlFromPackage(%d): expected %s, got %s", n, tt.expected, got)
		}
	}
}
