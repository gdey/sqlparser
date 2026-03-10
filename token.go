// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go tool goyacc -o sql.go sql.y
package sqlparser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gdey/sqlparser/sqltypes"
)

const EOFCHAR = 0x100

// TokenizerError is the struct used to hold the errors if there is any from
// Parsing. Line and Column are 1-based; Position is 0-based byte offset.
// Error() returns a string like "line N, column M (position P): message" and
// optionally " near \"token\"" when ErrToken is set.
type TokenizerError struct {
	Err      string
	Position int
	Line     int // 1-based line number
	Column   int // 1-based column number
	ErrToken []byte
}

func (te *TokenizerError) Error() string {
	if te == nil {
		return ""
	}
	msg := fmt.Sprintf("line %d, column %d (position %d): %s", te.Line, te.Column, te.Position, te.Err)
	if len(te.ErrToken) > 0 {
		msg += fmt.Sprintf(" near %q", te.ErrToken)
	}
	return msg
}

type CommentEntry struct {
	Position int
	Comment  []byte
}

// Tokenizer is the struct used to generate SQL
// tokens for the parser.
type Tokenizer struct {
	InStream                 *strings.Reader
	Input                    string // original SQL, used to compute line/column for errors
	AllowComments            bool
	CommentsTable            []CommentEntry
	ForceEOF                 bool
	lastChar                 uint16
	Position                 int
	errorToken               []byte
	lastLexError             string // specific message when Scan returns LEX_ERROR; used by Lex() to set LastError
	LastError                *TokenizerError
	posVarIndex              int
	ParseTree                Statement
	statementStartStack      []int // position when a top-level statement-start token was seen
	nextTokenStartsStatement bool  // true after ';' or at start; next SELECT/INSERT/etc. starts a statement
}

