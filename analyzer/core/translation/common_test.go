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
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"

	"github.com/stretchr/testify/require"
)

const _metaTestDataPrefix = "../../testdata/translation/"

// reflectVisit is a helper function that uses reflect package to recursively visit (i.e., call the
// visitFunc function parameter) all fields of the node structure. Note that only the following
// types of fields will be visited, everything else will be ignored:
// (1) Non-nil Interface or Pointer;
// (2) Slice or Array;
// (3) Struct.
// T is the type for all nodes in the AST.
// This function is used to serve as a ground truth for our testing.
func reflectVisit[T any](node reflect.Value, visitFunc func(T)) error {
	kind := node.Kind()

	// Early return if node is nil
	if kind == reflect.Pointer || kind == reflect.Interface || kind == reflect.Slice || kind == reflect.Array {
		if node.IsNil() {
			return nil
		}
	}

	// Process the visit for different types
	switch kind {
	case reflect.Slice, reflect.Array:
		// Recursively call reflectVisit for each of the element for slice / array types
		for j := 0; j < node.Len(); j++ {
			if err := reflectVisit(node.Index(j), visitFunc); err != nil {
				return err
			}
		}

		return nil

	case reflect.Pointer, reflect.Interface, reflect.Struct:
		// There could be fields attached to the AST structures that are unrelated to nodes, e.g.,
		// a "DebugInfo" struct fields. So here we _try_ to cast the potential node types to the
		// desired node type T. If
		// (1) true  -> call visitFunc and further inspect its fields for recursion;
		// (2) false -> it is an unrelated node, and we should simply ignore it.
		castNode, ok := node.Interface().(T)
		if !ok {
			return nil
		}

		// Call the visit function
		visitFunc(castNode)

		// Now we visit all fields of the node. First, we need to (iteratively) unwrap the
		// interface / pointer type to expose the innermost struct type for finding its fields,
		// e.g., "*ASTNode" -> "ASTNode".
		for node.Kind() == reflect.Pointer || node.Kind() == reflect.Interface {
			node = node.Elem()
		}

		// If after unwrapping we find out it is not a struct type (e.g., a "**string" -> "string"),
		// it means reflectVisit is called with a wrong generic type T.
		if node.Kind() != reflect.Struct {
			return fmt.Errorf(
				"struct type expected for innermost type of node type, got %s",
				node.Type().Name(),
			)
		}

		// Otherwise, we reflectVisit each of the field.
		for i := 0; i < node.Type().NumField(); i++ {
			field := node.Field(i)
			if err := reflectVisit(field, visitFunc); err != nil {
				return err
			}
		}

		return nil

	// Other types are not interesting to us, including maps, since they are usually not part of
	// the AST structure.
	default:
		return nil
	}
}

func TestUnsupportedExtension(t *testing.T) {
	t.Run("Test translating an unsupported file extension", func(t *testing.T) {
		// create a root golang node
		node := &ts.Node{
			Type: "source_file",
		}

		// a random suffix is added just to be safe.
		mastNode, err := Run(node, ".unsupported_B751DD")
		require.Nil(t, mastNode)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported file extension")
	})
}

func TestCorruptedNode(t *testing.T) {
	t.Run("Test translating a corrupted ts.Node", func(t *testing.T) {
		// create a corrupted golang ts.Node
		node := &ts.Node{
			Type: "source_file",
			Children: []*ts.Node{
				{
					Type: "package_clause",
					// package_clause can only have one child: package_identifier,
					// here we create two package_identifier nodes.
					Children: []*ts.Node{
						{
							Type: "package_identifier",
							Name: "test",
						},
						{
							Type: "package_identifier",
							Name: "test2",
						},
					},
				},
			},
		}

		mastNode, err := Run(node, ts.GoExt)
		require.Nil(t, mastNode)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected number of children")
	})
}

