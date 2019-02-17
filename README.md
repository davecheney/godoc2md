

# godoc2md
`import "github.com/WillAbides/godoc2md"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package godoc2md creates a markdown representation of a package's godoc.

This is forked from <a href="https://github.com/davecheney/godoc2md">https://github.com/davecheney/godoc2md</a>.  The primary difference being that this version is
a library that can be used by other packages.




## <a name="pkg-index">Index</a>
* [func Godoc2md(args []string, out io.Writer, config *Config)](#Godoc2md)
* [func ToMD(w io.Writer, text string)](#ToMD)
* [type Config](#Config)


#### <a name="pkg-files">Package files</a>
[commandline.go](/src/github.com/WillAbides/godoc2md/commandline.go) [comment.go](/src/github.com/WillAbides/godoc2md/comment.go) [main.go](/src/github.com/WillAbides/godoc2md/main.go) [template.go](/src/github.com/WillAbides/godoc2md/template.go) 





## <a name="Godoc2md">func</a> [Godoc2md](/src/target/main.go?s=3426:3485#L132)
``` go
func Godoc2md(args []string, out io.Writer, config *Config)
```
Godoc2md turns your godoc into markdown



## <a name="ToMD">func</a> [ToMD](/src/target/comment.go?s=4298:4333#L194)
``` go
func ToMD(w io.Writer, text string)
```
ToMD converts comment text to formatted Markdown.
The comment was prepared by DocReader,
so it is known not to have leading, trailing blank lines
nor to have trailing spaces at the end of lines.
The comment markers have already been removed.

Each span of unindented non-blank lines is converted into
a single paragraph. There is one exception to the rule: a span that
consists of a single line, is followed by another paragraph span,
begins with a capital letter, and contains no punctuation
is formatted as a heading.

A span of indented lines is converted into a <pre> block,
with the common indent prefix removed.

URLs in the comment text are converted into links.




## <a name="Config">type</a> [Config](/src/target/main.go?s=985:1254#L43)
``` go
type Config struct {
    TabWidth          int
    ShowTimestamps    bool
    AltPkgTemplate    string
    ShowPlayground    bool
    ShowExamples      bool
    DeclLinks         bool
    SrcLinkHashFormat string
    SrcLinkFormat     string
    Goroot            string
    Verbose           bool
}

```
Config contains config options for Godoc2md














