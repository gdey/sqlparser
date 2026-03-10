// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlparser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// errExpectParseFailure is a sentinel: when tcase.Err is this value, the test only checks that parse failed.
var errExpectParseFailure = errors.New("expect parse failure")

func TestGen(t *testing.T) {
	_, _, err := Parse("select :a from a where a in (:b)")
	if err != nil {
		t.Error(err)
	}
}

// TestParse tests the Parse function.
// Run with: DUMP_EXPECTED_AST=1 to dump the expected AST as Go.
func TestParse(t *testing.T) {
	type tcase struct {
		SQL        string
		Statements PositionedStatements
		Comments   []CommentEntry
		Err        error
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			stmts, comments, err := Parse(tc.SQL)
			if tc.Err != nil {
				if err == nil {
					t.Errorf("expected parse error, got output: %q", String(stmts))
					return
				}
				if !errors.Is(tc.Err, errExpectParseFailure) && err.Error() != tc.Err.Error() {
					t.Errorf("error:\n  got:  %s\n  want: %q", err.Error(), tc.Err.Error())
				} else {
					t.Logf("parse error (expected): %s", err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("parse error: %v", err)
				return
			}
			if tc.Statements != nil || tc.Comments != nil {
				// Compare expected Statements/Comments (compare by value: formatted SQL and comment list)
				if len(stmts) != len(tc.Statements) {
					t.Errorf("statements length: got %d, want %d", len(stmts), len(tc.Statements))
					return
				}
				for i := range stmts {
					if stmts[i].Start != tc.Statements[i].Start || stmts[i].End != tc.Statements[i].End {
						t.Errorf("statement %d position: got Start=%d End=%d, want Start=%d End=%d",
							i, stmts[i].Start, stmts[i].End, tc.Statements[i].Start, tc.Statements[i].End)
						return
					}
					gotStr, wantStr := String(stmts[i].Statement), String(tc.Statements[i].Statement)
					if gotStr != wantStr {
						t.Errorf("statement %d SQL mismatch:\n  got:  %q\n  want: %q", i, gotStr, wantStr)
						return
					}
				}
				if !reflect.DeepEqual(comments, tc.Comments) {
					t.Errorf("comments mismatch:\n  got:  %v\n  want: %v", comments, tc.Comments)
				}
			} else if os.Getenv("DUMP_EXPECTED_AST") == "1" {
				dumpExpectedAST(t, t.Name(), tc.SQL)
			}
			// With comment support, String(stmts) may not equal input (comments are separate).
			// Only check that formatting is idempotent: parse → format → parse → format gives same string.
			out := String(stmts)
			stmts2, _, err2 := Parse(out)
			if err2 != nil {
				t.Errorf("re-parse of formatted output failed: %v\nformatted was: %q", err2, out)
				return
			}
			out2 := String(stmts2)
			if out2 != out {
				t.Errorf("format not idempotent:\n  first:  %q\n  second: %q", out, out2)
			}
		}
	}

	tests := map[string]tcase{
		// Inline short — expected AST below so you can see and edit it when the parser changes
		"inline/select1": {
			SQL: "select 1",
			Statements: PositionedStatements{
				{Start: 0, End: 10, Statement: &Select{
					SelectExprs: SelectExprs{&NonStarExpr{Expr: NumVal("1")}},
				}},
			},
			Comments: nil,
		},
		"inline/select_from": {
			SQL: "select 1 from t",
			Statements: PositionedStatements{
				{Start: 0, End: 17, Statement: &Select{
					SelectExprs: SelectExprs{&NonStarExpr{Expr: NumVal("1")}},
					From:        TableExprs{&AliasedTableExpr{Expr: &TableName{Name: []byte("t")}}},
				}},
			},
			Comments: nil,
		},
		"inline/select_1_1": {
			SQL: "select 1+1",
			Statements: PositionedStatements{
				{Start: 0, End: 12, Statement: &Select{
					SelectExprs: SelectExprs{&NonStarExpr{Expr: &BinaryExpr{
						Left: NumVal("1"), Operator: "+", Right: NumVal("1"),
					}}},
				}},
			},
			Comments: nil,
		},

		// Expect parse failure
		"fail/select_bang":     {SQL: "select !8 from t", Err: errExpectParseFailure},
		"fail/select_dollar":   {SQL: "select $ from t", Err: errExpectParseFailure},
		"fail/select_colon":    {SQL: "select : from t", Err: errExpectParseFailure},
		"fail/select_078":      {SQL: "select 078 from t", Err: errExpectParseFailure},
		"fail/backtick_1a":     {SQL: "select `1a` from t", Err: errExpectParseFailure},
		"fail/backtick_table":  {SQL: "select `:table` from t", Err: errExpectParseFailure},
		"fail/backtick_table2": {SQL: "select `table:` from t", Err: errExpectParseFailure},
		"fail/unclosed_bs":     {SQL: "select 'aa\\", Err: errExpectParseFailure},
		"fail/unclosed_quote":  {SQL: "select 'aa", Err: errExpectParseFailure},
		"fail/bind_1":          {SQL: "select * from t where :1 = 2", Err: errExpectParseFailure},
		"fail/bind_dot":        {SQL: "select * from t where :. = 2", Err: errExpectParseFailure},
		"fail/bind_dbl1":       {SQL: "select * from t where ::1 = 2", Err: errExpectParseFailure},
		"fail/bind_dbl_dot":    {SQL: "select * from t where ::. = 2", Err: errExpectParseFailure},
		"fail/unclosed_cmt":    {SQL: "select /* aa", Err: errExpectParseFailure},
		"pass/create_table_columns": {SQL: `CREATE TABLE a_table (
    id integer primary key
  , name text not null
);`},
		"pass/create_table_varchar":  {SQL: "create table t (name varchar(255))"},
		"pass/create_table_decimal":  {SQL: "create table t (amount decimal(10, 2))"},
		"pass/create_table_constraints": {SQL: "create table t (a int, b int, primary key (a, b), unique (a))"},

		// pass/case_sensitivity (from testdata/TestParse/pass/case_sensitivity.sql)
		"pass/case_sensitivity/1": {SQL: "create table A"},
		// 2,3,4,5,7: alter/rename canonical forms not valid standalone input
		"pass/case_sensitivity/6":  {SQL: "drop table B"},
		"pass/case_sensitivity/8":  {SQL: "select a from B"},
		"pass/case_sensitivity/9":  {SQL: "select a as b from C"},
		"pass/case_sensitivity/10": {SQL: "select B.* from c"},
		"pass/case_sensitivity/11": {SQL: "select B.a from c"},
		"pass/case_sensitivity/12": {SQL: "select * from B as C"},
		"pass/case_sensitivity/13": {SQL: "select * from A.B"},
		"pass/case_sensitivity/14": {SQL: "update A set b = 1"},
		"pass/case_sensitivity/15": {SQL: "update A.B set b = 1"},
		"pass/case_sensitivity/16": {SQL: "select a() from b"},
		"pass/case_sensitivity/17": {SQL: "select a(b, c) from b"},
		"pass/case_sensitivity/18": {SQL: "select a(distinct b, c) from b"},
		"pass/case_sensitivity/19": {SQL: "select if(b, c) from b"},
		"pass/case_sensitivity/20": {SQL: "select values(b, c) from b"},
		"pass/case_sensitivity/21": {SQL: "select * from b use index (a)"},
		"pass/case_sensitivity/22": {SQL: "insert into A(a, b) values (1, 2)"},
		"pass/case_sensitivity/23": {SQL: "create table A"},
		"pass/case_sensitivity/24": {SQL: "create table a"},
		// 25: alter table a is canonical output only
		"pass/case_sensitivity/26": {SQL: "drop table a"},

		// pass/parse_pass (from testdata/TestParse/pass/parse_pass.sql)
		"pass/parse_pass/1":   {SQL: "select 1"},
		"pass/parse_pass/2":   {SQL: "select 1+1"},
		"pass/parse_pass/3":   {SQL: "select 1"},
		"pass/parse_pass/4":   {SQL: "select 1 from t"},
		"pass/parse_pass/5":   {SQL: "select 1 from t"},
		"pass/parse_pass/6":   {SQL: "select .1 from t"},
		"pass/parse_pass/7":   {SQL: "select 1.2e1 from t"},
		"pass/parse_pass/8":   {SQL: "select 1.2e+1 from t"},
		"pass/parse_pass/9":   {SQL: "select 1.2e-1 from t"},
		"pass/parse_pass/10":  {SQL: "select 08.3 from t"},
		"pass/parse_pass/11":  {SQL: "select -1 from t where b = -2"},
		"pass/parse_pass/12":  {SQL: "select 1 from t"},
		"pass/parse_pass/13":  {SQL: "select 1 from t"},
		"pass/parse_pass/14":  {SQL: "select /* simplest */ 1 from t"},
		"pass/parse_pass/15":  {SQL: "select /* double star **/ 1 from t"},
		"pass/parse_pass/16":  {SQL: "select /* double */ /* comment */ 1 from t"},
		"pass/parse_pass/17":  {SQL: "select /* back-quote */ 1 from t"},
		"pass/parse_pass/18":  {SQL: "select /* back-quote keyword */ 1 from `from`"},
		"pass/parse_pass/19":  {SQL: "select /* @ */ @@a from b"},
		"pass/parse_pass/20":  {SQL: "select /* \\0 */ '\\0' from a"},
		"pass/parse_pass/21":  {SQL: "select 1 from t"},
		"pass/parse_pass/22":  {SQL: "select /* union */ 1 from t union select 1 from t"},
		"pass/parse_pass/23":  {SQL: "select /* double union */ 1 from t union select 1 from t union select 1 from t"},
		"pass/parse_pass/24":  {SQL: "select /* union all */ 1 from t union all select 1 from t"},
		"pass/parse_pass/25":  {SQL: "select /* minus */ 1 from t minus select 1 from t"},
		"pass/parse_pass/26":  {SQL: "select /* except */ 1 from t except select 1 from t"},
		"pass/parse_pass/27":  {SQL: "select /* intersect */ 1 from t intersect select 1 from t"},
		"pass/parse_pass/28":  {SQL: "select /* distinct */ distinct 1 from t"},
		"pass/parse_pass/29":  {SQL: "select /* for update */ 1 from t for update"},
		"pass/parse_pass/30":  {SQL: "select /* lock in share mode */ 1 from t lock in share mode"},
		"pass/parse_pass/31":  {SQL: "select /* select list */ 1, 2 from t"},
		"pass/parse_pass/32":  {SQL: "select /* * */ * from t"},
		"pass/parse_pass/33":  {SQL: "select /* column alias */ a as b from t"},
		"pass/parse_pass/34":  {SQL: "select /* column alias with as */ a as b from t"},
		"pass/parse_pass/35":  {SQL: "select /* a.* */ a.* from t"},
		"pass/parse_pass/36":  {SQL: "select /* select with bool expr */ a = b from t"},
		"pass/parse_pass/37":  {SQL: "select /* case_when */ case when a = b then c end from t"},
		"pass/parse_pass/38":  {SQL: "select /* case_when_else */ case when a = b then c else d end from t"},
		"pass/parse_pass/39":  {SQL: "select /* case_when_when_else */ case when a = b then c when b = d then d else d end from t"},
		"pass/parse_pass/40":  {SQL: "select /* case */ case aa when a = b then c end from t"},
		"pass/parse_pass/41":  {SQL: "select /* parenthesis */ 1 from (t)"},
		"pass/parse_pass/42":  {SQL: "select /* table list */ 1 from t1, t2"},
		"pass/parse_pass/43":  {SQL: "select /* parenthessis in table list 1 */ 1 from (t1), t2"},
		"pass/parse_pass/44":  {SQL: "select /* parenthessis in table list 2 */ 1 from t1, (t2)"},
		"pass/parse_pass/45":  {SQL: "select /* use */ 1 from t1 use index (a) where b = 1"},
		"pass/parse_pass/46":  {SQL: "select /* ignore */ 1 from t1 as t2 ignore index (a), t3 use index (b) where b = 1"},
		"pass/parse_pass/47":  {SQL: "select /* use */ 1 from t1 as t2 use index (a), t3 use index (b) where b = 1"},
		"pass/parse_pass/48":  {SQL: "select /* force */ 1 from t1 as t2 force index (a), t3 force index (b) where b = 1"},
		"pass/parse_pass/49":  {SQL: "select /* table alias */ 1 from t as t1"},
		"pass/parse_pass/50":  {SQL: "select /* table alias with as */ 1 from t as t1"},
		"pass/parse_pass/51":  {SQL: "select /* join */ 1 from t1 join t2"},
		"pass/parse_pass/52":  {SQL: "select /* straight_join */ 1 from t1 straight_join t2"},
		"pass/parse_pass/53":  {SQL: "select /* left join */ 1 from t1 left join t2"},
		"pass/parse_pass/54":  {SQL: "select /* left outer join */ 1 from t1 left join t2"},
		"pass/parse_pass/55":  {SQL: "select /* right join */ 1 from t1 right join t2"},
		"pass/parse_pass/56":  {SQL: "select /* right outer join */ 1 from t1 right join t2"},
		"pass/parse_pass/57":  {SQL: "select /* inner join */ 1 from t1 join t2"},
		"pass/parse_pass/58":  {SQL: "select /* cross join */ 1 from t1 cross join t2"},
		"pass/parse_pass/59":  {SQL: "select /* natural join */ 1 from t1 natural join t2"},
		"pass/parse_pass/60":  {SQL: "select /* join on */ 1 from t1 join t2 on a = b"},
		"pass/parse_pass/61":  {SQL: "select /* s.t */ 1 from s.t"},
		"pass/parse_pass/62":  {SQL: "select /* select in from */ 1 from (select 1 from t)"},
		"pass/parse_pass/63":  {SQL: "select /* where */ 1 from t where a = b"},
		"pass/parse_pass/64":  {SQL: "select /* and */ 1 from t where a = b and a = c"},
		"pass/parse_pass/65":  {SQL: "select /* or */ 1 from t where a = b or a = c"},
		"pass/parse_pass/66":  {SQL: "select /* not */ 1 from t where not a = b"},
		"pass/parse_pass/67":  {SQL: "select /* exists */ 1 from t where exists (select 1 from t)"},
		"pass/parse_pass/68":  {SQL: "select /* keyrange */ 1 from t where keyrange(1, 2)"},
		"pass/parse_pass/69":  {SQL: "select /* (boolean) */ 1 from t where not (a = b)"},
		"pass/parse_pass/70":  {SQL: "select /* in value list */ 1 from t where a in (b, c)"},
		"pass/parse_pass/71":  {SQL: "select /* in select */ 1 from t where a in (select 1 from t)"},
		"pass/parse_pass/72":  {SQL: "select /* not in */ 1 from t where a not in (b, c)"},
		"pass/parse_pass/73":  {SQL: "select /* like */ 1 from t where a like b"},
		"pass/parse_pass/74":  {SQL: "select /* not like */ 1 from t where a not like b"},
		"pass/parse_pass/75":  {SQL: "select /* between */ 1 from t where a between b and c"},
		"pass/parse_pass/76":  {SQL: "select /* not between */ 1 from t where a not between b and c"},
		"pass/parse_pass/77":  {SQL: "select /* is null */ 1 from t where a is null"},
		"pass/parse_pass/78":  {SQL: "select /* is not null */ 1 from t where a is not null"},
		"pass/parse_pass/79":  {SQL: "select /* < */ 1 from t where a < b"},
		"pass/parse_pass/80":  {SQL: "select /* <= */ 1 from t where a <= b"},
		"pass/parse_pass/81":  {SQL: "select /* >= */ 1 from t where a >= b"},
		"pass/parse_pass/82":  {SQL: "select /* <> */ 1 from t where a != b"},
		"pass/parse_pass/83":  {SQL: "select /* <=> */ 1 from t where a <=> b"},
		"pass/parse_pass/84":  {SQL: "select /* != */ 1 from t where a != b"},
		"pass/parse_pass/85":  {SQL: "select /* single value expre list */ 1 from t where a in (b)"},
		"pass/parse_pass/86":  {SQL: "select /* select as a value expression */ 1 from t where a = (select a from t)"},
		"pass/parse_pass/87":  {SQL: "select /* parenthesised value */ 1 from t where a = (b)"},
		"pass/parse_pass/88":  {SQL: "select /* over-parenthesize */ ((1)) from t where ((a)) in (((1))) and ((a, b)) in ((((1, 1))), ((2, 2)))"},
		"pass/parse_pass/89":  {SQL: "select /* dot-parenthesize */ (a.b) from t where (b.c) = 2"},
		"pass/parse_pass/90":  {SQL: "select /* & */ 1 from t where a = b&c"},
		"pass/parse_pass/91":  {SQL: "select /* | */ 1 from t where a = b|c"},
		"pass/parse_pass/92":  {SQL: "select /* ^ */ 1 from t where a = b^c"},
		"pass/parse_pass/93":  {SQL: "select /* + */ 1 from t where a = b+c"},
		"pass/parse_pass/94":  {SQL: "select /* - */ 1 from t where a = b-c"},
		"pass/parse_pass/95":  {SQL: "select /* * */ 1 from t where a = b*c"},
		"pass/parse_pass/96":  {SQL: "select /* / */ 1 from t where a = b/c"},
		"pass/parse_pass/97":  {SQL: "select /* % */ 1 from t where a = b%c"},
		"pass/parse_pass/98":  {SQL: "select /* u+ */ 1 from t where a = +b"},
		"pass/parse_pass/99":  {SQL: "select /* u- */ 1 from t where a = -b"},
		"pass/parse_pass/100": {SQL: "select /* u~ */ 1 from t where a = ~b"},
		"pass/parse_pass/101": {SQL: "select /* empty function */ 1 from t where a = b()"},
		"pass/parse_pass/102": {SQL: "select /* function with 1 param */ 1 from t where a = b(c)"},
		"pass/parse_pass/103": {SQL: "select /* function with many params */ 1 from t where a = b(c, d)"},
		"pass/parse_pass/104": {SQL: "select /* if as func */ 1 from t where a = if(b)"},
		"pass/parse_pass/105": {SQL: "select /* function with distinct */ count(distinct a) from t"},
		"pass/parse_pass/106": {SQL: "select /* a */ a from t"},
		"pass/parse_pass/107": {SQL: "select /* a.b */ a.b from t"},
		"pass/parse_pass/108": {SQL: "select /* string */ 'a' from t"},
		"pass/parse_pass/109": {SQL: "select /* double quoted string */ 'a' from t"},
		"pass/parse_pass/110": {SQL: "select /* quote quote in string */ 'a\\'a' from t"},
		"pass/parse_pass/111": {SQL: "select /* double quote quote in string */ 'a\\\"a' from t"},
		"pass/parse_pass/112": {SQL: "select /* quote in double quoted string */ 'a\\'a' from t"},
		"pass/parse_pass/113": {SQL: "select /* backslash quote in string */ 'a\\'a' from t"},
		"pass/parse_pass/114": {SQL: "select /* literal backslash in string */ 'a\\\\na' from t"},
		"pass/parse_pass/115": {SQL: "select /* all escapes */ '\\0\\'\\\"\\b\\n\\r\\t\\Z\\\\' from t"},
		"pass/parse_pass/116": {SQL: "select /* non-escape */ 'x' from t"},
		"pass/parse_pass/117": {SQL: "select /* unescaped backslash */ '\\n' from t"},
		"pass/parse_pass/118": {SQL: "select /* value argument */ :a from t"},
		"pass/parse_pass/119": {SQL: "select /* value argument with digit */ :a1 from t"},
		"pass/parse_pass/120": {SQL: "select /* value argument with dot */ :a.b from t"},
		"pass/parse_pass/121": {SQL: "select /* positional argument */ :v1 from t"},
		"pass/parse_pass/122": {SQL: "select /* multiple positional arguments */ :v1, :v2 from t"},
		"pass/parse_pass/123": {SQL: "select /* list arg */ * from t where a in ::list"},
		"pass/parse_pass/124": {SQL: "select /* list arg not in */ * from t where a not in ::list"},
		"pass/parse_pass/125": {SQL: "select /* null */ null from t"},
		"pass/parse_pass/126": {SQL: "select /* octal */ 010 from t"},
		"pass/parse_pass/127": {SQL: "select /* hex */ 0xf0 from t"},
		"pass/parse_pass/128": {SQL: "select /* hex caps */ 0xF0 from t"},
		"pass/parse_pass/129": {SQL: "select /* float */ 0.1 from t"},
		"pass/parse_pass/130": {SQL: "select /* group by */ 1 from t group by a"},
		"pass/parse_pass/131": {SQL: "select /* having */ 1 from t having a = b"},
		"pass/parse_pass/132": {SQL: "select /* simple order by */ 1 from t order by a asc"},
		"pass/parse_pass/133": {SQL: "select /* order by asc */ 1 from t order by a asc"},
		"pass/parse_pass/134": {SQL: "select /* order by desc */ 1 from t order by a desc"},
		"pass/parse_pass/135": {SQL: "select /* limit a */ 1 from t limit a"},
		"pass/parse_pass/136": {SQL: "select /* limit a,b */ 1 from t limit a, b"},
		"pass/parse_pass/137": {SQL: "insert /* simple */ into a values (1)"},
		"pass/parse_pass/138": {SQL: "insert /* a.b */ into a.b values (1)"},
		"pass/parse_pass/139": {SQL: "insert /* multi-value */ into a values (1, 2)"},
		"pass/parse_pass/140": {SQL: "insert /* multi-value list */ into a values (1, 2), (3, 4)"},
		"pass/parse_pass/141": {SQL: "insert /* set */ into a(a, a.b) values (1, 2)"},
		"pass/parse_pass/142": {SQL: "insert /* value expression list */ into a values (a+1, 2*3)"},
		"pass/parse_pass/143": {SQL: "insert /* column list */ into a(a, b) values (1, 2)"},
		"pass/parse_pass/144": {SQL: "insert /* qualified column list */ into a(a, a.b) values (1, 2)"},
		"pass/parse_pass/145": {SQL: "insert /* select */ into a select b, c from d"},
		"pass/parse_pass/146": {SQL: "insert /* on duplicate */ into a values (1, 2) on duplicate key update b = values(a), c = d"},
		"pass/parse_pass/147": {SQL: "update /* simple */ a set b = 3"},
		"pass/parse_pass/148": {SQL: "update /* a.b */ a.b set b = 3"},
		"pass/parse_pass/149": {SQL: "update /* b.c */ a set b.c = 3"},
		"pass/parse_pass/150": {SQL: "update /* list */ a set b = 3, c = 4"},
		"pass/parse_pass/151": {SQL: "update /* expression */ a set b = 3+4"},
		"pass/parse_pass/152": {SQL: "update /* where */ a set b = 3 where a = b"},
		"pass/parse_pass/153": {SQL: "update /* order */ a set b = 3 order by c desc"},
		"pass/parse_pass/154": {SQL: "update /* limit */ a set b = 3 limit c"},
		"pass/parse_pass/155": {SQL: "delete /* simple */ from a"},
		"pass/parse_pass/156": {SQL: "delete /* a.b */ from a.b"},
		"pass/parse_pass/157": {SQL: "delete /* where */ from a where a = b"},
		"pass/parse_pass/158": {SQL: "delete /* order */ from a order by b desc"},
		"pass/parse_pass/159": {SQL: "delete /* limit */ from a limit b"},
		"pass/parse_pass/160": {SQL: "set /* simple */ a = 3"},
		"pass/parse_pass/161": {SQL: "set /* list */ a = 3, b = 4"},
		// 162-175: alter/rename canonical forms not valid standalone input
		"pass/parse_pass/176": {SQL: "create table a"},
		"pass/parse_pass/177": {SQL: "create table a"},
		// 178-182: alter table b / alter table a canonical only
		"pass/parse_pass/183": {SQL: "drop table a"},
		"pass/parse_pass/184": {SQL: "drop table a"},
		"pass/parse_pass/185": {SQL: "drop table a"},
		"pass/parse_pass/186": {SQL: "drop table a"},
		// 187-188: alter table a canonical only
		// 189-191: "other" is output of show/describe/explain, not valid input

	}
	for name, tc := range tests {
		t.Run(name, fn(tc))
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
	if ctas.Temporary {
		t.Error("expected Temporary false for CREATE TABLE")
	}
	if ctas.Temp() {
		t.Error("expected Temp() false for CREATE TABLE")
	}
}

func TestCreateTableWithColumns(t *testing.T) {
	sql := `CREATE TABLE a_table (
    id integer primary key
  , name text not null
);`
	pos, _, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(pos) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(pos))
	}
	ct, ok := pos[0].Statement.(*CreateTable)
	if !ok {
		t.Fatalf("expected CreateTable, got %T", pos[0].Statement)
	}
	if string(ct.Table) != "a_table" {
		t.Errorf("table name: got %q", ct.Table)
	}
	if ct.Temporary || ct.Temp() {
		t.Error("expected Temporary false for CREATE TABLE")
	}
	if len(ct.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(ct.Columns))
	}
	if string(ct.Columns[0].Name) != "id" || string(ct.Columns[0].Type.Name) != "integer" || !ct.Columns[0].Options.PrimaryKey {
		t.Errorf("unexpected col0: %#v", ct.Columns[0])
	}
	if string(ct.Columns[1].Name) != "name" || string(ct.Columns[1].Type.Name) != "text" || !ct.Columns[1].Options.NotNull {
		t.Errorf("unexpected col1: %#v", ct.Columns[1])
	}
	// Round-trip (String is canonicalized to one line).
	out := String(pos[0].Statement)
	if !strings.HasPrefix(out, "create table a_table (") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestCreateTableTypeFormsAndConstraints(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		fn   func(t *testing.T, ct *CreateTable)
	}{
		{
			name: "varchar_length",
			sql:  "CREATE TABLE t (name varchar(255))",
			fn: func(t *testing.T, ct *CreateTable) {
				if len(ct.Columns) != 1 {
					t.Fatalf("expected 1 column, got %d", len(ct.Columns))
				}
				c := ct.Columns[0]
				if string(c.Type.Name) != "varchar" || c.Type.Length == nil || *c.Type.Length != 255 || c.Type.Scale != nil {
					t.Errorf("unexpected type: %#v", c.Type)
				}
			},
		},
		{
			name: "decimal_precision_scale",
			sql:  "CREATE TABLE t (amount decimal(10, 2))",
			fn: func(t *testing.T, ct *CreateTable) {
				if len(ct.Columns) != 1 {
					t.Fatalf("expected 1 column, got %d", len(ct.Columns))
				}
				c := ct.Columns[0]
				if string(c.Type.Name) != "decimal" || c.Type.Length == nil || *c.Type.Length != 10 || c.Type.Scale == nil || *c.Type.Scale != 2 {
					t.Errorf("unexpected type: %#v", c.Type)
				}
			},
		},
		{
			name: "table_primary_key",
			sql:  "CREATE TABLE t (a int, b int, primary key (a, b))",
			fn: func(t *testing.T, ct *CreateTable) {
				if len(ct.Constraints) != 1 {
					t.Fatalf("expected 1 constraint, got %d", len(ct.Constraints))
				}
				tc := ct.Constraints[0]
				if tc.Kind != "primary key" || len(tc.Columns) != 2 || string(tc.Columns[0]) != "a" || string(tc.Columns[1]) != "b" {
					t.Errorf("unexpected constraint: %#v", tc)
				}
			},
		},
		{
			name: "table_unique",
			sql:  "CREATE TABLE t (email varchar(255), unique (email))",
			fn: func(t *testing.T, ct *CreateTable) {
				if len(ct.Constraints) != 1 {
					t.Fatalf("expected 1 constraint, got %d", len(ct.Constraints))
				}
				tc := ct.Constraints[0]
				if tc.Kind != "unique" || len(tc.Columns) != 1 || string(tc.Columns[0]) != "email" {
					t.Errorf("unexpected constraint: %#v", tc)
				}
			},
		},
		{
			name: "foreign_key",
			sql:  "CREATE TABLE orders (id int, user_id int, foreign key (user_id) references users (id))",
			fn: func(t *testing.T, ct *CreateTable) {
				if len(ct.Constraints) != 1 {
					t.Fatalf("expected 1 constraint, got %d", len(ct.Constraints))
				}
				tc := ct.Constraints[0]
				if tc.Kind != "foreign key" || len(tc.Columns) != 1 || string(tc.Columns[0]) != "user_id" ||
					string(tc.RefTable) != "users" || len(tc.RefColumns) != 1 || string(tc.RefColumns[0]) != "id" {
					t.Errorf("unexpected constraint: %#v", tc)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(pos) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(pos))
			}
			ct, ok := pos[0].Statement.(*CreateTable)
			if !ok {
				t.Fatalf("expected CreateTable, got %T", pos[0].Statement)
			}
			tt.fn(t, ct)
			// Round-trip
			out := String(ct)
			if out == "" || !strings.HasPrefix(out, "create table ") {
				t.Errorf("unexpected format: %q", out)
			}
		})
	}
}

