package main

import (
	"bytes"
	"fmt"
	"go/printer"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/godoc"
)

func exampleLinkFunc(funcName string) string {
	i := strings.LastIndex(funcName, "_")
	if 0 <= i && i < len(funcName)-1 && !startsWithUppercase(funcName[i+1:]) {
		name := strings.ToLower(funcName[:i])
		suffix := strings.ToLower(funcName[i+1:])
		return fmt.Sprintf("%s-%s", name, suffix)
	}
	return strings.ToLower(funcName)
}

// Based on example_textFunc from
// https://github.com/golang/tools/blob/master/godoc/godoc.go
func exampleMdFunc(info *godoc.PageInfo, funcName string) string {
	if !*showExamples {
		return ""
	}

	var buf bytes.Buffer
	first := true
	for _, eg := range info.Examples {
		name := stripExampleSuffix(eg.Name)
		if name != funcName {
			continue
		}

		if !first {
			buf.WriteString("\n")
		}
		first = false

		// print code
		cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
		config := &printer.Config{Mode: printer.UseSpaces, Tabwidth: pres.TabWidth}
		var buf1 bytes.Buffer
		config.Fprint(&buf1, info.FSet, cnode)
		code := buf1.String()
		output := strings.Trim(eg.Output, "\n")
		output = replaceLeadingIndentation(output, strings.Repeat(" ", pres.TabWidth), "")

		// Additional formatting if this is a function body. Unfortunately, we
		// can't print statements individually because we would lose comments
		// on later statements.
		if n := len(code); n >= 2 && code[0] == '{' && code[n-1] == '}' {
			// remove surrounding braces
			code = code[1 : n-1]
			// unindent
			code = replaceLeadingIndentation(code, strings.Repeat(" ", pres.TabWidth), "")
		}
		code = strings.Trim(code, "\n")
		name, suffix := splitExampleName(eg.Name)
		title := fmt.Sprintf("##### Example %s%s:\n", name, suffix)
		buf.WriteString(title)
		if len(eg.Doc) > 0 {
			buf.WriteString(eg.Doc)
			buf.WriteString("\n")
		}
		buf.WriteString("``` go\n")
		buf.WriteString(code)
		buf.WriteString("\n```\n\n")
		if len(output) > 0 {
			buf.WriteString("Output:\n")
			buf.WriteString("\n```\n")
			buf.WriteString(output)
			buf.WriteString("\n```\n\n")
		}
	}
	return buf.String()
}

// Copy/pasted from https://github.com/golang/tools/blob/master/godoc/godoc.go
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

// Copy/pasted from https://github.com/golang/tools/blob/master/godoc/godoc.go#L786
func stripExampleSuffix(name string) string {
	if i := strings.LastIndex(name, "_"); i != -1 {
		if i < len(name)-1 && !startsWithUppercase(name[i+1:]) {
			name = name[:i]
		}
	}
	return name
}

// Copy/pasted from https://github.com/golang/tools/blob/master/godoc/godoc.go#L777
func startsWithUppercase(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

// Copy/pasted from https://github.com/golang/tools/blob/master/godoc/godoc.go
func replaceLeadingIndentation(body, oldIndent, newIndent string) string {
	// Handle indent at the beginning of the first line. After this, we handle
	// indentation only after a newline.
	var buf bytes.Buffer
	if strings.HasPrefix(body, oldIndent) {
		buf.WriteString(newIndent)
		body = body[len(oldIndent):]
	}

	// Use a state machine to keep track of whether we're in a string or
	// rune literal while we process the rest of the code.
	const (
		codeState = iota
		runeState
		interpretedStringState
		rawStringState
	)
	searchChars := []string{
		"'\"`\n", // codeState
		`\'`,     // runeState
		`\"`,     // interpretedStringState
		"`\n",    // rawStringState
		// newlineState does not need to search
	}
	state := codeState
	for {
		i := strings.IndexAny(body, searchChars[state])
		if i < 0 {
			buf.WriteString(body)
			break
		}
		c := body[i]
		buf.WriteString(body[:i+1])
		body = body[i+1:]
		switch state {
		case codeState:
			switch c {
			case '\'':
				state = runeState
			case '"':
				state = interpretedStringState
			case '`':
				state = rawStringState
			case '\n':
				if strings.HasPrefix(body, oldIndent) {
					buf.WriteString(newIndent)
					body = body[len(oldIndent):]
				}
			}

		case runeState:
			switch c {
			case '\\':
				r, size := utf8.DecodeRuneInString(body)
				buf.WriteRune(r)
				body = body[size:]
			case '\'':
				state = codeState
			}

		case interpretedStringState:
			switch c {
			case '\\':
				r, size := utf8.DecodeRuneInString(body)
				buf.WriteRune(r)
				body = body[size:]
			case '"':
				state = codeState
			}

		case rawStringState:
			switch c {
			case '`':
				state = codeState
			case '\n':
				buf.WriteString(newIndent)
			}
		}
	}
	return buf.String()
}
