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
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
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

const (
	targetPath     = "/target"
	cmdPathPrefix  = "cmd/"
	srcPathPrefix  = "src/"
	toolsPath      = "golang.org/x/tools/cmd/"
	builtinPkgPath = "builtin"
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: godoc2md package [name ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	pres *godoc.Presentation
	fs   = vfs.NameSpace{}

	funcs = map[string]interface{}{
		"comment_md":  commentMdFunc,
		"base":        pathpkg.Base,
		"md":          mdFunc,
		"pre":         preFunc,
		"kebab":       kebabFunc,
		"bitscape":    bitscapeFunc, //Escape [] for bitbucket confusion
		"trim_prefix": strings.TrimPrefix,
	}
)

func commentMdFunc(comment string) string {
	var buf bytes.Buffer
	ToMD(&buf, comment)
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
	s = pathpkg.Clean("/" + s)
	if !strings.HasPrefix(s, "/src/") {
		s = "/src" + s
	}
	return s
}

// Removed code line that always substracted 10 from the value of `line`.
// Made format for the source link hash configurable to support source control platforms other than Github.
// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L540
func srcPosLinkFunc(s string, line, low, high int) string {
	if *srcLinkFormat != "" {
		return fmt.Sprintf(*srcLinkFormat, s, line, low, high)
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
		fmt.Fprintf(&buf, *srcLinkHashFormat, line) // no need for URL escaping
	}
	return buf.String()
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

func main() {
	flag.Usage = usage
	flag.Parse()

	// Check usage
	if flag.NArg() == 0 {
		usage()
	}

	// use file system of underlying OS
	fs.Bind("/", vfs.OS(*goroot), "/", vfs.BindReplace)

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src/pkg", vfs.OS(p), "/src", vfs.BindAfter)
	}

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *verbose

	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = *tabWidth
	pres.ShowTimestamps = *showTimestamps
	pres.ShowPlayground = *showPlayground
	pres.DeclLinks = *declLinks
	pres.SrcMode = false
	pres.HTMLMode = false
	pres.URLForSrcPos = srcPosLinkFunc

	var tmpl *template.Template

	if *altPkgTemplate != "" {
		buf, err := ioutil.ReadFile(*altPkgTemplate)
		if err != nil {
			log.Fatal(err)
		}
		tmpl = readTemplate("package.txt", string(buf))
	} else {
		tmpl = readTemplate("package.txt", pkgTemplate)
	}

	if err := writeOutput(os.Stdout, fs, pres, flag.Args(), tmpl); err != nil {
		log.Print(err)
	}
}

// writeOutpur returns godoc results to w.
// Note that it may add a /target path to fs.
func writeOutput(w io.Writer, fs vfs.NameSpace, pres *godoc.Presentation, args []string, packageText *template.Template) error {
	path := args[0]
	srcMode := pres.SrcMode
	cmdMode := strings.HasPrefix(path, cmdPathPrefix)
	if strings.HasPrefix(path, srcPathPrefix) {
		path = strings.TrimPrefix(path, srcPathPrefix)
		srcMode = true
	}
	var abspath, relpath string
	if cmdMode {
		path = strings.TrimPrefix(path, cmdPathPrefix)
	} else {
		abspath, relpath = paths(fs, pres, path)
	}

	var mode godoc.PageInfoMode
	if relpath == builtinPkgPath {
		// the fake built-in package contains unexported identifiers
		mode = godoc.NoFiltering | godoc.NoTypeAssoc
	}
	if pres.AllMode {
		mode |= godoc.NoFiltering
	}
	if srcMode {
		// only filter exports if we don't have explicit command-line filter arguments
		if len(args) > 1 {
			mode |= godoc.NoFiltering
		}
		mode |= godoc.ShowSource
	}

	// First, try as package unless forced as command.
	var info *godoc.PageInfo
	if !cmdMode {
		info = pres.GetPkgPageInfo(abspath, relpath, mode)
	}

	// Second, try as command (if the path is not absolute).
	var cinfo *godoc.PageInfo
	if !filepath.IsAbs(path) {
		// First try go.tools/cmd.
		abspath = pathpkg.Join(pres.PkgFSRoot(), toolsPath+path)
		cinfo = pres.GetCmdPageInfo(abspath, relpath, mode)
		if cinfo.IsEmpty() {
			// Then try $GOROOT/src/cmd.
			abspath = pathpkg.Join(pres.CmdFSRoot(), cmdPathPrefix, path)
			cinfo = pres.GetCmdPageInfo(abspath, relpath, mode)
		}
	}

	// determine what to use
	if info == nil || info.IsEmpty() {
		if cinfo != nil && !cinfo.IsEmpty() {
			// only cinfo exists - switch to cinfo
			info = cinfo
		}
	} else if cinfo != nil && !cinfo.IsEmpty() {
		// both info and cinfo exist - use cinfo if info
		// contains only subdirectory information
		if info.PAst == nil && info.PDoc == nil {
			info = cinfo
		} else if relpath != targetPath {
			// The above check handles the case where an operating system path
			// is provided (see documentation for paths below).  In that case,
			// relpath is set to "/target" (in anticipation of accessing packages there),
			// and is therefore not expected to match a command.
			fmt.Fprintf(w, "use 'godoc %s%s' for documentation on the %s command \n\n", cmdPathPrefix, relpath, relpath)
		}
	}

	if info == nil {
		return fmt.Errorf("%s: no such directory or package", args[0])
	}
	if info.Err != nil {
		return info.Err
	}

	if info.PDoc != nil && info.PDoc.ImportPath == targetPath {
		// Replace virtual /target with actual argument from command line.
		info.PDoc.ImportPath = args[0]
	}

	// If we have more than one argument, use the remaining arguments for filtering.
	if len(args) > 1 {
		info.IsFiltered = true
		filterInfo(args[1:], info)
	}

	if err := packageText.Execute(w, info); err != nil {
		return err
	}
	return nil
}

// paths determines the paths to use.
//
// If we are passed an operating system path like . or ./foo or /foo/bar or c:\mysrc,
// we need to map that path somewhere in the fs name space so that routines
// like getPageInfo will see it.  We use the arbitrarily-chosen virtual path "/target"
// for this.  That is, if we get passed a directory like the above, we map that
// directory so that getPageInfo sees it as /target.
// Returns the absolute and relative paths.
func paths(fs vfs.NameSpace, pres *godoc.Presentation, path string) (abspath, relpath string) {
	if filepath.IsAbs(path) {
		fs.Bind(targetPath, vfs.OS(path), "/", vfs.BindReplace)
		return targetPath, targetPath
	}
	if build.IsLocalImport(path) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("error while getting working directory: %v", err)
		}
		path = filepath.Join(cwd, path)
		fs.Bind(targetPath, vfs.OS(path), "/", vfs.BindReplace)
		return targetPath, targetPath
	}
	bp, err := build.Import(path, "", build.FindOnly)
	if err != nil {
		log.Printf("error while importing build package: %v", err)
	}
	if bp.Dir != "" && bp.ImportPath != "" {
		fs.Bind(targetPath, vfs.OS(bp.Dir), "/", vfs.BindReplace)
		return targetPath, bp.ImportPath
	}
	return pathpkg.Join(pres.PkgFSRoot(), path), path
}