// _allMASTNodeTypes is the set of all MAST node types. In order to get reflect.Type of each MAST
// node, we instantiate a typed nil pointer with each of MAST nodes (to avoid actual allocation of
// any structs) and call reflect.TypeOf to get their type. Then later in the actual test we can
// simply compare the visited nodes and this set to find if we are missing anything in the testing.
var _allMASTNodeTypes = map[reflect.Type]bool{
	reflect.TypeOf((*mast.Root)(nil)):                               true,
	reflect.TypeOf((*mast.Block)(nil)):                              true,
	reflect.TypeOf((*mast.PackageDeclaration)(nil)):                 true,
	reflect.TypeOf((*mast.ImportDeclaration)(nil)):                  true,
	reflect.TypeOf((*mast.ExpressionStatement)(nil)):                true,
	reflect.TypeOf((*mast.DeclarationStatement)(nil)):               true,
	reflect.TypeOf((*mast.ContinueStatement)(nil)):                  true,
	reflect.TypeOf((*mast.BreakStatement)(nil)):                     true,
	reflect.TypeOf((*mast.ReturnStatement)(nil)):                    true,
	reflect.TypeOf((*mast.SwitchStatement)(nil)):                    true,
	reflect.TypeOf((*mast.SwitchCase)(nil)):                         true,
	reflect.TypeOf((*mast.IfStatement)(nil)):                        true,
	reflect.TypeOf((*mast.LabelStatement)(nil)):                     true,
	reflect.TypeOf((*mast.Identifier)(nil)):                         true,
	reflect.TypeOf((*mast.ParenthesizedExpression)(nil)):            true,
	reflect.TypeOf((*mast.UnaryExpression)(nil)):                    true,
	reflect.TypeOf((*mast.BinaryExpression)(nil)):                   true,
	reflect.TypeOf((*mast.IndexExpression)(nil)):                    true,
	reflect.TypeOf((*mast.AccessPath)(nil)):                         true,
	reflect.TypeOf((*mast.CallExpression)(nil)):                     true,
	reflect.TypeOf((*mast.NullLiteral)(nil)):                        true,
	reflect.TypeOf((*mast.BooleanLiteral)(nil)):                     true,
	reflect.TypeOf((*mast.IntLiteral)(nil)):                         true,
	reflect.TypeOf((*mast.FloatLiteral)(nil)):                       true,
	reflect.TypeOf((*mast.StringLiteral)(nil)):                      true,
	reflect.TypeOf((*mast.CharacterLiteral)(nil)):                   true,
	reflect.TypeOf((*mast.UpdateExpression)(nil)):                   true,
	reflect.TypeOf((*mast.AssignmentExpression)(nil)):               true,
	reflect.TypeOf((*mast.ParameterDeclaration)(nil)):               true,
	reflect.TypeOf((*mast.VariableDeclaration)(nil)):                true,
	reflect.TypeOf((*mast.ForStatement)(nil)):                       true,
	reflect.TypeOf((*mast.KeyValuePair)(nil)):                       true,
	reflect.TypeOf((*mast.LiteralValue)(nil)):                       true,
	reflect.TypeOf((*mast.EntityCreationExpression)(nil)):           true,
	reflect.TypeOf((*mast.FunctionLiteral)(nil)):                    true,
	reflect.TypeOf((*mast.FieldDeclaration)(nil)):                   true,
	reflect.TypeOf((*mast.FunctionDeclaration)(nil)):                true,
	reflect.TypeOf((*mast.CastExpression)(nil)):                     true,
	reflect.TypeOf((*mast.GoSliceExpression)(nil)):                  true,
	reflect.TypeOf((*mast.GoEllipsisExpression)(nil)):               true,
	reflect.TypeOf((*mast.GoImaginaryLiteral)(nil)):                 true,
	reflect.TypeOf((*mast.GoPointerType)(nil)):                      true,
	reflect.TypeOf((*mast.GoArrayType)(nil)):                        true,
	reflect.TypeOf((*mast.GoMapType)(nil)):                          true,
	reflect.TypeOf((*mast.GoParenthesizedType)(nil)):                true,
	reflect.TypeOf((*mast.GoChannelType)(nil)):                      true,
	reflect.TypeOf((*mast.GoFunctionType)(nil)):                     true,
	reflect.TypeOf((*mast.GoTypeAssertionExpression)(nil)):          true,
	reflect.TypeOf((*mast.GoTypeSwitchHeaderExpression)(nil)):       true,
	reflect.TypeOf((*mast.GoDeferStatement)(nil)):                   true,
	reflect.TypeOf((*mast.GoGotoStatement)(nil)):                    true,
	reflect.TypeOf((*mast.GoFallthroughStatement)(nil)):             true,
	reflect.TypeOf((*mast.GoSendStatement)(nil)):                    true,
	reflect.TypeOf((*mast.GoGoStatement)(nil)):                      true,
	reflect.TypeOf((*mast.GoTypeDeclaration)(nil)):                  true,
	reflect.TypeOf((*mast.GoStructType)(nil)):                       true,
	reflect.TypeOf((*mast.GoInterfaceType)(nil)):                    true,
	reflect.TypeOf((*mast.GoFieldDeclarationFields)(nil)):           true,
	reflect.TypeOf((*mast.GoForRangeStatement)(nil)):                true,
	reflect.TypeOf((*mast.GoSelectStatement)(nil)):                  true,
	reflect.TypeOf((*mast.GoCommunicationCase)(nil)):                true,
	reflect.TypeOf((*mast.GoFunctionDeclarationFields)(nil)):        true,
	reflect.TypeOf((*mast.JavaTernaryExpression)(nil)):              true,
	reflect.TypeOf((*mast.Annotation)(nil)):                         true,
	reflect.TypeOf((*mast.JavaAnnotatedType)(nil)):                  true,
	reflect.TypeOf((*mast.JavaGenericType)(nil)):                    true,
	reflect.TypeOf((*mast.JavaWildcard)(nil)):                       true,
	reflect.TypeOf((*mast.JavaArrayType)(nil)):                      true,
	reflect.TypeOf((*mast.JavaDimension)(nil)):                      true,
	reflect.TypeOf((*mast.JavaInstanceOfExpression)(nil)):           true,
	reflect.TypeOf((*mast.JavaLiteralModifier)(nil)):                true,
	reflect.TypeOf((*mast.JavaTryStatement)(nil)):                   true,
	reflect.TypeOf((*mast.JavaCatchClause)(nil)):                    true,
	reflect.TypeOf((*mast.JavaCatchFormalParameter)(nil)):           true,
	reflect.TypeOf((*mast.JavaFinallyClause)(nil)):                  true,
	reflect.TypeOf((*mast.JavaWhileStatement)(nil)):                 true,
	reflect.TypeOf((*mast.JavaThrowStatement)(nil)):                 true,
	reflect.TypeOf((*mast.JavaAssertStatement)(nil)):                true,
	reflect.TypeOf((*mast.JavaSynchronizedStatement)(nil)):          true,
	reflect.TypeOf((*mast.JavaDoStatement)(nil)):                    true,
	reflect.TypeOf((*mast.JavaParameterDeclarationFields)(nil)):     true,
	reflect.TypeOf((*mast.JavaEnhancedForStatement)(nil)):           true,
	reflect.TypeOf((*mast.JavaModuleDeclaration)(nil)):              true,
	reflect.TypeOf((*mast.JavaModuleDirective)(nil)):                true,
	reflect.TypeOf((*mast.JavaTypeParameter)(nil)):                  true,
	reflect.TypeOf((*mast.JavaClassDeclaration)(nil)):               true,
	reflect.TypeOf((*mast.JavaInterfaceDeclaration)(nil)):           true,
	reflect.TypeOf((*mast.JavaEnumDeclaration)(nil)):                true,
	reflect.TypeOf((*mast.JavaEnumConstantDeclaration)(nil)):        true,
	reflect.TypeOf((*mast.JavaClassInitializer)(nil)):               true,
	reflect.TypeOf((*mast.JavaFunctionDeclarationFields)(nil)):      true,
	reflect.TypeOf((*mast.JavaCallExpressionFields)(nil)):           true,
	reflect.TypeOf((*mast.JavaAnnotationDeclaration)(nil)):          true,
	reflect.TypeOf((*mast.JavaAnnotationElementDeclaration)(nil)):   true,
	reflect.TypeOf((*mast.JavaMethodReference)(nil)):                true,
	reflect.TypeOf((*mast.JavaClassLiteral)(nil)):                   true,
	reflect.TypeOf((*mast.JavaEntityCreationExpressionFields)(nil)): true,
	reflect.TypeOf((*mast.JavaVariableDeclarationFields)(nil)):      true,
}

