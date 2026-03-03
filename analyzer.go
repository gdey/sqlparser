// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlparser

// analyzer.go contains utility analysis functions.

import (
	"fmt"

	"github.com/gdey/sqlparser/sqltypes"
)

// GetTableName returns the table name from the SimpleTableExpr
// only if it's a simple expression. Otherwise, it returns "".
func GetTableName(node SimpleTableExpr) string {
	if n, ok := node.(*TableName); ok && n.Qualifier == nil {
		return string(n.Name)
	}
	// sub-select or '.' expression
	return ""
}

// GetColName returns the column name, only if
// it's a simple expression. Otherwise, it returns "".
func GetColName(node Expr) string {
	if n, ok := node.(*ColName); ok {
		return string(n.Name)
	}
	return ""
}

// IsColName returns true if the ValExpr is a *ColName.
func IsColName(node ValExpr) bool {
	_, ok := node.(*ColName)
	return ok
}

// IsValue returns true if the ValExpr is a string, number or value arg.
// NULL is not considered to be a value.
func IsValue(node ValExpr) bool {
	switch node.(type) {
	case StrVal, NumVal, ValArg:
		return true
	}
	return false
}

// HasINCaluse returns true if any of the conditions has an IN clause.
func HasINClause(conditions []BoolExpr) bool {
	for _, node := range conditions {
		if c, ok := node.(*ComparisonExpr); ok && c.Operator == AST_IN {
			return true
		}
	}
	return false
}

// IsSimpleTuple returns true if the ValExpr is a ValTuple that
// contains simple values or if it's a list arg.
func IsSimpleTuple(node ValExpr) bool {
	switch vals := node.(type) {
	case ValTuple:
		for _, n := range vals {
			if !IsValue(n) {
				return false
			}
		}
		return true
	case ListArg:
		return true
	}
	// It's a subquery
	return false
}

// AsInterface converts the ValExpr to an interface. It converts
// ValTuple to []interface{}, ValArg to string, StrVal to sqltypes.String,
// NumVal to sqltypes.Numeric, NullVal to nil.
// Otherwise, it returns an error.
func AsInterface(node ValExpr) (interface{}, error) {
	switch node := node.(type) {
	case ValTuple:
		vals := make([]interface{}, 0, len(node))
		for _, val := range node {
			v, err := AsInterface(val)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}
		return vals, nil
	case ValArg:
		return string(node), nil
	case ListArg:
		return string(node), nil
	case StrVal:
		return sqltypes.MakeString(node), nil
	case NumVal:
		n, err := sqltypes.BuildNumeric(string(node))
		if err != nil {
			return nil, fmt.Errorf("type mismatch: %s", err)
		}
		return n, nil
	case *NullVal:
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected node %v", node)
}

// StringIn is a convenience function that returns
// true if str matches any of the values.
func StringIn(str string, values ...string) bool {
	for _, val := range values {
		if str == val {
			return true
		}
	}
	return false
}

// TableNamesFromTableExprs returns the set of table names referenced in the
// FROM clause (and in any nested JOINs or subqueries). Names are returned
// as single identifiers (e.g. "parcels_data"); qualified names use the Name
// only (the rightmost identifier).
func TableNamesFromTableExprs(from TableExprs) map[string]struct{} {
	out := make(map[string]struct{})
	for _, e := range from {
		tableNamesFromTableExpr(e, out)
	}
	return out
}

func tableNamesFromTableExpr(expr TableExpr, out map[string]struct{}) {
	switch e := expr.(type) {
	case *AliasedTableExpr:
		tableNamesFromSimpleTableExpr(e.Expr, out)
	case *JoinTableExpr:
		tableNamesFromTableExpr(e.LeftExpr, out)
		tableNamesFromTableExpr(e.RightExpr, out)
	case *ParenTableExpr:
		tableNamesFromTableExpr(e.Expr, out)
	}
}

func tableNamesFromSimpleTableExpr(expr SimpleTableExpr, out map[string]struct{}) {
	switch s := expr.(type) {
	case *TableName:
		out[string(s.Name)] = struct{}{}
	case *Subquery:
		tableNamesFromSelectStatement(s.Select, out)
	}
}

func tableNamesFromSelectStatement(sel SelectStatement, out map[string]struct{}) {
	switch s := sel.(type) {
	case *Select:
		if s.From != nil {
			for _, e := range s.From {
				tableNamesFromTableExpr(e, out)
			}
		}
	case *Union:
		tableNamesFromSelectStatement(s.Left, out)
		tableNamesFromSelectStatement(s.Right, out)
	}
}

// SelectStatementReferencesAny returns true if the given SelectStatement's
// FROM clause (and nested subqueries) references any of the given table names.
func SelectStatementReferencesAny(sel SelectStatement, tableNames map[string]struct{}) bool {
	refs := make(map[string]struct{})
	tableNamesFromSelectStatement(sel, refs)
	for name := range refs {
		if _, ok := tableNames[name]; ok {
			return true
		}
	}
	return false
}