// positionToLineColumn returns 1-based line and column for the given 0-based byte position in input.
func positionToLineColumn(input string, pos int) (line, col int) {
	line = 1
	col = 1
	if pos > len(input) {
		pos = len(input)
	}
	for i := 0; i < pos; i++ {
		if input[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

// NewStringTokenizer creates a new Tokenizer for the
// sql string.
func NewStringTokenizer(sql string) *Tokenizer {
	return &Tokenizer{
		InStream:                 strings.NewReader(sql),
		Input:                    sql,
		nextTokenStartsStatement: true,
	}
}

// statementStartTokens are token types that begin a top-level statement.
// When we see one and nextTokenStartsStatement is true, we push Position.
var statementStartTokens = map[int]bool{
	SELECT: true, INSERT: true, UPDATE: true, DELETE: true, SET: true,
	CREATE: true, ALTER: true, RENAME: true, DROP: true, ANALYZE: true,
	BEGIN: true, COMMIT: true,
	WITH: true,
	SHOW: true, DESCRIBE: true, EXPLAIN: true,
}

// GetAndPopStatementStart returns the position at which the current statement
// started and clears it. Call from the grammar when reducing a "command".
// Returns 0 if the stack is empty.
func (tkn *Tokenizer) GetAndPopStatementStart() int {
	if len(tkn.statementStartStack) == 0 {
		return 0
	}
	pos := tkn.statementStartStack[len(tkn.statementStartStack)-1]
	tkn.statementStartStack = tkn.statementStartStack[:len(tkn.statementStartStack)-1]
	return pos
}

var keywords = map[string]int{
	"all":           ALL,
	"add":          ADD,
	"alter":        ALTER,
	"analyze":      ANALYZE,
	"and":          AND,
	"begin":        BEGIN,
	"commit":       COMMIT,
	"cast":         CAST,
	"column":       COLUMN,
	"as":            AS,
	"asc":           ASC,
	"between":       BETWEEN,
	"by":            BY,
	"case":          CASE,
	"create":        CREATE,
	"cross":         CROSS,
	"default":       DEFAULT,
	"delete":        DELETE,
	"desc":          DESC,
	"describe":      DESCRIBE,
	"distinct":      DISTINCT,
	"drop":          DROP,
	"duplicate":     DUPLICATE,
	"else":          ELSE,
	"end":           END,
	"except":        EXCEPT,
	"exists":        EXISTS,
	"explain":       EXPLAIN,
	"for":           FOR,
	"force":         FORCE,
	"from":          FROM,
	"group":         GROUP,
	"having":        HAVING,
	"if":            IF,
	"ignore":        IGNORE,
	"in":            IN,
	"index":         INDEX,
	"inner":         INNER,
	"insert":        INSERT,
	"intersect":     INTERSECT,
	"into":          INTO,
	"is":            IS,
	"join":          JOIN,
	"key":           KEY,
	"keyrange":      KEYRANGE,
	"left":          LEFT,
	"like":          LIKE,
	"limit":         LIMIT,
	"lock":          LOCK,
	"minus":         MINUS,
	"natural":       NATURAL,
	"not":           NOT,
	"null":          NULL,
	"on":            ON,
	"or":            OR,
	"order":         ORDER,
	"outer":         OUTER,
	"foreign":     FOREIGN,
	"primary":      PRIMARY,
	"references":   REFERENCES,
	"rename":       RENAME,
	"right":        RIGHT,
	"select":        SELECT,
	"set":           SET,
	"show":          SHOW,
	"straight_join": STRAIGHT_JOIN,
	"table":         TABLE,
	"temp":          TEMP,
	"temporary":     TEMPORARY,
	"then":          THEN,
	"to":            TO,
	"union":         UNION,
	"unique":        UNIQUE,
	"update":        UPDATE,
	"use":           USE,
	"using":         USING,
	"values":        VALUES,
	"view":          VIEW,
	"when":          WHEN,
	"where":         WHERE,
	"with":          WITH,
}

// Lex returns the next token form the Tokenizer.
// This function is used by go yacc.
func (tkn *Tokenizer) Lex(lval *yySymType) int {
	startPos := tkn.Position
	typ, val := tkn.Scan()
	for typ == COMMENT {
		// Let's add this comment to our comment table so that we can
		// aline the comment afterwords
		tkn.CommentsTable = append(tkn.CommentsTable, CommentEntry{
			Position: startPos,
			Comment:  val,
		})
		if tkn.AllowComments {
			break
		}
		startPos = tkn.Position
		typ, val = tkn.Scan()
	}
	lval.position = startPos
	if typ == LEX_ERROR {
		msg := tkn.lastLexError
		if msg == "" {
			msg = "invalid token"
		}
		tkn.lastLexError = ""
		line, col := positionToLineColumn(tkn.Input, startPos)
		tkn.LastError = &TokenizerError{
			Err:      msg,
			Position: startPos,
			Line:     line,
			Column:   col,
			ErrToken: val,
		}
	}
	switch typ {
	case ID, STRING, NUMBER, VALUE_ARG, LIST_ARG, COMMENT:
		lval.bytes = val
	}
	// Track top-level statement start for position-based comment interleaving.
	if typ == ';' {
		tkn.nextTokenStartsStatement = true
	} else if statementStartTokens[typ] && tkn.nextTokenStartsStatement {
		tkn.statementStartStack = append(tkn.statementStartStack, startPos)
		tkn.nextTokenStartsStatement = false
	}
	tkn.errorToken = val
	return typ
}

// Error is called by go yacc when the parser hits a syntax error. If the lexer
// already set LastError (e.g. for LEX_ERROR), we keep it so the more specific
// message is preserved.
func (tkn *Tokenizer) Error(err string) {
	if tkn.LastError != nil {
		return
	}
	line, col := positionToLineColumn(tkn.Input, tkn.Position)
	tkn.LastError = &TokenizerError{
		Err:      err,
		Position: tkn.Position,
		Line:     line,
		Column:   col,
		ErrToken: tkn.errorToken,
	}
}

// Scan scans the tokenizer for the next token and returns
// the token type and an optional value.
func (tkn *Tokenizer) Scan() (int, []byte) {
	if tkn.ForceEOF {
		return 0, nil
	}

	if tkn.lastChar == 0 {
		tkn.next()
	}
	tkn.skipBlank()
	switch ch := tkn.lastChar; {
	case isLetter(ch):
		return tkn.scanIdentifier()
	case isDigit(ch):
		return tkn.scanNumber(false)
	case ch == ':':
		return tkn.scanBindVar()
	default:
		tkn.next()
		switch ch {
		case EOFCHAR:
			return 0, nil
		case '=', ',', ';', '(', ')', '+', '*', '%', '&', '^', '~', '$':
			return int(ch), nil
		case '|':
			if tkn.lastChar == '|' {
				tkn.next()
				return CONCAT, nil
			}
			return int(ch), nil
		case '?':
			tkn.posVarIndex++
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, ":v%d", tkn.posVarIndex)
			return VALUE_ARG, buf.Bytes()
		case '.':
			if isDigit(tkn.lastChar) {
				return tkn.scanNumber(true)
			} else {
				return int(ch), nil
			}
		case '/':
			switch tkn.lastChar {
			case '/':
				tkn.next()
				return tkn.scanCommentType1("//")
			case '*':
				tkn.next()
				return tkn.scanCommentType2()
			default:
				return int(ch), nil
			}
		case '-':
			if tkn.lastChar == '-' {
				tkn.next()
				return tkn.scanCommentType1("--")
			}
			if tkn.lastChar == '>' {
				tkn.next()
				if tkn.lastChar == '>' {
					tkn.next()
					return JSON_EXTRACT_TEXT, nil
				}
				return JSON_EXTRACT, nil
			}
			return int(ch), nil
		case '#':
			return tkn.scanCommentType1("#")
		case '<':
			switch tkn.lastChar {
			case '>':
				tkn.next()
				return NE, nil
			case '=':
				tkn.next()
				switch tkn.lastChar {
				case '>':
					tkn.next()
					return NULL_SAFE_EQUAL, nil
				default:
					return LE, nil
				}
			default:
				return int(ch), nil
			}
		case '>':
			if tkn.lastChar == '=' {
				tkn.next()
				return GE, nil
			} else {
				return int(ch), nil
			}
		case '!':
			if tkn.lastChar == '=' {
				tkn.next()
				return NE, nil
			}
			tkn.lastLexError = "unexpected '!'"
			return LEX_ERROR, []byte("!")
		case '\'', '"':
			return tkn.scanString(ch, STRING)
		case '`':
			return tkn.scanLiteralIdentifier()
		default:
			tkn.lastLexError = "invalid character"
			return LEX_ERROR, []byte{byte(ch)}
		}
	}
}

func (tkn *Tokenizer) skipBlank() {
	ch := tkn.lastChar
	for ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
		tkn.next()
		ch = tkn.lastChar
	}
}

func (tkn *Tokenizer) scanIdentifier() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(tkn.lastChar))
	for tkn.next(); isLetter(tkn.lastChar) || isDigit(tkn.lastChar); tkn.next() {
		buffer.WriteByte(byte(tkn.lastChar))
	}
	lowered := bytes.ToLower(buffer.Bytes())
	if keywordId, found := keywords[string(lowered)]; found {
		return keywordId, lowered
	}
	return ID, buffer.Bytes()
}

