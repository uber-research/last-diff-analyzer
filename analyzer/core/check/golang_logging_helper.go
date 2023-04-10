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

package check

import (
	"analyzer/core/mast"
)

// _supportedZapMethods contains the names of logging methods that can
// be subject to auto-approval. These calls are also interchangeable
// from the point of view of auto-approval (i.e., a Debug call can be
// replaced with a Warn call).
var _supportedZapMethods = map[string]bool{
	"Debug": true,
	"Warn":  true,
	"Info":  true,
	"Error": true,
}

// _helperZapFunctions contains the names of logging helper functions
// whose calls are "safe" (free of side-effects) to use as arguments
// to the logging methods.
var _helperZapFunctions = map[string]bool{
	"Any":         true,
	"Object":      true,
	"Array":       true,
	"Bool":        true,
	"Boolp":       true,
	"Bools":       true,
	"Complex128":  true,
	"Complex128p": true,
	"Complex128s": true,
	"Complex64":   true,
	"Complex64p":  true,
	"Complex64s":  true,
	"Float64":     true,
	"Float64p":    true,
	"Float64s":    true,
	"Float32":     true,
	"Float32p":    true,
	"Float32s":    true,
	"Int":         true,
	"Intp":        true,
	"Ints":        true,
	"Int64":       true,
	"Int64p":      true,
	"Int64s":      true,
	"Int32":       true,
	"Int32p":      true,
	"Int32s":      true,
	"Int16":       true,
	"Int16p":      true,
	"Int16s":      true,
	"Int8":        true,
	"Int8p":       true,
	"Int8s":       true,
	"String":      true,
	"Stringp":     true,
	"Strings":     true,
	"Uint":        true,
	"Uintp":       true,
	"Uints":       true,
	"Uint64":      true,
	"Uint64p":     true,
	"Uint64s":     true,
	"Uint32":      true,
	"Uint32p":     true,
	"Uint32s":     true,
	"Uint16":      true,
	"Uint16p":     true,
	"Uint16s":     true,
	"Uint8":       true,
	"Uint8p":      true,
	"Uint8s":      true,
	"Binary":      true,
	"Uintptr":     true,
	"Uintptrp":    true,
	"Uintptrs":    true,
	"Time":        true,
	"Timep":       true,
	"Times":       true,
	"Duration":    true,
	"Durationp":   true,
	"Durations":   true,
	"NamedError":  true,
	"Errors":      true,
	"Stringer":    true,
	"Reflect":     true,
}

// _internalSafeFunctions contains names of internal "safe" (free of
// side-effects) functions.
var _internalSafeFunctions = map[string]bool{
	"len":    true,
	"append": true,
}

const (
	// _zapPath is the zap package path
	_zapPath = "\"go.uber.org/zap\""
	// _zapPath is the zap package name
	_zapName = "zap"
)

// isLogImportClean determines if logging imports in a given root of a
// forest were "clean", that it if they imported the right logging
// library and were not aliased (which could make auto-approval of
// related logging changes incorrect).
func (g *GoChecker) isLogImportClean(root *mast.Root) bool {
	for _, decl := range root.Declarations {
		imp, ok := decl.(*mast.ImportDeclaration)
		if !ok {
			// not an import declaration
			continue
		}
		pkg, ok := imp.Package.(*mast.StringLiteral)
		if !ok {
			// should not happen, but just in case
			return false
		}
		if pkg.Value != _zapPath {
			// does not match zap package path
			continue
		}
		if imp.Alias == nil || imp.Alias.Name == _zapName {
			// correct zap package path and (optional) alias
			continue
		}
		// correct zap package path but import is incorrectly aliased
		return false
	}
	return true
}

// callIgnored verifies if the call can be ignored when comparing
// different files. It also returns a (possibly empty) list of
// expressions representing arguments and access path expression, both
// (or either) of which could prevent this call from being ignored due
// to side-effects. The verification is applied to either outer calls
// (actual logging calls) or to calls used in arguments to (e.g.,
// helper functions to form logging call arguments).
func (g *GoChecker) callIgnored(expr mast.Expression, outer bool) (bool, []mast.Expression) {
	callExpr, ok := expr.(*mast.CallExpression)
	if !ok || callExpr == nil {
		// not a call
		return false, nil
	}

	if outer {
		// check for actual logging calls
		if !g.isZapCallSignature(callExpr) {
			return false, nil
		}
	} else {
		// check for helper calls used in arguments to logging calls
		if !isZapHelperSignature(callExpr) && !isOtherSafeSignature(callExpr) {
			return false, nil
		}
	}

	ignore := true
	safe, unsafeExprs := isSafeAccessPath(callExpr.Function, nil)
	if !safe {
		ignore = false
	}

	for _, a := range callExpr.Arguments {
		if g.isSafeArg(a) {
			continue
		}
		unsafeExprs = append(unsafeExprs, a)
		ignore = false
	}

	return ignore, unsafeExprs
}

