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

package gomod

import (
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// astEq checks if two go.mod ASTs are equivalent.
func astEq(baseAst *modfile.File, lastAst *modfile.File) bool {
	if !modEq(baseAst.Module, lastAst.Module) || !goEq(baseAst.Go, lastAst.Go) {
		return false
	}

	if len(baseAst.Require) != len(lastAst.Require) {
		return false
	}
	for i, baseReq := range baseAst.Require {
		if !reqEq(baseReq, lastAst.Require[i]) {
			return false
		}
	}

	if len(baseAst.Exclude) != len(lastAst.Exclude) {
		return false
	}
	for i, baseExcl := range baseAst.Exclude {
		if !exclEq(baseExcl, lastAst.Exclude[i]) {
			return false
		}
	}

	if len(baseAst.Replace) != len(lastAst.Replace) {
		return false
	}
	for i, baseRep := range baseAst.Replace {
		if !repEq(baseRep, lastAst.Replace[i]) {
			return false
		}
	}

	if len(baseAst.Retract) != len(lastAst.Retract) {
		return false
	}
	for i, baseRet := range baseAst.Retract {
		if !retEq(baseRet, lastAst.Retract[i]) {
			return false
		}
	}

	// The modfile package documentations says that FileSyntax
	// "represents an entire go.mod file" so verification of the AST
	// nodes below may be redundant, but we do it just to be safe.
	return syntaxEq(baseAst.Syntax, lastAst.Syntax)
}

// modEq checks if two Module AST components are equivalent.
func modEq(base *modfile.Module, last *modfile.Module) bool {
	return modVersionEq(base.Mod, last.Mod) && lineEq(base.Syntax, last.Syntax)
}

// modVersionEq checks if two module.Version AST components are equivalent.
func modVersionEq(base module.Version, last module.Version) bool {
	return base.Path == last.Path && base.Version == last.Version
}

// goEq checks if two Go AST components are equivalent.
func goEq(base *modfile.Go, last *modfile.Go) bool {
	return base.Version == last.Version && lineEq(base.Syntax, last.Syntax)
}

// reqEq checks if two require AST components are equivalent.
func reqEq(base *modfile.Require, last *modfile.Require) bool {
	return modVersionEq(base.Mod, last.Mod) &&
		base.Indirect == last.Indirect &&
		lineEq(base.Syntax, last.Syntax)
}

// exclEq checks if two Exclude AST components are equivalent.
func exclEq(base *modfile.Exclude, last *modfile.Exclude) bool {
	return modVersionEq(base.Mod, last.Mod) && lineEq(base.Syntax, last.Syntax)
}

// repEq checks if two Replace AST components are equivalent.
func repEq(base *modfile.Replace, last *modfile.Replace) bool {
	return modVersionEq(base.Old, last.Old) &&
		modVersionEq(base.New, last.New) &&
		lineEq(base.Syntax, last.Syntax)
}

// retEq checks if two Retract AST components are equivalent.
func retEq(base *modfile.Retract, last *modfile.Retract) bool {
	return base.Low == last.Low &&
		base.High == last.High &&
		base.Rationale == last.Rationale &&
		lineEq(base.Syntax, last.Syntax)
}

// syntaxEq checks if two FileSyntax AST components are equivalent.
func syntaxEq(base *modfile.FileSyntax, last *modfile.FileSyntax) bool {
	// FileSyntax also contains strings representing paths to the
	// files represented by ASTs but we don't compare those here (as
	// they would never match for tests). The file names themselves
	// are still compared when constructing the ASTs (in astBuild) as
	// a sanity check.
	if len(base.Stmt) != len(last.Stmt) {
		return false
	}
	for i, baseExpr := range base.Stmt {
		if !exprEq(baseExpr, last.Stmt[i]) {
			return false
		}
	}

	return true
}

// exprEq checks if two Expr AST components are equivalent.
func exprEq(base modfile.Expr, last modfile.Expr) bool {
	switch baseExpr := base.(type) {
	case *modfile.CommentBlock:
		if lastExpr, ok := last.(*modfile.CommentBlock); ok {
			return commentBlockEq(baseExpr, lastExpr)
		}
	case *modfile.LParen:
		if lastExpr, ok := last.(*modfile.LParen); ok {
			return lParenEq(baseExpr, lastExpr)
		}
	case *modfile.Line:
		if lastExpr, ok := last.(*modfile.Line); ok {
			return lineEq(baseExpr, lastExpr)
		}
	case *modfile.LineBlock:
		if lastExpr, ok := last.(*modfile.LineBlock); ok {
			return lineBlockEq(baseExpr, lastExpr)
		}
	case *modfile.RParen:
		if lastExpr, ok := last.(*modfile.RParen); ok {
			return rParenEq(baseExpr, lastExpr)
		}
	}
	return false
}

// commentBlockEq checks if two CommentBlock AST components are equivalent.
func commentBlockEq(base *modfile.CommentBlock, last *modfile.CommentBlock) bool {
	// we ignore these (they contain only comments and posistions) so
	// they are always equivalent
	return true
}

// lParebEq checks if two LParen AST components are equivalent.
func lParenEq(base *modfile.LParen, last *modfile.LParen) bool {
	// we ignore these (they contain only comments and posistions) so
	// they are always equivalent
	return true
}

// lineEq checks if two Line AST components are equivalent.
func lineEq(base *modfile.Line, last *modfile.Line) bool {
	// a bit surprisingly, Token is actually an array of string-s
	return strArrayEq(base.Token, last.Token) && base.InBlock == last.InBlock
}

// lineBlockEq checks if two LineBlock AST components are equivalent.
func lineBlockEq(base *modfile.LineBlock, last *modfile.LineBlock) bool {
	if !lParenEq(&base.LParen, &last.LParen) ||
		// a bit surprisingly, Token is actually an array of string-s
		!strArrayEq(base.Token, last.Token) ||
		!rParenEq(&base.RParen, &last.RParen) {
		return false
	}

	if len(base.Line) != len(last.Line) {
		return false
	}
	for i, baseLine := range base.Line {
		if !lineEq(baseLine, last.Line[i]) {
			return false
		}
	}

	return true
}

// rParebEq checks if two RParen AST components are equivalent.
func rParenEq(base *modfile.RParen, last *modfile.RParen) bool {
	// we ignore these (they contain only comments and posistions) so
	// they are always equivalent
	return true
}

// strArrayEq checks if two string arrays are equivalent.
func strArrayEq(base []string, last []string) bool {
	if len(base) != len(last) {
		return false
	}
	for i, baseStr := range base {
		if baseStr != last[i] {
			return false
		}
	}
	return true
}
