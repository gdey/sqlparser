/*
sqlfmt — Will take an sql file and print out a pretty version, if it can not parse the sql file. It will
then show where the error is and exit with an error code of 1
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gdey/sqlparser"
)

const SytnaxErrorCode = 1

func x(n int, vv string) (v string) {
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

// FormatSQL will walk the tree and produce "pretty" SQL.
// Given a statement like: select a,b from a where a = 1
// Return
//    SELECT a, b
//    FROM a
//    WHERE a = 1;
func formatSQL(tree sqlparser.Statement) (string, error) {

}
func main() {
	// We will have flags later.
	flag.Parse()

	for i := 1; i < len(os.Args); i++ {
		file := os.Args[i]
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		buff := string(buf)
		statement, err := sqlparser.Parse(buff)
		if err != nil {
			tokenError, _ := err.(*sqlparser.TokenizerError)
			// So, I need to figure out how to highlight where the error is.
			// I know that it's going to be at file position x, but there could be
			// new lines in there, and when I display the file I just want to display
			// The line that has the issue.

			idx := strings.IndexRune(buff[tokenError.Position-1:], '\n')
			if idx == -1 {
				idx = len(buff)
			} else {
				log.Println("idx", idx)
				idx += tokenError.Position - 1
			}
			ln, pos := PosTrans(buff, tokenError.Position)
			fmt.Fprintf(os.Stdout, "In file %s line %v %s:\n%v\n", file, ln+1, tokenError.Err, buff[:idx])
			fmt.Fprintf(os.Stdout, "%s%s\n", x(pos-1, "—"), "^")
			fmt.Fprintf(os.Stdout, "%s\n", tokenError.Err)
			os.Exit(SytnaxErrorCode)

		}
		fmt.Println(sqlparser.String(statement))
	}
}
