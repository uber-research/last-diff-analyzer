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
	"fmt"

	"analyzer/core/mast"
	"analyzer/core/mast/mastutil"
)

// GenericSymbolTableBuilder is the generic builder for building symbol information. It does all the bookkeeping of the
// underlying symbol table and scope manager. Note that the language-specific builder should embed this builder and
// "override" the Process* methods.
type GenericSymbolTableBuilder struct {
	// lang is the language-specific symbol table builder for the generic to call back to.
	lang SymbolTableBuilder
	// symbolTable is the pointer to the built symbol table which will be populated during the traversal of the MAST
	// tree.
	symbolTable *SymbolTable
	// scopes contains the stack of scopes as well as the package-level scopes, it provides access methods to
	// conveniently modify the scopes.
	scopes *scopeManager

	// ignoreNodes keeps track of the nodes to ignore. The children of some nodes will be properly processed when we are
	// visiting their parents. For example, the mast.Identifier nodes in the chain of mast.AccessPath will be properly
	// processed when we are visiting the outermost mast.AccessPath. However, the Walk function will still traverse the
	// children of it, meaning the builder will see the children nodes again later. Therefore, we keep a record of
	// nodes to be ignored to avoid double processing in this case.
	ignoreNodes map[mast.Node]bool

	// inactiveDecls is a stack of inactive symbol table entries
	// pushed to when starting to process declaration node introducing
	// inactive declarations and popped from after such declaration is
	// processed.
	inactiveDecls [][]*SymbolTableEntry
}

// SymbolTable returns the built symbol table.
func (g *GenericSymbolTableBuilder) SymbolTable() *SymbolTable {
	return g.symbolTable
}

// NewGenericSymbolTableBuilder returns a properly initialized *GenericSymbolTableBuilder. It takes a language specific
// symbol table builder to call back to. If it is nil the generic symbol table will fall back to itself. This
// is safe from cycles, since the language specific symbol table builder is only ever called with the _children_
// of the current node.
func NewGenericSymbolTableBuilder(lang SymbolTableBuilder) *GenericSymbolTableBuilder {
	builder := &GenericSymbolTableBuilder{
		symbolTable: newSymbolTable(),
		scopes:      newScopeManager(),
		ignoreNodes: make(map[mast.Node]bool),
	}

	// If no language specific builder is specified, we will fall back to itself. This eliminates the needs for nilness
	// testings during traversal.
	if lang == nil {
		builder.lang = builder
	} else {
		builder.lang = lang
	}
	return builder
}

// createDeclarationEntry is a helper function designed exclusively for
// ProcessDeclaration that creates a declaration entry on the _current_ scope.
// It will determine the privateness of the declaration based on the privateness
// of the outer scope and the language-specific characteristics from calling
// IsDeclarationPrivate method.
// Note that a new entry will _not_ be created if there exists an entry with the
// same identifier and kind within the _current_ scope in order to avoid
// creating duplicate declaration entries. Instead, the existing entry will be
// returned.
func (g *GenericSymbolTableBuilder) createDeclarationEntry(node mast.Node, ident *mast.Identifier, active bool) (*SymbolTableEntry, error) {
	// We could be processing declaration nodes more than one time (e.g., due to
	// pre-processing declarations nodes in mast.Root node for package-level
	// visibility or similar scenarios), we will only create a new declaration
	// entry when it really is new.
	existingEntry, err := g.scopes.FindDeclarationEntry(ident, CurrentOnly /* option */, false /* activeOnly */)
	if err != nil {
		return nil, err
	}
	if existingEntry != nil {
		// But here we respect the active flag and change it to make it
		// transparent for the caller as if it is a newly-created entry.
		// The caller will be responsible to change it back in later traversal.
		existingEntry.IsActive = active
		return existingEntry, nil
	}

	// If we reached here, we need to create a new declaration entry.

	// If the current scope is private, everything in it will also be private [1, 2]. Otherwise, we try to determine the
	// privateness from language-specific IsDeclarationPrivate method.
	// [1] Java: "A member (class, interface, field, or method) of a reference (class, interface, or array) type or a
	// constructor of a class type is accessible only if the type is accessible and the member or constructor is
	// declared to permit access" from https://docs.oracle.com/javase/specs/jls/se7/html/jls-6.html#jls-6.6.1
	// [2] For Go, the above statement also holds, the only "public" scope is the top-level scope (which may contain
	// private and public declarations). All other scopes are "private" and declarations there are private as well.
	// The only exception in Go are fields inside a struct declaration (we can have a public field inside a private struct), but these are handled in language-specific manner.
	isPrivate := g.scopes.IsCurrentScopePrivate()
	if !isPrivate {
		var err error
		isPrivate, err = g.lang.IsDeclarationPrivate(node, ident)
		if err != nil {
			return nil, err
		}
	}

	entry := &SymbolTableEntry{
		Identifier: ident,
		Node:       node,
		IsPrivate:  isPrivate,
		IsActive:   active,
	}
	if err := g.scopes.AddDeclarationEntry(entry, ident.Kind); err != nil {
		return nil, err
	}
	return entry, nil
}

