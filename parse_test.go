// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGen(t *testing.T) {
	_, _, err := Parse("select :a from a where a in (:b)")
	if err != nil {
		t.Error(err)
	}
}

func TestParse(t *testing.T) {
	for tcase := range iterateFiles("sqlparser_test/*.sql") {
		if tcase.output == "" {
			tcase.output = tcase.input
		}
		tree, _, err := Parse(tcase.input)
		var out string
		if err != nil {
			out = err.Error()
		} else {
			out = String(tree)
		}
		if out != tcase.output {
			t.Errorf("File:%s Line:%v\n%q\n%q", tcase.file, tcase.lineno, tcase.output, out)
		}
	}
}

// TestFullFileParse will test to make sure we can handle files with multiple sql statments, and comments.
func TestFullFilePassParse(t *testing.T) {
	// the test case are the files in this directory. The expected file is located in the expected
	// directory.
	for _, name := range glob("sqlparser_test/full_file/pass/*.sql") {
		sql, err := os.ReadFile(name)
		if err != nil {
			t.Errorf("Skipping test %v, got error trying to load it: %v", name, err)
			continue
		}
		expectedName := "sqlparser_test/full_file/pass/expected/" + filepath.Base(name)
		esql, err := os.ReadFile(expectedName)
		if err != nil {
			t.Errorf("Skipping test %v, got error trying to load expected file(%v): %v", name, expectedName, err)
			continue
		}
		tree, commentEntries, err := Parse(string(sql))
		if err != nil {
			t.Errorf("Failed test %v: Got error: %v", t.Name(), err)
		}
		out := FormatWithComments(tree, commentEntries)
		if out != string(esql) {
			sesql := string(esql)
			t.Errorf("Failed test %v:\nexpected(%v):\n[%v]\ngot(%v):\n[%v]\n", name, len(sesql), sesql, len(out), out)
			fmt.Printf("Failed test %v:\nexpected(%v):\n[%v]\ngot(%v):\n[%v]\n", name, len(sesql), sesql, len(out), out)
		}
	}
}

// TestFullFileParse will test to make sure we can handle files with multiple sql statments, and comments.
func TestFullFileFailParse(t *testing.T) {
	// the test case are the files in this directory. The expected file is located in the expected
	// directory.
	for _, name := range glob("sqlparser_test/full_file/fail/*.sql") {
		sql, err := os.ReadFile(name)
		if err != nil {
			t.Errorf("Skipping test %v, got error trying to load it: %v", name, err)
			continue
		}
		/*
			expectedName := "sqlparser_test/full_file/fail/expected/" + filepath.Base(name)
			esql, err := os.ReadFile(expectedName)
			if err != nil {
				t.Errorf("Skipping test %v, got error trying to load expected file(%v): %v", name, expectedName, err)
				continue
			}
		*/
		tree, _, err := Parse(string(sql))
		if err == nil {
			out := String(tree)
			t.Errorf("Failed test %v , expected error Got(%v)\n%v", name, len(out), out)
		}
		t.Logf("Got expected error (%v)", err)
	}
}

