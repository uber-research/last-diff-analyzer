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

package mastutil

import (
	"fmt"
	"strings"

	"analyzer/core/mast"
)

// ExtractAccessPath is a helper function to extract the _operand_ Identifier
// nodes in the AccessPath in order. (i.e., access path "a.b[c].d" will return
// ["a", "b", "d"]). While extracting the Identifier, it also searches for a
// specific identifier (there should be only one such identifier on the path)
// and returns its prefix.
func ExtractAccessPath(path *mast.AccessPath, searchedIDs map[string]bool) (chain []*mast.Identifier, ignoreFirst bool, foundPrefix mast.Expression, err error) {
	// We iteratively unwrap the AccessPath and store the identifier nodes
	// to chain slice until we hit the innermost Identifier.

	// We will use the current variable to keep track of the iteration (by iteratively assigning Operand field of
	// AccessPath to the current variable).
	var current mast.Expression = path
	shouldContinue := true

	// In access paths starting with an identifier, the first identifier
	// is treated specially. This flag keeps track if it should or
	// should not be treated as such.
	ignoreFirst = false
	for shouldContinue {
		switch n := current.(type) {
		case *mast.AccessPath:
			chain = append(chain, n.Field)
			current = n.Operand

			// Here (and only here) we also look for special
			// identifiers (such as Java's "this" and "super"). The
			// reason for it is that this is where we identify the
			// access path's prefix with respect to a given
			// identifier. AccessPath is a recursive structure
			// where Field is the last identifier of the path and
			// Operand is its prefix, so with Field being some
			// identifier foo and Operand being [...] (an arbitrarily
			// complicated prefix expression), the access path is
			// [...].foo.

			if searchedIDs[n.Field.Name] {
				// found special identifier
				if foundPrefix != nil {
					return nil, ignoreFirst, nil, fmt.Errorf("incorrect additional special identifer %s in access path %T", n.Field.Name, path)
				}
				foundPrefix = n.Operand
			}

		case *mast.CallExpression:
			current = n.Function
		case *mast.IndexExpression:
			current = n.Operand

		// The AccessPath node structure is nested: the first identifier of the access path is actually nested in
		// the innermost level. Therefore, if we encounter an Identifier node, we should break the loop.
		case *mast.Identifier:
			chain = append(chain, n)
			// we cannot go further so we break the loop
			shouldContinue = false

		// For expressions like "(&logger).a.b.c", we could be eventually hitting a ParenthesizedExpression or
		// UnaryExpression, therefore we properly handle them (by simply unwrapping them) here.
		case *mast.ParenthesizedExpression:
			current = n.Expr
		case *mast.UnaryExpression:
			current = n.Expr
		case *mast.EntityCreationExpression:
			// If we reached this node (which is essentially a
			// constructor), it means that we are at the beginning of
			// the access path (as a constructor cannot be part of
			// access path happening in the middle or at the
			// end). Consequently, this path does not start with an
			// identifier, so the first (but not starting) identifier
			// on the path should not be treated specially (should be
			// ignored).
			ignoreFirst = true
			// we cannot go further so we break the loop
			shouldContinue = false
		default:
			return nil, ignoreFirst, nil, fmt.Errorf("unhandled node type %T in access path", current)
		}
	}

	// The chain stores the identifier nodes in reverse order, so we reverse it back.
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	return chain, ignoreFirst, foundPrefix, err
}

// JoinAccessPath joins the _operands_ of a mast.AccessPath node and returns the string
// representation of it. For example, mast.AccessPath("a.b[c].d") will converted to "a.b.d".
func JoinAccessPath(path *mast.AccessPath) (string, error) {
	chain, _, _, err := ExtractAccessPath(path, nil /* searchedIDs */)
	if err != nil {
		return "", err
	}

	elements := make([]string, 0, len(chain))
	for _, e := range chain {
		elements = append(elements, e.Name)
	}
	return strings.Join(elements, "."), nil
}
