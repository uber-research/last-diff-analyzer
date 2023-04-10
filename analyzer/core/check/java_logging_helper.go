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
	"analyzer/core/mast/mastutil"
	"analyzer/core/symbolication"
)

// _supportedLoggingClasses marks the supported logging import classes.
var _supportedLoggingClasses = map[string]bool{
	"org.slf4j.Logger":         true,
	"org.apache.log4j.Logger":  true,
	"java.util.logging.Logger": true,
}

// _loggerMethods stores the set of methods that can be called on a logger variable.
var _loggerMethods = map[string]bool{
	"trace":   true,
	"debug":   true,
	"info":    true,
	"warn":    true,
	"warning": true,
	"severe":  true,
	"error":   true,
	"fatal":   true,
}

// _safeLoggerHelpers marks the helper functions that we consider to be side-effect-free. They
// may be safely ignored when used in logger functions.
var _safeLoggerHelpers = map[string]bool{
	"String.format": true,
}

// findImportedLogger examines the import declarations and return imported the logger class name
// from supported libraries (e.g., slf4j).
func (j *JavaChecker) findImportedLogger(root *mast.Root) string {
	for _, decl := range root.Declarations {
		i, ok := decl.(*mast.ImportDeclaration)
		if !ok {
			continue
		}
		path, ok := i.Package.(*mast.AccessPath)
		if !ok {
			continue
		}

		// Check if the imported class is one of the supported logging class.
		p, err := mastutil.JoinAccessPath(path)
		if err != nil {
			return ""
		}
		if _supportedLoggingClasses[p] {
			return path.Field.Name
		}
	}

	return ""
}

// loggingStmtIgnore checks if a statement is a logging-related statement and can be ignored.
func (j *JavaChecker) loggingStmtIgnore(stmt mast.Statement, symbols *symbolication.SymbolTable) (bool, []mast.Expression) {
	// Must have an imported logger class.
	if len(j.importedLoggerClass) == 0 {
		return false, nil
	}

	// Check if the statement is a method call.
	exprStmt, ok := stmt.(*mast.ExpressionStatement)
	if !ok || exprStmt.Expr == nil {
		return false, nil
	}
	call, ok := exprStmt.Expr.(*mast.CallExpression)
	if !ok {
		return false, nil
	}
	path, ok := call.Function.(*mast.AccessPath)
	if !ok {
		return false, nil
	}
	// Logger calls will not have type arguments (e.g., "logger.<A, B...>debug()") associated with
	// it, so we early return if there are.
	if call.LangFields != nil {
		return false, nil
	}

	// The call should end with a logger method (e.g., "XX.XX.XX.debug(...)").
	if !_loggerMethods[path.Field.Name] {
		return false, nil
	}

	// We currently only support directly calling logger.debug instead of loggerFunc().debug. So
	// the operand must be a simple identifier.
	first, ok := path.Operand.(*mast.Identifier)
	if !ok {
		return false, nil
	}

	// The declaration entry for the variable must be a private and final variable declaration.
	entry, err := symbols.DeclarationEntry(first)
	if err != nil || entry == nil {
		return false, nil
	}
	decl, ok := entry.Node.(*mast.VariableDeclaration)
	if !ok {
		return false, nil
	}
	modifiers := decl.LangFields.(*mast.JavaVariableDeclarationFields).Modifiers
	if !entry.IsPrivate || !mastutil.HasJavaModifier(modifiers, mast.FinalMod) {
		return false, nil
	}

	// Check if the type of the variable declaration is the imported logger class.
	if ident, ok := decl.Type.(*mast.Identifier); !ok || ident.Name != j.importedLoggerClass {
		return false, nil
	}

	// Recursively collect all expressions that could have side effects. The equality will depend
	// on the equalities of them.
	var exprs []mast.Expression
	for _, arg := range call.Arguments {
		exprs = append(exprs, extractSideEffectExprs(arg)...)
	}

	// If the arguments of the logger call is side-effect-free, we can simply ignore this statement.
	if len(exprs) == 0 {
		return true, nil
	}

	return false, exprs
}

// extractSideEffectExprs recursively extracts and returns a slice of MAST expressions that
// potentially have side effects in expr.
func extractSideEffectExprs(expr mast.Expression) []mast.Expression {
	switch e := expr.(type) {
	// Literal nodes and simple identifiers have no side effects.
	case *mast.NullLiteral, *mast.BooleanLiteral, *mast.IntLiteral, *mast.FloatLiteral,
		*mast.StringLiteral, *mast.CharacterLiteral, *mast.Identifier:
		return nil

	case *mast.BinaryExpression:
		var exprs []mast.Expression
		exprs = append(exprs, extractSideEffectExprs(e.Left)...)
		exprs = append(exprs, extractSideEffectExprs(e.Right)...)
		return exprs

	case *mast.UnaryExpression:
		return extractSideEffectExprs(e.Expr)

	case *mast.CallExpression:
		// Convert the function to string representation to be checked against our allow-list.
		name := ""
		switch a := e.Function.(type) {
		case *mast.Identifier:
			name = a.Name
		case *mast.AccessPath:
			if p, err := mastutil.JoinAccessPath(a); err == nil {
				name = p
			}
		}

		// Check if it is a safe logger helper function.
		if !_safeLoggerHelpers[name] {
			return []mast.Expression{e}
		}

		// It is a safe logger helper function, but we still need to make sure the arguments do
		// not contain expressions that could have side effects.
		var exprs []mast.Expression
		for _, arg := range e.Arguments {
			exprs = append(exprs, extractSideEffectExprs(arg)...)
		}

		return exprs

	default:
		return []mast.Expression{expr}
	}
}
