

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
* [type Config](#Config)


#### <a name="pkg-files">Package files</a>
[commandline.go](/src/github.com/WillAbides/godoc2md/commandline.go) [comment.go](/src/github.com/WillAbides/godoc2md/comment.go) [godoc2md.go](/src/github.com/WillAbides/godoc2md/godoc2md.go) [template.go](/src/github.com/WillAbides/godoc2md/template.go) 





## <a name="Godoc2md">func</a> [Godoc2md](/src/target/godoc2md.go?s=3439:3498#L132)
``` go
func Godoc2md(args []string, out io.Writer, config *Config)
```
Godoc2md turns your godoc into markdown




## <a name="Config">type</a> [Config](/src/target/godoc2md.go?s=985:1254#L43)
``` go
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

```
Config contains config options for Godoc2md














