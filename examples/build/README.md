
# build
    import "go/build"

Package build gathers information about Go packages.

### Go Path
The Go path is a list of directory trees containing Go source code.
It is consulted to resolve imports that cannot be found in the standard
Go tree.  The default path is the value of the GOPATH environment
variable, interpreted as a path list appropriate to the operating system
(on Unix, the variable is a colon-separated string;
on Windows, a semicolon-separated string;
on Plan 9, a list).

Each directory listed in the Go path must have a prescribed structure:

The src/ directory holds source code.  The path below 'src' determines
the import path or executable name.

The pkg/ directory holds installed package objects.
As in the Go tree, each target operating system and
architecture pair has its own subdirectory of pkg
(pkg/GOOS_GOARCH).

If DIR is a directory listed in the Go path, a package with
source in DIR/src/foo/bar can be imported as "foo/bar" and
has its compiled form installed to "DIR/pkg/GOOS_GOARCH/foo/bar.a"
(or, for gccgo, "DIR/pkg/gccgo/foo/libbar.a").

The bin/ directory holds compiled commands.
Each command is named for its source directory, but only
using the final element, not the entire path.  That is, the
command with source in DIR/src/foo/quux is installed into
DIR/bin/quux, not DIR/bin/foo/quux.  The foo/ is stripped
so that you can add DIR/bin to your PATH to get at the
installed commands.

Here's an example directory layout:


	GOPATH=/home/user/gocode
	
	/home/user/gocode/
	    src/
	        foo/
	            bar/               (go code in package bar)
	                x.go
	            quux/              (go code in package main)
	                y.go
	    bin/
	        quux                   (installed command)
	    pkg/
	        linux_amd64/
	            foo/
	                bar.a          (installed package object)

### Build Constraints
A build constraint is a line comment beginning with the directive +build
that lists the conditions under which a file should be included in the package.
Constraints may appear in any kind of source file (not just Go), but
they must appear near the top of the file, preceded
only by blank lines and other line comments.

To distinguish build constraints from package documentation, a series of
build constraints must be followed by a blank line.

A build constraint is evaluated as the OR of space-separated options;
each option evaluates as the AND of its comma-separated terms;
and each term is an alphanumeric word or, preceded by !, its negation.
That is, the build constraint:


	// +build linux,386 darwin,!cgo

corresponds to the boolean formula:


	(linux AND 386) OR (darwin AND (NOT cgo))

A file may have multiple build constraints. The overall constraint is the AND
of the individual constraints. That is, the build constraints:


	// +build linux darwin
	// +build 386

corresponds to the boolean formula:


	(linux OR darwin) AND 386

During a particular build, the following words are satisfied:


	- the target operating system, as spelled by runtime.GOOS
	- the target architecture, as spelled by runtime.GOARCH
	- the compiler being used, either "gc" or "gccgo"
	- "cgo", if ctxt.CgoEnabled is true
	- "go1.1", from Go version 1.1 onward
	- "go1.2", from Go version 1.2 onward
	- any additional words listed in ctxt.BuildTags

If a file's name, after stripping the extension and a possible _test suffix,
matches any of the following patterns:


	*_GOOS
	*_GOARCH
	*_GOOS_GOARCH

(example: source_windows_amd64.go) or the literals:


	GOOS
	GOARCH

(example: windows.go) where GOOS and GOARCH represent any known operating
system and architecture values respectively, then the file is considered to
have an implicit build constraint requiring those terms.

To keep a file from being considered for the build:


	// +build ignore

(any other unsatisfied word will work as well, but ``ignore'' is conventional.)

To build a file only when using cgo, and only on Linux and OS X:


	// +build linux,cgo darwin,cgo

Such a file is usually paired with another file implementing the
default functionality for other systems, which in this case would
carry the constraint:


	// +build !linux,!darwin !cgo

Naming a file dns_windows.go will cause it to be included only when
building the package for Windows; similarly, math_386.s will be included
only when building the package for 32-bit x86.






## Variables

