package main

import (
	"bytes"
	"go/doc"
	"go/printer"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/godoc"
)

type example struct {
	Name   string
	Doc    string
	Code   string
	Output string
}

func examplesFunc(info *godoc.PageInfo, name string) []*example {
	if !*showExamples {
		return nil
	}
	var egs []*example
	for _, eg := range info.Examples {
		if name != "*" && stripExampleSuffix(eg.Name) != name {
			continue
		}
		doc := eg.Doc
		out := eg.Output
		code, wholeFile := exampleCode(info, eg)
		if wholeFile {
			doc = ""
			out = ""
		}
		egs = append(egs, &example{
			Name:   eg.Name,
			Doc:    doc,
			Code:   code,
			Output: out,
		})
	}
	sort.Slice(egs, func(i int, j int) bool {
		ni, si := splitExampleName(egs[i].Name)
		nj, sj := splitExampleName(egs[j].Name)
		if ni == nj {
			return si < sj
		}
		return ni < nj
	})
	return egs
}

var exampleOutputRx = regexp.MustCompile(`(?i)//[[:space:]]*(unordered )?output:`)

func exampleCode(info *godoc.PageInfo, eg *doc.Example) (code string, wholeFile bool) {
	// Print code
	var buf bytes.Buffer
	cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
	config := &printer.Config{Mode: printer.UseSpaces, Tabwidth: *tabWidth}
	config.Fprint(&buf, info.FSet, cnode)
	code = strings.Trim(buf.String(), "\n")
	wholeFile = true

	if n := len(code); n >= 2 && code[0] == '{' && code[n-1] == '}' {
		wholeFile = false
		// Remove surrounding braces.
		code = strings.Trim(code[1:n-1], "\n")
		// Remove output from code.
		if loc := exampleOutputRx.FindStringIndex(code); loc != nil {
			code = strings.TrimRightFunc(code[:loc[0]], unicode.IsSpace)
		}
		// Unindent code.
		lines := strings.Split(code, "\n")
		unindent(lines)
		code = strings.Join(lines, "\n")
	}

	return code, wholeFile
}

func splitExampleName(s string) (name, suffix string) {
	i := strings.LastIndex(s, "_")
	if 0 <= i && i < len(s)-1 && !startsWithUppercase(s[i+1:]) {
		name = s[:i]
		suffix = " (" + strings.Title(s[i+1:]) + ")"
		return
	}
	name = s
	return
}

// stripExampleSuffix strips lowercase braz in Foo_braz or Foo_Bar_braz from name
// while keeping uppercase Braz in Foo_Braz.
func stripExampleSuffix(name string) string {
	if i := strings.LastIndex(name, "_"); i != -1 {
		if i < len(name)-1 && !startsWithUppercase(name[i+1:]) {
			name = name[:i]
		}
	}
	return name
}

func startsWithUppercase(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