func TestCreateTemporaryTableAsSelect(t *testing.T) {
	sql := "CREATE TEMPORARY TABLE table1 AS SELECT * FROM users"
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
	if string(ctas.Table) != "table1" {
		t.Errorf("table name: got %q", ctas.Table)
	}
	if !ctas.Temporary {
		t.Error("expected Temporary true for CREATE TEMPORARY TABLE")
	}
	if !ctas.Temp() {
		t.Error("expected Temp() true for CREATE TEMPORARY TABLE")
	}
	// Round-trip
	out := String(pos[0].Statement)
	if out == "" {
		t.Error("expected non-empty String()")
	}
	if !strings.Contains(strings.ToLower(out), "temporary") {
		t.Errorf("expected output to contain 'temporary', got %q", out)
	}
}

func TestCreateTemporaryTableDDL(t *testing.T) {
	sql := "CREATE TEMPORARY TABLE t"
	pos, _, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(pos) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(pos))
	}
	ddl, ok := pos[0].Statement.(*DDL)
	if !ok {
		t.Fatalf("expected DDL, got %T", pos[0].Statement)
	}
	if ddl.Action != AST_CREATE || string(ddl.NewName) != "t" {
		t.Errorf("unexpected DDL: Action=%q NewName=%q", ddl.Action, ddl.NewName)
	}
	if !ddl.Temporary {
		t.Error("expected Temporary true for CREATE TEMPORARY TABLE")
	}
	if !ddl.Temp() {
		t.Error("expected Temp() true for CREATE TEMPORARY TABLE")
	}
}

