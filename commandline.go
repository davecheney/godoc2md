package godoc2md

import (
	"fmt"
	"go/ast"
	"go/build"
	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"io"
	"log"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	target         = "/target"
	cmdPrefix      = "cmd/"
	srcPrefix      = "src/"
	toolsPath      = "golang.org/x/tools/cmd/"
	builtinPkgPath = "builtin"
)

// commandLine returns godoc results to w.
// Note that it may add a /target path to fs.
func commandLine(w io.Writer, fs vfs.NameSpace, pres *godoc.Presentation, tmpl *template.Template, args []string) error {
	path := args[0]
	srcMode := pres.SrcMode
	cmdMode := strings.HasPrefix(path, cmdPrefix)
	if strings.HasPrefix(path, srcPrefix) {
		path = strings.TrimPrefix(path, srcPrefix)
		srcMode = true
	}
	var abspath, relpath string
	if cmdMode {
		path = strings.TrimPrefix(path, cmdPrefix)
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
			abspath = pathpkg.Join(pres.CmdFSRoot(), cmdPrefix, path)
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
		} else if relpath != target {
			// The above check handles the case where an operating system path
			// is provided (see documentation for paths below).  In that case,
			// relpath is set to "/target" (in anticipation of accessing packages there),
			// and is therefore not expected to match a command.
			fmt.Fprintf(w, "use 'godoc %s%s' for documentation on the %s command \n\n", cmdPrefix, relpath, relpath)
		}
	}

	if info == nil {
		return fmt.Errorf("%s: no such directory or package", args[0])
	}
	if info.Err != nil {
		return info.Err
	}

	if info.PDoc != nil && info.PDoc.ImportPath == target {
		// Replace virtual /target with actual argument from command line.
		info.PDoc.ImportPath = args[0]
	}

	// If we have more than one argument, use the remaining arguments for filtering.
	if len(args) > 1 {
		info.IsFiltered = true
		filterInfo(args[1:], info)
	}

	if err := tmpl.Execute(w, info); err != nil {
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
		fs.Bind(target, vfs.OS(path), "/", vfs.BindReplace)
		return target, target
	}
	if build.IsLocalImport(path) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("error while getting working directory: %v", err)
		}
		path = filepath.Join(cwd, path)
		fs.Bind(target, vfs.OS(path), "/", vfs.BindReplace)
		return target, target
	}
	bp, err := build.Import(path, "", build.FindOnly)
	if err != nil {
		log.Printf("error while importing build package: %v", err)
	}
	if bp.Dir != "" && bp.ImportPath != "" {
		fs.Bind(target, vfs.OS(bp.Dir), "/", vfs.BindReplace)
		return target, bp.ImportPath
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

	filter := func(s string) bool {
		fmt.Fprintf(os.Stderr, "filtering on string: %s\n", s)

		return rx.MatchString(s)
	}
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
	fmt.Fprintf(os.Stderr, "regex string: %s\n", s)
	return regexp.Compile(s)
}
