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

package bazel

import "github.com/bazelbuild/buildtools/build"

// depsEqResult defines result of comparing dependencies.
type depsEqResult int

// The following are dependencies comparison results.
const (
	_depsEqual = iota
	_depsAdded
	_depsRemoved
	_depsModified
	_depsIncomparable
)

// _analyzableRules contains names of rules that can be analyzed for
// potential auto-approval (any change to non-analyzable rule can ever
// be approved).
var _analyzableRules = [...]string{
	// Go-specific rules.
	"go_library", "go_binary", "go_test",
	// Java-specific rules.
	"java_library", "java_binary", "java_test",
	// Misc rules.
	"starlark_library"}

// _testCalls defines names of test-related calls in Bazel build
// files.
var _testCalls = [...]string{"go_test", "java_test"}

// astEq checks if two Bazel ASTs are equivalent.
func astEq(baseAst *build.File, lastAst *build.File) bool {
	if baseAst.Pkg != lastAst.Pkg {
		return false
	}
	if baseAst.Label != lastAst.Label {
		return false
	}
	if baseAst.WorkspaceRoot != lastAst.WorkspaceRoot {
		return false
	}
	if baseAst.Type != lastAst.Type {
		return false
	}
	j := 0
	i := 0
	// iterate over all top-level statements, but omit those related
	// to testing
	for i < len(baseAst.Stmt) && j < len(lastAst.Stmt) {
		if ignoreTopExpr(baseAst.Stmt[i]) {
			// ignore and reloop to check if more statements available
			i++
			continue
		}
		if ignoreTopExpr(lastAst.Stmt[j]) {
			// ignore and reloop to check if more statements available
			j++
			continue
		}
		if !topExprEq(baseAst.Stmt[i], lastAst.Stmt[j]) {
			return false
		}
		i++
		j++

	}
	if i < len(baseAst.Stmt) {
		// base file has more expressions to analyze
		return areRemainingIgnored(i, baseAst.Stmt)
	}
	if j < len(lastAst.Stmt) {
		// last file has more expressions to analyze
		return areRemainingIgnored(j, lastAst.Stmt)
	}

	// otherwise all expressions have been analyzed already
	return true

	// build.File has another field (Path), but we do not compare this
	// one as:
	// - in production, it is already guaranteed that we are comparing
	// two versions of the same file
	// - when testing, by neccessity, files to compare are on
	// different paths
}

// ignoreTopExpr determines if a given expression represents a
// test-related call in a Bazel build file.
func ignoreTopExpr(expr build.Expr) bool {
	return isTestCall(expr) || isCommentBlock(expr)
}