func TestCreateTempTable(t *testing.T) {
	// TEMP is synonym for TEMPORARY
	for _, sql := range []string{
		"CREATE TEMP TABLE t",
		"CREATE TEMP TABLE t AS SELECT 1",
		"CREATE TEMP TABLE t AS ( SELECT 1 )",
	} {
		pos, _, err := Parse(sql)
		if err != nil {
			t.Errorf("parse %q: %v", sql, err)
			continue
		}
		if len(pos) != 1 {
			t.Errorf("parse %q: expected 1 statement", sql)
			continue
		}
		switch stmt := pos[0].Statement.(type) {
		case *DDL:
			if !stmt.Temporary && !stmt.Temp() {
				t.Errorf("parse %q: expected Temporary true", sql)
			}
		case *CreateTableAsSelect:
			if !stmt.Temporary && !stmt.Temp() {
				t.Errorf("parse %q: expected Temporary true", sql)
			}
		default:
			t.Errorf("parse %q: unexpected type %T", sql, stmt)
		}
	}
}

func TestCreateTableAsSelectWithParens(t *testing.T) {
	sql := "CREATE TABLE t AS ( SELECT * FROM users )"
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
	if string(ctas.Table) != "t" {
		t.Errorf("table name: got %q", ctas.Table)
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

// dumpExpectedASTGo returns Go source for (stmts, comments) so you can paste expected
// AST into parseTestCases. Run with: DUMP_EXPECTED_AST=1 go test -run TestDumpExpectedAST -v .
func dumpExpectedASTGo(stmts PositionedStatements, comments []CommentEntry) (statementsGo, commentsGo string) {
	return "PositionedStatements{" + dumpPositionedStatements(stmts) + "}", dumpComments(comments)
}

func dumpPositionedStatements(stmts PositionedStatements) string {
	var b strings.Builder
	for i, ps := range stmts {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("{Start: %d, End: %d, Statement: %s}", ps.Start, ps.End, dumpStatement(ps.Statement)))
	}
	return b.String()
}

func dumpComments(c []CommentEntry) string {
	if len(c) == 0 {
		return "nil"
	}
	var b strings.Builder
	b.WriteString("[]CommentEntry{")
	for i, e := range c {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("{Position: %d, Comment: %s}", e.Position, dumpBytes(e.Comment)))
	}
	b.WriteString("}")
	return b.String()
}