// IsDeclarationPrivate directly returns true for all declaration nodes since it does not check language-specific
// details of the declaration node. This method is meant to be "overridden" by language-specific builders.
func (g *GenericSymbolTableBuilder) IsDeclarationPrivate(node mast.Node, identifier *mast.Identifier) (bool, error) {
	return false, nil
}

// ProcessDeclaration adds the declaration entry to the current scope.
func (g *GenericSymbolTableBuilder) ProcessDeclaration(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.PackageDeclaration:
		// The package declaration will not be used within the same file, therefore we do not register it in the scope.
		// "The package clause is not a declaration; the package name does not appear in any scope. Its purpose is to
		// identify the files belonging to the same package and to specify the default package name for import
		// declarations."
		// [1] https://golang.org/ref/spec#Declarations_and_scope
		// TODO: Actually in Java it is allowed to use the package name, we need to re-visit this when supporting
		//       full resolving of access paths.
		if err := g.nilifyNodeIDs(n.Name); err != nil {
			return err
		}
		return nil

	case *mast.VariableDeclaration:
		// We need a special case here for creating inactive declaration entries, please refer to
		// the comments for "SymbolTableEntry.IsActive" for more explanations.

		// We create a stack entry to contain the inactive declarations (refer to comments for
		// "SymbolTableEntry.IsActive" to better understand what inactive declarations mean). We
		// always create it (even if not actual declarations are present) so that we can
		// unilaterally pop it in PostProcessDeclaration.
		g.inactiveDecls = append(g.inactiveDecls, []*SymbolTableEntry{})

		for _, ident := range n.Names {
			entry, err := g.createDeclarationEntry(node, ident, false /* active */)
			if err != nil {
				return err
			}
			l := len(g.inactiveDecls)
			g.inactiveDecls[l-1] = append(g.inactiveDecls[l-1], entry)
			// add a "self" link for the declared identifier and mark
			// the identifier node as already handled
			if err := g.symbolTable.AddLink(ident, entry); err != nil {
				return err
			}
			g.ignoreNodes[ident] = true
		}
		return nil

	case *mast.ParameterDeclaration:
		// The Name field in mast.ParameterDeclaration is optional, so we skip if it is nil.
		if n.Name == nil {
			return nil
		}
		// Parameter declaration will always be private to the current file.
		entry := &SymbolTableEntry{Identifier: n.Name, Node: n, IsPrivate: true, IsActive: true}
		return g.scopes.AddDeclarationEntry(entry, mast.Blanket /* kind */)
	case *mast.LabelStatement:
		entry := &SymbolTableEntry{Identifier: n.Label, Node: n, IsPrivate: true, IsActive: true}
		return g.scopes.AddDeclarationEntry(entry, mast.Label /* kind */)
	}

	return nil
}

