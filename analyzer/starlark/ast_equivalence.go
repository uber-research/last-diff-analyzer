//  Copyright (c) 2023 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package starlark

import (
	"bytes"
	"encoding/json"

	"go.starlark.net/syntax"
)

// astEq checks if two starlark ASTs are equivalent.
func astEq(baseAst, lastAst *syntax.File) (bool, error) {
	// First strip the ASTs
	syntax.Walk(baseAst, strip)
	syntax.Walk(lastAst, strip)

	// Since starlark/syntax library does not provide an easy way to unparse the ASTs, we simply
	// fall back to using json.Marshal to compare if two stripped ASTs are equal.
	baseBytes, err := json.Marshal(baseAst.Stmts)
	if err != nil {
		return false, err
	}
	lastBytes, err := json.Marshal(lastAst.Stmts)
	if err != nil {
		return false, err
	}

	a, b := string(baseBytes), string(lastBytes)
	print(a, b)

	return bytes.Equal(baseBytes, lastBytes), nil
}

// _emptyPos is an empty Position struct to be used to clear the position information for each
// starlark AST node.
var _emptyPos syntax.Position

// strip clears out position information and docstrings in the starlark AST. We assume the AST is
// parsed without RetainComments option (see go.starlark.net/syntax/parse.go), so we do _not_
// strip comments in this function (docstrings will still be properly stripped, since they are
// technically raw strings).
func strip(node syntax.Node) bool {
	switch n := node.(type) {
	case *syntax.File:
		n.Stmts = stripDocstring(n.Stmts)

	// Statements
	case *syntax.AssignStmt:
		n.OpPos = _emptyPos
	case *syntax.BranchStmt:
		n.TokenPos = _emptyPos
	case *syntax.DefStmt:
		n.Def = _emptyPos
		n.Function = nil
		n.Body = stripDocstring(n.Body)
	case *syntax.ExprStmt:
		// no-op
	case *syntax.ForStmt:
		n.For = _emptyPos
		n.Body = stripDocstring(n.Body)
	case *syntax.WhileStmt:
		n.While = _emptyPos
		n.Body = stripDocstring(n.Body)
	case *syntax.IfStmt:
		n.If = _emptyPos
		n.ElsePos = _emptyPos
		n.True = stripDocstring(n.True)
		n.False = stripDocstring(n.False)
	case *syntax.LoadStmt:
		n.Load = _emptyPos
		n.Rparen = _emptyPos
	case *syntax.ReturnStmt:
		n.Return = _emptyPos

	// Expressions
	case *syntax.BinaryExpr:
		n.OpPos = _emptyPos
	case *syntax.CallExpr:
		n.Rparen = _emptyPos
		n.Lparen = _emptyPos
	case *syntax.Comprehension:
		n.Lbrack = _emptyPos
		n.Rbrack = _emptyPos
	case *syntax.CondExpr:
		n.If = _emptyPos
		n.ElsePos = _emptyPos
	case *syntax.DictEntry:
		n.Colon = _emptyPos
	case *syntax.DictExpr:
		n.Rbrace = _emptyPos
		n.Lbrace = _emptyPos
	case *syntax.DotExpr:
		n.Dot = _emptyPos
		n.NamePos = _emptyPos
	case *syntax.Ident:
		n.NamePos = _emptyPos
	case *syntax.IndexExpr:
		n.Rbrack = _emptyPos
		n.Lbrack = _emptyPos
	case *syntax.LambdaExpr:
		n.Lambda = _emptyPos
	case *syntax.ListExpr:
		n.Rbrack = _emptyPos
		n.Lbrack = _emptyPos
	case *syntax.Literal:
		n.TokenPos = _emptyPos
	case *syntax.ParenExpr:
		n.Rparen = _emptyPos
		n.Lparen = _emptyPos
	case *syntax.SliceExpr:
		n.Rbrack = _emptyPos
		n.Lbrack = _emptyPos
	case *syntax.TupleExpr:
		n.Rparen = _emptyPos
		n.Lparen = _emptyPos
	case *syntax.UnaryExpr:
		n.OpPos = _emptyPos

	// Clauses
	case *syntax.ForClause:
		n.For = _emptyPos
		n.In = _emptyPos
	case *syntax.IfClause:
		n.If = _emptyPos
	}

	return true
}

// stripDocstring removes any docstrings that is wrapped in ExprStmt in the statement slice and
// return the updated slice.
func stripDocstring(stmts []syntax.Stmt) []syntax.Stmt {
	var newStmts []syntax.Stmt
	for _, stmt := range stmts {

		if exprStmt, ok := stmt.(*syntax.ExprStmt); ok {
			if docstring, ok := exprStmt.X.(*syntax.Literal); ok && docstring.Token == syntax.STRING {
				continue
			}
		}
		newStmts = append(newStmts, stmt)
	}

	return newStmts
}