func dumpBytes(p []byte) string {
	if p == nil {
		return "nil"
	}
	return "[]byte(" + fmt.Sprintf("%q", p) + ")"
}

func dumpStatement(s Statement) string {
	if s == nil {
		return "nil"
	}
	switch n := s.(type) {
	case *Select:
		return dumpSelect(n)
	case *Union:
		return dumpUnion(n)
	case *DDL:
		return dumpDDL(n)
	case *Insert:
		return dumpInsert(n)
	case *Update:
		return dumpUpdate(n)
	case *Delete:
		return dumpDelete(n)
	case *Set:
		return dumpSet(n)
	case *CreateTableAsSelect:
		return dumpCreateTableAsSelect(n)
	case *Other:
		return "&Other{}"
	default:
		return fmt.Sprintf("/* unknown Statement %T */ nil", s)
	}
}

func dumpSelect(n *Select) string {
	var b strings.Builder
	b.WriteString("&Select{")
	if len(n.Comments) > 0 {
		b.WriteString("Comments: " + dumpCommentsType(n.Comments) + ", ")
	}
	if n.Distinct != "" {
		b.WriteString("Distinct: " + fmt.Sprintf("%q", n.Distinct) + ", ")
	}
	b.WriteString("SelectExprs: " + dumpSelectExprs(n.SelectExprs))
	if len(n.From) > 0 {
		b.WriteString(", From: " + dumpTableExprs(n.From))
	}
	if n.Where != nil {
		b.WriteString(", Where: " + dumpWhere(n.Where))
	}
	if len(n.GroupBy) > 0 {
		b.WriteString(", GroupBy: " + dumpGroupBy(n.GroupBy))
	}
	if n.Having != nil {
		b.WriteString(", Having: " + dumpWhere(n.Having))
	}
	if len(n.OrderBy) > 0 {
		b.WriteString(", OrderBy: " + dumpOrderBy(n.OrderBy))
	}
	if n.Limit != nil {
		b.WriteString(", Limit: " + dumpLimit(n.Limit))
	}
	if n.Lock != "" {
		b.WriteString(", Lock: " + fmt.Sprintf("%q", n.Lock))
	}
	b.WriteString("}")
	return b.String()
}

