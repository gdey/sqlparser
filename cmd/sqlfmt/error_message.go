package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/gdey/sqlparser"
	"github.com/mgutz/ansi"
)

func repeat(n int, vv string) (v string) {
	for i := 0; i < n; i++ {
		v += vv
	}
	return
}

// PosTrans will scan through the string up to n looking for newlines. For each
// newline it countiers it will increment the ln count and reset pos to zero.
func PosTrans(s string, n int) (ln, pos int) {
	b := []rune(s)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if b[i] == '\n' {
			ln++
			pos = 0
			continue
		}
		pos++
	}
	return ln, pos
}

func printWithNumbers(w io.Writer, buff string, hlLine int) {
	line := 1
	start := 0
	nb := ansi.ColorFunc("black+h")
	hlnb := ansi.ColorFunc("white+h:black+h")
	hl := ansi.ColorFunc("yellow:black+h")
	nofmt := ansi.ColorFunc("")
	nbfmt := nb
	valfmt := nofmt

	var num, val string
	for i := range buff {
		if buff[i] == '\n' {
			if line == hlLine {
				nbfmt = hlnb
				valfmt = hl
			} else {
				nbfmt = nb
				valfmt = nofmt
			}
			num = nbfmt(fmt.Sprintf("%03v ", line))
			val = valfmt(fmt.Sprintf("%v", buff[start:i]))
			fmt.Fprintln(w, ansi.Reset, num, ansi.Reset, val)
			start = i + 1
			line++
		}
	}
	if line == hlLine {
		nbfmt = hlnb
		valfmt = hl
	} else {
		nbfmt = nb
		valfmt = nofmt
	}
	if start < len(buff) {
		num = nbfmt(fmt.Sprintf("%03v ", line))
		val = valfmt(fmt.Sprintf("%v", buff[start:]))
		fmt.Fprintln(w, ansi.Reset, num, ansi.Reset, val)
	}
}

func FormatErrorMessage(w io.Writer, err *sqlparser.TokenizerError, file, buff string, interactive bool) {
	// So, I need to figure out how to highlight where the error is.
	// I know that it's going to be at file position x, but there could be
	// new lines in there, and when I display the file I just want to display
	// The line that has the issue.

	idx := strings.IndexRune(buff[err.Position-1:], '\n')
	if idx == -1 {
		idx = len(buff)
	} else {
		idx += err.Position - 1
	}
	ln, pos := PosTrans(buff, err.Position)
	fmt.Fprintf(w, "%sIn file %s on line %v we had a %s.\n", ansi.Reset, file, ln+1, err.Err)
	printWithNumbers(w, buff[:idx], ln+1)
	fmt.Fprintf(w, " %s%s\n", repeat(pos+4, "—"), "^")
	if err.ErrToken == nil {
		fmt.Fprintf(w, "%s\n", err.Err)
	} else {
		fmt.Fprintf(w, "%s near %s\n", err.Err, err.ErrToken)
	}
}
