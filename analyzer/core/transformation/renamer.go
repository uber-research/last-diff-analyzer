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

package transformation

import (
	"strconv"

	"analyzer/core/mast"
	"analyzer/core/symbolication"
)

// Renamer replaces each relevant (i.e., whose renaming we currently handle) identifier in each AST with a freshly
// generated one that depends only on the tree structure (i.e., it is independent on the source code). Two ASTs for
// files that differ only in identifier names will be equal after such treatment.
type Renamer struct {
	// symbolTable is the generated symbol table from the symbolicaton process, which will be referenced for renaming.
	symbolTable *symbolication.SymbolTable
}

// NewRenamer returns a properly-initialized Renamer.
func NewRenamer(symbolTable *symbolication.SymbolTable) *Renamer {
	return &Renamer{
		symbolTable: symbolTable,
	}
}

// Transform renames all identifiers in the node and returns the transformed node.
func (r *Renamer) Transform(forest []mast.Node) ([]mast.Node, error) {
	symbols := r.symbolTable.OrderedSymbols()

	renamed := map[*mast.Identifier]bool{}

	count := 0
	for _, symbol := range symbols {
		declEntry, err := r.symbolTable.DeclarationEntry(symbol)
		if err != nil {
			return nil, err
		}
		if declEntry == nil {
			continue
		}
		if vdecl, ok := declEntry.Node.(*mast.VariableDeclaration); ok && vdecl.IsConst {
			// do not rename constants as they are handled differently
			// to allow for constant addition and removal
			continue
		}
		// not renamed yet
		if !renamed[declEntry.Identifier] {
			// skip renaming the public declarations
			if !declEntry.IsPrivate {
				continue
			}
			// rename with a $ prefix to avoid collisions with existing names
			newName := "$_renamed_declaration_" + strconv.Itoa(count)
			declEntry.Identifier.Name = newName
			count++
			renamed[declEntry.Identifier] = true
		}
		symbol.Name = declEntry.Identifier.Name
	}

	// We have actually modified the node from referencing the symbol table, so here we directly return itself.
	return forest, nil
}
