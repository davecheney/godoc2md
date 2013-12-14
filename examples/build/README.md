
# build

Package build gathers information about Go packages.

<h3 id="hdr-Go_Path">Go Path</h3>
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

<h3 id="hdr-Build_Constraints">Build Constraints</h3>
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

<pre>var ToolDir = filepath.Join(runtime.GOROOT(), &#34;pkg/tool/&#34;+runtime.GOOS+&#34;_&#34;+runtime.GOARCH)</pre>
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
    GOARCH      string <span class="comment">// target architecture</span>
    GOOS        string <span class="comment">// target operating system</span>
    GOROOT      string <span class="comment">// Go root</span>
    GOPATH      string <span class="comment">// Go path</span>
    CgoEnabled  bool   <span class="comment">// whether cgo can be used</span>
    UseAllFiles bool   <span class="comment">// use files regardless of +build lines, file names</span>
    Compiler    string <span class="comment">// compiler to assume when computing target paths</span>

    <span class="comment">// The build and release tags specify build constraints</span>
    <span class="comment">// that should be considered satisfied when processing +build lines.</span>
    <span class="comment">// Clients creating a new context may customize BuildTags, which</span>
    <span class="comment">// defaults to empty, but it is usually an error to customize ReleaseTags,</span>
    <span class="comment">// which defaults to the list of Go releases the current release is compatible with.</span>
    <span class="comment">// In addition to the BuildTags and ReleaseTags, build constraints</span>
    <span class="comment">// consider the values of GOARCH and GOOS as satisfied tags.</span>
    BuildTags   []string
    ReleaseTags []string

    <span class="comment">// The install suffix specifies a suffix to use in the name of the installation</span>
    <span class="comment">// directory. By default it is empty, but custom builds that need to keep</span>
    <span class="comment">// their outputs separate can set InstallSuffix to do so. For example, when</span>
    <span class="comment">// using the race detector, the go command uses InstallSuffix = &#34;race&#34;, so</span>
    <span class="comment">// that on a Linux/386 system, packages are written to a directory named</span>
    <span class="comment">// &#34;linux_386_race&#34; instead of the usual &#34;linux_386&#34;.</span>
    InstallSuffix string

    <span class="comment">// JoinPath joins the sequence of path fragments into a single path.</span>
    <span class="comment">// If JoinPath is nil, Import uses filepath.Join.</span>
    JoinPath func(elem ...string) string

    <span class="comment">// SplitPathList splits the path list into a slice of individual paths.</span>
    <span class="comment">// If SplitPathList is nil, Import uses filepath.SplitList.</span>
    SplitPathList func(list string) []string

    <span class="comment">// IsAbsPath reports whether path is an absolute path.</span>
    <span class="comment">// If IsAbsPath is nil, Import uses filepath.IsAbs.</span>
    IsAbsPath func(path string) bool

    <span class="comment">// IsDir reports whether the path names a directory.</span>
    <span class="comment">// If IsDir is nil, Import calls os.Stat and uses the result&#39;s IsDir method.</span>
    IsDir func(path string) bool

    <span class="comment">// HasSubdir reports whether dir is a subdirectory of</span>
    <span class="comment">// (perhaps multiple levels below) root.</span>
    <span class="comment">// If so, HasSubdir sets rel to a slash-separated path that</span>
    <span class="comment">// can be joined to root to produce a path equivalent to dir.</span>
    <span class="comment">// If HasSubdir is nil, Import uses an implementation built on</span>
    <span class="comment">// filepath.EvalSymlinks.</span>
    HasSubdir func(root, dir string) (rel string, ok bool)

    <span class="comment">// ReadDir returns a slice of os.FileInfo, sorted by Name,</span>
    <span class="comment">// describing the content of the named directory.</span>
    <span class="comment">// If ReadDir is nil, Import uses ioutil.ReadDir.</span>
    ReadDir func(dir string) (fi []os.FileInfo, err error)

    <span class="comment">// OpenFile opens a file (not a directory) for reading.</span>
    <span class="comment">// If OpenFile is nil, Import uses os.Open.</span>
    OpenFile func(path string) (r io.ReadCloser, err error)
}</pre>
A Context specifies the supporting context for a build.






<pre>var Default Context = defaultContext()</pre>
Default is the default Context for builds.
It uses the GOARCH, GOOS, GOROOT, and GOPATH environment variables
if set, or else the compiled code's GOARCH, GOOS, and GOROOT.










### func (*Context) Import
<pre>func (ctxt *Context) Import(path string, srcDir string, mode ImportMode) (*Package, error)</pre>
<p>
Import returns details about the Go package named by the import path,
interpreting local import paths relative to the srcDir directory.
If the path is a local import path naming a package that can be imported
using a standard import path, the returned package will set p.ImportPath
to that path.
</p>
<p>
In the directory containing the package, .go, .c, .h, and .s files are
considered part of the package except for:
</p>
<pre>- .go files in package documentation
- files starting with _ or . (likely editor temporary files)
- files with build constraints not satisfied by the context
</pre>
<p>
If an error occurs, Import returns a non-nil error and a non-nil
*Package containing partial information.
</p>