// filterInfo updates info to include only the nodes that match the given
// filter args.
func filterInfo(args []string, info *godoc.PageInfo) {
	rx, err := makeRx(args)
	if err != nil {
		log.Fatalf("illegal regular expression from %v: %v", args, err)
	}

	filter := func(s string) bool { return rx.MatchString(s) }
	switch {
	case info.PAst != nil:
		newPAst := map[string]*ast.File{}
		for name, a := range info.PAst {
			cmap := ast.NewCommentMap(info.FSet, a, a.Comments)
			a.Comments = []*ast.CommentGroup{} // remove all comments.
			ast.FilterFile(a, filter)
			if len(a.Decls) > 0 {
				newPAst[name] = a
			}
			for _, d := range a.Decls {
				// add back the comments associated with d only
				comments := cmap.Filter(d).Comments()
				a.Comments = append(a.Comments, comments...)
			}
		}
		info.PAst = newPAst // add only matching files.
	case info.PDoc != nil:
		info.PDoc.Filter(filter)
	}
}

// Does s look like a regular expression?
func isRegexp(s string) bool {
	return strings.ContainsAny(s, ".(|)*+?^$[]")
}

// Make a regular expression of the form
// names[0]|names[1]|...names[len(names)-1].
// Returns an error if the regular expression is illegal.
func makeRx(names []string) (*regexp.Regexp, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no expression provided")
	}
	s := ""
	for i, name := range names {
		if i > 0 {
			s += "|"
		}
		if isRegexp(name) {
			s += name
		} else {
			s += "^" + name + "$" // must match exactly
		}
	}
	return regexp.Compile(s)
}
