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

// SymbolicatedChecker takes symbolication information when comparing
// MAST forests.
type SymbolicatedChecker struct {
	GenericChecker
	baseSymbols *symbolication.SymbolTable
	lastSymbols *symbolication.SymbolTable
}

// CheckNode compares "generic" nodes.
func (s *SymbolicatedChecker) CheckNode(base mast.Node, last mast.Node, c LangChecker) bool {
	switch n1 := base.(type) {
	case *mast.Identifier:
		if n2, ok := last.(*mast.Identifier); ok {
			return s.IdentifierEq(n1, n2, c)
		}
	}
	eq := s.GenericChecker.CheckNode(base, last, c)
	if !eq {
		if ident, ok := last.(*mast.Identifier); ok {
			lastConst, _ := getConstDecl(ident, s.lastSymbols)
			if lastConst != nil {
				// check if an expression has not been replaced by a
				// constant with value the same as the expression
				return c.CheckNode(base, lastConst.Value, c)
			}
		}
	}
	return eq
}

// CheckDeclarationList is a helper function that check if every
// declaration in the given two lists is equal. All elements in the
// list must be non-nil.
func (s *SymbolicatedChecker) CheckDeclarationList(base []mast.Declaration, last []mast.Declaration, c LangChecker) bool {

	i := 0
	j := 0
	for i < len(base) && j < len(last) {
		if c.CheckNode(base[i], last[j], c) {
			// top base are the same
			i++
			j++
			continue
		}
		// In most cases we can ignore constant definitions
		// themselves, as even if they are added or removed, code
		// equivalence is decided at the point of the constant being
		// used (e.g. if constant is added and used instead of a
		// literal with the same value, it's OK). On the other hand,
		// there are cases where the constant value is modified
		// (without name change) and the uses are not in the diff - we
		// should prevent auto-approval in this case as this value
		// change may affect a different part of the service.
		if !s.sameNameOtherwiseModified(base[i], last[j], c) {
			return false
		}
		if ignoreDecl(base[i], s.baseSymbols) {
			// declaration in base diff is for a global const that
			// might have been removed in the last diff - ignore it
			// and keep comparing
			i++
			continue
		}
		if ignoreDecl(last[j], s.lastSymbols) {
			// declaration in last diff is for a global const that
			// might have been added in the last diff - ignore it and
			// keep comparing
			j++
			continue
		}
		return false
	}

	if i < len(base) {
		// base file has more declarations to analyze
		return areRemainingDeclIgnored(i, base, s.baseSymbols)
	}
	if j < len(last) {
		// last file has more declarations to analyze
		return areRemainingDeclIgnored(j, last, s.lastSymbols)
	}
	// otherwise all declarations have been analyzed already
	return true
}

// CheckStatementList is a helper function that check if every
// statement in the given two lists is equal. All elements in the list
// must be non-nil.
func (s *SymbolicatedChecker) CheckStatementList(base []mast.Statement, last []mast.Statement, c LangChecker) bool {
	i := 0
	j := 0
	// The statement comparison algorithm below works as follows:
	// - check if statements are equivalent
	// - for statements that are not equivalent according to "standard" algorithm
	//   - check if any of the statements can be ignored due to being
	//   irrelevant from the point of view of comparing two diffs
	//   (e.g, constant declarations)
	// - if at the end of the statements iteration, any statements
	// remain that have not been analyzed, check if the remaining
	// statements can be ignored
	for i < len(base) && j < len(last) {
		if c.CheckNode(base[i], last[j], c) {
			// statements are the same
			i++
			j++
			continue
		}

		// skip over constant declaration statements that can be ignored
		if declStmt, ok := base[i].(*mast.DeclarationStatement); ok && castConstDecl(declStmt.Decl) != nil && ignoreDecl(declStmt.Decl, s.baseSymbols) {
			// declaration in base diff is for a local const that
			// might have been removed in the last diff - ignore it
			// and keep comparing
			i++
			continue
		}
		if declStmt, ok := last[j].(*mast.DeclarationStatement); ok && castConstDecl(declStmt.Decl) != nil && ignoreDecl(declStmt.Decl, s.lastSymbols) {
			// declaration in last diff is for a local const that
			// might have been added in the last diff - ignore it and
			// keep comparing
			j++
			continue
		}

		ignore, baseSubExpressions := c.StmtIgnored(base[i], s.baseSymbols)
		if ignore {
			i++
			continue
		}

		ignore, lastSubExpressions := c.StmtIgnored(last[j], s.lastSymbols)
		if ignore {
			j++
			continue
		}

		if baseSubExpressions != nil && lastSubExpressions != nil {
			if !c.CheckExpressionList(baseSubExpressions, lastSubExpressions, c) {
				return false
			}
			// expressions are equivalent
			i++
			j++
			continue
		}

		// statements cannot be ignored and are also not equivalent
		return false
	}
	if i < len(base) {
		// base file has more declarations to analyze
		return areRemainingStmtIgnored(i, base, s.baseSymbols, c)
	}
	if j < len(last) {
		// last file has more declarations to analyze
		return areRemainingStmtIgnored(j, last, s.lastSymbols, c)
	}
	// otherwise all statements have been analyzed already
	return true
}

