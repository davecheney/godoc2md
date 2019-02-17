// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// godoc2md converts godoc formatted package documentation into Markdown format.
//
//
// Usage
//
//    godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/WillAbides/godoc2md"
)

var (
	verbose = flag.Bool("v", false, "verbose mode")

	// file system roots
	// TODO(gri) consider the invariant that goroot always end in '/'
	goroot = flag.String("goroot", runtime.GOROOT(), "Go root directory")

	// layout control
	tabWidth       = flag.Int("tabwidth", 4, "tab width")
	showTimestamps = flag.Bool("timestamps", false, "show timestamps with directory listings")
	altPkgTemplate = flag.String("template", "", "path to an alternate template file")
	showPlayground = flag.Bool("play", false, "enable playground in web interface")
	showExamples   = flag.Bool("ex", false, "show examples in command line mode")
	declLinks      = flag.Bool("links", true, "link identifiers to their declarations")

	// The hash format for Github is the default `#L%d`; but other source control platforms do not
	// use the same format. For example Bitbucket Enterprise uses `#%d`. This option provides the
	// user the option to switch the format as needed and still remain backwards compatible.
	srcLinkHashFormat = flag.String("hashformat", "#L%d", "source link URL hash format")

	srcLinkFormat = flag.String("srclink", "", "if set, format for entire source link")
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: godoc2md [options] package\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	pkgName := flag.Arg(0)

	config := &godoc2md.Config{
		TabWidth:          *tabWidth,
		ShowTimestamps:    *showTimestamps,
		AltPkgTemplate:    *altPkgTemplate,
		ShowPlayground:    *showPlayground,
		ShowExamples:      *showExamples,
		DeclLinks:         *declLinks,
		Goroot:            *goroot,
		SrcLinkHashFormat: *srcLinkHashFormat,
		SrcLinkFormat:     *srcLinkFormat,
		Verbose:           *verbose,
	}

	godoc2md.Godoc2md([]string{pkgName}, os.Stdout, config)
}