<pre>var ToolDir = filepath.Join(runtime.GOROOT(), "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH)</pre>
ToolDir is the directory containing build tools.







## func ArchChar
<pre>func ArchChar(goarch string) (string, error)</pre>
ArchChar returns the architecture character for the given goarch.
For example, ArchChar("amd64") returns "6".






## func IsLocalImport
<pre>func IsLocalImport(path string) bool</pre>
IsLocalImport reports whether the import path is
a local import path, like ".", "..", "./foo", or "../foo".







## type Context
<pre>type Context struct {
    GOARCH      string // target architecture
    GOOS        string // target operating system
    GOROOT      string // Go root
    GOPATH      string // Go path
    CgoEnabled  bool   // whether cgo can be used
    UseAllFiles bool   // use files regardless of +build lines, file names
    Compiler    string // compiler to assume when computing target paths

    // The build and release tags specify build constraints
    // that should be considered satisfied when processing +build lines.
    // Clients creating a new context may customize BuildTags, which
    // defaults to empty, but it is usually an error to customize ReleaseTags,
    // which defaults to the list of Go releases the current release is compatible with.
    // In addition to the BuildTags and ReleaseTags, build constraints
    // consider the values of GOARCH and GOOS as satisfied tags.
    BuildTags   []string
    ReleaseTags []string

    // The install suffix specifies a suffix to use in the name of the installation
    // directory. By default it is empty, but custom builds that need to keep
    // their outputs separate can set InstallSuffix to do so. For example, when
    // using the race detector, the go command uses InstallSuffix = "race", so
    // that on a Linux/386 system, packages are written to a directory named
    // "linux_386_race" instead of the usual "linux_386".
    InstallSuffix string

    // JoinPath joins the sequence of path fragments into a single path.
    // If JoinPath is nil, Import uses filepath.Join.
    JoinPath func(elem ...string) string

    // SplitPathList splits the path list into a slice of individual paths.
    // If SplitPathList is nil, Import uses filepath.SplitList.
    SplitPathList func(list string) []string

    // IsAbsPath reports whether path is an absolute path.
    // If IsAbsPath is nil, Import uses filepath.IsAbs.
    IsAbsPath func(path string) bool

    // IsDir reports whether the path names a directory.
    // If IsDir is nil, Import calls os.Stat and uses the result's IsDir method.
    IsDir func(path string) bool

    // HasSubdir reports whether dir is a subdirectory of
    // (perhaps multiple levels below) root.
    // If so, HasSubdir sets rel to a slash-separated path that
    // can be joined to root to produce a path equivalent to dir.
    // If HasSubdir is nil, Import uses an implementation built on
    // filepath.EvalSymlinks.
    HasSubdir func(root, dir string) (rel string, ok bool)

    // ReadDir returns a slice of os.FileInfo, sorted by Name,
    // describing the content of the named directory.
    // If ReadDir is nil, Import uses ioutil.ReadDir.
    ReadDir func(dir string) (fi []os.FileInfo, err error)

    // OpenFile opens a file (not a directory) for reading.
    // If OpenFile is nil, Import uses os.Open.
    OpenFile func(path string) (r io.ReadCloser, err error)
}</pre>
A Context specifies the supporting context for a build.






<pre>var Default Context = defaultContext()</pre>
Default is the default Context for builds.
It uses the GOARCH, GOOS, GOROOT, and GOPATH environment variables
if set, or else the compiled code's GOARCH, GOOS, and GOROOT.










### func (*Context) Import

    func (ctxt *Context) Import(path string, srcDir string, mode ImportMode) (*Package, error)

Import returns details about the Go package named by the import path,
interpreting local import paths relative to the srcDir directory.
If the path is a local import path naming a package that can be imported
using a standard import path, the returned package will set p.ImportPath
to that path.

In the directory containing the package, .go, .c, .h, and .s files are
considered part of the package except for:


	- .go files in package documentation
	- files starting with _ or . (likely editor temporary files)
	- files with build constraints not satisfied by the context

If an error occurs, Import returns a non-nil error and a non-nil
*Package containing partial information.






### func (*Context) ImportDir

    func (ctxt *Context) ImportDir(dir string, mode ImportMode) (*Package, error)

ImportDir is like Import but processes the Go package found in
the named directory.






### func (*Context) MatchFile

    func (ctxt *Context) MatchFile(dir, name string) (match bool, err error)

MatchFile reports whether the file with the given name in the given directory
matches the context and would be included in a Package created by ImportDir
of that directory.

MatchFile considers the name of the file and may use ctxt.OpenFile to
read some or all of the file's content.






### func (*Context) SrcDirs

    func (ctxt *Context) SrcDirs() []string

SrcDirs returns a list of package source root directories.
It draws from the current Go root and Go path but omits directories
that do not exist.








## type ImportMode
<pre>type ImportMode uint</pre>
An ImportMode controls the behavior of the Import method.




<pre>const (
    // If FindOnly is set, Import stops after locating the directory
    // that should contain the sources for a package.  It does not
    // read any files in the directory.
    FindOnly ImportMode = 1 << iota

    // If AllowBinary is set, Import can be satisfied by a compiled
    // package object without corresponding sources.
    AllowBinary
)</pre>













## type NoGoError
<pre>type NoGoError struct {
    Dir string
}</pre>
NoGoError is the error used by Import to describe a directory
containing no buildable Go source files. (It may still contain
test files, files hidden by build tags, and so on.)













### func (*NoGoError) Error

    func (e *NoGoError) Error() string








## type Package
<pre>type Package struct {
    Dir         string   // directory containing package sources
    Name        string   // package name
    Doc         string   // documentation synopsis
    ImportPath  string   // import path of package ("" if unknown)
    Root        string   // root of Go tree where this package lives
    SrcRoot     string   // package source root directory ("" if unknown)
    PkgRoot     string   // package install root directory ("" if unknown)
    BinDir      string   // command install directory ("" if unknown)
    Goroot      bool     // package found in Go root
    PkgObj      string   // installed .a file
    AllTags     []string // tags that can influence file selection in this directory
    ConflictDir string   // this directory shadows Dir in $GOPATH

    // Source files
    GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
    CgoFiles       []string // .go source files that import "C"
    IgnoredGoFiles []string // .go source files ignored for this build
    CFiles         []string // .c source files
    CXXFiles       []string // .cc, .cpp and .cxx source files
    HFiles         []string // .h, .hh, .hpp and .hxx source files
    SFiles         []string // .s source files
    SwigFiles      []string // .swig files
    SwigCXXFiles   []string // .swigcxx files
    SysoFiles      []string // .syso system object files to add to archive

    // Cgo directives
    CgoCFLAGS    []string // Cgo CFLAGS directives
    CgoCPPFLAGS  []string // Cgo CPPFLAGS directives
    CgoCXXFLAGS  []string // Cgo CXXFLAGS directives
    CgoLDFLAGS   []string // Cgo LDFLAGS directives
    CgoPkgConfig []string // Cgo pkg-config directives

    // Dependency information
    Imports   []string                    // imports from GoFiles, CgoFiles
    ImportPos map[string][]token.Position // line information for Imports

    // Test information
    TestGoFiles    []string                    // _test.go files in package
    TestImports    []string                    // imports from TestGoFiles
    TestImportPos  map[string][]token.Position // line information for TestImports
    XTestGoFiles   []string                    // _test.go files outside package
    XTestImports   []string                    // imports from XTestGoFiles
    XTestImportPos map[string][]token.Position // line information for XTestImports
}</pre>
A Package describes the Go package found in a directory.











### func Import

    func Import(path, srcDir string, mode ImportMode) (*Package, error)

Import is shorthand for Default.Import.





### func ImportDir

    func ImportDir(dir string, mode ImportMode) (*Package, error)

ImportDir is shorthand for Default.ImportDir.







### func (*Package) IsCommand

    func (p *Package) IsCommand() bool

IsCommand reports whether the package is considered a
command to be installed (not just a library).
Packages named "main" are treated as commands.












- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)