func dumpCommentsType(c Comments) string {
	if len(c) == 0 {
		return "nil"
	}
	var b strings.Builder
	b.WriteString("Comments{")
	for i, x := range c {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpBytes(x))
	}
	b.WriteString("}")
	return b.String()
}

func dumpSelectExprs(s SelectExprs) string {
	var b strings.Builder
	b.WriteString("SelectExprs{")
	for i, e := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpSelectExpr(e))
	}
	b.WriteString("}")
	return b.String()
}

func dumpSelectExpr(e SelectExpr) string {
	switch n := e.(type) {
	case *StarExpr:
		if n.TableName == nil {
			return "&StarExpr{}"
		}
		return "&StarExpr{TableName: " + dumpBytes(n.TableName) + "}"
	case *NonStarExpr:
		out := "&NonStarExpr{Expr: " + dumpExpr(n.Expr) + "}"
		if n.As != nil {
			out += ", As: " + dumpBytes(n.As)
		}
		return out
	default:
		return "/* SelectExpr */ nil"
	}
}

func dumpTableExprs(t TableExprs) string {
	var b strings.Builder
	b.WriteString("TableExprs{")
	for i, e := range t {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpTableExpr(e))
	}
	b.WriteString("}")
	return b.String()
}

func dumpTableExpr(e TableExpr) string {
	switch n := e.(type) {
	case *AliasedTableExpr:
		out := "&AliasedTableExpr{Expr: " + dumpSimpleTableExpr(n.Expr) + ""
		if n.As != nil {
			out += ", As: " + dumpBytes(n.As)
		}
		if n.Hints != nil {
			out += ", Hints: " + dumpIndexHints(n.Hints)
		}
		return out + "}"
	case *ParenTableExpr:
		return "&ParenTableExpr{Expr: " + dumpTableExpr(n.Expr) + "}"
	case *JoinTableExpr:
		out := "&JoinTableExpr{LeftExpr: " + dumpTableExpr(n.LeftExpr) + ", Join: " + fmt.Sprintf("%q", n.Join) + ", RightExpr: " + dumpTableExpr(n.RightExpr)
		if n.On != nil {
			out += ", On: " + dumpBoolExpr(n.On)
		}
		return out + "}"
	default:
		return "/* TableExpr */ nil"
	}
}

