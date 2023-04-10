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

// errNoScopesAvailable error is returned whenever there is no scope on the scope stack but access methods such as
// findDeclaration is called.
var errNoScopesAvailable error = errors.New("no scopes available")

// scopeEntry describes an entry in a given scope (multiple entries
// with the same name but different kinds can exist).
type scopeEntry map[mast.NameKind]*SymbolTableEntry

// scope represents a variable scope. Since the languages we support do not allow redeclarations of variables inside the
// _same_ scope [1], internally we use a map from identifier name to the identifier pointer to represent the data
// structure for faster lookup. Note that an error will be returned if a redeclaration is found inside the same scope.
// [1] Golang: https://golang.org/ref/spec#Declarations_and_scope
//
//	Java: https://docs.oracle.com/javase/specs/jls/se11/html/jls-6.html#jls-6.4
//
// TODO: consider the possibility of redeclarations when more languages are supported.
type scope struct {
	// table stores the symbol entries for this particular scope.
	table map[string]scopeEntry
	// isPrivate indicates whether this scope is private or not.
	isPrivate bool
}

// newScope creates a properly-initialized scope.
func newScope(isPrivate bool) *scope {
	return &scope{
		table:     make(map[string]scopeEntry),
		isPrivate: isPrivate,
	}
}

// scopeManager manages the stack of scopes, as well as package-level scopes. It provides convenient methods for easy
// accesses on the scopes.
type scopeManager struct {
	// scopeStack is a stack that keeps the stack of scopes with a slice of (ordered) declarations, where the current
	// scope is at the tail (i.e., at index len(scoped)-1).
	stack []*scope

	// packageScopes stores the mapping from the string representation of the package name (i.e., "pkg.a.b.c") to the
	// scope for that particular scope.  Note that the scope stack does _not_ store the package-level scope. The
	// management of scopes is two-level:
	// - scopeStack[N]
	// - ...
	// - scopeStack[1]
	// - scopeStack[0]          <-- this is the individual scope for function, class etc.
	// - currentPackageScope    <-- this is the current package scope
	// In other words, the `currentPackageScope` always stays at the bottom of the stack, and will be swapped when
	// analyzing a new file.
	packageScopes map[string]*scope
	// currentPackageScope keeps the scope for the current package, it will be created or retrieved from the
	// packageScopes map when visiting a mast.Root. Since the symbolicators current only accept mast.Root nodes in the
	// forest, this field should never be nil.
	currentPackageScope *scope

	// thisScopes keeps track of the scope of "this" keyword, so that we can properly resolve access paths like
	// "this.xx.xx".
	thisScopes []*scope
}

// newScopeManager creates a properly initialized scope manager.
func newScopeManager() *scopeManager {
	manager := &scopeManager{
		packageScopes: make(map[string]*scope),
	}
	return manager
}

// CreatePackageScope creates/retrieves a package-level scope and properly sets the current package scope based on the
// input mast.PackageDeclaration node. If the package declaration node is nil, a standalone package level scope will be
// created and isolated from other package-level scopes in the forest.
func (s *scopeManager) CreatePackageScope(pkg *mast.PackageDeclaration) error {
	if pkg == nil {
		// package scope is never private, instead, the privateness of the declarations inside it will be determined
		// individually.
		s.currentPackageScope = newScope(false /* isPrivate */)
		return nil
	}

	// convert the package declaration node to string representation (do not include annotation nodes of the package
	// declaration for Java since "@SomeAnnotation pkg" and "pkg" refers to the same package)

	// store the strings for every identifier along the access path in chain
	chain := []string{}
	// The Name field of a mast.PackageDeclaration could either be an mast.AccessPath (of mast.Identifier nodes) or a
	// plain mast.Identifier node. Therefore we iteratively unwrap the mast.AccessPath and collect the names of all
	// mast.Identifier nodes along the way, until we hit a mast.Identifier.
	var current mast.Expression = pkg.Name
	shouldContinue := true
	for shouldContinue {
		switch n := current.(type) {
		case *mast.AccessPath:
			current = n.Operand
			chain = append(chain, n.Field.Name)
		case *mast.Identifier:
			// we cannot go further so we break the loop
			chain = append(chain, n.Name)
			shouldContinue = false
		default:
			return fmt.Errorf("unexpected node type %T in package declaration", current)
		}
	}
	// The names are stored in reverse order in chain, so first reverse it.
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	// join the names to form the string representation of the package name
	name := strings.Join(chain, ".")
	if _, exists := s.packageScopes[name]; !exists {
		// package scope is never private
		packageScope := newScope(false /* isPrivate */)
		s.packageScopes[name] = packageScope
	}
	s.currentPackageScope = s.packageScopes[name]

	return nil
}