// ProcessScope creates a new scope for the scope-related nodes if we are entering the node, otherwise pops the current
// scope.
func (g *GenericSymbolTableBuilder) ProcessScope(node mast.Node, onEnter bool) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch node.(type) {
	case *mast.Block, *mast.IfStatement, *mast.ForStatement, *mast.SwitchStatement, *mast.SwitchCase,
		*mast.FunctionDeclaration, *mast.FunctionLiteral:
		// We should create private scopes for the nodes above.
		return g.handleScope(onEnter, true /* isPrivate */)
	}
	return nil
}

// ProcessUse searches the identifier in the stack of scopes and adds the appropriate information to the symbol table.
func (g *GenericSymbolTableBuilder) ProcessUse(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.Identifier:
		// Try to match the identifier with a declaration identifier. Note that we will not be able
		// to find a declaration entry for built-in identifiers (e.g., "int" and "make") or
		// the corresponding declarations for this identifier are in other files outside the
		// forests being analyzed.
		entry, err := g.scopes.FindDeclarationEntry(n, EntireStack, true /* activeOnly */)
		if err != nil {
			return err
		}
		return g.symbolTable.AddLink(n, entry)
	case *mast.AccessPath:
		// The logic for handling AccessPath is this:
		// (1) extract the identifier nodes in AccessPath in order with the slice of processed nodes along the way. This
		//     is done by using the helper function mast.ExtractAccessPath(path);
		// (2) process the chain of identifier nodes (implemented here), language-specific builder can "override" the
		//     handling here and do some language-specific stuff (e.g., Java builder will process "this" keyword
		//     properly);
		// (3) put the processed nodes in ignore set (after the full chain has been processed).

		// Note that (2) and (3) cannot be swapped, otherwise we cannot re-use ProcessIdentifier(ident) in (2) since
		// "ident" is already in the ignore set.

		// For simplicity, we only handle the first (the innermost) component of the AccessPath, the rest of the
		// mast.Identifier nodes along the way will be linked to a nil declaration identifier. For example, in the
		// access path "a.b.c", we will only try to link the identifier "a", the identifiers "b" and "c" will be linked
		// to nil.

		// call processAccessPath to get the unwrapped chain of identifiers in the mast.AccessPath node
		identifierChain, ignoreFirst, _, err := mastutil.ExtractAccessPath(n, nil)
		if err != nil {
			return err
		}

		// mark all identifiers in the returned chain as "unresolved"
		// (linked to nil) starting with this index:
		//
		// - start with 0 if path does not start with an identifier, e.g., (tmp3{[]int{42}}).a.b
		// - start with 1 if path starts with an identifier, e.g., a.b.c (resolve "a")
		unresolvedIndex := 0

		if !ignoreFirst {
			// resolve the first identifier in the access path
			err = g.lang.ProcessUse(identifierChain[0])
			if err != nil {
				return err
			}
			unresolvedIndex = 1
		}

		// The rest of the mast.Identifier nodes in the chain will simply be linked to a nil declaration.
		if err := g.nilifyIDs(identifierChain, unresolvedIndex); err != nil {
			return err
		}
	}

	return nil
}