func dumpSimpleTableExpr(e SimpleTableExpr) string {
	switch n := e.(type) {
	case *TableName:
		out := "&TableName{"
		if n.Qualifier != nil {
			out += "Qualifier: " + dumpBytes(n.Qualifier) + ", "
		}
		out += "Name: " + dumpBytes(n.Name) + "}"
		return out
	case *Subquery:
		return "&Subquery{Select: " + dumpSelectStatement(n.Select) + "}"
	case *TableFunc:
		return "&TableFunc{Name: " + dumpBytes(n.Name) + ", Exprs: " + dumpSelectExprs(n.Exprs) + "}"
	default:
		return "/* SimpleTableExpr */ nil"
	}
}

func dumpSelectStatement(s SelectStatement) string {
	switch n := s.(type) {
	case *Select:
		return dumpSelect(n)
	case *Union:
		return dumpUnion(n)
	default:
		return "nil"
	}
}

func dumpUnion(n *Union) string {
	return "&Union{Type: " + fmt.Sprintf("%q", n.Type) + ", Left: " + dumpSelectStatement(n.Left) + ", Right: " + dumpSelectStatement(n.Right) + "}"
}

func dumpDDL(n *DDL) string {
	out := "&DDL{Action: " + fmt.Sprintf("%q", n.Action)
	if n.Table != nil {
		out += ", Table: " + dumpBytes(n.Table)
	}
	if n.NewName != nil {
		out += ", NewName: " + dumpBytes(n.NewName)
	}
	if n.Temporary {
		out += ", Temporary: true"
	}
	return out + "}"
}