func (tkn *Tokenizer) scanLiteralIdentifier() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(tkn.lastChar))
	if !isLetter(tkn.lastChar) {
		tkn.lastLexError = "backtick identifier must start with a letter"
		return LEX_ERROR, buffer.Bytes()
	}
	for tkn.next(); isLetter(tkn.lastChar) || isDigit(tkn.lastChar); tkn.next() {
		buffer.WriteByte(byte(tkn.lastChar))
	}
	if tkn.lastChar != '`' {
		tkn.lastLexError = "unterminated backtick identifier"
		return LEX_ERROR, buffer.Bytes()
	}
	tkn.next()
	return ID, buffer.Bytes()
}

func (tkn *Tokenizer) scanBindVar() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(tkn.lastChar))
	token := VALUE_ARG
	tkn.next()
	if tkn.lastChar == ':' {
		token = LIST_ARG
		buffer.WriteByte(byte(tkn.lastChar))
		tkn.next()
	}
	if !isLetter(tkn.lastChar) {
		tkn.lastLexError = "invalid bind variable"
		return LEX_ERROR, buffer.Bytes()
	}
	for isLetter(tkn.lastChar) || isDigit(tkn.lastChar) || tkn.lastChar == '.' {
		buffer.WriteByte(byte(tkn.lastChar))
		tkn.next()
	}
	return token, buffer.Bytes()
}

