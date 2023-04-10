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
	ts "analyzer/core/treesitter"
)

func entryNotExistError(node mast.Node, name *mast.Identifier) error {
	return fmt.Errorf("declaration entry does not exist for %T: %v", node, name.Name)
}

func langFieldsNodeNotExistsError(node mast.Node) error {
	return fmt.Errorf("lang field node does not exist for %T", node)
}

// SymbolTableBuilder is the interface that all symbol builders must implement. It consists of four processors which
// will be called sequentially during the MAST node traversal. See the comments of each method for detailed explanations.
type SymbolTableBuilder interface {
	// IsDeclarationPrivate returns whether a MAST declaration node is declared private or not. For example, in Go it
	// will return true if the first letter of the identifier is upper-case and in Java it will return true if the
	// declaration does not have "public" or "protected" modifiers.
	// Since some declaration nodes will declare several variables, each identifier in the declaration will be queried
	// via the given identifier node.
	IsDeclarationPrivate(node mast.Node, identifier *mast.Identifier) (bool, error)
	// ProcessDeclaration handles declaration nodes when they are visited.
	ProcessDeclaration(node mast.Node) error
	// ProcessScope for scoping-related nodes: append a scope to the scope stack; onEnter indicates whether the
	// traversal is entering the node or not.
	ProcessScope(node mast.Node, onEnter bool) error
	// ProcessUse for "use" nodes (i.e., mast.AccessPath and mast.Identifier etc.): search the required
	// information in the scope stack and build the symbol information in the symbol table.
	ProcessUse(node mast.Node) error
	// ProcessOther for doing any extra work: for example, for nodes that have a global scope (i.e., visible throughout
	// the entire package), we need to process all the declaration nodes in it _before_ the traversal continues.
	ProcessOther(node mast.Node) error
	// SymbolTable returns the constructed symbol table.
	SymbolTable() *SymbolTable
	// PostProcessDeclaration handles declaration node after they have been visited.
	PostProcessDeclaration(node mast.Node) error
	// PostSymbolicationFixup is responsible for any actions that need
	// to be performed after the whole symbolication process is
	// finished.
	PostSymbolicationFixup() error
}

// driver is the driver visitor for the symbol table builder, it implements the required interfaces for mast.Visitor and
// dispatches the calls to the corresponding SymbolTableBuilder methods in defined order.
type driver struct {
	// builder is the builder pointer to the actual symbol table builder.
	builder SymbolTableBuilder
}

// Pre implements the required Pre interface for mast.Visitor and calls the corresponding SymbolTableBuilder methods in
// defined order.
func (v *driver) Pre(node mast.Node) error {
	// (1) call ProcessDeclaration method for registering the declaration to the current scope
	if err := v.builder.ProcessDeclaration(node); err != nil {
		return err
	}
	// (2) call ProcessScope to add scope
	if err := v.builder.ProcessScope(node, true /* onEnter */); err != nil {
		return err
	}
	// (3) call ProcessUse to add information to the symbol table
	if err := v.builder.ProcessUse(node); err != nil {
		return err
	}
	// (4) call ProcessOther to do any additional work, such as pre-registering the declaration nodes for the global
	//     scope.
	if err := v.builder.ProcessOther(node); err != nil {
		return err
	}
	return nil
}

// Post implements the required Post interface for mast.Visitor and only calls ProcessScope method to handle leaving the
// scope.
func (v *driver) Post(node mast.Node) error {
	if err := v.builder.PostProcessDeclaration(node); err != nil {
		return err
	}
	if err := v.builder.ProcessScope(node, false /* onEnter */); err != nil {
		return err
	}
	return nil
}

// Run is the main driver of the symbolication process. It traverses the given MAST forests and returns a map that links
// from the Identifier node to its corresponding declaration node. All nodes in the forests must be *mast.Root node.
func Run(forest []mast.Node, suffix string) (*SymbolTable, error) {
	var builder SymbolTableBuilder
	switch suffix {
	case ts.JavaExt:
		builder = NewJavaSymbolTableBuilder()
	case ts.GoExt:
		builder = NewGoSymbolTableBuilder()
	default:
		return nil, fmt.Errorf("unsupported file extension %q during symbolication", suffix)
	}

	// create the driver visitor and add the builder to it
	driver := &driver{builder: builder}

	// Here we use the entire forest to build the symbol table for cross-file visibility of package-level declarations.
	for _, node := range forest {
		root, ok := node.(*mast.Root)
		if !ok {
			return nil, fmt.Errorf("non-root node in the forest for symbolication is currently not supported: %T", node)
		}
		if err := builder.ProcessOther(root); err != nil {
			return nil, err
		}
	}

	// do the actual walk
	for _, node := range forest {
		err := mast.Walk(driver, node)
		if err != nil {
			return nil, err
		}
	}

	if err := builder.PostSymbolicationFixup(); err != nil {
		return nil, err
	}

	return builder.SymbolTable(), nil
}

// processTypeScope is a helper function for processing type-related
// scopes (e.g., interface scope).
func processTypeScope(s *scopeManager, name *mast.Identifier, node mast.Node, onEnter bool) error {
	// special case for exiting the scope
	if !onEnter {
		return s.PopScope()
	}

	// first retrieve the symbol entry for the class declaration for later use
	entry, err := s.FindDeclarationEntry(name, CurrentOnly, false /* activeOnly */)
	if err != nil {
		return err
	}
	if entry == nil {
		return entryNotExistError(node, name)
	}

	// The privateness of the scope for the class will be determined by the privateness of the declaration itself.
	s.CreateNewScope(entry.IsPrivate /* isPrivate */)
	return nil
}