### func (*Context) ImportDir
<pre>func (ctxt *Context) ImportDir(dir string, mode ImportMode) (*Package, error)</pre>
<p>
ImportDir is like Import but processes the Go package found in
the named directory.
</p>





### func (*Context) MatchFile
<pre>func (ctxt *Context) MatchFile(dir, name string) (match bool, err error)</pre>
<p>
MatchFile reports whether the file with the given name in the given directory
matches the context and would be included in a Package created by ImportDir
of that directory.
</p>
<p>
MatchFile considers the name of the file and may use ctxt.OpenFile to
read some or all of the file&#39;s content.
</p>





### func (*Context) SrcDirs
<pre>func (ctxt *Context) SrcDirs() []string</pre>
<p>
SrcDirs returns a list of package source root directories.
It draws from the current Go root and Go path but omits directories
that do not exist.
</p>







## type ImportMode
<pre>type ImportMode uint</pre>
An ImportMode controls the behavior of the Import method.




<pre>const (
    <span class="comment">// If FindOnly is set, Import stops after locating the directory</span>
    <span class="comment">// that should contain the sources for a package.  It does not</span>
    <span class="comment">// read any files in the directory.</span>
    FindOnly ImportMode = 1 &lt;&lt; iota

    <span class="comment">// If AllowBinary is set, Import can be satisfied by a compiled</span>
    <span class="comment">// package object without corresponding sources.</span>
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
<pre>func (e *NoGoError) Error() string</pre>







## type Package
<pre>type Package struct {
    Dir         string   <span class="comment">// directory containing package sources</span>
    Name        string   <span class="comment">// package name</span>
    Doc         string   <span class="comment">// documentation synopsis</span>
    ImportPath  string   <span class="comment">// import path of package (&#34;&#34; if unknown)</span>
    Root        string   <span class="comment">// root of Go tree where this package lives</span>
    SrcRoot     string   <span class="comment">// package source root directory (&#34;&#34; if unknown)</span>
    PkgRoot     string   <span class="comment">// package install root directory (&#34;&#34; if unknown)</span>
    BinDir      string   <span class="comment">// command install directory (&#34;&#34; if unknown)</span>
    Goroot      bool     <span class="comment">// package found in Go root</span>
    PkgObj      string   <span class="comment">// installed .a file</span>
    AllTags     []string <span class="comment">// tags that can influence file selection in this directory</span>
    ConflictDir string   <span class="comment">// this directory shadows Dir in $GOPATH</span>

    <span class="comment">// Source files</span>
    GoFiles        []string <span class="comment">// .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)</span>
    CgoFiles       []string <span class="comment">// .go source files that import &#34;C&#34;</span>
    IgnoredGoFiles []string <span class="comment">// .go source files ignored for this build</span>
    CFiles         []string <span class="comment">// .c source files</span>
    CXXFiles       []string <span class="comment">// .cc, .cpp and .cxx source files</span>
    HFiles         []string <span class="comment">// .h, .hh, .hpp and .hxx source files</span>
    SFiles         []string <span class="comment">// .s source files</span>
    SwigFiles      []string <span class="comment">// .swig files</span>
    SwigCXXFiles   []string <span class="comment">// .swigcxx files</span>
    SysoFiles      []string <span class="comment">// .syso system object files to add to archive</span>

    <span class="comment">// Cgo directives</span>
    CgoCFLAGS    []string <span class="comment">// Cgo CFLAGS directives</span>
    CgoCPPFLAGS  []string <span class="comment">// Cgo CPPFLAGS directives</span>
    CgoCXXFLAGS  []string <span class="comment">// Cgo CXXFLAGS directives</span>
    CgoLDFLAGS   []string <span class="comment">// Cgo LDFLAGS directives</span>
    CgoPkgConfig []string <span class="comment">// Cgo pkg-config directives</span>

    <span class="comment">// Dependency information</span>
    Imports   []string                    <span class="comment">// imports from GoFiles, CgoFiles</span>
    ImportPos map[string][]token.Position <span class="comment">// line information for Imports</span>

    <span class="comment">// Test information</span>
    TestGoFiles    []string                    <span class="comment">// _test.go files in package</span>
    TestImports    []string                    <span class="comment">// imports from TestGoFiles</span>
    TestImportPos  map[string][]token.Position <span class="comment">// line information for TestImports</span>
    XTestGoFiles   []string                    <span class="comment">// _test.go files outside package</span>
    XTestImports   []string                    <span class="comment">// imports from XTestGoFiles</span>
    XTestImportPos map[string][]token.Position <span class="comment">// line information for XTestImports</span>
}</pre>
A Package describes the Go package found in a directory.











### func Import
<pre>func Import(path, srcDir string, mode ImportMode) (*Package, error)</pre>
Import is shorthand for Default.Import.





### func ImportDir
<pre>func ImportDir(dir string, mode ImportMode) (*Package, error)</pre>
ImportDir is shorthand for Default.ImportDir.







### func (*Package) IsCommand
<pre>func (p *Package) IsCommand() bool</pre>
<p>
IsCommand reports whether the package is considered a
command to be installed (not just a library).
Packages named &#34;main&#34; are treated as commands.
</p>











- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)