func (tkn *Tokenizer) scanMantissa(base int, buffer *bytes.Buffer) {
	for digitVal(tkn.lastChar) < base {
		tkn.ConsumeNext(buffer)
	}
}

func (tkn *Tokenizer) scanNumber(seenDecimalPoint bool) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	if seenDecimalPoint {
		buffer.WriteByte('.')
		tkn.scanMantissa(10, buffer)
		goto exponent
	}

	if tkn.lastChar == '0' {
		// int or float
		tkn.ConsumeNext(buffer)
		if tkn.lastChar == 'x' || tkn.lastChar == 'X' {
			// hexadecimal int
			tkn.ConsumeNext(buffer)
			tkn.scanMantissa(16, buffer)
		} else {
			// octal int or float
			seenDecimalDigit := false
			tkn.scanMantissa(8, buffer)
			if tkn.lastChar == '8' || tkn.lastChar == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				tkn.scanMantissa(10, buffer)
			}
			if tkn.lastChar == '.' || tkn.lastChar == 'e' || tkn.lastChar == 'E' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				tkn.lastLexError = "invalid octal number"
				return LEX_ERROR, buffer.Bytes()
			}
		}
		goto exit
	}

	// decimal int or float
	tkn.scanMantissa(10, buffer)

fraction:
	if tkn.lastChar == '.' {
		tkn.ConsumeNext(buffer)
		tkn.scanMantissa(10, buffer)
	}

exponent:
	if tkn.lastChar == 'e' || tkn.lastChar == 'E' {
		tkn.ConsumeNext(buffer)
		if tkn.lastChar == '+' || tkn.lastChar == '-' {
			tkn.ConsumeNext(buffer)
		}
		tkn.scanMantissa(10, buffer)
	}

exit:
	return NUMBER, buffer.Bytes()
}

func (tkn *Tokenizer) scanString(delim uint16, typ int) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	for {
		ch := tkn.lastChar
		tkn.next()
		if ch == delim {
			if tkn.lastChar == delim {
				tkn.next()
			} else {
				break
			}
		} else if ch == '\\' {
			if tkn.lastChar == EOFCHAR {
				tkn.lastLexError = "unterminated string"
				return LEX_ERROR, buffer.Bytes()
			}
			if decodedChar := sqltypes.SqlDecodeMap[byte(tkn.lastChar)]; decodedChar == sqltypes.DONTESCAPE {
				ch = tkn.lastChar
			} else {
				ch = uint16(decodedChar)
			}
			tkn.next()
		}
		if ch == EOFCHAR {
			tkn.lastLexError = "unterminated string"
			return LEX_ERROR, buffer.Bytes()
		}
		buffer.WriteByte(byte(ch))
	}
	return typ, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType1(prefix string) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteString(prefix)
	for tkn.lastChar != EOFCHAR {
		if tkn.lastChar == '\n' {
			tkn.ConsumeNext(buffer)
			break
		}
		tkn.ConsumeNext(buffer)
	}
	return COMMENT, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType2() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteString("/*")
	for {
		if tkn.lastChar == '*' {
			tkn.ConsumeNext(buffer)
			if tkn.lastChar == '/' {
				tkn.ConsumeNext(buffer)
				break
			}
			continue
		}
		if tkn.lastChar == EOFCHAR {
			tkn.lastLexError = "unterminated block comment"
			return LEX_ERROR, buffer.Bytes()
		}
		tkn.ConsumeNext(buffer)
	}
	return COMMENT, buffer.Bytes()
}

func (tkn *Tokenizer) ConsumeNext(buffer *bytes.Buffer) {
	if tkn.lastChar == EOFCHAR {
		// This should never happen.
		panic("unexpected EOF")
	}
	buffer.WriteByte(byte(tkn.lastChar))
	tkn.next()
}

func (tkn *Tokenizer) next() {
	if ch, err := tkn.InStream.ReadByte(); err != nil {
		// Only EOF is possible.
		tkn.lastChar = EOFCHAR
	} else {
		tkn.lastChar = uint16(ch)
	}
	tkn.Position++
}

func isLetter(ch uint16) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '@'
}

func digitVal(ch uint16) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

func isDigit(ch uint16) bool {
	return '0' <= ch && ch <= '9'
}