// isTestCall determines if a given expression represents a
// test-related call in a Bazel build file.
func isTestCall(expr build.Expr) bool {
	call, ok := expr.(*build.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.X.(*build.Ident)
	if !ok {
		return false
	}
	for _, testCall := range _testCalls {
		if ident.Name == testCall {
			return true
		}
	}
	return false
}

// isCommentBlock determines if a given expression represents a comment block
// in a Bazel build file.
func isCommentBlock(expr build.Expr) bool {
	_, ok := expr.(*build.CommentBlock)
	return ok
}

// areRemainingIgnored checks if remaining expressions in the array
// can be ignored (excluded from comparison).
func areRemainingIgnored(startIdx int, remaining []build.Expr) bool {
	for i := startIdx; i < len(remaining); i++ {
		if !ignoreTopExpr(remaining[i]) {
			return false
		}
	}
	return true
}

// topExprEq checks if two top-level Expr AST components are equivalent.
func topExprEq(base build.Expr, last build.Expr) bool {
	// We are looking for "deps" arguments in top-level calls
	// (CallExpr-s) to see if any deps have been removed. This custom
	// check involves determining that both expressions are indeed
	// calls, that they are equivalent other than their "deps"-s, and
	// proceding with custom comparison of "deps"-s. The fall-back at
	// each step (e.g., expression types are not CallExpr) is to do
	// the "default" comparison.
	baseCall, ok := base.(*build.CallExpr)
	if !ok {
		// not a CallExpr - default comparison
		return exprEq(base, last)
	}

	lastCall, ok := last.(*build.CallExpr)
	if !ok {
		// see comment directly above
		return exprEq(base, last)
	}

	var baseDep *build.Expr
	// extract (remove) "deps" argument from the base CallExpr
	baseDep, baseCall.List = extractDeps(baseCall)

	var lastDep *build.Expr
	// extract (remove) "deps" argument from the last CallExpr
	lastDep, lastCall.List = extractDeps(lastCall)

	// check if the call expressions match with the dep" argument removed
	if !callExprEq(baseCall, lastCall) {
		return false
	}

	// check the case of no deps in base diff call
	if baseDep == nil {
		if lastDep == nil {
			// no deps in either call
			return true
		}
		// Deps added - adding dependencies can only be auto-approved for a certain
		// set of rules. It's enough to check name of a rule in one
		// diff as equivalence of rules in both diffs (along with
		// their names) has already been checked earlier in this
		// function.
		return isAnalyzableRule(lastCall)
	}

	if lastDep == nil {
		// all deps in last diff removed
		return true
	}

	// check if dependencies lists are equivalent
	switch depsEq(*baseDep, *lastDep) {
	case _depsEqual, _depsRemoved:
		return true
	case _depsAdded:
		return isAnalyzableRule(lastCall)
	}
	return false
}

// isAnalyzableRule checks if a given rule is analyzable, that is if
// its modifications could potentially be auto-approved. No changes to
// non-analyzable rules can ever be auto-approved.
func isAnalyzableRule(callExpr *build.CallExpr) bool {
	ident, ok := callExpr.X.(*build.Ident)
	if !ok {
		return false
	}
	for _, n := range _analyzableRules {
		if ident.Name == n {
			return true
		}
	}
	return false
}

// extractDeps returns "deps" argument for a call (if any) and the
// arguments array with the "deps" argument removed (if found).
func extractDeps(callExpr *build.CallExpr) (*build.Expr, []build.Expr) {
	// newList ultimately stores all the same arguments as the
	// original callExpr.List with the exception of the "deps"
	// argument (if any)
	var newList []build.Expr
	for i, e := range callExpr.List {
		assign, ok := e.(*build.AssignExpr)
		if !ok {
			newList = append(newList, e)
			continue
		}
		ident, ok := assign.LHS.(*build.Ident)
		if !ok {
			newList = append(newList, e)
			continue
		}
		if ident.Name != "deps" {
			newList = append(newList, e)
			continue
		}

		// found "deps" argument
		l := len(callExpr.List)
		if i != l-1 {
			// not the last argument - skip the current one and append
			// the rest of arguments rather than looping unnecessarily
			// until the end of the list is reached
			newList = append(newList, callExpr.List[i+1:l]...)
		}
		return &assign.RHS, newList
	}
	return nil, newList
}

// depsEq performs equivalence check between two sets of
// dependencies.
func depsEq(base build.Expr, last build.Expr) depsEqResult {
	baseList, baseOK := base.(*build.ListExpr)
	lastList, lastOK := last.(*build.ListExpr)
	if !baseOK || !lastOK {
		// not a ListExpr - default comparison (unlikely to happen,
		// but going for default comparison will arguably handle the
		// situation more gracefully than a failed type check)
		if exprEq(base, last) {
			return _depsEqual
		}
		return _depsModified
	}

	baseStrings := depsStrings(baseList)
	lastStrings := depsStrings(lastList)

	if baseStrings == nil || lastStrings == nil {
		// at least one of the "deps" sets contains unrecognized
		// elements - bail out
		return _depsIncomparable
	}

	if len(baseStrings) < len(lastStrings) {
		// The last diff contains dependencies not present on the base
		// one one. We must distinguish between new dependencies added
		// ONLY (which is safe for some rules), and other
		// modifications.
		for dep := range baseStrings {
			if !lastStrings[dep] {
				// dependency changed between base and last diffs
				return _depsModified
			}
		}
		// at this point it is guaranteed that all dependencies
		// present in the base diff were also present in last diff (i.e.,
		// last diff only contains dependency additions)
		return _depsAdded
	}

	// The base diff may contain dependencies not present in the last
	// one. We must distinguish between dependencies removed ONLY
	// (which is safe for all rules), and other modifications.
	for dep := range lastStrings {
		if !baseStrings[dep] {
			// dependency changed between base and last diffs
			return _depsModified
		}
	}

	// at this point it is guaranteed that all dependencies present in
	// the last diff were also present in base diff

	if len(baseStrings) > len(lastStrings) {
		// some dependencies present in base diff were removed in last
		// diff
		return _depsRemoved
	}
	// dependency lists are identical
	return _depsEqual
}

// depsStrings constructs a set of strings representing dependencies
// (bails out and returns null if not all dependencies are strings).
func depsStrings(deps *build.ListExpr) map[string]bool {
	strs := make(map[string]bool)
	for _, e := range deps.List {
		if strExpr, ok := e.(*build.StringExpr); ok {
			strs[strExpr.Value] = true
		} else {
			return nil
		}
	}
	return strs
}

// exprEq checks if two Expr AST components are equivalent.
func exprEq(base build.Expr, last build.Expr) bool {
	switch baseExpr := base.(type) {
	case *build.AssignExpr:
		if lastExpr, ok := last.(*build.AssignExpr); ok {
			return assignExprEq(baseExpr, lastExpr)
		}
	case *build.BinaryExpr:
		if lastExpr, ok := last.(*build.BinaryExpr); ok {
			return binaryExprEq(baseExpr, lastExpr)
		}
	case *build.BranchStmt:
		if lastExpr, ok := last.(*build.BranchStmt); ok {
			return branchStmtEq(baseExpr, lastExpr)
		}
	case *build.CallExpr:
		if lastExpr, ok := last.(*build.CallExpr); ok {
			return callExprEq(baseExpr, lastExpr)
		}
	case *build.CommentBlock:
		if lastExpr, ok := last.(*build.CommentBlock); ok {
			return commentBlockEq(baseExpr, lastExpr)
		}
	case *build.Comprehension:
		if lastExpr, ok := last.(*build.Comprehension); ok {
			return comprehensionEq(baseExpr, lastExpr)
		}
	case *build.ConditionalExpr:
		if lastExpr, ok := last.(*build.ConditionalExpr); ok {
			return conditionalExprEq(baseExpr, lastExpr)
		}
	case *build.DefStmt:
		if lastExpr, ok := last.(*build.DefStmt); ok {
			return defStmtEq(baseExpr, lastExpr)
		}
	case *build.DictExpr:
		if lastExpr, ok := last.(*build.DictExpr); ok {
			return dictExprEq(baseExpr, lastExpr)
		}
	case *build.DotExpr:
		if lastExpr, ok := last.(*build.DotExpr); ok {
			return dotExprEq(baseExpr, lastExpr)
		}
	case *build.End:
		if lastExpr, ok := last.(*build.End); ok {
			return endEq(baseExpr, lastExpr)
		}
	case *build.ForClause:
		if lastExpr, ok := last.(*build.ForClause); ok {
			return forClauseEq(baseExpr, lastExpr)
		}
	case *build.ForStmt:
		if lastExpr, ok := last.(*build.ForStmt); ok {
			return forStmtEq(baseExpr, lastExpr)
		}
	case *build.Function:
		if lastExpr, ok := last.(*build.Function); ok {
			return fnEq(baseExpr, lastExpr)
		}
	case *build.Ident:
		if lastExpr, ok := last.(*build.Ident); ok {
			return identEq(baseExpr, lastExpr)
		}
	case *build.IfClause:
		if lastExpr, ok := last.(*build.IfClause); ok {
			return ifClauseEq(baseExpr, lastExpr)
		}
	case *build.IfStmt:
		if lastExpr, ok := last.(*build.IfStmt); ok {
			return ifStmtEq(baseExpr, lastExpr)
		}
	case *build.IndexExpr:
		if lastExpr, ok := last.(*build.IndexExpr); ok {
			return indexExprEq(baseExpr, lastExpr)
		}
	case *build.KeyValueExpr:
		if lastExpr, ok := last.(*build.KeyValueExpr); ok {
			return keyValExprEq(baseExpr, lastExpr)
		}
	case *build.LambdaExpr:
		if lastExpr, ok := last.(*build.LambdaExpr); ok {
			return lambdaExprEq(baseExpr, lastExpr)
		}
	case *build.ListExpr:
		if lastExpr, ok := last.(*build.ListExpr); ok {
			return listExprEq(baseExpr, lastExpr)
		}
	case *build.LiteralExpr:
		if lastExpr, ok := last.(*build.LiteralExpr); ok {
			return literalExprEq(baseExpr, lastExpr)
		}
	case *build.LoadStmt:
		if lastExpr, ok := last.(*build.LoadStmt); ok {
			return loadStmtEq(baseExpr, lastExpr)
		}
	case *build.ParenExpr:
		if lastExpr, ok := last.(*build.ParenExpr); ok {
			return parenExprEq(baseExpr, lastExpr)
		}
	case *build.ReturnStmt:
		if lastExpr, ok := last.(*build.ReturnStmt); ok {
			return returnStmtEq(baseExpr, lastExpr)
		}
	case *build.SetExpr:
		if lastExpr, ok := last.(*build.SetExpr); ok {
			return setExprEq(baseExpr, lastExpr)
		}
	case *build.SliceExpr:
		if lastExpr, ok := last.(*build.SliceExpr); ok {
			return sliceExprEq(baseExpr, lastExpr)
		}
	case *build.StringExpr:
		if lastExpr, ok := last.(*build.StringExpr); ok {
			return stringExprEq(baseExpr, lastExpr)
		}
	case *build.TupleExpr:
		if lastExpr, ok := last.(*build.TupleExpr); ok {
			return tupleExprEq(baseExpr, lastExpr)
		}
	case *build.UnaryExpr:
		if lastExpr, ok := last.(*build.UnaryExpr); ok {
			return unaryExprEq(baseExpr, lastExpr)
		}
	}
	return false

}

// assignExprEq checks if two AssignExpr AST components are equivalent.
func assignExprEq(base *build.AssignExpr, last *build.AssignExpr) bool {
	if !exprEq(base.LHS, last.LHS) {
		return false
	}
	if base.Op != last.Op {
		return false
	}
	if !exprEq(base.RHS, last.RHS) {
		return false
	}
	return true
}

// binaryExprEq checks if two BinaryExpr AST components are equivalent.
func binaryExprEq(base *build.BinaryExpr, last *build.BinaryExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if base.Op != last.Op {
		return false
	}
	if !exprEq(base.Y, last.Y) {
		return false
	}
	return true
}

// branchStmtEq checks if two BranchStmt AST components are equivalent.
func branchStmtEq(base *build.BranchStmt, last *build.BranchStmt) bool {
	if base.Token != last.Token {
		return false
	}
	return true
}

// callExprEq checks if two CallExpr AST components are equivalent.
func callExprEq(base *build.CallExpr, last *build.CallExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if !exprArrayEq(base.List, last.List) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// commentBlockEq checks if two CommentBlock AST components are equivalent.
func commentBlockEq(base *build.CommentBlock, last *build.CommentBlock) bool {
	// we ignore positions and comments, so just return true
	return true
}

// comprehensionEq checks if two Comprehension AST components are equivalent.
func comprehensionEq(base *build.Comprehension, last *build.Comprehension) bool {
	if !exprEq(base.Body, last.Body) {
		return false
	}
	if !exprArrayEq(base.Clauses, last.Clauses) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// conditionalExprEq checks if two ConditionalExpr AST components are equivalent.
func conditionalExprEq(base *build.ConditionalExpr, last *build.ConditionalExpr) bool {
	if !exprEq(base.Then, last.Then) {
		return false
	}
	if !exprEq(base.Test, last.Test) {
		return false
	}
	if !exprEq(base.Else, last.Else) {
		return false
	}
	return true
}

// defStmtEq checks if two DefStmt AST components are equivalent.
func defStmtEq(base *build.DefStmt, last *build.DefStmt) bool {
	if !fnEq(&base.Function, &last.Function) {
		return false
	}
	if base.Name != last.Name {
		return false
	}
	return true
}

// dictExprEq checks if two DictExpr AST components are equivalent.
func dictExprEq(base *build.DictExpr, last *build.DictExpr) bool {
	if len(base.List) != len(last.List) {
		return false
	}
	for i := 0; i < len(base.List); i++ {
		if !keyValExprEq(base.List[i], last.List[i]) {
			return false
		}
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// dotExprEq checks if two DotExpr AST components are equivalent.
func dotExprEq(base *build.DotExpr, last *build.DotExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if base.Name != last.Name {
		return false
	}
	return true
}

// endEq checks if two End AST components are equivalent.
func endEq(base *build.End, last *build.End) bool {
	// we ignore positions and comments, so just return true
	return true
}

// forClauseEq checks if two ForClause AST components are equivalent.
func forClauseEq(base *build.ForClause, last *build.ForClause) bool {
	if !exprEq(base.Vars, last.Vars) {
		return false
	}
	if !exprEq(base.X, last.X) {
		return false
	}
	return true
}

// forStmtEq checks if two ForStmt AST components are equivalent.
func forStmtEq(base *build.ForStmt, last *build.ForStmt) bool {
	if !fnEq(&base.Function, &last.Function) {
		return false
	}
	if !exprEq(base.Vars, last.Vars) {
		return false
	}
	if !exprEq(base.X, last.X) {
		return false
	}
	if !exprArrayEq(base.Body, last.Body) {
		return false
	}
	return true
}

// fnEq checks if two Function AST components are equivalent.
func fnEq(base *build.Function, last *build.Function) bool {
	if !exprArrayEq(base.Params, last.Params) {
		return false
	}
	if !exprArrayEq(base.Body, last.Body) {
		return false
	}
	return true
}

// identEq checks if two Ident AST components are equivalent.
func identEq(base *build.Ident, last *build.Ident) bool {
	if base.Name != last.Name {
		return false
	}
	return true
}

// ifClauseEq checks if two IfClause AST components are equivalent.
func ifClauseEq(base *build.IfClause, last *build.IfClause) bool {
	if !exprEq(base.Cond, last.Cond) {
		return false
	}
	return true
}

// ifStmtEq checks if two IfStmt AST components are equivalent.
func ifStmtEq(base *build.IfStmt, last *build.IfStmt) bool {
	if !exprEq(base.Cond, last.Cond) {
		return false
	}
	if !exprArrayEq(base.True, last.True) {
		return false
	}
	if !exprArrayEq(base.False, last.False) {
		return false
	}
	return true
}

// indexExprEq checks if two IndexExpr AST components are equivalent.
func indexExprEq(base *build.IndexExpr, last *build.IndexExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if !exprEq(base.Y, last.Y) {
		return false
	}
	return true
}

// keyValExprEq checks if two KeyValueExpr AST components are equivalent.
func keyValExprEq(base *build.KeyValueExpr, last *build.KeyValueExpr) bool {
	if !exprEq(base.Key, last.Key) {
		return false
	}
	if !exprEq(base.Value, last.Value) {
		return false
	}
	return true
}

// lambdaExprEq checks if two LambdaExpr AST components are equivalent.
func lambdaExprEq(base *build.LambdaExpr, last *build.LambdaExpr) bool {
	if !fnEq(&base.Function, &last.Function) {
		return false
	}
	return true
}

// listExprEq checks if two ListExpr AST components are equivalent.
func listExprEq(base *build.ListExpr, last *build.ListExpr) bool {
	if !exprArrayEq(base.List, last.List) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// literalExprEq checks if two LiteralExpr AST components are equivalent.
func literalExprEq(base *build.LiteralExpr, last *build.LiteralExpr) bool {
	if base.Token != last.Token {
		return false
	}
	return true
}

// loadStmtEq checks if two LoadStmt AST components are equivalent.
func loadStmtEq(base *build.LoadStmt, last *build.LoadStmt) bool {
	if !stringExprEq(base.Module, last.Module) {
		return false
	}
	if !identArrayEq(base.From, last.From) {
		return false
	}
	if !identArrayEq(base.To, last.To) {
		return false
	}
	if !endEq(&base.Rparen, &last.Rparen) {
		return false
	}
	return true
}

// parenExprEq checks if two ParenExpr AST components are equivalent.
func parenExprEq(base *build.ParenExpr, last *build.ParenExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// returnStmtEq checks if two ReturnStmt AST components are equivalent.
func returnStmtEq(base *build.ReturnStmt, last *build.ReturnStmt) bool {
	if !exprEq(base.Result, last.Result) {
		return false
	}
	return true
}

// ruleEq checks if two Rule AST components are equivalent.
func ruleEq(base *build.Rule, last *build.Rule) bool {
	if !callExprEq(base.Call, last.Call) {
		return false
	}
	if base.ImplicitName != last.ImplicitName {
		return false
	}
	return true
}

// setExprEq checks if two SetExpr AST components are equivalent.
func setExprEq(base *build.SetExpr, last *build.SetExpr) bool {
	if !exprArrayEq(base.List, last.List) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// sliceExprEq checks if two sliceExpr AST components are equivalent.
func sliceExprEq(base *build.SliceExpr, last *build.SliceExpr) bool {
	if !exprEq(base.X, last.X) {
		return false
	}
	if !exprEq(base.From, last.From) {
		return false
	}
	if !exprEq(base.To, last.To) {
		return false
	}
	if !exprEq(base.Step, last.Step) {
		return false
	}
	return true
}

// stringExprEq checks if two StringExpr AST components are equivalent.
func stringExprEq(base *build.StringExpr, last *build.StringExpr) bool {
	if base.Value != last.Value {
		return false
	}
	if base.Token != last.Token {
		return false
	}
	return true
}

// tupleExprEq checks if two TupleExpr AST components are equivalent.
func tupleExprEq(base *build.TupleExpr, last *build.TupleExpr) bool {
	if !exprArrayEq(base.List, last.List) {
		return false
	}
	if !endEq(&base.End, &last.End) {
		return false
	}
	return true
}

// unaryExprEq checks if two UnaryExpr AST components are equivalent.
func unaryExprEq(base *build.UnaryExpr, last *build.UnaryExpr) bool {
	if base.Op != last.Op {
		return false
	}
	if !exprEq(base.X, last.X) {
		return false
	}
	return true
}

// exprArrayEq checks if two Expr arrays are equivalent.
func exprArrayEq(base []build.Expr, last []build.Expr) bool {
	if len(base) != len(last) {
		return false
	}
	for i := 0; i < len(base); i++ {
		if !exprEq(base[i], last[i]) {
			return false
		}
	}
	return true
}

// identArrayEq checks if two Ident arrays are equivalent.
func identArrayEq(base []*build.Ident, last []*build.Ident) bool {
	if len(base) != len(last) {
		return false
	}
	for i := 0; i < len(base); i++ {
		if !identEq(base[i], last[i]) {
			return false
		}
	}
	return true
}