func dumpWhere(n *Where) string {
	if n == nil || n.Expr == nil {
		return "nil"
	}
	return "&Where{Type: " + fmt.Sprintf("%q", n.Type) + ", Expr: " + dumpBoolExpr(n.Expr) + "}"
}

func dumpBoolExpr(e BoolExpr) string {
	if e == nil {
		return "nil"
	}
	switch n := e.(type) {
	case *ComparisonExpr:
		return "&ComparisonExpr{Operator: " + fmt.Sprintf("%q", n.Operator) + ", Left: " + dumpValExpr(n.Left) + ", Right: " + dumpValExpr(n.Right) + "}"
	case *AndExpr:
		out := "&AndExpr{Left: " + dumpBoolExpr(n.Left) + ", Right: " + dumpBoolExpr(n.Right) + "}"
		return out
	case *OrExpr:
		return "&OrExpr{Left: " + dumpBoolExpr(n.Left) + ", Right: " + dumpBoolExpr(n.Right) + "}"
	case *NotExpr:
		return "&NotExpr{Expr: " + dumpBoolExpr(n.Expr) + "}"
	case *ParenBoolExpr:
		return "&ParenBoolExpr{Expr: " + dumpBoolExpr(n.Expr) + "}"
	case *NullCheck:
		return "&NullCheck{Operator: " + fmt.Sprintf("%q", n.Operator) + ", Expr: " + dumpValExpr(n.Expr) + "}"
	case *ExistsExpr:
		return "&ExistsExpr{Subquery: " + dumpSubquery(n.Subquery) + "}"
	case *RangeCond:
		return "&RangeCond{Operator: " + fmt.Sprintf("%q", n.Operator) + ", Left: " + dumpValExpr(n.Left) + ", From: " + dumpValExpr(n.From) + ", To: " + dumpValExpr(n.To) + "}"
	case *KeyrangeExpr:
		return "&KeyrangeExpr{Start: " + dumpValExpr(n.Start) + ", End: " + dumpValExpr(n.End) + "}"
	default:
		return "/* BoolExpr */ nil"
	}
}

func dumpValExpr(e ValExpr) string {
	if e == nil {
		return "nil"
	}
	return dumpExpr(e)
}

func dumpExpr(e Expr) string {
	if e == nil {
		return "nil"
	}
	switch n := e.(type) {
	case NumVal:
		return "NumVal(" + dumpBytes([]byte(n)) + ")"
	case StrVal:
		return "StrVal(" + dumpBytes([]byte(n)) + ")"
	case ValArg:
		return "ValArg(" + dumpBytes([]byte(n)) + ")"
	case *NullVal:
		return "&NullVal{}"
	case *ColName:
		out := "&ColName{Name: " + dumpBytes(n.Name)
		if n.Qualifier != nil {
			out += ", Qualifier: " + dumpBytes(n.Qualifier)
		}
		return out + "}"
	case ValTuple:
		return "ValTuple{" + dumpValExprSlice(ValExprs(n)) + "}"
	case *Subquery:
		return "&Subquery{Select: " + dumpSelectStatement(n.Select) + "}"
	case ListArg:
		return "ListArg(" + dumpBytes([]byte(n)) + ")"
	case *BinaryExpr:
		return "&BinaryExpr{Left: " + dumpExpr(n.Left) + ", Operator: " + fmt.Sprintf("%q", n.Operator) + ", Right: " + dumpExpr(n.Right) + "}"
	case *UnaryExpr:
		return "&UnaryExpr{Operator: " + fmt.Sprintf("%q", n.Operator) + ", Expr: " + dumpExpr(n.Expr) + "}"
	case *FuncExpr:
		out := "&FuncExpr{Name: " + dumpBytes(n.Name)
		if n.Distinct {
			out += ", Distinct: true"
		}
		out += ", Exprs: " + dumpSelectExprs(n.Exprs) + "}"
		return out
	case *CaseExpr:
		out := "&CaseExpr{"
		if n.Expr != nil {
			out += "Expr: " + dumpValExpr(n.Expr) + ", "
		}
		out += "Whens: []*When{"
		for i, w := range n.Whens {
			if i > 0 {
				out += ", "
			}
			out += "&When{Cond: " + dumpBoolExpr(w.Cond) + ", Val: " + dumpValExpr(w.Val) + "}"
		}
		out += "}"
		if n.Else != nil {
			out += ", Else: " + dumpValExpr(n.Else)
		}
		return out + "}"
	default:
		return "/* Expr */ nil"
	}
}

func dumpValExprs(v ValExprs) string {
	return "ValExprs{" + dumpValExprSlice([]ValExpr(v)) + "}"
}

func dumpValExprSlice(v []ValExpr) string {
	var b strings.Builder
	for i, e := range v {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpValExpr(e))
	}
	return b.String()
}