func BenchmarkParse1(b *testing.B) {
	sql := "select 'abcd', 20, 30.0, eid from a where 1=eid and name='3'"
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(sql)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse2(b *testing.B) {
	sql := "select aaaa, bbb, ccc, ddd, eeee, ffff, gggg, hhhh, iiii from tttt, ttt1, ttt3 where aaaa = bbbb and bbbb = cccc and dddd+1 = eeee group by fff, gggg having hhhh = iiii and iiii = jjjj order by kkkk, llll limit 3, 4"
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(sql)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type testCase struct {
	file   string
	lineno int
	input  string
	output string
}

func iterateFiles(pattern string) (testCaseIterator chan testCase) {
	names := glob(pattern)
	testCaseIterator = make(chan testCase)
	go func() {
		defer close(testCaseIterator)
		for _, name := range names {
			fd, err := os.OpenFile(name, os.O_RDONLY, 0)
			if err != nil {
				panic(fmt.Sprintf("Could not open file %s", name))
			}

			r := bufio.NewReader(fd)
			lineno := 0
			for {
				line, err := r.ReadString('\n')
				lines := strings.Split(strings.TrimRight(line, "\n"), "#")
				lineno++
				if err != nil {
					if err != io.EOF {
						panic(fmt.Sprintf("Error reading file %s: %s", name, err.Error()))
					}
					break
				}
				input := lines[0]
				output := ""
				if len(lines) > 1 {
					output = lines[1]
				}
				if input == "" {
					continue
				}
				testCaseIterator <- testCase{name, lineno, input, output}
			}
		}
	}()
	return testCaseIterator
}

func TestCaseWhenBareExpression(t *testing.T) {
	// WHEN with bare value (truth test) and WHEN with boolean_expression
	sql := `UPDATE t SET x = CASE WHEN ST_Contains(a, b) THEN 1 WHEN col IS NOT NULL THEN 2 END`
	tree, _, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	_ = tree
	// Round-trip: parse and re-format
	out := String(tree)
	if out == "" {
		t.Error("expected non-empty String(tree)")
	}
}

func TestCreateTableAsSelect(t *testing.T) {
	sql := "CREATE TABLE addresses AS SELECT 1 FROM parcels_data"
	pos, _, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(pos) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(pos))
	}
	ctas, ok := pos[0].Statement.(*CreateTableAsSelect)
	if !ok {
		t.Fatalf("expected CreateTableAsSelect, got %T", pos[0].Statement)
	}
	if string(ctas.Table) != "addresses" {
		t.Errorf("table name: got %q", ctas.Table)
	}
	refs := TableNamesFromTableExprs(ctas.Select.(*Select).From)
	if _, ok := refs["parcels_data"]; !ok {
		t.Errorf("expected parcels_data in FROM refs, got %v", refs)
	}
}

// parseMultiStatements splits sql by ";\n" and parses each segment so that
// multi-statement input can be parsed (each segment yields one statement).
func parseMultiStatements(sql string) (PositionedStatements, error) {
	sql = strings.ReplaceAll(sql, "\r\n", "\n")
	parts := strings.Split(sql, ";\n")
	var out PositionedStatements
	pos := 0
	for _, seg := range parts {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		ps, _, err := Parse(seg)
		if err != nil {
			return nil, fmt.Errorf("parse segment at pos %d: %w", pos, err)
		}
		if len(ps) != 1 {
			return nil, fmt.Errorf("expected single statement at pos %d", pos)
		}
		ps[0].Start = pos
		pos += len(seg) + 2 // +2 for ";\n"
		ps[0].End = pos
		out = append(out, ps[0])
	}
	return out, nil
}

func TestParseLACountyAddressesTransform(t *testing.T) {
	sql, err := os.ReadFile("testdata/la_county_addresses_transform.sql")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	pos, err := parseMultiStatements(string(sql))
	if err != nil {
		t.Fatalf("parse LA County addresses transform SQL: %v", err)
	}
	// Should have multiple statements including at least one CreateTableAsSelect (addresses).
	var foundCTAS bool
	for _, ps := range pos {
		if _, ok := ps.Statement.(*CreateTableAsSelect); ok {
			foundCTAS = true
			break
		}
	}
	if !foundCTAS {
		types := make([]string, 0, len(pos))
		for _, ps := range pos {
			if ps.Statement != nil {
				types = append(types, fmt.Sprintf("%T", ps.Statement))
			}
		}
		t.Errorf("expected at least one CreateTableAsSelect (CREATE TABLE addresses AS SELECT ...); got %d statements: %v", len(pos), types)
	}
}

func TestParseVarSQL(t *testing.T) {
	filename := "testdata/sql_with_var.sql"
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %s with error: %v", filename, err)
	}
	stmts, commentEntries, err := Parse(string(fileContent))
	if err != nil {
		t.Errorf("Failed to parse %s, with error: %v", filename, err)
	}
	_ = stmts
	_ = commentEntries
	//spew.Dump(commentEntries)
	spew.Dump(stmts)

}

func glob(pattern string) []string {
	out, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	return out
}