// IsCurrentScopePrivate indicates whether the current scope is private or not.
func (s *scopeManager) IsCurrentScopePrivate() bool {
	return s.currentScope().isPrivate
}

// CreateNewScope puts an empty scope onto the scope stack. isPrivate indicates whether this scope is private or not.
func (s *scopeManager) CreateNewScope(isPrivate bool) {
	// create a new scope
	s.stack = append(s.stack, newScope(isPrivate))
}

// PopScope pops a scope from the scope stack.
func (s *scopeManager) PopScope() error {
	if len(s.stack) == 0 {
		return fmt.Errorf("pop scope: %v", errNoScopesAvailable)
	}
	// pop the current scope
	s.stack = s.stack[:len(s.stack)-1]

	return nil
}

// MarkThis marks the current scope as the scope for "this" keyword.
func (s *scopeManager) MarkThis() error {
	// this keyword can never appear at the package-level declarations, so the scope stack must have at least one scope.
	if len(s.stack) == 0 {
		return fmt.Errorf("MarkThis: %v", errNoScopesAvailable)
	}
	s.thisScopes = append(s.thisScopes, s.stack[len(s.stack)-1])
	return nil
}

// ClearThis clears the current scope for "this" keyword.
func (s *scopeManager) ClearThis() error {
	if len(s.thisScopes) == 0 {
		return errors.New(`ClearThis called, but "this" scope does not exist`)
	}
	s.thisScopes = s.thisScopes[:len(s.thisScopes)-1]
	return nil
}

// SearchOption indicates the different search strategy for the findDeclaredIdentifier.
type SearchOption int

// TODO: Although right now we only have two possible options, in the future when there are more strategies needed we
// can always expand it. Moreover, we can represent them as bit masks so that we can easily combine the options like
// "EntireStack | IgnoreBlank". See https://golang.org/ref/spec#Iota for details.
const (
	// CurrentOnly indicates that the search only happens in the current scope instead of the entire stack of scopes.
	CurrentOnly SearchOption = iota
	// EntireStack indicates the search happens in the entire stack of scopes in reverse order.
	EntireStack
	// This indicates the search happens in the scope of "this" keyword.
	This
)

// FindDeclarationEntry searches the stack of scopes in reverse order, and returns the found entry.
// An option for search strategy can be provided, e.g., CurrentOnly or EntireStack, for searching
// entries only in the current scope or in the entire scope stack. Additionally, activeOnly can be
// set to search only active declaration entries (please see the comments for
// "SymbolTableEntry.IsActive" for more explanations on the activeness of a declaration entry).
// Nil is returned if the target identifier is not found with the given options.
func (s *scopeManager) FindDeclarationEntry(target *mast.Identifier, option SearchOption, activeOnly bool) (*SymbolTableEntry, error) {
	switch option {
	// special case for current-scope-only search
	case CurrentOnly:
		return s.findEntryInScope(target, s.currentScope(), activeOnly)

	case EntireStack:
		// iterate in reverse order of the scopes to find the declaration due to shadowing
		for i := len(s.stack) - 1; i >= 0; i-- {
			// check if the declaration exists in the current scope
			entry, err := s.findEntryInScope(target, s.stack[i], activeOnly)
			if err != nil {
				return nil, err
			}
			if entry != nil {
				return entry, nil
			}
		}

		// if still not found, check the package-level scope
		entry, err := s.findEntryInScope(target, s.currentPackageScope, activeOnly)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			return entry, nil
		}

		return nil, nil

	case This:
		if len(s.thisScopes) == 0 {
			return nil, errors.New(`no scope registered for "this"`)
		}

		entry, err := s.findEntryInScope(target, s.thisScopes[len(s.thisScopes)-1], activeOnly)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			return entry, nil
		}

		return nil, nil

	default:
		return nil, fmt.Errorf("unhandled search option: %v", option)
	}

}