// ProcessOther does the additional work for Root node to pre-register the child declaration nodes in the current scope
// for visibility throughout the entire package. Language-specific builders will "override" this method to handle any
// other nodes that require special handling for the particular language.
func (g *GenericSymbolTableBuilder) ProcessOther(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.Root:
		// The declarations in the package level should be visible in the entire package. That is, you can call function
		// "foo()" before it is declared later in the package. Therefore, we iterate through the declaration list and
		// put the declarations in the current scope. We do not go any deeper here; it will be done during the traversal
		// of its children.
		// A side-effect is that the declaration will exist in the current scope when we are actually visiting them
		// during children traversal, which might seem like a forbidden re-declaration. However, this is ok since a
		// re-declaration error is returned only when the declaration node and the one in the scope are different.
		// To do this, we will visit the declaration nodes before the actual traversal begins, which will put the
		// declared identifiers properly in the current scope. Note that the declaration nodes will be visited again
		// during the actual traversal, but it is fine since the createDeclarationEntry method is designed to ignore
		// identical declared entry.

		// skip if no declarations are in the root node
		if len(n.Declarations) == 0 {
			return nil
		}

		// Here we determine the package for this root node and assign the currentPackageScope properly by either getting
		// an existing package scope from the packageScope map or creating a new one if no scopes have been created for
		// this package yet.
		pkg, ok := n.Declarations[0].(*mast.PackageDeclaration)
		if ok {
			if err := g.scopes.CreatePackageScope(pkg); err != nil {
				return err
			}
		} else {
			// If no package declaration is present, we create a standalone scope for this particular file, which will
			// simply be dropped when we move to the next MAST in the forest. Therefore we pass a nil node to
			// CreatePackageScope to properly handle that. Note that it is still required to do so, otherwise the
			// scope manager will not create an empty package-level scope for this file. This is mainly designed for
			// unnamed packages [1] in Java and for the languages that do not have a "package" concept.
			// [1] https://docs.oracle.com/javase/specs/jls/se11/html/jls-7.html#jls-7.4
			if err := g.scopes.CreatePackageScope(nil); err != nil {
				return err
			}
		}

		// process the package-level declarations
		for _, decl := range n.Declarations {
			if err := g.lang.ProcessDeclaration(decl); err != nil {
				return err
			}
			if err := g.lang.PostProcessDeclaration(decl); err != nil {
				return err
			}
		}

	}

	return nil
}

// handleScope is a helper function that simply creates a new scope on enter and pops a scope otherwise.
func (g *GenericSymbolTableBuilder) handleScope(onEnter bool, isPrivate bool) error {
	// create/pop the scope based on onEnter
	if onEnter {
		g.scopes.CreateNewScope(isPrivate)
	} else {
		if err := g.scopes.PopScope(); err != nil {
			return err
		}
	}
	return nil
}

// nilifyNodeIDs links all identifiers in a given node nil.
func (g *GenericSymbolTableBuilder) nilifyNodeIDs(node mast.Node) error {
	switch n := node.(type) {
	case *mast.Identifier:
		g.ignoreNodes[n] = true
		if err := g.symbolTable.AddLink(n, nil); err != nil {
			return err
		}
	case *mast.AccessPath:
		identifierChain, _, _, err := mastutil.ExtractAccessPath(n, nil)
		if err != nil {
			return err
		}
		if err := g.nilifyIDs(identifierChain, 0); err != nil {
			return err
		}

	default:
		// just in case - we want to nilify only "known" nodes so
		// let's signal if this does not hold
		return fmt.Errorf("unhandled node %T when nilifying identifiers", node)
	}
	return nil
}

// nilifyIds links all identifiers in a given array to nil.
func (g *GenericSymbolTableBuilder) nilifyIDs(ids []*mast.Identifier, start int) error {
	for i := start; i < len(ids); i++ {
		err := g.symbolTable.AddLink(ids[i], nil)
		if err != nil {
			return err
		}
		// ignore the processed node
		g.ignoreNodes[ids[i]] = true
	}

	return nil
}

// PostProcessDeclaration handles declaration node after they have been visited.
func (g *GenericSymbolTableBuilder) PostProcessDeclaration(node mast.Node) error {
	if g.ignoreNodes[node] {
		return nil
	}

	switch node.(type) {
	case *mast.VariableDeclaration:
		// We need to "activate" this declaration so that it's
		// available for all statements following this one (please see
		// mast.VariableDeclaration handling in ProcessDeclaration for
		// additional explanation).
		g.activateDeclarations()
	}
	return nil
}

// activateDeclarations changes state of a declaration from "inactive"
// (when it should not be available yet, for example to its own RHS)
// to "active".
func (g *GenericSymbolTableBuilder) activateDeclarations() {
	l := len(g.inactiveDecls)
	inactive := g.inactiveDecls[l-1]
	// pop the inactive declarations stack
	g.inactiveDecls = g.inactiveDecls[:l-1]
	for _, entry := range inactive {
		entry.IsActive = true
	}
}

// PostSymbolicationFixup is responsible for any actions that need to
// be performed after the whole symbolication process is finished.
func (g *GenericSymbolTableBuilder) PostSymbolicationFixup() error {
	// nothing to do for the generic case
	return nil
}
