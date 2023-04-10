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
	"fmt"

	"analyzer/core/mast"
	"analyzer/core/symbolication"
	ts "analyzer/core/treesitter"
)

// CommonChecker is the interface that all checkers (language-specific
// and language-agnostic) must implement.
type CommonChecker interface {
	// The following methods are responsible for comparing individual
	// nodes or lists of nodes. The reason they take the interface as a
	// parameter is that we want to be able to "override" them with
	// respect to the generic implementation (otherwise, when calling
	// recursively using the receiver, the "override" would not work).

	// CheckNode compares individual nodes.
	CheckNode(node mast.Node, other mast.Node, c LangChecker) bool
	// CheckDeclarationList is a helper function that check if every
	// declaration in the given two lists is equal.
	CheckDeclarationList(declarations []mast.Declaration, other []mast.Declaration, c LangChecker) bool
	// CheckStatementList is a helper function that check if every
	// statement in the given two lists is equal.
	CheckStatementList(statements []mast.Statement, other []mast.Statement, c LangChecker) bool
	// CheckExpressionList is a helper function that check if every
	// expression in the given two lists is equal.
	CheckExpressionList(expressions []mast.Expression, other []mast.Expression, c LangChecker) bool
	// CheckDimensions is a helper function that check if two arrays of
	// dimensions are equal. All elements in the list must be non-nil.
	CheckDimensions(dimensions []*mast.JavaDimension, other []*mast.JavaDimension, c LangChecker) bool
}

// LangChecker is the interface that all language-specific checkers
// must implement.
type LangChecker interface {
	CommonChecker
	// StmtIgnored checks if a given statement can be ignored when
	// checking for equivalence. If it can be ignored altogether, this
	// method returns true. If it can be ignored under the condition
	// that some of its sub-expressions are equal, this method returns
	// false and a non-nil list of sub-expressions. If it cannot be
	// ignored under any circumstances, this method returns false and
	// an empty (or nil) list of sub-expressions.
	//
	// An example here could be logging calls: ignored completely if
	// logging side-effect free arguments (e.g. ints), ignored
	// conditionally if some of their arguments have side-effects but
	// are equal in both diffs, and not ignored if they have a
	// different number of arguments.
	StmtIgnored(stmt mast.Statement, symbols *symbolication.SymbolTable) (bool, []mast.Expression)
}

// Checker is the interface that all "top-level" (language-agnostic)
// checkers must implement.
type Checker interface {
	CommonChecker
	// Equal checks whether the given forest and other forest are
	// functionally equivalent. It is the main driver method of the
	// checker.
	Equal(forest []mast.Node, other []mast.Node, c LangChecker) (bool, error)
	// Lang returns the specific language-specifc checker.
}

// Run is the main driver of the checker. It traverses the given MAST forest and returns whether two MAST forests are
// functionally equivalent.
func Run(base []mast.Node, last []mast.Node, baseSymbols *symbolication.SymbolTable, lastSymbols *symbolication.SymbolTable, suffix string, loggingOn bool) (bool, error) {
	var checker Checker
	var langChecker LangChecker
	var err error
	if baseSymbols != nil && lastSymbols != nil {
		symbolicatedChecker := SymbolicatedChecker{
			GenericChecker: GenericChecker{},
			baseSymbols:    baseSymbols,
			lastSymbols:    lastSymbols,
		}
		switch suffix {
		case ts.GoExt:
			langChecker = NewGoSymbolicatedChecker(&symbolicatedChecker, loggingOn)
		case ts.JavaExt:
			langChecker = NewJavaSymbolicatedChecker(&symbolicatedChecker, loggingOn)
		default:
			langChecker, err = newLangChecker(&symbolicatedChecker, suffix, loggingOn)
		}
		checker = &symbolicatedChecker
	} else {
		checker = &GenericChecker{}
		langChecker, err = newLangChecker(checker, suffix, loggingOn)
	}
	if err != nil {
		return false, err
	}
	isEqual, err := checker.Equal(base, last, langChecker)
	if err != nil {
		return false, err
	}

	return isEqual, nil
}

// newLangChecker creates and returns a new language-specific checker.
func newLangChecker(checker Checker, suffix string, loggingOn bool) (LangChecker, error) {

	var langChecker LangChecker
	switch suffix {
	case ts.JavaExt:
		langChecker = NewJavaChecker(checker, loggingOn)
	case ts.GoExt:
		langChecker = NewGoChecker(checker, loggingOn)
	default:
		return nil, fmt.Errorf("unsupported file extension %q during equivalence checking", suffix)
	}
	return langChecker, nil
}
