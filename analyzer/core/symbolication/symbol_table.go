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

package symbolication

import (
	"errors"
	"fmt"
	"strings"

	"analyzer/core/mast"
)

// SymbolTableEntry is a wrapper that wraps different kinds of information that a symbol is
// associated with.
type SymbolTableEntry struct {
	// Identifier is the identifier in the declaration node for the symbol.
	Identifier *mast.Identifier
	// Node represents declaration itself.
	Node mast.Node
	// IsPrivate marks whether this declaration is private or not. For example, whether "private"
	// is present for the declaration in Java, or whether the first letter of the identifier is
	// upper-case in Go.
	IsPrivate bool
	// IsActive determines if the declaration is active, this is specially designed to allow
	// temporarily "deactivating" a declaration entry so that it can be ignored during the search
	// to link "use" to "def". For example, when we have a declaration with LHS and RHS having the
	// same identifiers:
	// tmp := &tmp{}
	// it is a perfectly valid declaration where the RHS refers to some struct declared earlier.
	// However, due to the nature of the traversal we will be first visiting the variable
	// declaration in symbolication, adding a declaration entry onto the current stack, and _later_
	// visiting the identifiers in the RHS. At which point we will be erroneously linking the "tmp"
	// in RHS to this variable declaration. To solve this, we introduce this field and add an
	// inactive entry when visiting the variable declaration, the search for entry for "tmp" in RHS
	// will then search for active entries only. The entry is finally "activated" when the
	// traversal for the entire variable declaration finishes.
	IsActive bool
}

// SymbolTable contains mappings from identifier to information associated with it, such as links to its declaration.
type SymbolTable struct {
	// table is a map that maps from symbol to multiple information associated with it (wrapped in symbolInfo struct).
	table map[*mast.Identifier]*SymbolTableEntry
	// symbols keeps a slice of all symbols in the table, it preserves the insertion order so that the iteration on the
	// map will be deterministic.
	symbols []*mast.Identifier
}

func (s *SymbolTable) String() string {
	var b strings.Builder
	for id, entry := range s.table {
		if entry != nil {
			fmt.Fprintf(&b, "%s: %p => %p\n", id.Name, id, entry.Identifier)
		} else {
			fmt.Fprintf(&b, "%s: %p => nil\n", id.Name, id)
		}
	}
	return b.String()
}

// newSymbolTable creates and returns a properly initialized *SymbolTable.
func newSymbolTable() *SymbolTable {
	return &SymbolTable{
		table: make(map[*mast.Identifier]*SymbolTableEntry),
	}
}

// DeclarationEntry returns the declaration entry linked for the identifier node, an error is returned if the input
// identifier is nil, or the identifier does not exist in the symbol table.
func (s *SymbolTable) DeclarationEntry(identifier *mast.Identifier) (*SymbolTableEntry, error) {
	if identifier == nil {
		return nil, errors.New("input identifier node is nil")
	}
	result, ok := s.table[identifier]
	if !ok {
		return nil, fmt.Errorf("%q does not exist in the symbol table", identifier.Name)
	}
	return result, nil
}

// OrderedSymbols returns all the symbols stored in the symbol table.
func (s *SymbolTable) OrderedSymbols() []*mast.Identifier {
	symbols := make([]*mast.Identifier, len(s.symbols))
	copy(symbols, s.symbols)
	return symbols
}

// AddLink adds a link from the identifier to the declaration entry.
func (s *SymbolTable) AddLink(symbol *mast.Identifier, entry *SymbolTableEntry) error {
	if symbol == nil {
		return errors.New("adding a nil link")
	}

	// During the translation, we could be re-using the translated type node for a multi-variable declaration. For
	// example, "int a, b;" would be translated to two mast.VariableDeclaration nodes, with a shared identifier node
	// "int" as their types. This means that we could be adding the same link twice for "int" node. So we return an
	// error only if the entry to be stored is different from the existing entry.
	existingEntry, ok := s.table[symbol]
	if ok && existingEntry != entry {
		return fmt.Errorf("AddLink: link exists for identifier %q", symbol.Name)
	}

	// do nothing if the entry already in the table
	if ok {
		return nil
	}

	// add the link to the symbol info
	s.table[symbol] = entry
	s.symbols = append(s.symbols, symbol)

	return nil
}

// RemoveLink removes a link for a given identifier. It's relatively
// expensive so should be used sparingly.
func (s *SymbolTable) RemoveLink(symbol *mast.Identifier) error {
	if symbol == nil {
		return errors.New("removing a nil link")
	}
	delete(s.table, symbol)
	newSymbols := []*mast.Identifier{}
	for _, sym := range s.symbols {
		if sym != symbol {
			newSymbols = append(newSymbols, sym)
		}
	}
	s.symbols = newSymbols
	return nil
}

// ReplaceLink replaces a link for a given identifier with the new
// entry (the link must already exist).
func (s *SymbolTable) ReplaceLink(symbol *mast.Identifier, entry *SymbolTableEntry) error {
	if symbol == nil {
		return errors.New("replacing a nil link")
	}
	if _, exists := s.table[symbol]; !exists {
		return errors.New("replacing a non-existing link")
	}
	s.table[symbol] = entry
	return nil
}
