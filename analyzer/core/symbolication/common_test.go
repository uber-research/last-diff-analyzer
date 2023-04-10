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
	"strconv"
	"testing"

	"analyzer/core/mast"
	"analyzer/core/translation"
	ts "analyzer/core/treesitter"

	"github.com/stretchr/testify/require"
)

const _metaTestDataPrefix = "../../testdata/symbolication/"

// _recorderVisitor simply records all the identifier nodes that has a number suffix during traversal.
type _recorderVisitor struct {
	// Visited is a slice of identifier nodes visited during traversal.
	Visited []*mast.Identifier
}

// Pre simply puts any identifier nodes that has a number suffix into its Visited slice.
func (r *_recorderVisitor) Pre(node mast.Node) error {
	if n, ok := node.(*mast.Identifier); ok && len(n.Name) != 0 {
		// check if the identifier's name ends with a digit
		_, err := strconv.ParseInt(string(n.Name[len(n.Name)-1]), 10 /* base */, 0 /* bitSize */)
		if err == nil {
			r.Visited = append(r.Visited, n)
		}
	}
	return nil
}

// Post is implemented so that _recorderVisitor satisfies the mast.Visitor interface, it does nothing and returns nil.
func (r *_recorderVisitor) Post(node mast.Node) error { return nil }

// recordNames records the identifiers with a number suffix in a map for easy access. Each identifier node may appear multiple
// times in the test files, so we use a slice to keep track of all identifier nodes. The length of each slice must match
// the length given in the expectedLengths.
func recordNames(t *testing.T, node mast.Node, expectedLengths map[string]int) map[string][]*mast.Identifier {
	// create a names map based on the expectedLengths map
	names := make(map[string][]*mast.Identifier)
	for name := range expectedLengths {
		names[name] = nil
	}

	// traverse the MAST node and records all identifier nodes
	recorder := &_recorderVisitor{}
	err := mast.Walk(recorder, node)
	require.NoError(t, err)

	for _, ident := range recorder.Visited {
		_, exists := names[ident.Name]
		require.True(t, exists, "unexpected identifier %q", ident.Name)
		names[ident.Name] = append(names[ident.Name], ident)
	}

	// check if the lengths match
	for name, identifiers := range names {
		require.Equal(t, expectedLengths[name], len(identifiers), "%q has unexpected length of identifiers", name)
	}

	return names
}

// getNames performs the same actions as recordNames but takes source
// file information as parameters instead of an already translated
// MAST node (the expectedLength parameter serves the same role as in
// recordNames - it maps identifier names to the number of times they
// appear in the test file).
func getNames(t *testing.T, fname string, ext string, expectedLength map[string]int) (*SymbolTable, map[string][]*mast.Identifier) {
	node, err := ts.ParseFile(fname)
	require.NoError(t, err)
	actual, err := translation.Run(node, ext)
	require.NoError(t, err)

	symbolTable, err := Run([]mast.Node{actual}, ext)
	require.NoError(t, err)
	require.NotNil(t, symbolTable)

	return symbolTable, recordNames(t, actual, expectedLength)
}

// verifySymbols verifies if symbolic information is correct.  The
// expectedLinks parameter maps each identifier name (as a string) to
// a map of uses of identifiers with the same name to their
// corresponding declarations, which are expected to match those in
// the symbolTable parameter.
func verifySymbols(t *testing.T, symbolTable *SymbolTable, expectedLinks map[string]map[*mast.Identifier]*mast.Identifier) {
	for name, link := range expectedLinks {
		ind := 0
		for use, def := range link {
			defInTable, err := symbolTable.DeclarationEntry(use)
			require.NoError(t, err)
			if defInTable == nil {
				require.True(t, def == nil, "identifier %q (%p) at index %d has incorrect link %p, expected %p", name, use, ind, defInTable, def)
			} else {
				require.True(t, def == defInTable.Identifier, "identifier %q (%p) at index %d has incorrect link %p, expected %p", name, use, ind, defInTable, def)
			}
			ind++
		}
	}

}
