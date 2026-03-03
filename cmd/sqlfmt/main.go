/*
sqlfmt — Will take an sql file and print out a pretty version, if it can not parse the sql file. It will
then show where the error is and exit with an error code of 1
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gdey/sqlparser"
)

const SytnaxErrorCode = 1

// formatSQL will walk the tree and produce "pretty" SQL.
// Given a statement like: select a,b from a where a = 1
// Return
//
//	SELECT a, b
//	FROM a
//	WHERE a = 1;
func formatSQL(tree sqlparser.Statement) (string, error) {
	buf := sqlparser.NewTrackedBuffer(prettyFormat)
	buf.Myprintf("%v", tree)
	s := buf.String()
	if s != "" && !strings.HasSuffix(s, ";") {
		s += ";"
	}
	return s, nil
}

// formatWithDefault formats a node using the default (single-line) formatter.
func formatWithDefault(node sqlparser.SQLNode) string {
	if node == nil {
		return ""
	}
	b := sqlparser.NewTrackedBuffer(nil)
	b.Myprintf("%v", node)
	return b.String()
}

// prettyFormat formats nodes with newlines and uppercase keywords where applicable.
func prettyFormat(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	switch n := node.(type) {
	case sqlparser.Statements:
		for i, st := range n {
			if i > 0 {
				buf.WriteString(" ;\n\n")
			}
			prettyFormat(buf, st)
		}
	case *sqlparser.Select:
		formatPrettySelect(buf, n)
	case *sqlparser.Union:
		buf.WriteString(formatWithDefault(n.Left))
		buf.WriteString(" ")
		buf.WriteString(strings.ToUpper(n.Type))
		buf.WriteString(" ")
		buf.WriteString(formatWithDefault(n.Right))
	case sqlparser.ValArg:
		// Parser normalizes $1 to :v$1; preserve original $n form in output.
		s := string(n)
		if strings.HasPrefix(s, ":v$") {
			buf.WriteString("$" + s[3:])
		} else {
			buf.WriteString(s)
		}
	default:
		// Delegate to node's Format so sub-nodes (e.g. ValArg) use this formatter.
		node.Format(buf)
	}
}

func formatPrettySelect(buf *sqlparser.TrackedBuffer, n *sqlparser.Select) {
	buf.WriteString("SELECT ")
	buf.Myprintf("%v", n.Comments)
	buf.WriteString(strings.ToUpper(n.Distinct))
	buf.Myprintf("%v", n.SelectExprs)
	if len(n.From) > 0 {
		buf.WriteString("\nFROM ")
		buf.Myprintf("%v", n.From)
		buf.Myprintf("%v", n.FromComments)
	}
	if n.Where != nil {
		buf.WriteString("\n")
		buf.WriteString(strings.ToUpper(n.Where.Type))
		buf.WriteString(" ")
		buf.Myprintf("%v", n.Where.Expr)
	}
	if len(n.GroupBy) > 0 {
		buf.WriteString("\nGROUP BY ")
		for i, e := range n.GroupBy {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Myprintf("%v", e)
		}
	}
	if n.Having != nil {
		buf.WriteString("\nHAVING ")
		buf.Myprintf("%v", n.Having.Expr)
	}
	if len(n.OrderBy) > 0 {
		buf.WriteString("\nORDER BY ")
		for i, o := range n.OrderBy {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Myprintf("%v", o)
		}
	}
	if n.Limit != nil {
		buf.WriteString("\nLIMIT ")
		if n.Limit.Offset != nil {
			buf.Myprintf("%v", n.Limit.Offset)
			buf.WriteString(", ")
		}
		buf.Myprintf("%v", n.Limit.Rowcount)
	}
	if n.Lock != "" {
		buf.WriteString("\n")
		buf.WriteString(strings.TrimSpace(strings.ToUpper(n.Lock)))
	}
}
func main() {
	// We will have flags later.
	flag.Parse()
	statusCode := 0

	for i := 1; i < len(os.Args); i++ {
		file := os.Args[i]
		buf, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		buff := string(buf)
		statement, _, err := sqlparser.Parse(buff)
		if err != nil {
			tokenError, _ := err.(*sqlparser.TokenizerError)
			FormatErrorMessage(os.Stderr, tokenError, file, buff, true)
			statusCode = SytnaxErrorCode
			continue
		}
		formatted, err := formatSQL(statement)
		if err != nil {
			panic(err)
		}
		fmt.Println(formatted)
	}
	os.Exit(statusCode)
}