// _recordVisitor is a simple visitor which records the types of nodes visited for testing the Walk
// function.
type _recordVisitor struct {
	// Visited keeps track of all MAST nodes the visitor visits.
	Visited map[mast.Node]bool
}

// Pre is a method that simply records the type of the visited node.
func (d *_recordVisitor) Pre(node mast.Node) error {
	// mark the node as visited
	d.Visited[node] = true
	return nil
}

// Post method does nothing and directly return a nil error. This is implemented so that
// _recordVisitor implements the mast.Visitor interface.
func (*_recordVisitor) Post(mast.Node) error { return nil }

func TestWalk(t *testing.T) {
	t.Run("Test walking a translated node", func(t *testing.T) {
		// Walk through all the test files which are supposed to have full coverage on the MAST
		// nodes

		visitedNodeTypes := make(map[reflect.Type]bool)
		err := filepath.WalkDir(_metaTestDataPrefix, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failure accessing a path %q: %v", path, err)
			}
			if d.IsDir() {
				return nil
			}
			tsNode, err := ts.ParseFile(path)
			require.NoError(t, err)
			mastNode, err := Run(tsNode, filepath.Ext(path))
			require.NoError(t, err)
			require.NotNil(t, mastNode)

			// Pass a recorder visitor to record all the visited nodes
			visitor := &_recordVisitor{Visited: make(map[mast.Node]bool)}
			err = mast.Walk(visitor, mastNode)
			require.NoError(t, err)

			// Compare it against our ground truth that is collected by the reflect package
			reflectVisited := make(map[mast.Node]bool)
			err = reflectVisit(reflect.ValueOf(mastNode), func(node mast.Node) {
				reflectVisited[node] = true
			})
			require.NoError(t, err)

			// We use cmp.Diff here since the diffing algorithm in require.Equal is not powerful
			// enough to give clear error messages.
			if diff := cmp.Diff(reflectVisited, visitor.Visited); diff != "" {
				require.FailNow(t, "mismatch (-expected +actual)", diff)
			}

			// Also test mast.Inspect - a functional version of Walk.
			inspectVisited := make(map[mast.Node]bool)
			err = mast.Inspect(mastNode, func(n mast.Node) { inspectVisited[n] = true })
			require.NoError(t, err)
			if diff := cmp.Diff(reflectVisited, inspectVisited); diff != "" {
				require.FailNow(t, "mismatch (-expected +actual)", diff)
			}

			// Now merge the set of visited node types for this test case into the global set of
			// node types visited by all tests to be checked later.
			for node := range visitor.Visited {
				visitedNodeTypes[reflect.TypeOf(node)] = true
			}

			return nil
		})
		require.NoError(t, err)

		// Now we check if all MAST node types have been visited and tested after all tests.
		// We use cmp.Diff here since the diffing algorithm in require.Equal is not powerful enough
		// to give clear error messages.
		if diff := cmp.Diff(_allMASTNodeTypes, visitedNodeTypes); diff != "" {
			require.Fail(t, "mismatch (-expected +actual)", diff)
		}
	})
}
