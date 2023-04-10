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
	"analyzer/core/symbolication"
)

// GoSymbolicatedChecker represents a Go-specific symbolicated
// equivalence checker.
type GoSymbolicatedChecker struct {
	// GoChecker is a regular (non-symbolicated) Go checker.
	GoChecker

	// Symbolicated is the language-agnostic symbolicated checker.
	Symbolicated *SymbolicatedChecker
}

// NewGoSymbolicatedChecker creates and returns Go-specific
// symbolicated equivalence checker.
func NewGoSymbolicatedChecker(checker *SymbolicatedChecker, loggingOn bool) *GoSymbolicatedChecker {
	goChecker := NewGoChecker(checker, loggingOn)
	return &GoSymbolicatedChecker{GoChecker: *goChecker, Symbolicated: checker}
}

// CheckNode compares "generic" nodes.
func (g *GoSymbolicatedChecker) CheckNode(node mast.Node, other mast.Node, c LangChecker) bool {
	if baseRoot, ok := node.(*mast.Root); ok {
		if lastRoot, ok := other.(*mast.Root); ok {
			g.LogImportClean = g.isLogImportClean(baseRoot) && g.isLogImportClean(lastRoot)
		} else {
			// should not happen
			return false
		}
	}

	switch n1 := node.(type) {
	case *mast.EntityCreationExpression:
		if n2, ok := other.(*mast.EntityCreationExpression); ok {
			return g.EntityCreationExpressionEq(n1, n2, c)
		}
	}
	// fall back on a "regular" Go checker
	return g.GoChecker.CheckNode(node, other, c)
}

// EntityCreationExpressionEq compares EntityCreationExpression-s.
func (g *GoSymbolicatedChecker) EntityCreationExpressionEq(n1 *mast.EntityCreationExpression, n2 *mast.EntityCreationExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if !c.CheckNode(n1.Type, n2.Type, c) || !c.CheckNode(n1.LangFields, n2.LangFields, c) {
		return false
	}

	if n1.Value == nil && n2.Value == nil {
		return true
	}

	if n1.Value == nil || n2.Value == nil {
		return false
	}

	values1 := n1.Value.Values
	values2 := n2.Value.Values

	if len(values1) != len(values2) {
		return false
	}

	// lastFields contains a list of struct fields (if any) for the
	// constructor represented by the n2 parameter (coming from the
	// last diff), if this constructor is for a struct type. This
	// variable is initialized lazily to avoid unnecessary overheads.
	var lastFields []*mast.FieldDeclaration

	for i := 0; i < len(values1); i++ {
		if baseExpr, ok := values1[i].(*mast.KeyValuePair); ok {
			// expression in the in base diff is (example):
			// someStruct{fieldOne: 7, fieldTwo: 42}
			if lastExpr, ok := values2[i].(*mast.KeyValuePair); ok {
				// expression in last diff also is
				// (example): someStruct{fieldOne: 7, fieldTwo: 42}
				// (compare them directly)
				if !c.CheckNode(baseExpr, lastExpr, c) {
					return false
				}
			} else {
				// if expression is a key-value expressions in the
				// base diff but not in the last diff, the literals
				// are incomparable
				return false
			}
		} else if lastExpr, ok := values2[i].(*mast.KeyValuePair); ok {
			// expression in last diff is (example):
			// someStruct{fieldOne: 7, fieldTwo: 42}
			//
			// we may be converting from a simple list of expressions:
			// someStruct{7, 42}

			// field values must still match in the base diff and in the last diff
			if !c.CheckNode(values1[i], lastExpr.Value, c) {
				return false
			}
			// field names in the last diff must match the field order
			// in the base diff
			lastFields = cacheStructFields(lastFields, n2.Type, g.Symbolicated.lastSymbols)
			if lastFields == nil {
				// cannot get to a field declaration - bail out
				return false
			}
			if !isFieldNameMatch(lastExpr.Key, i, lastFields) {
				return false
			}
		} else if !c.CheckNode(values1[i], values2[i], c) {
			return false
		}
	}

	return true
}

// cacheStructFields returns the list of fields defined in a struct,
// or nil if type represents a different type or no fields are
// defined. The list of fields is computed only if the fields
// parameter is nil.
func cacheStructFields(fields []*mast.FieldDeclaration, typ mast.Expression, symbolTable *symbolication.SymbolTable) []*mast.FieldDeclaration {
	if fields != nil {
		return fields
	}
	ident, ok := typ.(*mast.Identifier)
	if !ok {
		// if we can't get to the type declaration, bail out
		return nil
	}
	declEntry, err := symbolTable.DeclarationEntry(ident)
	if err != nil || declEntry == nil {
		return nil
	}
	decl, ok := declEntry.Node.(*mast.GoTypeDeclaration)
	if !ok {
		return nil
	}
	if structType, ok := decl.Type.(*mast.GoStructType); ok {
		return structType.Declarations
	}
	return nil
}

// isFieldNameMatch checks if a given field is a match in a given list at
// a given position.
func isFieldNameMatch(fieldKey mast.Expression, fieldPos int, fields []*mast.FieldDeclaration) bool {
	fieldIdent, ok := fieldKey.(*mast.Identifier)
	if !ok {
		return false
	}
	if fieldPos >= len(fields) {
		return false
	}
	return fieldIdent.Name == fields[fieldPos].Name.Name
}