// IdentifierEq compares Identifier-s.
func (s *SymbolicatedChecker) IdentifierEq(n1 *mast.Identifier, n2 *mast.Identifier, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}
	baseConst, baseEntry := getConstDecl(n1, s.baseSymbols)
	lastConst, lastEntry := getConstDecl(n2, s.lastSymbols)
	if baseConst == nil || lastConst == nil {
		// just a regular check
		return s.GenericChecker.IdentifierEq(n1, n2, c)
	}

	if (baseEntry.IsPrivate && !lastEntry.IsPrivate) ||
		(!baseEntry.IsPrivate && lastEntry.IsPrivate) {
		// both or neither should be private (this should already have
		// been checked and flagged when comparing const declaratons
		// but just in case)
		return false
	}
	if baseEntry.IsPrivate && lastEntry.IsPrivate &&
		baseConst.Value != nil && lastConst.Value != nil {
		// ignore names (as they could have been renamed) and instead
		// compare the actual constant values (do it only if values
		// are non-nil - if they are nil then actual value cannot be
		// determined as it depends on the context of the constant's
		// declaration).
		return c.CheckNode(baseConst.Value, lastConst.Value, c)
	}

	// just a regular check
	return s.GenericChecker.IdentifierEq(n1, n2, c)
}

// getConstDecl returns const declaration and its respective symbol
// table entry if the identifer refers to this declaration, otherwise
// nils.
func getConstDecl(ident *mast.Identifier, symbolTable *symbolication.SymbolTable) (*mast.VariableDeclaration, *symbolication.SymbolTableEntry) {
	declEntry, err := symbolTable.DeclarationEntry(ident)
	if err != nil || declEntry == nil {
		return nil, nil
	}
	if vdecl, ok := declEntry.Node.(*mast.VariableDeclaration); ok && vdecl.IsConst {
		return vdecl, declEntry
	}
	return nil, nil
}

// sameNameOtherwiseModified checks if declarations represent
// constants with the same name and, (only) if so, checks if the
// constant values (and other constant declaration components are the
// same). The function returns false if declarations are for constants
// with the same name but differing in some other aspect, otherwise it
// returns true (e.g., if declarations are for differerent language
// components, or if other constant declaration components are the
// same).
func (s *SymbolicatedChecker) sameNameOtherwiseModified(base mast.Declaration, last mast.Declaration, c LangChecker) bool {
	baseDecl := castConstDecl(base)
	if baseDecl == nil {
		return true
	}
	lastDecl := castConstDecl(last)
	if lastDecl == nil {
		return true
	}
	if len(baseDecl.Names) != len(lastDecl.Names) {
		return true
	}

	for i, n := range baseDecl.Names {
		if n.Name != lastDecl.Names[i].Name {
			return true
		}
	}

	return c.CheckNode(baseDecl.Type, lastDecl.Type, c) &&
		c.CheckNode(baseDecl.Value, lastDecl.Value, c) &&
		c.CheckNode(baseDecl.LangFields, lastDecl.LangFields, c)
}

// areRemainingDeclIgnored checks if remaining declarations in the array
// can be ignored (excluded from comparison).
func areRemainingDeclIgnored(startIdx int, remaining []mast.Declaration, symbolTable *symbolication.SymbolTable) bool {
	for i := startIdx; i < len(remaining); i++ {
		if !ignoreDecl(remaining[i], symbolTable) {
			return false
		}
	}
	return true
}

// areRemainingStmtIgnored checks if remaining statement in the array
// can be ignored (excluded from comparison).
func areRemainingStmtIgnored(startIdx int, remaining []mast.Statement, symbolTable *symbolication.SymbolTable, c LangChecker) bool {
	for i := startIdx; i < len(remaining); i++ {
		switch s := remaining[i].(type) {
		case *mast.DeclarationStatement:
			if ignoreDecl(s.Decl, symbolTable) {
				continue
			}
		default:
			if ignore, _ := c.StmtIgnored(s, symbolTable); ignore {
				continue
			}
		}
		return false
	}
	return true
}

// ignoreDecl determines if a given declaration can be excluded from
// comparison (ignored).
func ignoreDecl(decl mast.Declaration, symbolTable *symbolication.SymbolTable) bool {
	// it must be a constant to be ignored
	cdecl := castConstDecl(decl)
	if cdecl == nil {
		return false
	}

	// it can be ignored only if it has a value (otherwise value
	// depends on context and cannot be easily determined)
	if cdecl.Value == nil {
		return false
	}

	// it can be ignored only if all declared names are renameable
	// (otherwise cannot be modified and should be checked for
	// "standard" equivalence)
	for _, ident := range cdecl.Names {
		declEntry, err := symbolTable.DeclarationEntry(ident)
		if err != nil || declEntry == nil || !declEntry.IsPrivate {
			// cannot be ignored and will be compared structurally (as
			// if symbolication information was not available)
			return false
		}
	}
	return true
}

// castConstDecl returns const declaration or nil if the declaration is
// not for a const.
func castConstDecl(decl mast.Declaration) *mast.VariableDeclaration {
	if vdecl, ok := decl.(*mast.VariableDeclaration); ok && vdecl.IsConst {
		return vdecl
	}
	return nil
}