// isZapCallSignature verifies if a call is to a function with a
// supported signature.
func (g *GoChecker) isZapCallSignature(callExpr *mast.CallExpression) bool {
	// We use a simplified algorithm to check the signature:
	// - function has to have supported name in the Zap package
	// - the first argument must be a string literal
	// - if there are more than one arguments, at least one must be a
	// call to a supported Zap package helper function
	name := getCalleeName(callExpr.Function)
	if !_supportedZapMethods[name] {
		// wrong name
		return false
	}

	args := callExpr.Arguments

	if len(args) == 0 {
		// no arguments
		return false
	}

	if _, ok := args[0].(*mast.StringLiteral); !ok {
		// first argument is not a string literal
		return false
	}

	if len(args) == 1 {
		// only one argument - checking done
		return true
	}

	if !g.LogImportClean {
		// no need to even check for calls to supported Zap package
		// helper functions as the Zap package import itself is
		// incorrect
		return false
	}

	for i := 1; i < len(args); i++ {
		if callExpr, ok := args[i].(*mast.CallExpression); ok && isZapHelperSignature(callExpr) {
			// argument is a a call to one of the Zap helper functions
			// (we need only one for verification)
			return true
		}
	}

	return false
}

// getCalleeName returns the name of the function/method being called
// (empty string if name cannot be computed).
func getCalleeName(expr mast.Expression) string {
	switch n := expr.(type) {
	case *mast.Identifier:
		return n.Name
	case *mast.AccessPath:
		return n.Field.Name
	}
	return ""
}

// isOtherSafeSignature checks if the call represents a known "safe"
// (no side-effects) function, other than those defined in the Zap
// package.
func isOtherSafeSignature(callExpr *mast.CallExpression) bool {
	ident, ok := callExpr.Function.(*mast.Identifier)
	if !ok {
		// at this point "other" safe signatures only include internal
		// functions whose calls do not involve selectors (e.g. to
		// specify a package name).
		return false
	}
	return _internalSafeFunctions[ident.Name]
}

// isZapHelperSignature returns true if the call expression represents
// a call to one of Zap's helper functions which are "safe" arguments
// to Zap logging calls. Otherwise it returns false.
func isZapHelperSignature(callExpr *mast.CallExpression) bool {
	ap, ok := callExpr.Function.(*mast.AccessPath)
	if !ok {
		// call is not represented by access path (as it should be for
		// a zap helper function, whose access path should consist of
		// the "zap" identifier and a function name)
		return false
	}
	if pkg, ok := ap.Operand.(*mast.Identifier); !ok || pkg.Name != _zapName {
		// function name is preceded but something else than the "zap"
		// identifier
		return false
	}
	return _helperZapFunctions[ap.Field.Name]
}

// isSafeAccessPath checks if the access path of the logging call, for
// example foo.bar in foo.bar.Error("error"), is safe (free of
// side-effects). Returns the list of expressions on the access path.
func isSafeAccessPath(expr mast.Expression, exprList []mast.Expression) (bool, []mast.Expression) {
	// at this point, the access path is valid if it contains
	// identifiers only (to avoid potential side-effects on the access
	// path)
	switch n := expr.(type) {
	case *mast.Identifier:
		exprList = append(exprList, n)
		return true, exprList
	case *mast.AccessPath:
		exprList = append(exprList, n.Field)
		return isSafeAccessPath(n.Operand, exprList)
	}
	exprList = append(exprList, expr)
	return false, exprList
}

// isSafeArg verifies if the expression represents a "safe" (no
// side-effects) argument to a method or function used in a logging
// call.
func (g *GoChecker) isSafeArg(arg mast.Expression) bool {
	switch e := arg.(type) {
	case *mast.NullLiteral, *mast.BooleanLiteral, *mast.IntLiteral, *mast.FloatLiteral, *mast.StringLiteral, *mast.CharacterLiteral:
		return true
	case *mast.Identifier:
		// a simple identifier
		return true
	case *mast.UnaryExpression:
		if e.Operator == "&" {
			// an identifier whose address is being taken (no
			// side-effects even if taking address of the nil value)
			return true
		}
	case *mast.CallExpression:
		ignore, _ := g.callIgnored(e, false)
		return ignore
	}
	return false
}
