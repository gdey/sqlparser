// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlparser

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGen(t *testing.T) {
	_, err := Parse("select :a from a where a in (:b)")
	if err != nil {
		t.Error(err)
	}
}

func TestParse(t *testing.T) {
	for tcase := range iterateFiles("sqlparser_test/*.sql") {
		if tcase.output == "" {
			tcase.output = tcase.input
		}
		tree, err := Parse(tcase.input)
		var out string
		if err != nil {
			out = err.Error()
		} else {
			out = String(tree)
		}
		if out != tcase.output {
			t.Error(fmt.Sprintf("File:%s Line:%v\n%q\n%q", tcase.file, tcase.lineno, tcase.output, out))
		}
	}
}

// TestFullFileParse will test to make sure we can handle files with multiple sql statments, and comments.
func TestFullFileParse(t *testing.T) {
	// the test case are the files in this directory. The expected file is located in the expected
	// directory.
	for _, name := range glob("sqlparser_test/full_file/pass/*.sql") {
		expectedName := "sqlparser_test/full_file/pass/expected/" + filepath.Base(name)
		sql, err := ioutil.ReadFile(name)
		if err != nil {
			t.Errorf("Skipping test %v, got error trying to load it: %v", name, err)
			continue
		}
		esql, err := ioutil.ReadFile(expectedName)
		if err != nil {
			t.Errorf("Skipping test %v, got error trying to load expected file(%v): %v", name, expectedName, err)
			continue
		}
		tree, err := Parse(string(sql))
		out := String(tree)
		if out != string(esql) {
			sesql := string(esql)
			t.Errorf("Failed test %v:\nexpected(%v):\n[%v]\ngot(%v):\n[%v]\n", name, len(sesql), sesql, len(out), out)
			fmt.Printf("Failed test %v:\nexpected(%v):\n[%v]\ngot(%v):\n[%v]\n", name, len(sesql), sesql, len(out), out)
		}
	}
}

func BenchmarkParse1(b *testing.B) {
	sql := "select 'abcd', 20, 30.0, eid from a where 1=eid and name='3'"
	for i := 0; i < b.N; i++ {
		_, err := Parse(sql)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse2(b *testing.B) {
	sql := "select aaaa, bbb, ccc, ddd, eeee, ffff, gggg, hhhh, iiii from tttt, ttt1, ttt3 where aaaa = bbbb and bbbb = cccc and dddd+1 = eeee group by fff, gggg having hhhh = iiii and iiii = jjjj order by kkkk, llll limit 3, 4"
	for i := 0; i < b.N; i++ {
		_, err := Parse(sql)
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

func glob(pattern string) []string {
	out, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	return out
}
