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

package translation

import (
	"fmt"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"
)

// Translator is an interface that all language translators must implement.
type Translator interface {
	// Translate translates the tree-sitter node to MAST node. The translated
	// MAST node (mast.Node) or nil is returned if the translation does not
	// generate any node. Specifically, a mast.TempGroupNode is returned
	// indicating the ts.Node is translated into multiple mast.Nodes.
	// See mast.TempGroupNode for details.
	Translate(node *ts.Node) (mast.Node, error)
}

// Run is the main driver of the translation process. It translates the given ts.Node into a mast.Node using a
// language-specific translator by the given suffix.
func Run(node *ts.Node, suffix string) (mast.Node, error) {
	// initialize translate function by suffix
	var translator Translator
	switch suffix {
	case ts.GoExt:
		translator = &GoTranslator{}
	case ts.JavaExt:
		translator = &JavaTranslator{}
	default:
		return nil, fmt.Errorf("unsupported file extension %q", suffix)
	}

	n, err := translator.Translate(node)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// translateNodes is a helper function that simply translates all nodes in the input slice and return the translated
// node in another slice.
// Nil values will be part of the return slice if and only if the corresponding nodes in the _input_ slice are nil.
// If a translation would result in a non-nil input node being converted into a nil in the output mast.Node array,
// this method will instead return an error.
// If shouldUngroup is set to true, it will also un-group any mast.TempGroupNode returned from the lower-level
// translations, making the length of the result slice potentially different from the original slice.
// Otherwise, the lengths will be the same (an error will be returned if not).
func translateNodes(t Translator, nodes []*ts.Node, shouldUngroup bool) ([]mast.Node, error) {
	// the length of the result slice cannot be determined due to potential TempGroupNode nodes, so here we create an
	// array with a capacity of len(nodes) to reduce copying.
	result := make([]mast.Node, 0, len(nodes))

	// translate each node and append the translated node to result
	for _, node := range nodes {
		// put nil in result if the input node is nil
		if node == nil {
			result = append(result, nil)
			continue
		}
		translated, err := t.Translate(node)
		if err != nil {
			return nil, err
		}
		// translation should never return nil.
		if translated == nil {
			return nil, fmt.Errorf("translation returned nil for node %q", node.Type)
		}
		// Translate might return a TempGroupNode wrapping a group of translated node, so here we un-group it and append
		// the nodes to result if shouldUngroup flag is set to true.
		n, ok := translated.(*mast.TempGroupNode)
		if ok {
			if !shouldUngroup {
				return nil, fmt.Errorf("translation for %q returned TempGroupNode, but shouldUngroup flag is set to false", node.Type)
			}
			result = append(result, n.Nodes...)
		} else {
			result = append(result, translated)
		}
	}
	return result, nil
}

// toIdentifiers casts nodes to *mast.Identifier nodes, an error will be returned if any of the casts fail.
// "nil"s in the input slice will be preserved in the result slice.
func toIdentifiers(nodes []mast.Node) ([]*mast.Identifier, error) {
	// apply type assertions on the nodes and return error if any assertion fails
	result := make([]*mast.Identifier, len(nodes))
	for i, node := range nodes {
		// nil is allowed in result slice if the input is also nil
		if node == nil && nodes[i] == nil {
			continue
		}
		n, ok := node.(*mast.Identifier)
		if !ok {
			return nil, nodeTypeError(node)
		}
		result[i] = n
	}

	return result, nil
}

// toExpressions casts nodes to mast.Expression nodes, an error will be returned if any of the casts fail.
// "nil"s in the input slice will be preserved in the result slice.
func toExpressions(nodes []mast.Node) ([]mast.Expression, error) {
	// apply type assertions on the nodes and return error if any assertion fails
	result := make([]mast.Expression, len(nodes))
	for i, node := range nodes {
		// nil is allowed in result slice if the input is also nil.
		if node == nil {
			continue
		}
		n, ok := node.(mast.Expression)
		if !ok {
			return nil, nodeTypeError(node)
		}
		result[i] = n
	}

	return result, nil
}

// toStatements casts nodes to mast.Statement nodes, an error will be returned if any of the casts fail.
// "nil"s in the input slice will be preserved in the result slice.
func toStatements(nodes []mast.Node) ([]mast.Statement, error) {
	// apply type assertions on the nodes and return error if any assertion fails
	result := make([]mast.Statement, len(nodes))
	for i, node := range nodes {
		// nil is allowed in result slice if the input is also nil.
		if node == nil {
			continue
		}
		// tree-sitter does not generate expression_statement properly for Go AST, instead it directly generates
		// expression node even if the expression is in a block of statements. So here we do a post-fix for all
		// translated mast.Expression nodes and wrap around them with mast.ExpressionStatement.
		// Similarly, we do the same for DeclarationStatement.
		switch n := node.(type) {
		case mast.Statement:
			result[i] = n
		case mast.Expression:
			wrapped := &mast.ExpressionStatement{Expr: n}
			result[i] = wrapped
		case mast.Declaration:
			wrapped := &mast.DeclarationStatement{Decl: n}
			result[i] = wrapped
		default:
			// an error is returned if any other type of node is found.
			return nil, nodeTypeError(node)
		}
	}

	return result, nil
}

// toDeclarations casts nodes to mast.Declaration nodes, an error will be returned if any of the casts fail.
// "nil"s in the input slice will be preserved in the result slice.
func toDeclarations(nodes []mast.Node) ([]mast.Declaration, error) {
	// apply type assertions on the nodes and return error if any assertion fails
	result := make([]mast.Declaration, len(nodes))
	for i, node := range nodes {
		// nil is allowed in result slice if the input is also nil.
		if node == nil {
			continue
		}
		n, ok := node.(mast.Declaration)
		if !ok {
			return nil, nodeTypeError(node)
		}
		result[i] = n
	}

	return result, nil
}
