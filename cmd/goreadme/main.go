// goreadme creates a README.md from your godoc
//
//
// Usage
//
//    goreadme [-out path/to/README.md] $PACKAGE
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/WillAbides/godoc2md"
)

//go:generate ../../bin/goreadme github.com/WillAbides/godoc2md/cmd/goreadme

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: goreadme [options] package\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	out := flag.String("out", filepath.FromSlash("./README.md"), "path to README.md")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	pkgName := flag.Arg(0)

	var buf bytes.Buffer
	config := &godoc2md.Config{
		TabWidth:          4,
		DeclLinks:         true,
		Goroot:            runtime.GOROOT(),
		SrcLinkHashFormat: "#L%d",
	}

	godoc2md.Godoc2md([]string{pkgName}, &buf, config)
	mdContent := buf.String()
	mdContent = strings.Replace(mdContent, `/src/target/`, `./`, -1)
	mdContent = strings.Replace(mdContent, fmt.Sprintf("/src/%s/", pkgName), `./`, -1)

	err := ioutil.WriteFile(*out, []byte(mdContent), 0640)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed writing to %s\n", *out)
		os.Exit(1)
	}
}