// findEntryInScope finds entry in a scope for a given identifier
// based on the identifier kind.
func (s *scopeManager) findEntryInScope(target *mast.Identifier, curScope *scope, activeOnly bool) (*SymbolTableEntry, error) {
	existingScopeEntry, exists := curScope.table[target.Name]
	if !exists {
		return nil, nil
	}

	labelEntryExists := false
	var result *SymbolTableEntry
	for existingKind, existingEntry := range existingScopeEntry {
		if existingKind == mast.Label {
			labelEntryExists = true
		}
		if existingKind == target.Kind {
			// exact kind match
			result = existingEntry
		} else if target.Kind == mast.Blanket && existingKind != mast.Label {
			// identifier is of blanket (unknown) kind so it can be
			// matched with all declarations other than label
			result = existingEntry
		}
	}
	// sanity check - at this point we should not have more than one
	// declaration with the same name in a scope, unless the
	// additional declaration is for a label
	if labelEntryExists {
		if len(existingScopeEntry) > 2 {
			return nil, fmt.Errorf("too many (%d) entries in scope for identifier %q (including label)", len(existingScopeEntry), target.Name)
		}
	} else {
		if len(existingScopeEntry) > 1 {
			return nil, fmt.Errorf("too many (%d) entries in scope for identifier %q", len(existingScopeEntry), target.Name)
		}
	}

	if activeOnly && result != nil && !result.IsActive {
		// we found an inactive entry but should return only active ones
		return nil, nil
	}
	return result, nil
}

// AddDeclarationEntry adds the declaration entry to the current scope.
func (s *scopeManager) AddDeclarationEntry(entry *SymbolTableEntry, kind mast.NameKind) error {
	if entry == nil {
		return errors.New("adding a nil declaration entry")
	}

	// Here we add the declaration entry to the scope, if the scope stack is empty, the package scope will then be used.
	curScope := s.currentScope()

	// We might be re-visiting the same declaration node for top-level declarations since we did a quick pass on them
	// when we are handling mast.Root node for package level visibility. So here we only return an error if the
	// declaration already exists and it is _not_ the same as the one being added.
	existingScopeEntry, exists := curScope.table[entry.Identifier.Name]

	if exists {
		// make sure that an entry in given scope with the same kind does not exist
		if existingEntry, exists := existingScopeEntry[kind]; exists && existingEntry.Identifier != entry.Identifier {
			return fmt.Errorf("identifier %q already exists in the current scope", entry.Identifier.Name)
		}
	} else {
		existingScopeEntry = make(scopeEntry)
	}

	// put the declaration entry to the scope
	existingScopeEntry[kind] = entry
	curScope.table[entry.Identifier.Name] = existingScopeEntry

	return nil
}

// currentScope is a helper function that returns the current scope.
func (s *scopeManager) currentScope() *scope {
	// The current scope is either the top scope in the stack or the package level scope.
	var curScope *scope
	if len(s.stack) == 0 {
		curScope = s.currentPackageScope
	} else {
		curScope = s.stack[len(s.stack)-1]
	}
	return curScope
}

// String returns a string representation of the scopes
func (s *scopeManager) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "PACKAGE SCOPES:\n")
	for name, scope := range s.packageScopes {
		fmt.Fprintf(&b, "\tPKG %s -- %s\n", name, scopeString(scope))
	}
	if len(s.stack) == 0 {
		return b.String()
	}
	fmt.Fprintf(&b, "OTHER SCOPES:\n")
	for i := len(s.stack) - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "\t%s\n", scopeString(s.stack[i]))
	}
	return b.String()
}

// scopeString returns string representation of the scope.
func scopeString(s *scope) string {
	var b strings.Builder
	fmt.Fprintf(&b, "DECLS: ")
	for name, existingScopeEntry := range s.table {
		for kind := range existingScopeEntry {
			fmt.Fprintf(&b, "%s (%d) ", name, kind)
		}
	}
	return b.String()
}