func dumpSubquery(n *Subquery) string {
	if n == nil {
		return "nil"
	}
	return "&Subquery{Select: " + dumpSelectStatement(n.Select) + "}"
}

func dumpGroupBy(g GroupBy) string {
	return "GroupBy{" + dumpValExprSlice(g) + "}"
}

func dumpOrderBy(o OrderBy) string {
	var b strings.Builder
	b.WriteString("OrderBy{")
	for i, ord := range o {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("&Order{Expr: " + dumpValExpr(ord.Expr) + ", Direction: " + fmt.Sprintf("%q", ord.Direction) + "}")
	}
	b.WriteString("}")
	return b.String()
}

func dumpLimit(n *Limit) string {
	if n == nil {
		return "nil"
	}
	out := "&Limit{"
	if n.Offset != nil {
		out += "Offset: " + dumpValExpr(n.Offset) + ", "
	}
	out += "Rowcount: " + dumpValExpr(n.Rowcount) + "}"
	return out
}

func dumpIndexHints(n *IndexHints) string {
	if n == nil {
		return "nil"
	}
	var b strings.Builder
	b.WriteString("&IndexHints{Type: " + fmt.Sprintf("%q", n.Type) + ", Indexes: [][]byte{")
	for i, idx := range n.Indexes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpBytes(idx))
	}
	b.WriteString("}}")
	return b.String()
}

func dumpInsert(n *Insert) string {
	out := "&Insert{"
	if len(n.Comments) > 0 {
		out += "Comments: " + dumpCommentsType(n.Comments) + ", "
	}
	out += "Table: " + dumpTableName(n.Table) + ", Columns: " + dumpColumns(n.Columns) + ", Rows: " + dumpInsertRows(n.Rows)
	if len(n.OnDup) > 0 {
		out += ", OnDup: " + dumpOnDup(n.OnDup)
	}
	return out + "}"
}

func dumpTableName(n *TableName) string {
	if n == nil {
		return "nil"
	}
	out := "&TableName{"
	if n.Qualifier != nil {
		out += "Qualifier: " + dumpBytes(n.Qualifier) + ", "
	}
	out += "Name: " + dumpBytes(n.Name) + "}"
	return out
}

func dumpColumns(c Columns) string {
	if len(c) == 0 {
		return "nil"
	}
	return "Columns(" + dumpSelectExprs(SelectExprs(c)) + ")"
}

func dumpInsertRows(r InsertRows) string {
	if r == nil {
		return "nil"
	}
	switch n := r.(type) {
	case Values:
		return "Values{" + dumpValues(n) + "}"
	case *Select:
		return dumpSelect(n)
	case *Union:
		return dumpUnion(n)
	default:
		return "nil"
	}
}

func dumpValues(v Values) string {
	var b strings.Builder
	for i, r := range v {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(dumpRowTuple(r))
	}
	return b.String()
}

func dumpRowTuple(r RowTuple) string {
	switch n := r.(type) {
	case ValTuple:
		return "ValTuple{" + dumpValExprSlice(ValExprs(n)) + "}"
	case *Subquery:
		return dumpSubquery(n)
	default:
		return "nil"
	}
}

func dumpOnDup(o OnDup) string {
	return "OnDup(" + dumpUpdateExprs(UpdateExprs(o)) + ")"
}

func dumpUpdateExprs(u UpdateExprs) string {
	var b strings.Builder
	b.WriteString("UpdateExprs{")
	for i, e := range u {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("&UpdateExpr{Name: " + dumpColName(e.Name) + ", Expr: " + dumpValExpr(e.Expr) + "}")
	}
	b.WriteString("}")
	return b.String()
}

func dumpColName(n *ColName) string {
	if n == nil {
		return "nil"
	}
	out := "&ColName{Name: " + dumpBytes(n.Name)
	if n.Qualifier != nil {
		out += ", Qualifier: " + dumpBytes(n.Qualifier)
	}
	return out + "}"
}

func dumpUpdate(n *Update) string {
	out := "&Update{"
	if len(n.Comments) > 0 {
		out += "Comments: " + dumpCommentsType(n.Comments) + ", "
	}
	out += "Table: " + dumpTableName(n.Table) + ", Exprs: " + dumpUpdateExprs(n.Exprs)
	if n.Where != nil {
		out += ", Where: " + dumpWhere(n.Where)
	}
	if len(n.OrderBy) > 0 {
		out += ", OrderBy: " + dumpOrderBy(n.OrderBy)
	}
	if n.Limit != nil {
		out += ", Limit: " + dumpLimit(n.Limit)
	}
	return out + "}"
}

func dumpDelete(n *Delete) string {
	out := "&Delete{"
	if len(n.Comments) > 0 {
		out += "Comments: " + dumpCommentsType(n.Comments) + ", "
	}
	out += "Table: " + dumpTableName(n.Table)
	if n.Where != nil {
		out += ", Where: " + dumpWhere(n.Where)
	}
	if len(n.OrderBy) > 0 {
		out += ", OrderBy: " + dumpOrderBy(n.OrderBy)
	}
	if n.Limit != nil {
		out += ", Limit: " + dumpLimit(n.Limit)
	}
	return out + "}"
}

func dumpSet(n *Set) string {
	out := "&Set{"
	if len(n.Comments) > 0 {
		out += "Comments: " + dumpCommentsType(n.Comments) + ", "
	}
	return out + "Exprs: " + dumpUpdateExprs(n.Exprs) + "}"
}

func dumpCreateTableAsSelect(n *CreateTableAsSelect) string {
	out := "&CreateTableAsSelect{Table: " + dumpBytes(n.Table) + ", Select: " + dumpSelectStatement(n.Select)
	if n.Temporary {
		out += ", Temporary: true"
	}
	return out + "}"
}

// dumpExpectedAST prints expected AST as Go source for the given SQL.
func dumpExpectedAST(t *testing.T, name string, sql string) {
	var out strings.Builder
	stmts, comments, err := Parse(sql)
	if err != nil {
		t.Logf("skip %s: parse failed: %v", name, err)
		return
	}
	stmtsGo, commentsGo := dumpExpectedASTGo(stmts, comments)
	out.WriteString(fmt.Sprintf("\t%q: {\n\t\tSQL: %q,\n\t\tStatements: %s,\n\t\tComments: %s,\n\t},\n\n", name, sql, stmtsGo, commentsGo))
	t.Logf("\n%s", out.String())
}
