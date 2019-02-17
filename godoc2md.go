// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package godoc2md creates a markdown representation of a package's godoc.
//
// This is forked from https://github.com/davecheney/godoc2md.  The primary difference being that this version is
// a library that can be used by other packages.
package godoc2md

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
)

var (
	pres *godoc.Presentation
	fs   = vfs.NameSpace{}

	funcs = map[string]interface{}{
		"comment_md":  commentMdFunc,
		"base":        path.Base,
		"md":          mdFunc,
		"pre":         preFunc,
		"kebab":       kebabFunc,
		"bitscape":    bitscapeFunc, //Escape [] for bitbucket confusion
		"trim_prefix": strings.TrimPrefix,
	}
)

//Config contains config options for Godoc2md
type Config struct {
	AltPkgTemplate    string
	SrcLinkHashFormat string
	SrcLinkFormat     string
	Goroot            string
	TabWidth          int
	ShowTimestamps    bool
	ShowPlayground    bool
	ShowExamples      bool
	DeclLinks         bool
	Verbose           bool
}

func commentMdFunc(comment string) string {
	var buf bytes.Buffer
	toMd(&buf, comment)
	return buf.String()
}

func mdFunc(text string) string {
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)
	return text
}

func preFunc(text string) string {
	return "``` go\n" + text + "\n```"
}

// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L562
func srcLinkFunc(s string) string {
	s = path.Clean("/" + s)
	if !strings.HasPrefix(s, "/src/") {
		s = "/src" + s
	}
	return s
}

// Removed code line that always substracted 10 from the value of `line`.
// Made format for the source link hash configurable to support source control platforms other than Github.
// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L540
func genSrcPosLinkFunc(srcLinkFormat, srcLinkHashFormat string) func(s string, line, low, high int) string {
	return func(s string, line, low, high int) string {
		if srcLinkFormat != "" {
			return fmt.Sprintf(srcLinkFormat, s, line, low, high)
		}

		s = srcLinkFunc(s)
		var buf bytes.Buffer
		template.HTMLEscape(&buf, []byte(s))
		// selection ranges are of form "s=low:high"
		if low < high {
			fmt.Fprintf(&buf, "?s=%d:%d", low, high) // no need for URL escaping
			if line < 1 {
				line = 1
			}
		}
		// line id's in html-printed source are of the
		// form "L%d" (on Github) where %d stands for the line number
		if line > 0 {
			fmt.Fprintf(&buf, srcLinkHashFormat, line) // no need for URL escaping
		}
		return buf.String()
	}
}

func readTemplate(name, data string) *template.Template {
	// be explicit with errors (for app engine use)
	t, err := template.New(name).Funcs(pres.FuncMap()).Funcs(funcs).Parse(data)
	if err != nil {
		log.Fatal("readTemplate: ", err)
	}
	return t
}

func kebabFunc(text string) string {
	s := strings.Replace(strings.ToLower(text), " ", "-", -1)
	s = strings.Replace(s, ".", "-", -1)
	s = strings.Replace(s, "\\*", "42", -1)
	return s
}

func bitscapeFunc(text string) string {
	s := strings.Replace(text, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	return s
}

//Godoc2md turns your godoc into markdown
func Godoc2md(args []string, out io.Writer, config *Config) {
	// use file system of underlying OS
	fs.Bind("/", vfs.OS(config.Goroot), "/", vfs.BindReplace)
	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src/pkg", vfs.OS(p), "/src", vfs.BindAfter)
	}
	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = config.Verbose
	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = config.TabWidth
	pres.ShowTimestamps = config.ShowTimestamps
	pres.ShowPlayground = config.ShowPlayground
	pres.DeclLinks = config.DeclLinks
	pres.SrcMode = false
	pres.URLForSrcPos = genSrcPosLinkFunc(config.SrcLinkFormat, config.SrcLinkHashFormat)
	var tmpl *template.Template
	if config.AltPkgTemplate != "" {
		buf, err := ioutil.ReadFile(config.AltPkgTemplate)
		if err != nil {
			log.Fatal(err)
		}
		tmpl = readTemplate("package.txt", string(buf))
	} else {
		tmpl = readTemplate("package.txt", pkgTemplate)
	}
	if err := commandLine(out, fs, pres, tmpl, args); err != nil {
		log.Print(err)
	}
}
