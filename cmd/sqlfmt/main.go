/*
sqlfmt — Will take an sql file and print out a pretty version, if it can not parse the sql file. It will
then show where the error is and exit with an error code of 1
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gdey/sqlparser"
)

const SytnaxErrorCode = 1

const defaultIndent = "    "

// FormatOptions controls formatting behavior.
type FormatOptions struct {
	// MaxFieldsInline: when number of SELECT fields, GROUP BY, or ORDER BY items
	// exceeds this value, each item is placed on its own line with leading comma.
	// Zero means always keep on one line.
	MaxFieldsInline int
}

// formatSQL produces "pretty" SQL from the parse tree. When tree is PositionedStatements
// and comments is non-nil, comments are interleaved by source position.
func formatSQL(tree sqlparser.Statement, comments []sqlparser.CommentEntry, opts FormatOptions) (string, error) {
	if posStmts, ok := tree.(sqlparser.PositionedStatements); ok {
		return formatWithPositionedComments(posStmts, comments, opts)
	}
	// Fallback: format tree only (no position-based comment interleaving).
	formatter := makePrettyFormat(opts)
	buf := sqlparser.NewTrackedBuffer(formatter)
	buf.Myprintf("%v", tree)
	s := buf.String()
	if s != "" && !strings.HasSuffix(s, ";") && !endsWithCommentOnly(tree) {
		s += ";"
	}
	return s, nil
}

// formatWithPositionedComments merges statements and comments by position and writes formatted SQL.
// Comments inside a statement's [Start, End) are skipped (they are in the AST and appear in the formatted statement).
func formatWithPositionedComments(stmts sqlparser.PositionedStatements, comments []sqlparser.CommentEntry, opts FormatOptions) (string, error) {
	type event struct {
		pos       int
		end       int
		comment   []byte
		statement sqlparser.Statement
	}
	var events []event
	for _, ps := range stmts {
		events = append(events, event{pos: ps.Start, end: ps.End, statement: ps.Statement})
	}
	for _, c := range comments {
		events = append(events, event{pos: c.Position, comment: c.Comment})
	}
	sort.Slice(events, func(i, j int) bool { return events[i].pos < events[j].pos })

	var b strings.Builder
	for i, ev := range events {
		if len(ev.comment) > 0 {
			b.Write(ev.comment)
			if len(ev.comment) > 0 && ev.comment[len(ev.comment)-1] != '\n' {
				b.WriteByte('\n')
			}
		} else if ev.statement != nil {
			// Use no-comments formatter so all comments come from CommentsTable only (no duplicates).
			formatter := makePrettyFormatNoComments(opts)
			buf := sqlparser.NewTrackedBuffer(formatter)
			buf.Myprintf("%v", ev.statement)
			s := buf.String()
			b.WriteString(s)
			if s != "" && !strings.HasSuffix(s, ";") {
				b.WriteString(";")
			}
			if i+1 < len(events) {
				b.WriteString("\n\n")
			}
		}
	}
	return b.String(), nil
}

// endsWithCommentOnly reports whether the tree is only comments or ends with a Comments statement.
func endsWithCommentOnly(tree sqlparser.Statement) bool {
	if _, ok := tree.(sqlparser.Comments); ok {
		return true
	}
	sts, ok := tree.(sqlparser.Statements)
	if !ok || len(sts) == 0 {
		return false
	}
	_, lastIsComment := sts[len(sts)-1].(sqlparser.Comments)
	return lastIsComment
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

// makePrettyFormat returns a formatter that uses the given options (e.g. max fields inline, leading comma).
func makePrettyFormat(opts FormatOptions) func(*sqlparser.TrackedBuffer, sqlparser.SQLNode) {
	return func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		prettyFormat(buf, node, opts)
	}
}

// prettyFormat formats nodes with newlines and uppercase keywords where applicable.
func prettyFormat(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode, opts FormatOptions) {
	switch n := node.(type) {
	case sqlparser.PositionedStatements:
		for i, ps := range n {
			if i > 0 {
				buf.WriteString(" ;\n\n")
			}
			if ps.Statement != nil {
				prettyFormat(buf, ps.Statement, opts)
			}
		}
	case sqlparser.PositionStatement:
		if n.Statement != nil {
			prettyFormat(buf, n.Statement, opts)
		}
	case sqlparser.Statements:
		for i, st := range n {
			if i > 0 {
				_, prevIsComment := n[i-1].(sqlparser.Comments)
				_, currIsComment := st.(sqlparser.Comments)
				if prevIsComment && currIsComment {
					buf.WriteString("\n\n")
				} else if prevIsComment {
					buf.WriteString("\n\n")
				} else {
					buf.WriteString(" ;\n\n")
				}
			}
			prettyFormat(buf, st, opts)
		}
	case *sqlparser.Select:
		formatPrettySelect(buf, n, opts)
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
	case *sqlparser.AliasedTableExpr:
		buf.Myprintf("%v", n.Expr)
		if n.As != nil {
			buf.Myprintf(" as %s", n.As)
		}
		if len(n.Comments) > 0 {
			buf.WriteString(" ")
		}
		buf.Myprintf("%v", n.Comments)
		if n.Hints != nil {
			buf.Myprintf("%v", n.Hints)
		}
	default:
		node.Format(buf)
	}
}

func makePrettyFormatNoComments(opts FormatOptions) func(*sqlparser.TrackedBuffer, sqlparser.SQLNode) {
	return func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
		prettyFormatNoComments(buf, node, opts)
	}
}

func prettyFormatNoComments(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode, opts FormatOptions) {
	switch n := node.(type) {
	case sqlparser.PositionedStatements:
		for i, ps := range n {
			if i > 0 {
				buf.WriteString(" ;\n\n")
			}
			if ps.Statement != nil {
				prettyFormatNoComments(buf, ps.Statement, opts)
			}
		}
	case sqlparser.PositionStatement:
		if n.Statement != nil {
			prettyFormatNoComments(buf, n.Statement, opts)
		}
	case *sqlparser.Select:
		formatPrettySelectNoComments(buf, n, opts)
	case *sqlparser.Union:
		buf.WriteString(formatWithDefault(n.Left))
		buf.WriteString(" ")
		buf.WriteString(strings.ToUpper(n.Type))
		buf.WriteString(" ")
		buf.WriteString(formatWithDefault(n.Right))
	case sqlparser.ValArg:
		s := string(n)
		if strings.HasPrefix(s, ":v$") {
			buf.WriteString("$" + s[3:])
		} else {
			buf.WriteString(s)
		}
	case *sqlparser.AliasedTableExpr:
		buf.Myprintf("%v", n.Expr)
		if n.As != nil {
			buf.Myprintf(" as %s", n.As)
		}
		if n.Hints != nil {
			buf.Myprintf("%v", n.Hints)
		}
	default:
		node.Format(buf)
	}
}

// formatListWithLeadingComma writes each item on its own line with leading comma when count > maxInline.
func formatListWithLeadingComma(buf *sqlparser.TrackedBuffer, maxInline int, count int, formatItem func(i int)) {
	if maxInline > 0 && count > maxInline {
		for i := 0; i < count; i++ {
			if i > 0 {
				buf.WriteString("\n")
				buf.WriteString(defaultIndent)
				buf.WriteString(", ")
			} else {
				buf.WriteString("\n")
				buf.WriteString(defaultIndent)
			}
			formatItem(i)
		}
	} else {
		for i := 0; i < count; i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			formatItem(i)
		}
	}
}

func formatPrettySelectNoComments(buf *sqlparser.TrackedBuffer, n *sqlparser.Select, opts FormatOptions) {
	buf.WriteString("SELECT ")
	buf.WriteString(strings.ToUpper(n.Distinct))
	formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.SelectExprs), func(i int) {
		buf.Myprintf("%v", n.SelectExprs[i])
	})
	if len(n.From) > 0 {
		buf.WriteString("\nFROM ")
		buf.Myprintf("%v", n.From)
	}
	if n.Where != nil {
		buf.WriteString("\n")
		buf.WriteString(strings.ToUpper(n.Where.Type))
		buf.WriteString(" ")
		buf.Myprintf("%v", n.Where.Expr)
	}
	if len(n.GroupBy) > 0 {
		buf.WriteString("\nGROUP BY ")
		formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.GroupBy), func(i int) {
			buf.Myprintf("%v", n.GroupBy[i])
		})
	}
	if n.Having != nil {
		buf.WriteString("\nHAVING ")
		buf.Myprintf("%v", n.Having.Expr)
	}
	if len(n.OrderBy) > 0 {
		buf.WriteString("\nORDER BY ")
		formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.OrderBy), func(i int) {
			buf.Myprintf("%v", n.OrderBy[i])
		})
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

func formatPrettySelect(buf *sqlparser.TrackedBuffer, n *sqlparser.Select, opts FormatOptions) {
	buf.WriteString("SELECT ")
	buf.Myprintf("%v", n.Comments)
	buf.WriteString(strings.ToUpper(n.Distinct))
	formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.SelectExprs), func(i int) {
		buf.Myprintf("%v", n.SelectExprs[i])
	})
	if len(n.From) > 0 {
		buf.WriteString("\nFROM ")
		buf.Myprintf("%v", n.From)
		buf.Myprintf("%v", n.FromComments)
	}
	if n.Where != nil {
		buf.WriteString("\n")
		buf.WriteString(strings.ToUpper(n.Where.Type))
		buf.WriteString(" ")
		buf.Myprintf("%v", n.Where.Comments)
		buf.Myprintf("%v", n.Where.Expr)
	}
	if len(n.GroupBy) > 0 {
		buf.WriteString("\nGROUP BY ")
		formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.GroupBy), func(i int) {
			buf.Myprintf("%v", n.GroupBy[i])
		})
	}
	if n.Having != nil {
		buf.WriteString("\nHAVING ")
		buf.Myprintf("%v", n.Having.Expr)
	}
	if len(n.OrderBy) > 0 {
		buf.WriteString("\nORDER BY ")
		formatListWithLeadingComma(buf, opts.MaxFieldsInline, len(n.OrderBy), func(i int) {
			buf.Myprintf("%v", n.OrderBy[i])
		})
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
	maxFieldsInline := flag.Int("max-fields-inline", 0, "when SELECT/GROUP BY/ORDER BY item count exceeds this, list each on its own line with leading comma (0 = always inline)")
	flag.Parse()
	opts := FormatOptions{MaxFieldsInline: *maxFieldsInline}
	statusCode := 0

	for i := 0; i < flag.NArg(); i++ {
		file := flag.Arg(i)
		buf, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		buff := string(buf)
		statement, commentEntries, err := sqlparser.Parse(buff)
		if err != nil {
			tokenError, _ := err.(*sqlparser.TokenizerError)
			FormatErrorMessage(os.Stderr, tokenError, file, buff, true)
			statusCode = SytnaxErrorCode
			continue
		}
		formatted, err := formatSQL(statement, commentEntries, opts)
		if err != nil {
			panic(err)
		}
		fmt.Println(formatted)
	}
	os.Exit(statusCode)
}
