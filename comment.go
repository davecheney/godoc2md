// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Godoc comment extraction and comment -> Markdown formatting.

package main

import (
	"io"
	"regexp"
	"strings"
	"text/template" // for HTMLEscape
	"unicode"
	"unicode/utf8"
)

const (
	// Regexp for Go identifiers
	identRx = `[a-zA-Z_][a-zA-Z_0-9]*` // TODO(gri) ASCII only for now - fix this

	// Regexp for URLs
	protocol = `(https?|ftp|file|gopher|mailto|news|nntp|telnet|wais|prospero):`
	hostPart = `[a-zA-Z0-9_@\-]+`
	filePart = `[a-zA-Z0-9_?%#~&/\-+=]+`
	urlRx    = protocol + `//` + // http://
		hostPart + `([.:]` + hostPart + `)*/?` + // //www.google.com:8080/
		filePart + `([:.,]` + filePart + `)*`
)

var matchRx = regexp.MustCompile(`(` + urlRx + `)|(` + identRx + `)`)

var (
	htmlA    = []byte(`<a href="`)
	htmlAq   = []byte(`">`)
	htmlEnda = []byte("</a>")

	mdPre     = []byte("\t")
	mdNewline = []byte("\n")
	mdH3      = []byte("### ")
)

// Emphasize and escape a line of text for HTML. URLs are converted into links.
func emphasize(w io.Writer, line string) {
	for {
		m := matchRx.FindStringSubmatchIndex(line)
		if m == nil {
			break
		}
		// m >= 6 (two parenthesized sub-regexps in matchRx, 1st one is urlRx)

		// write text before match
		_, _ = w.Write([]byte(line[0:m[0]]))

		// analyze match
		match := line[m[0]:m[1]]

		// if URL then write as link
		if m[2] >= 0 {
			_, _ = w.Write(htmlA)
			template.HTMLEscape(w, []byte(match))
			_, _ = w.Write(htmlAq)
		}

		// write match
		_, _ = w.Write([]byte(match))

		if m[2] >= 0 {
			_, _ = w.Write(htmlEnda)
		}

		// advance
		line = line[m[1]:]
	}
	_, _ = w.Write([]byte(line))
}

func indentLen(s string) int {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i
}

func isBlank(s string) bool {
	return len(s) == 0 || (len(s) == 1 && s[0] == '\n')
}

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func unindent(block []string) {
	if len(block) == 0 {
		return
	}

	// compute maximum common white prefix
	prefix := block[0][0:indentLen(block[0])]
	for _, line := range block {
		if !isBlank(line) {
			prefix = commonPrefix(prefix, line[0:indentLen(line)])
		}
	}
	n := len(prefix)

	// remove
	for i, line := range block {
		if !isBlank(line) {
			block[i] = line[n:]
		}
	}
}

// heading returns the trimmed line if it passes as a section heading;
// otherwise it returns the empty string.
func heading(line string) string {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return ""
	}

	// a heading must start with an uppercase letter
	r, _ := utf8.DecodeRuneInString(line)
	if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
		return ""
	}

	// it must end in a letter or digit:
	r, _ = utf8.DecodeLastRuneInString(line)
	if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
		return ""
	}

	// exclude lines with illegal characters
	if strings.ContainsAny(line, ",.;:!?+*/=()[]{}_^°&§~%#@<\">\\") {
		return ""
	}

	// allow "'" for possessive "'s" only
	for b := line; ; {
		i := strings.IndexRune(b, '\'')
		if i < 0 {
			break
		}
		if i+1 >= len(b) || b[i+1] != 's' || (i+2 < len(b) && b[i+2] != ' ') {
			return "" // not followed by "s "
		}
		b = b[i+2:]
	}

	return line
}

type op int

const (
	opPara op = iota
	opHead
	opPre
)

type block struct {
	op    op
	lines []string
}

var nonAlphaNumRx = regexp.MustCompile(`[^a-zA-Z0-9]`)

func anchorID(line string) string {
	// Add a "hdr-" prefix to avoid conflicting with IDs used for package symbols.
	return "hdr-" + nonAlphaNumRx.ReplaceAllString(line, "_")
}

// ToMD converts comment text to formatted Markdown.
// The comment was prepared by DocReader,
// so it is known not to have leading, trailing blank lines
// nor to have trailing spaces at the end of lines.
// The comment markers have already been removed.
//
// Each span of unindented non-blank lines is converted into
// a single paragraph. There is one exception to the rule: a span that
// consists of a single line, is followed by another paragraph span,
// begins with a capital letter, and contains no punctuation
// is formatted as a heading.
//
// A span of indented lines is converted into a <pre> block,
// with the common indent prefix removed.
//
// URLs in the comment text are converted into links.
func ToMD(w io.Writer, text string) {
	for _, b := range blocks(text) {
		switch b.op {
		case opPara:
			for _, line := range b.lines {
				emphasize(w, line)
			}
			_, _ = w.Write(mdNewline) // trailing newline to emulate </p>
		case opHead:
			_, _ = w.Write(mdH3)
			id := ""
			for _, line := range b.lines {
				if id == "" {
					id = anchorID(line)
				}
				_, _ = w.Write([]byte(line))
			}
			_, _ = w.Write(mdNewline)
		case opPre:
			_, _ = w.Write(mdNewline)
			for _, line := range b.lines {
				_, _ = w.Write(mdPre)
				emphasize(w, line)
			}
			_, _ = w.Write(mdNewline)
		}
	}
}

func blocks(text string) []block {
	var (
		out  []block
		para []string

		lastWasBlank   = false
		lastWasHeading = false
	)

	close := func() {
		if para != nil {
			out = append(out, block{opPara, para})
			para = nil
		}
	}

	lines := strings.SplitAfter(text, "\n")
	unindent(lines)
	for i := 0; i < len(lines); {
		line := lines[i]
		if isBlank(line) {
			// close paragraph
			close()
			i++
			lastWasBlank = true
			continue
		}
		if indentLen(line) > 0 {
			// close paragraph
			close()

			// count indented or blank lines
			j := i + 1
			for j < len(lines) && (isBlank(lines[j]) || indentLen(lines[j]) > 0) {
				j++
			}
			// but not trailing blank lines
			for j > i && isBlank(lines[j-1]) {
				j--
			}
			pre := lines[i:j]
			i = j

			unindent(pre)

			// put those lines in a pre block
			out = append(out, block{opPre, pre})
			lastWasHeading = false
			continue
		}

		if lastWasBlank && !lastWasHeading && i+2 < len(lines) &&
			isBlank(lines[i+1]) && !isBlank(lines[i+2]) && indentLen(lines[i+2]) == 0 {
			// current line is non-blank, surrounded by blank lines
			// and the next non-blank line is not indented: this
			// might be a heading.
			if head := heading(line); head != "" {
				close()
				out = append(out, block{opHead, []string{head}})
				i += 2
				lastWasHeading = true
				continue
			}
		}

		// open paragraph
		lastWasBlank = false
		lastWasHeading = false
		para = append(para, lines[i])
		i++
	}
	close()

	return out
}
