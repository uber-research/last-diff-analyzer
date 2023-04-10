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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"

	"github.com/stretchr/testify/require"
)

// the Java grammar is listed here: https://github.com/tree-sitter/tree-sitter-java/blob/v0.19.1/src/grammar.json
// and all node types are listed here: https://github.com/tree-sitter/tree-sitter-java/blob/v0.19.1/src/node-types.json
// a node type will appear in the final AST if it satisfies the following conditions:
// (1) it is a named type (i.e., "named": true in the node-types.json file);
// (2) it is _not_ a supertype (i.e., no "supertypes" field in the node-types.json file).
// see https://tree-sitter.github.io/tree-sitter/using-parsers#supertype-nodes for explanations on the supertypes.

// _allJavaTSNodeTypes keeps track of the set of all tree-sitter nodes for verification.
var _allJavaTSNodeTypes = map[string]bool{
	"annotated_type":                      true,
	"annotation":                          true,
	"annotation_argument_list":            true,
	"annotation_type_body":                true,
	"annotation_type_declaration":         true,
	"annotation_type_element_declaration": true,
	"argument_list":                       true,
	"array_access":                        true,
	"array_creation_expression":           true,
	"array_initializer":                   true,
	"array_type":                          true,
	"assert_statement":                    true,
	"assignment_expression":               true,
	"asterisk":                            true,
	"binary_expression":                   true,
	"block":                               true,
	"break_statement":                     true,
	"cast_expression":                     true,
	"catch_clause":                        true,
	"catch_formal_parameter":              true,
	"catch_type":                          true,
	"class_body":                          true,
	"class_declaration":                   true,
	"constant_declaration":                true,
	"constructor_body":                    true,
	"constructor_declaration":             true,
	"continue_statement":                  true,
	"dimensions":                          true,
	"dimensions_expr":                     true,
	"do_statement":                        true,
	"element_value_array_initializer":     true,
	"element_value_pair":                  true,
	"enhanced_for_statement":              true,
	"enum_body":                           true,
	"enum_body_declarations":              true,
	"enum_constant":                       true,
	"enum_declaration":                    true,
	"explicit_constructor_invocation":     true,
	"expression_statement":                true,
	"extends_interfaces":                  true,
	"field_access":                        true,
	"field_declaration":                   true,
	"finally_clause":                      true,
	"floating_point_type":                 true,
	"for_statement":                       true,
	"formal_parameter":                    true,
	"formal_parameters":                   true,
	"generic_type":                        true,
	"if_statement":                        true,
	"import_declaration":                  true,
	"inferred_parameters":                 true,
	"instanceof_expression":               true,
	"integral_type":                       true,
	"interface_body":                      true,
	"interface_declaration":               true,
	"interface_type_list":                 true,
	"labeled_statement":                   true,
	"lambda_expression":                   true,
	"local_variable_declaration":          true,
	"marker_annotation":                   true,
	"method_declaration":                  true,
	"method_invocation":                   true,
	"method_reference":                    true,
	"modifiers":                           true,
	"module_body":                         true,
	"module_declaration":                  true,
	"module_directive":                    true,
	"object_creation_expression":          true,
	"package_declaration":                 true,
	"parenthesized_expression":            true,
	"program":                             true,
	"receiver_parameter":                  true,
	"requires_modifier":                   true,
	"resource":                            true,
	"resource_specification":              true,
	"return_statement":                    true,
	"scoped_identifier":                   true,
	"scoped_type_identifier":              true,
	"spread_parameter":                    true,
	"static_initializer":                  true,
	"super_interfaces":                    true,
	"superclass":                          true,
	"switch_block":                        true,
	"switch_label":                        true,
	"switch_statement":                    true,
	"synchronized_statement":              true,
	"ternary_expression":                  true,
	"throw_statement":                     true,
	"throws":                              true,
	"try_statement":                       true,
	"try_with_resources_statement":        true,
	"type_arguments":                      true,
	"type_bound":                          true,
	"type_parameter":                      true,
	"type_parameters":                     true,
	"unary_expression":                    true,
	"update_expression":                   true,
	"variable_declarator":                 true,
	"while_statement":                     true,
	"wildcard":                            true,
	"binary_integer_literal":              true,
	"boolean_type":                        true,
	"character_literal":                   true,
	"decimal_floating_point_literal":      true,
	"decimal_integer_literal":             true,
	"false":                               true,
	"hex_floating_point_literal":          true,
	"hex_integer_literal":                 true,
	"identifier":                          true,
	"null_literal":                        true,
	"octal_integer_literal":               true,
	"string_literal":                      true,
	"super":                               true,
	"this":                                true,
	"true":                                true,
	"type_identifier":                     true,
	"void_type":                           true,
	// The following TS nodes are _not_ included in this set:
	// ""class_literal"": class_literal nodes are not reproducible in tree-sitter
	// "comment": dropped in our tree-sitter wrapper
}

func TestJavaTranslation(t *testing.T) {
	// unit tests for different translation rules
	testCases := []struct {
		description string
		file        string
		expected    mast.Node
	}{
		{
			// comments containing numbers (e.g., // 0) indicate array
			// indexes and are there for debugging purposes - the path
			// diffing routine used for error reporting returns a path
			// to an error including array indices that are otherwise
			// difficult to track for large arrays

			description: "Test translating declaration nodes",
			file:        _metaTestDataPrefix + "java/declarations.java",
			expected: &mast.Root{
				Declarations: []mast.Declaration{
					// 0
					&mast.PackageDeclaration{
						Annotation: &mast.Annotation{
							Name:      &mast.Identifier{Name: "TestMarker", Kind: mast.Typ},
							Arguments: nil,
						},
						Name: &mast.Identifier{Name: "dummy"},
					},
					// 1
					&mast.PackageDeclaration{
						Annotation: &mast.Annotation{
							Name: &mast.Identifier{Name: "TestSingle", Kind: mast.Typ},
							Arguments: []mast.Expression{
								&mast.KeyValuePair{
									Key:   &mast.Identifier{Name: "value", Kind: mast.Method},
									Value: &mast.BooleanLiteral{Value: true},
								},
							},
						},
						Name: &mast.Identifier{Name: "dummy"},
					},
					// 2
					&mast.PackageDeclaration{
						Annotation: &mast.Annotation{
							Name: &mast.Identifier{Name: "TestMulti", Kind: mast.Typ},
							Arguments: []mast.Expression{
								&mast.KeyValuePair{
									Key:   &mast.Identifier{Name: "a", Kind: mast.Method},
									Value: &mast.IntLiteral{Value: "1"},
								},
								&mast.KeyValuePair{
									Key:   &mast.Identifier{Name: "b", Kind: mast.Method},
									Value: &mast.IntLiteral{Value: "2"},
								},
								&mast.KeyValuePair{
									Key: &mast.Identifier{Name: "c", Kind: mast.Method},
									Value: &mast.LiteralValue{
										Values: []mast.Expression{
											&mast.IntLiteral{Value: "1"},
											&mast.IntLiteral{Value: "2"},
											&mast.IntLiteral{Value: "3"},
											&mast.Annotation{
												Name: &mast.Identifier{Name: "Nested", Kind: mast.Typ},
												Arguments: []mast.Expression{
													&mast.KeyValuePair{
														Key:   &mast.Identifier{Name: "d", Kind: mast.Method},
														Value: &mast.IntLiteral{Value: "1"},
													},
												},
											},
											&mast.Annotation{
												Name:      &mast.Identifier{Name: "NestedMarker", Kind: mast.Typ},
												Arguments: nil,
											},
										},
									},
								},
							},
						},
						Name: &mast.Identifier{Name: "dummy"},
					},
					// 3
					&mast.PackageDeclaration{
						Annotation: nil,
						Name:       &mast.Identifier{Name: "dummy"},
					},
					// 4
					&mast.PackageDeclaration{
						Annotation: nil,
						Name: &mast.AccessPath{
							Operand: &mast.AccessPath{
								Operand:     &mast.Identifier{Name: "a"},
								Annotations: nil,
								Field:       &mast.Identifier{Name: "b"},
							},
							Annotations: nil,
							Field:       &mast.Identifier{Name: "c"},
						},
					},
					// 5
					&mast.ImportDeclaration{
						Alias: nil,
						Package: &mast.AccessPath{
							Operand: &mast.AccessPath{
								Operand: &mast.AccessPath{
									Operand:     &mast.Identifier{Name: "java"},
									Annotations: nil,
									Field:       &mast.Identifier{Name: "util"},
								},
								Annotations: nil,
								Field:       &mast.Identifier{Name: "jar"},
							},
							Annotations: nil,
							Field:       &mast.Identifier{Name: "*"},
						},
					},
					// 6
					&mast.ImportDeclaration{
						Alias:   nil,
						Package: &mast.Identifier{Name: "java"},
					},
					// 7
					&mast.JavaModuleDeclaration{
						Annotations: []*mast.Annotation{
							{
								Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
								Arguments: nil,
							},
							{
								Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
								Arguments: nil,
							},
							{
								Name: &mast.AccessPath{
									Operand:     &mast.Identifier{Name: "pkg", Kind: mast.Typ},
									Annotations: nil,
									Field:       &mast.Identifier{Name: "Test3", Kind: mast.Typ},
								},
							},
						},
						IsOpen: true,
						Name: &mast.AccessPath{
							Operand: &mast.AccessPath{
								Operand:     &mast.Identifier{Name: "test"},
								Annotations: nil,
								Field:       &mast.Identifier{Name: "a"},
							},
							Annotations: nil,
							Field:       &mast.Identifier{Name: "b"},
						},
						Directives: []*mast.JavaModuleDirective{
							{
								Keyword: "requires",
								Exprs: []mast.Expression{
									&mast.JavaLiteralModifier{
										Modifier: "transitive",
									},
									&mast.Identifier{Name: "a"},
								},
							},
							{
								Keyword: "exports",
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "b"},
								},
							},
							{
								Keyword: "opens",
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "c"},
									&mast.Identifier{Name: "d"},
									&mast.Identifier{Name: "f"},
								},
							},
							{
								Keyword: "uses",
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "g"},
								},
							},
							{
								Keyword: "provides",
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "h"},
									&mast.Identifier{Name: "i"},
								},
							},
						},
					},
					// 8
					&mast.JavaModuleDeclaration{
						Annotations: nil,
						IsOpen:      false,
						Name:        &mast.Identifier{Name: "mod"},
						Directives:  nil,
					},
					// 9
					&mast.JavaClassDeclaration{
						Modifiers: []mast.Expression{
							&mast.Annotation{
								Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
								Arguments: nil,
							},
						},
						Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
						TypeParameters: []*mast.JavaTypeParameter{
							{
								Annotations: nil,
								Type:        &mast.Identifier{Name: "A", Kind: mast.Typ},
								Extends:     nil,
							},
							{
								Annotations: []*mast.Annotation{
									{
										Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
										Arguments: nil,
									},
								},
								Type: &mast.Identifier{Name: "B", Kind: mast.Typ},
								Extends: []mast.Expression{
									&mast.Identifier{Name: "C", Kind: mast.Typ},
									&mast.Identifier{Name: "D", Kind: mast.Typ},
								},
							},
						},
						SuperClass: &mast.Identifier{Name: "E", Kind: mast.Typ},
						Interfaces: []mast.Expression{
							&mast.Identifier{Name: "F", Kind: mast.Typ},
							&mast.AccessPath{
								Operand:     &mast.Identifier{Name: "G", Kind: mast.Typ},
								Annotations: nil,
								Field:       &mast.Identifier{Name: "H", Kind: mast.Typ},
							},
						},
						Body: []mast.Declaration{
							// 9-0
							&mast.VariableDeclaration{
								Names: []*mast.Identifier{&mast.Identifier{Name: "a"}},
								Type:  &mast.Identifier{Name: "int"},
								Value: &mast.IntLiteral{Value: "1"},
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.Annotation{
											Name:      &mast.Identifier{Name: "Test3", Kind: mast.Typ},
											Arguments: nil,
										},
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									Dimensions: nil,
								},
							},
							// 9-1
							&mast.VariableDeclaration{
								Names: []*mast.Identifier{&mast.Identifier{Name: "b"}},
								Type:  &mast.Identifier{Name: "int"},
								Value: &mast.IntLiteral{Value: "2"},
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.Annotation{
											Name:      &mast.Identifier{Name: "Test3", Kind: mast.Typ},
											Arguments: nil,
										},
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									Dimensions: nil,
								},
							},
							// 9-2
							&mast.VariableDeclaration{
								Names: []*mast.Identifier{&mast.Identifier{Name: "c"}},
								Type:  &mast.Identifier{Name: "String", Kind: mast.Typ},
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.JavaLiteralModifier{Modifier: "private"},
									},
									Dimensions: []*mast.JavaDimension{
										{
											Length:      nil,
											Annotations: nil,
										},
									},
								},
							},
							// 9-3
							&mast.VariableDeclaration{
								Names: []*mast.Identifier{&mast.Identifier{Name: "d"}},
								Type:  &mast.Identifier{Name: "int"},
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.JavaLiteralModifier{Modifier: "static"},
									},
									Dimensions: nil,
								},
							},
							// 9-4
							&mast.JavaClassInitializer{
								IsStatic: true,
								Block: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function: &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments: []mast.Expression{
													&mast.Identifier{Name: "d"},
												},
												LangFields: nil,
											},
										},
									},
								},
							},
							// 9-5
							&mast.JavaClassInitializer{
								IsStatic: false,
								Block: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
									},
								},
							},
							// 9-6
							&mast.FunctionDeclaration{
								Name: &mast.Identifier{Name: "Test", Kind: mast.Method},
								Parameters: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "int"},
										Name:       &mast.Identifier{Name: "a"},
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Returns: nil,
								Statements: []mast.Statement{
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.Identifier{Name: "super", Kind: mast.Method},
											Arguments: []mast.Expression{
												&mast.IntLiteral{Value: "1"},
											},
											LangFields: &mast.JavaCallExpressionFields{
												TypeArguments: []mast.Expression{
													&mast.Identifier{Name: "A", Kind: mast.Typ},
													&mast.Identifier{Name: "B", Kind: mast.Typ},
												},
											},
										},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.Annotation{
											Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
											Arguments: nil,
										},
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									TypeParameters: []*mast.JavaTypeParameter{
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "A", Kind: mast.Typ},
											Extends:     nil,
										},
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "B", Kind: mast.Typ},
											Extends:     nil,
										},
									},
									Dimensions:  nil,
									Annotations: nil,
									Throws: []mast.Expression{
										&mast.Identifier{Name: "Ex", Kind: mast.Typ},
									},
								},
							},
							// 9-7
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "Test", Kind: mast.Method},
								Parameters: nil,
								Returns:    nil,
								Statements: []mast.Statement{
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand:     &mast.Identifier{Name: "pkg"},
													Annotations: nil,
													Field:       &mast.Identifier{Name: "A"},
												},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "super", Kind: mast.Method},
											},
											Arguments: nil,
											LangFields: &mast.JavaCallExpressionFields{
												TypeArguments: []mast.Expression{
													&mast.Identifier{Name: "A", Kind: mast.Typ},
													&mast.Identifier{Name: "B", Kind: mast.Typ},
												},
											},
										},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers:      nil,
									TypeParameters: nil,
									Dimensions:     nil,
									Annotations:    nil,
									Throws:         nil,
								},
							},
							// 9-8
							&mast.FunctionDeclaration{
								Name: &mast.Identifier{Name: "Test", Kind: mast.Method},
								Parameters: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "String", Kind: mast.Typ},
										Name:       &mast.Identifier{Name: "b"},
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Returns: nil,
								Statements: []mast.Statement{
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function:   &mast.Identifier{Name: "this", Kind: mast.Method},
											Arguments:  nil,
											LangFields: nil,
										},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{

									Modifiers:      nil,
									TypeParameters: nil,
									Dimensions:     nil,
									Annotations:    nil,
									Throws:         nil,
								},
							},
							// 9-9
							&mast.FunctionDeclaration{
								Name: &mast.Identifier{Name: "hello", Kind: mast.Method},
								Parameters: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "String", Kind: mast.Typ},
										Name:       &mast.Identifier{Name: "a"},
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type: &mast.JavaGenericType{
											Name: &mast.Identifier{Name: "Tuple", Kind: mast.Typ},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "A", Kind: mast.Typ},
												&mast.Identifier{Name: "B", Kind: mast.Typ},
											},
										},
										Name: nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Dimensions: nil,
											Modifiers:  nil,
										},
									},
								},
								Statements: []mast.Statement{
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
											Arguments:  nil,
											LangFields: nil,
										},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.Annotation{
											Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
											Arguments: nil,
										},
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									TypeParameters: []*mast.JavaTypeParameter{
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "A", Kind: mast.Typ},
											Extends:     nil,
										},
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "B", Kind: mast.Typ},
											Extends:     nil,
										},
									},
									Dimensions: []*mast.JavaDimension{
										{
											Length:      nil,
											Annotations: nil,
										},
										{
											Length:      nil,
											Annotations: nil,
										},
									},
									Annotations: []*mast.Annotation{
										{
											Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
											Arguments: nil,
										},
										{
											Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
											Arguments: nil,
										},
									},
									Throws: []mast.Expression{
										&mast.Identifier{Name: "C", Kind: mast.Typ},
										&mast.Identifier{Name: "D", Kind: mast.Typ},
									},
								},
							},
							// 9-10
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "hello2", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: nil,
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers:      nil,
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
						},
					},
					// 10
					&mast.JavaInterfaceDeclaration{
						Modifiers: []mast.Expression{
							&mast.Annotation{
								Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
								Arguments: nil,
							},
						},
						Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
						TypeParameters: []*mast.JavaTypeParameter{
							{
								Annotations: nil,
								Type:        &mast.Identifier{Name: "A", Kind: mast.Typ},
								Extends:     nil,
							},
							{
								Annotations: []*mast.Annotation{
									{
										Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
										Arguments: nil,
									},
								},
								Type: &mast.Identifier{Name: "B", Kind: mast.Typ},
								Extends: []mast.Expression{
									&mast.Identifier{Name: "C", Kind: mast.Typ},
									&mast.Identifier{Name: "D", Kind: mast.Typ},
								},
							},
						},
						Extends: []mast.Expression{
							&mast.Identifier{Name: "E", Kind: mast.Typ},
							&mast.AccessPath{
								Operand:     &mast.Identifier{Name: "F", Kind: mast.Typ},
								Annotations: nil,
								Field:       &mast.Identifier{Name: "G", Kind: mast.Typ},
							},
						},
						Body: []mast.Declaration{
							&mast.VariableDeclaration{
								Names:   []*mast.Identifier{&mast.Identifier{Name: "a"}},
								Type:    &mast.Identifier{Name: "int"},
								Value:   &mast.IntLiteral{Value: "1"},
								IsConst: true,
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.JavaLiteralModifier{Modifier: "public"},
										&mast.JavaLiteralModifier{Modifier: "static"},
										&mast.JavaLiteralModifier{Modifier: "final"},
									},
									Dimensions: nil,
								},
							},
							&mast.FunctionDeclaration{
								Name: &mast.Identifier{Name: "hello", Kind: mast.Method},
								Parameters: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "String", Kind: mast.Typ},
										Name:       &mast.Identifier{Name: "a"},
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type: &mast.JavaGenericType{
											Name: &mast.Identifier{Name: "Tuple", Kind: mast.Typ},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "A", Kind: mast.Typ},
												&mast.Identifier{Name: "B", Kind: mast.Typ},
											},
										},
										Name: nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Dimensions: nil,
											Modifiers:  nil,
										},
									},
								},
								Statements: nil,
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.Annotation{
											Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
											Arguments: nil,
										},
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									TypeParameters: []*mast.JavaTypeParameter{
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "A", Kind: mast.Typ},
											Extends:     nil,
										},
										{
											Annotations: nil,
											Type:        &mast.Identifier{Name: "B", Kind: mast.Typ},
											Extends:     nil,
										},
									},
									Dimensions: []*mast.JavaDimension{
										{
											Length:      nil,
											Annotations: nil,
										},
										{
											Length:      nil,
											Annotations: nil,
										},
									},
									Annotations: []*mast.Annotation{
										{
											Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
											Arguments: nil,
										},
										{
											Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
											Arguments: nil,
										},
									},
									Throws: []mast.Expression{
										&mast.Identifier{Name: "C", Kind: mast.Typ},
										&mast.Identifier{Name: "D", Kind: mast.Typ},
									},
								},
							},
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "hello2", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: nil,
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers:      nil,
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
						},
					},
					// 11
					&mast.JavaInterfaceDeclaration{
						Modifiers:      nil,
						Name:           &mast.Identifier{Name: "Test2", Kind: mast.Typ},
						TypeParameters: nil,
						Extends:        nil,
						Body:           nil,
					},
					// 12
					&mast.JavaEnumDeclaration{
						Modifiers:  nil,
						Name:       &mast.Identifier{Name: "Enum"},
						Interfaces: nil,
						Body: []mast.Declaration{
							&mast.JavaEnumConstantDeclaration{
								Modifiers: nil,
								Name:      &mast.Identifier{Name: "ONE"},
								Arguments: nil,
								Body:      nil,
							},
							&mast.JavaEnumConstantDeclaration{
								Modifiers: nil,
								Name:      &mast.Identifier{Name: "TWO"},
								Arguments: nil,
								Body:      nil,
							},
							&mast.JavaEnumConstantDeclaration{
								Modifiers: nil,
								Name:      &mast.Identifier{Name: "THREE"},
								Arguments: nil,
								Body:      nil,
							},
						},
					},
					// 13
					&mast.JavaEnumDeclaration{
						Modifiers: []mast.Expression{
							&mast.Annotation{
								Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
							},
						},
						Name: &mast.Identifier{Name: "T1"},
						Interfaces: []mast.Expression{
							&mast.Identifier{Name: "I1", Kind: mast.Typ},
							&mast.Identifier{Name: "I2", Kind: mast.Typ},
						},
						Body: []mast.Declaration{
							&mast.JavaEnumConstantDeclaration{
								Modifiers: []mast.Expression{
									&mast.Annotation{
										Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
									},
								},
								Name: &mast.Identifier{Name: "V1"},
								Arguments: []mast.Expression{
									&mast.IntLiteral{Value: "1"},
								},
								Body: []mast.Declaration{
									&mast.VariableDeclaration{
										Type:  &mast.Identifier{Name: "int"},
										Names: []*mast.Identifier{&mast.Identifier{Name: "a"}},
										LangFields: &mast.JavaVariableDeclarationFields{
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
									&mast.VariableDeclaration{
										Type:  &mast.Identifier{Name: "int"},
										Names: []*mast.Identifier{&mast.Identifier{Name: "b"}},
										LangFields: &mast.JavaVariableDeclarationFields{
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
									&mast.FunctionDeclaration{
										Name:       &mast.Identifier{Name: "m1", Kind: mast.Method},
										Parameters: nil,
										Returns: []mast.Declaration{
											&mast.ParameterDeclaration{
												IsVariadic: false,
												Type:       &mast.Identifier{Name: "void"},
												Name:       nil,
												LangFields: &mast.JavaParameterDeclarationFields{
													IsReceiver: false,
													Modifiers:  nil,
													Dimensions: nil,
												},
											},
										},
										Statements: []mast.Statement{},
										LangFields: &mast.JavaFunctionDeclarationFields{
											Modifiers:      nil,
											TypeParameters: nil,
											Annotations:    nil,
											Dimensions:     nil,
											Throws:         nil,
										},
									},
								},
							},
							&mast.JavaEnumConstantDeclaration{
								Modifiers: nil,
								Name:      &mast.Identifier{Name: "V2"},
								Arguments: nil,
								Body:      nil,
							},
							&mast.JavaEnumConstantDeclaration{
								Modifiers: nil,
								Name:      &mast.Identifier{Name: "V3"},
								Arguments: nil,
								Body:      nil,
							},
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "m2", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: []mast.Statement{},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "m3", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: []mast.Statement{},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers: []mast.Expression{
										&mast.JavaLiteralModifier{Modifier: "public"},
									},
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
						},
					},
					// 14
					&mast.JavaAnnotationDeclaration{
						Modifiers: []mast.Expression{
							&mast.Annotation{
								Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
								Arguments: nil,
							},
						},
						Name: &mast.Identifier{Name: "Anno"},
						Body: []mast.Declaration{
							&mast.JavaAnnotationElementDeclaration{
								Modifiers: []mast.Expression{
									&mast.Annotation{
										Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
										Arguments: nil,
									},
								},
								Type:       &mast.Identifier{Name: "int"},
								Name:       &mast.Identifier{Name: "A", Kind: mast.Method},
								Dimensions: nil,
								Value:      &mast.IntLiteral{Value: "1"},
							},
							&mast.JavaAnnotationElementDeclaration{
								Modifiers: nil,
								Type:      &mast.Identifier{Name: "int"},
								Name:      &mast.Identifier{Name: "B", Kind: mast.Method},
								Dimensions: []*mast.JavaDimension{
									{
										Length:      nil,
										Annotations: nil,
									},
									{
										Length:      nil,
										Annotations: nil,
									},
								},
								Value: nil,
							},
							&mast.VariableDeclaration{
								Type:  &mast.Identifier{Name: "int"},
								Names: []*mast.Identifier{&mast.Identifier{Name: "C"}},
								Value: &mast.IntLiteral{Value: "1"},
								LangFields: &mast.JavaVariableDeclarationFields{
									Modifiers:  nil,
									Dimensions: nil,
								},
							},
						},
					},
				},
			},
		},
		{
			description: "Test translating expressions",
			file:        _metaTestDataPrefix + "java/expressions.java",
			expected:
			// All expressions should be wrapped with ExpressionStatement nodes.
			&mast.Root{
				Declarations: []mast.Declaration{
					&mast.JavaClassDeclaration{
						Modifiers:      nil,
						Name:           &mast.Identifier{Name: "Root", Kind: mast.Typ},
						TypeParameters: nil,
						SuperClass:     nil,
						Interfaces:     nil,
						Body: []mast.Declaration{
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "root", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: []mast.Statement{
									// 0
									&mast.ExpressionStatement{Expr: &mast.NullLiteral{}},
									// 1
									&mast.ExpressionStatement{Expr: &mast.BooleanLiteral{Value: true}},
									// 2
									&mast.ExpressionStatement{Expr: &mast.BooleanLiteral{Value: false}},
									// 3
									&mast.ExpressionStatement{
										Expr: &mast.StringLiteral{
											IsRaw: false,
											Value: `"test\t"`,
										},
									},
									// 4
									&mast.ExpressionStatement{Expr: &mast.IntLiteral{Value: "123"}},
									// 5
									&mast.ExpressionStatement{Expr: &mast.IntLiteral{Value: "0x5"}},
									// 6
									&mast.ExpressionStatement{Expr: &mast.IntLiteral{Value: "0o5"}},
									// 7
									&mast.ExpressionStatement{Expr: &mast.IntLiteral{Value: "0b11"}},
									// 8
									&mast.ExpressionStatement{Expr: &mast.FloatLiteral{Value: "1.5"}},
									// 9
									&mast.ExpressionStatement{Expr: &mast.FloatLiteral{Value: "0x0.C90FDAP2f"}},
									// 10
									&mast.ExpressionStatement{Expr: &mast.CharacterLiteral{Value: "'a'"}},
									// 11
									&mast.ExpressionStatement{
										Expr: &mast.BinaryExpression{
											Operator: "+",
											Left: &mast.BinaryExpression{
												Operator: "*",
												Left:     &mast.Identifier{Name: "a"},
												Right:    &mast.Identifier{Name: "b"},
											},
											Right: &mast.Identifier{Name: "c"},
										},
									},
									// 12
									&mast.ExpressionStatement{
										Expr: &mast.UnaryExpression{
											Operator: "!",
											Expr:     &mast.Identifier{Name: "a"},
										},
									},
									// 13
									&mast.ExpressionStatement{
										Expr: &mast.UnaryExpression{
											Operator: "!",
											Expr: &mast.ParenthesizedExpression{
												Expr: &mast.Identifier{Name: "a"},
											},
										},
									},
									// 14
									&mast.ExpressionStatement{
										Expr: &mast.IndexExpression{
											Operand: &mast.Identifier{Name: "a"},
											Index:   &mast.Identifier{Name: "i"},
										},
									},
									// 15
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "a"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "b"},
											},
											Annotations: nil,
											Field:       &mast.Identifier{Name: "c"},
										},
									},
									// 16
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.CallExpression{
												Function: &mast.AccessPath{
													Operand:     &mast.Identifier{Name: "a"},
													Annotations: nil,
													Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
												},
												Arguments:  nil,
												LangFields: nil,
											},
											Field: &mast.Identifier{Name: "b"},
										},
									},
									// 17
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
											Field: &mast.Identifier{Name: "b"},
										},
									},
									// 18
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.AccessPath{
												Operand: &mast.CallExpression{
													Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
													Arguments:  nil,
													LangFields: nil,
												},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "a"},
											},
											Field: &mast.Identifier{Name: "b"},
										},
									},
									// 19
									&mast.ExpressionStatement{
										Expr: &mast.JavaTernaryExpression{
											Condition:   &mast.Identifier{Name: "a"},
											Consequence: &mast.Identifier{Name: "b"},
											Alternative: &mast.Identifier{Name: "c"},
										},
									},
									// 20
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.Identifier{Name: "add", Kind: mast.Method},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
												&mast.Identifier{Name: "b"},
											},
											LangFields: nil,
										},
									},
									// 21
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "foo"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
												&mast.Identifier{Name: "b"},
											},
											LangFields: nil,
										},
									},
									// 22
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "this"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
											},
											LangFields: nil,
										},
									},
									// 23
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "super"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
											},
											LangFields: nil,
										},
									},
									// 24
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand: &mast.AccessPath{
														Operand:     &mast.Identifier{Name: "a"},
														Annotations: nil,
														Field:       &mast.Identifier{Name: "b"},
													},
													Annotations: nil,
													Field:       &mast.Identifier{Name: "c"},
												},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
											},
											LangFields: nil,
										},
									},
									// 25
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.Identifier{Name: "void"},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 26
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.Identifier{Name: "int"},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 27
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.Identifier{Name: "float"},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 28
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.JavaGenericType{
													Name: &mast.AccessPath{
														Operand: &mast.AccessPath{
															Operand: &mast.AccessPath{
																Operand:     &mast.Identifier{Name: "a", Kind: mast.Typ},
																Annotations: nil,
																Field:       &mast.Identifier{Name: "b", Kind: mast.Typ},
															},
															Annotations: nil,
															Field:       &mast.Identifier{Name: "c", Kind: mast.Typ},
														},
														Annotations: nil,
														Field:       &mast.Identifier{Name: "Map", Kind: mast.Typ},
													},
													Arguments: []mast.Expression{
														&mast.JavaAnnotatedType{
															Annotations: []*mast.Annotation{
																{
																	Name:      &mast.Identifier{Name: "Even", Kind: mast.Typ},
																	Arguments: nil,
																},
															},
															Type: &mast.Identifier{Name: "int", Kind: mast.Typ},
														},
														&mast.Identifier{Name: "boolean"},
													},
												},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 29
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.JavaGenericType{
													Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
													Arguments: []mast.Expression{
														&mast.JavaAnnotatedType{
															Annotations: []*mast.Annotation{
																{
																	Name:      &mast.Identifier{Name: "NotNull", Kind: mast.Typ},
																	Arguments: nil,
																},
															},
															Type: &mast.JavaWildcard{
																Super:   nil,
																Extends: &mast.Identifier{Name: "A", Kind: mast.Typ},
															},
														},
													},
												},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 30
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.JavaGenericType{
													Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
													Arguments: []mast.Expression{
														&mast.JavaAnnotatedType{
															Annotations: []*mast.Annotation{
																{
																	Name:      &mast.Identifier{Name: "Foo", Kind: mast.Typ},
																	Arguments: nil,
																},
																{
																	Name:      &mast.Identifier{Name: "Bar", Kind: mast.Typ},
																	Arguments: nil,
																},
															},
															Type: &mast.JavaWildcard{
																Super: &mast.JavaAnnotatedType{
																	Annotations: []*mast.Annotation{
																		{
																			Name: &mast.Identifier{Name: "Foo", Kind: mast.Typ},
																		},
																	},
																	Type: &mast.Identifier{Name: "A", Kind: mast.Typ},
																},
																Extends: nil,
															},
														},
													},
												},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 31
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.JavaGenericType{
													Name: &mast.Identifier{Name: "Test", Kind: mast.Typ},
													Arguments: []mast.Expression{
														&mast.JavaWildcard{
															Super:   nil,
															Extends: nil,
														},
													},
												},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 32
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.JavaAnnotatedType{
													Annotations: []*mast.Annotation{
														{
															Name:      &mast.Identifier{Name: "NotNull", Kind: mast.Typ},
															Arguments: nil,
														},
													},
													Type: &mast.JavaArrayType{
														Name: &mast.Identifier{Name: "String", Kind: mast.Typ},
														Dimensions: []*mast.JavaDimension{
															{
																Length:      nil,
																Annotations: nil,
															},
															{
																Annotations: []*mast.Annotation{
																	{
																		Name:      &mast.Identifier{Name: "Foo", Kind: mast.Typ},
																		Arguments: nil,
																	},
																},
															},
															{
																Annotations: []*mast.Annotation{
																	{
																		Name:      &mast.Identifier{Name: "Bar", Kind: mast.Typ},
																		Arguments: nil,
																	},
																	{
																		Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
																		Arguments: nil,
																	},
																},
															},
														},
													},
												},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 33
									&mast.ExpressionStatement{
										Expr: &mast.CastExpression{
											Types: []mast.Expression{
												&mast.Identifier{Name: "T1", Kind: mast.Typ},
												&mast.Identifier{Name: "T2", Kind: mast.Typ},
											},
											Operand: &mast.Identifier{Name: "a"},
										},
									},
									// 34
									&mast.ExpressionStatement{
										Expr: &mast.JavaInstanceOfExpression{
											Operand: &mast.Identifier{Name: "a"},
											Type:    &mast.Identifier{Name: "B", Kind: mast.Typ},
										},
									},
									// 35
									&mast.ExpressionStatement{
										Expr: &mast.UpdateExpression{
											OperatorSide: mast.OperatorAfter,
											Operator:     "++",
											Operand:      &mast.Identifier{Name: "a"},
										},
									},
									// 36
									&mast.ExpressionStatement{
										Expr: &mast.UpdateExpression{
											OperatorSide: mast.OperatorAfter,
											Operator:     "--",
											Operand:      &mast.Identifier{Name: "a"},
										},
									},
									// 37
									&mast.ExpressionStatement{
										Expr: &mast.UpdateExpression{
											OperatorSide: mast.OperatorBefore,
											Operator:     "++",
											Operand:      &mast.Identifier{Name: "a"},
										},
									},
									// 38
									&mast.ExpressionStatement{
										Expr: &mast.UpdateExpression{
											OperatorSide: mast.OperatorBefore,
											Operator:     "--",
											Operand:      &mast.Identifier{Name: "a"},
										},
									},
									// 39
									&mast.ExpressionStatement{
										Expr: &mast.FunctionLiteral{
											Parameters: nil,
											Returns:    nil,
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 40
									&mast.ExpressionStatement{
										Expr: &mast.FunctionLiteral{
											Parameters: nil,
											Returns:    nil,
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 41
									&mast.ExpressionStatement{
										Expr: &mast.FunctionLiteral{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "ClassName", Kind: mast.Typ},
													Name:       &mast.Identifier{Name: "name"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.Annotation{
																Name:      &mast.Identifier{Name: "T1", Kind: mast.Typ},
																Arguments: nil,
															},
															&mast.Annotation{
																Name:      &mast.Identifier{Name: "T2", Kind: mast.Typ},
																Arguments: nil,
															},
														},
														Dimensions: nil,
													},
												},
											},
											Returns: nil,
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.Identifier{Name: "a"},
												},
											},
										},
									},
									// 42
									&mast.ExpressionStatement{
										Expr: &mast.FunctionLiteral{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "ClassName", Kind: mast.Typ},
													Name:       &mast.Identifier{Name: "name"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.Annotation{
																Name:      &mast.Identifier{Name: "T1", Kind: mast.Typ},
																Arguments: nil,
															},
															&mast.Annotation{
																Name:      &mast.Identifier{Name: "T2", Kind: mast.Typ},
																Arguments: nil,
															},
														},
														Dimensions: nil,
													},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "Test1", Kind: mast.Typ},
													Name:       &mast.Identifier{Name: "t1"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers:  nil,
														Dimensions: nil,
													},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "Test2", Kind: mast.Typ},
													Name:       &mast.Identifier{Name: "t2"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.JavaLiteralModifier{Modifier: "final"},
														},
														Dimensions: []*mast.JavaDimension{
															{
																Length:      nil,
																Annotations: nil,
															},
															{
																Annotations: []*mast.Annotation{
																	{
																		Name:      &mast.Identifier{Name: "Anno", Kind: mast.Typ},
																		Arguments: nil,
																	},
																},
															},
														},
													},
												},
												&mast.ParameterDeclaration{
													IsVariadic: true,
													Type:       &mast.Identifier{Name: "Test3", Kind: mast.Typ},
													Name:       &mast.Identifier{Name: "a"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers:  nil,
														Dimensions: nil,
													},
												},
											},
											Returns: nil,
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.Identifier{Name: "a"},
												},
											},
										},
									},
									// 43
									&mast.ExpressionStatement{
										Expr: &mast.FunctionLiteral{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       nil,
													Name:       &mast.Identifier{Name: "a"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers:  nil,
														Dimensions: nil,
													},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       nil,
													Name:       &mast.Identifier{Name: "b"},
													LangFields: &mast.JavaParameterDeclarationFields{
														Modifiers:  nil,
														Dimensions: nil,
													},
												},
											},
											Returns: nil,
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.Identifier{Name: "a"},
												},
											},
										},
									},
									// 44
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   nil,
										},
									},
									// 45
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "final"},
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: true,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   nil,
										},
									},
									// 46
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   &mast.IntLiteral{Value: "2"},
										},
									},
									// 47
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "b"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   nil,
										},
									},
									// 48
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "c"}},
											Type:    &mast.Identifier{Name: "int"},
											Value: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
									},
									// 49
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "String", Kind: mast.Typ},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: nil,
										},
									},
									// 50
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "b"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "String", Kind: mast.Typ},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: nil,
										},
									},
									// 51
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.Annotation{
														Name:      &mast.Identifier{Name: "Test1", Kind: mast.Typ},
														Arguments: nil,
													},
													&mast.Annotation{
														Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
														Arguments: nil,
													},
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "String", Kind: mast.Typ},
											Value:   nil,
										},
									},
									// 52
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.Annotation{
														Name:      &mast.Identifier{Name: "Test1", Kind: mast.Typ},
														Arguments: nil,
													},
													&mast.Annotation{
														Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
														Arguments: nil,
													},
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "b"}},
											Type:    &mast.Identifier{Name: "String", Kind: mast.Typ},
											Value:   nil,
										},
									},
									// 53
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "int"},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: &mast.EntityCreationExpression{
												Object: nil,
												Type:   &mast.Identifier{Name: "int"},
												Value: &mast.LiteralValue{
													Values: []mast.Expression{
														&mast.IntLiteral{Value: "1"},
														&mast.IntLiteral{Value: "2"},
													},
												},
												LangFields: &mast.JavaEntityCreationExpressionFields{
													Dimensions: []*mast.JavaDimension{
														{
															Length:      nil,
															Annotations: nil,
														},
													},
													Body: nil,
												},
											},
										},
									},
									// 54
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "int"},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: &mast.EntityCreationExpression{
												Object: nil,
												Type:   &mast.Identifier{Name: "int"},
												Value:  nil,
												LangFields: &mast.JavaEntityCreationExpressionFields{
													Dimensions: []*mast.JavaDimension{
														{
															Length:      &mast.IntLiteral{Value: "5"},
															Annotations: nil,
														},
													},
													Body: nil,
												},
											},
										},
									},
									// 55
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "int"},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: &mast.EntityCreationExpression{
												Type:  &mast.Identifier{Name: "int"},
												Value: nil,
												LangFields: &mast.JavaEntityCreationExpressionFields{
													Dimensions: []*mast.JavaDimension{
														{
															Length:      &mast.Identifier{Name: "b"},
															Annotations: nil,
														},
														{
															Length:      &mast.Identifier{Name: "c"},
															Annotations: nil,
														},
													},
													Body: nil,
												},
											},
										},
									},
									// 56
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type: &mast.JavaArrayType{
												Name: &mast.Identifier{Name: "int"},
												Dimensions: []*mast.JavaDimension{
													{
														Length:      nil,
														Annotations: nil,
													},
													{
														Length:      nil,
														Annotations: nil,
													},
												},
											},
											Value: &mast.EntityCreationExpression{
												Object: nil,
												Type:   &mast.Identifier{Name: "int"},
												Value: &mast.LiteralValue{
													Values: []mast.Expression{
														&mast.LiteralValue{
															Values: []mast.Expression{
																&mast.IntLiteral{Value: "1"},
																&mast.IntLiteral{Value: "2"},
															},
														},
														&mast.LiteralValue{
															Values: []mast.Expression{
																&mast.IntLiteral{Value: "3"},
																&mast.IntLiteral{Value: "4"},
															},
														},
													},
												},
												LangFields: &mast.JavaEntityCreationExpressionFields{
													Dimensions: []*mast.JavaDimension{
														{
															Length:      nil,
															Annotations: nil,
														},
														{
															Length:      nil,
															Annotations: nil,
														},
													},
													Body: nil,
												},
											},
										},
									},
									// 57
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "String", Kind: mast.Typ},
											Value: &mast.EntityCreationExpression{
												Object: nil,
												Type:   &mast.Identifier{Name: "String", Kind: mast.Method},
												Value: &mast.LiteralValue{
													Values: []mast.Expression{
														&mast.StringLiteral{Value: `"test"`},
													},
												},
											},
										},
									},
									// 58
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.Identifier{Name: "foo", Kind: mast.Method},
											Arguments: []mast.Expression{
												&mast.JavaClassLiteral{
													Type: &mast.JavaArrayType{
														Name: &mast.Identifier{Name: "String", Kind: mast.Typ},
														Dimensions: []*mast.JavaDimension{
															{
																Length:      nil,
																Annotations: nil,
															},
														},
													},
												},
											},
											LangFields: nil,
										},
									},
									// 59
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.Identifier{Name: "foo", Kind: mast.Method},
											Arguments: []mast.Expression{
												&mast.AccessPath{
													Operand:     &mast.Identifier{Name: "String"},
													Annotations: nil,
													Field:       &mast.Identifier{Name: "class"},
												},
											},
											LangFields: nil,
										},
									},
									// 60
									&mast.ExpressionStatement{
										Expr: &mast.EntityCreationExpression{
											Object: nil,
											Type: &mast.JavaGenericType{
												Name: &mast.Identifier{Name: "ConcurrentHashMap", Kind: mast.Method},
												Arguments: []mast.Expression{
													&mast.Identifier{Name: "MethodDeclaration", Kind: mast.Typ},
													&mast.AccessPath{
														Operand: &mast.Identifier{Name: "Annotations", Kind: mast.Typ},
														Annotations: []*mast.Annotation{
															{
																Name:      &mast.Identifier{Name: "NotNull", Kind: mast.Typ},
																Arguments: nil,
															},
														},
														Field: &mast.Identifier{Name: "Measure", Kind: mast.Typ},
													},
												},
											},
											Value: nil,
										},
									},
									// 61
									&mast.ExpressionStatement{
										Expr: &mast.EntityCreationExpression{
											Object: &mast.Identifier{Name: "outer"},
											Type:   &mast.Identifier{Name: "T", Kind: mast.Method},
											Value:  nil,
										},
									},
									// 62
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "T", Kind: mast.Typ},
											Value: &mast.EntityCreationExpression{
												Object: &mast.Identifier{Name: "outer"},
												Type:   &mast.Identifier{Name: "T", Kind: mast.Method},
												Value:  nil,
												LangFields: &mast.JavaEntityCreationExpressionFields{
													Body: []mast.Declaration{
														&mast.FunctionDeclaration{
															Name:       &mast.Identifier{Name: "hello", Kind: mast.Method},
															Parameters: nil,
															Returns: []mast.Declaration{
																&mast.ParameterDeclaration{
																	IsVariadic: false,
																	Type:       &mast.Identifier{Name: "void"},
																	Name:       nil,
																	LangFields: &mast.JavaParameterDeclarationFields{
																		IsReceiver: false,
																		Modifiers:  nil,
																		Dimensions: nil,
																	},
																},
															},
															Statements: []mast.Statement{
																&mast.ReturnStatement{
																	Exprs: nil,
																},
															},
															LangFields: &mast.JavaFunctionDeclarationFields{
																Modifiers: []mast.Expression{
																	&mast.JavaLiteralModifier{Modifier: "public"},
																},
																TypeParameters: nil,
																Dimensions:     nil,
																Annotations:    nil,
																Throws:         nil,
															},
														},
													},
												},
											},
										},
									},
									// 63
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "Clazz"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
											},
											Arguments: nil,
											LangFields: &mast.JavaCallExpressionFields{
												TypeArguments: []mast.Expression{
													&mast.Identifier{Name: "String", Kind: mast.Typ},
												},
											},
										},
									},
									// 64
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "Clazz"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
											},
											Arguments: nil,
											LangFields: &mast.JavaCallExpressionFields{
												TypeArguments: []mast.Expression{
													&mast.Identifier{Name: "String", Kind: mast.Typ},
													&mast.Identifier{Name: "Integer", Kind: mast.Typ},
												},
											},
										},
									},
									// 65
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand:     &mast.Identifier{Name: "Clazz"},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "foo", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
												&mast.Identifier{Name: "b"},
											},
											LangFields: &mast.JavaCallExpressionFields{
												TypeArguments: []mast.Expression{
													&mast.Identifier{Name: "String", Kind: mast.Typ},
													&mast.Identifier{Name: "Integer", Kind: mast.Typ},
												},
											},
										},
									},
									// 66
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.AccessPath{
												Operand: &mast.Identifier{Name: "a"},
												Field:   &mast.Identifier{Name: "this"},
											},
											Field: &mast.Identifier{Name: "bar"},
										},
									},
									// 67
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.Identifier{Name: "super"},
											Field:   &mast.Identifier{Name: "bar"},
										},
									},
									// 68
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.AccessPath{
												Operand: &mast.Identifier{Name: "a"},
												Field:   &mast.Identifier{Name: "super"},
											},
											Field: &mast.Identifier{Name: "bar"},
										},
									},
									// 69
									&mast.ExpressionStatement{
										Expr: &mast.AccessPath{
											Operand: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand: &mast.Identifier{Name: "a"},
													Field:   &mast.Identifier{Name: "b"},
												},
												Field: &mast.Identifier{Name: "super"},
											},
											Field: &mast.Identifier{Name: "bar"},
										},
									},
									// 70
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand: &mast.Identifier{Name: "a"},
													Field:   &mast.Identifier{Name: "this"},
												},
												Field: &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											LangFields: nil,
										},
									},
									// 71
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.Identifier{Name: "super"},
												Field:   &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											LangFields: nil,
										},
									},
									// 72
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand: &mast.Identifier{Name: "a"},
													Field:   &mast.Identifier{Name: "super"},
												},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											LangFields: nil,
										},
									},
									// 73
									&mast.ExpressionStatement{
										Expr: &mast.CallExpression{
											Function: &mast.AccessPath{
												Operand: &mast.AccessPath{
													Operand: &mast.Identifier{Name: "a"},
													Field:   &mast.Identifier{Name: "super"},
												},
												Annotations: nil,
												Field:       &mast.Identifier{Name: "bar", Kind: mast.Method},
											},
											Arguments: []mast.Expression{
												&mast.Identifier{Name: "a"},
												&mast.Identifier{Name: "b"},
											},
											LangFields: nil,
										},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers:      nil,
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
						},
					},
				},
			},
		},
		{
			description: "Test translating statements",
			file:        _metaTestDataPrefix + "java/statements.java",
			expected: &mast.Root{
				Declarations: []mast.Declaration{
					&mast.JavaClassDeclaration{
						Modifiers:      nil,
						Name:           &mast.Identifier{Name: "Root", Kind: mast.Typ},
						TypeParameters: nil,
						SuperClass:     nil,
						Interfaces:     nil,
						Body: []mast.Declaration{
							&mast.FunctionDeclaration{
								Name:       &mast.Identifier{Name: "root", Kind: mast.Method},
								Parameters: nil,
								Returns: []mast.Declaration{
									&mast.ParameterDeclaration{
										IsVariadic: false,
										Type:       &mast.Identifier{Name: "void"},
										Name:       nil,
										LangFields: &mast.JavaParameterDeclarationFields{
											IsReceiver: false,
											Modifiers:  nil,
											Dimensions: nil,
										},
									},
								},
								Statements: []mast.Statement{
									// 0
									&mast.ContinueStatement{
										Label: nil,
									},
									// 1
									&mast.ContinueStatement{
										Label: &mast.Identifier{Name: "here", Kind: mast.Label},
									},
									// 2
									&mast.BreakStatement{
										Label: nil,
									},
									// 3
									&mast.BreakStatement{
										Label: &mast.Identifier{Name: "there", Kind: mast.Label},
									},
									// 4
									&mast.ReturnStatement{
										Exprs: nil,
									},
									// 5
									&mast.ReturnStatement{
										Exprs: []mast.Expression{
											&mast.Identifier{Name: "a"},
										},
									},
									// 6
									&mast.ExpressionStatement{
										Expr: &mast.AssignmentExpression{
											IsShortVarDeclaration: false,
											Left: []mast.Expression{
												&mast.Identifier{Name: "a"},
											},
											Right: []mast.Expression{
												&mast.Identifier{Name: "b"},
											},
										},
									},
									// 7
									&mast.JavaWhileStatement{
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "a"}},
										Body: &mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
									},
									// 8
									&mast.JavaWhileStatement{
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "a"}},
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 9
									&mast.SwitchStatement{
										Initializer: nil,
										Value:       &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "a"}},
										Cases: []*mast.SwitchCase{
											{
												Values: []mast.Expression{
													&mast.StringLiteral{Value: `"1"`},
												},
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "bar", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
											{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "2"},
												},
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "test", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
											{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "3"},
												},
												Statements: nil,
											},
											{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "4"},
												},
												Statements: nil,
											},
											{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "5"},
												},
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
											{
												Values:     nil,
												Statements: nil,
											},
										},
									},
									// 10
									&mast.SwitchStatement{
										Initializer: nil,
										Value: &mast.ParenthesizedExpression{
											Expr: &mast.Identifier{Name: "a"},
										},
										Cases: []*mast.SwitchCase{
											{
												Values: nil,
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
										},
									},
									// 11
									&mast.JavaThrowStatement{
										Expr: &mast.Identifier{Name: "ex"},
									},
									// 12
									&mast.JavaAssertStatement{
										Condition: &mast.BinaryExpression{
											Operator: "==",
											Left:     &mast.Identifier{Name: "a"},
											Right:    &mast.Identifier{Name: "b"},
										},
										ErrorString: nil,
									},
									// 13
									&mast.JavaAssertStatement{
										Condition: &mast.BinaryExpression{
											Operator: "==",
											Left:     &mast.Identifier{Name: "a"},
											Right:    &mast.Identifier{Name: "b"},
										},
										ErrorString: &mast.StringLiteral{
											IsRaw: false,
											Value: `"ERROR"`,
										},
									},
									// 14
									&mast.JavaSynchronizedStatement{
										Expr: &mast.ParenthesizedExpression{
											Expr: &mast.Identifier{Name: "a"},
										},
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.UpdateExpression{
														OperatorSide: mast.OperatorAfter,
														Operator:     "++",
														Operand:      &mast.Identifier{Name: "a"},
													},
												},
											},
										},
									},
									// 15
									&mast.IfStatement{
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "a"}},
										Consequence: &mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
										Alternative: nil,
									},
									// 16
									&mast.IfStatement{
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "a"}},
										Consequence: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "t1", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										Alternative: &mast.IfStatement{
											Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "b"}},
											Consequence: &mast.Block{
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "t2", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
											Alternative: &mast.Block{
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "t3", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
										},
									},
									// 17
									&mast.LabelStatement{
										Label: &mast.Identifier{Name: "hello", Kind: mast.Label},
									},
									// 18
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "a"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   nil,
										},
									},
									// 19
									&mast.DeclarationStatement{
										Decl: &mast.VariableDeclaration{
											LangFields: &mast.JavaVariableDeclarationFields{
												Modifiers: []mast.Expression{
													&mast.JavaLiteralModifier{Modifier: "private"},
												},
											},
											IsConst: false,
											Names:   []*mast.Identifier{{Name: "b"}},
											Type:    &mast.Identifier{Name: "int"},
											Value:   nil,
										},
									},
									// 20
									&mast.JavaDoStatement{
										Body: &mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "bar"}},
									},
									// 21
									&mast.JavaDoStatement{
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										Condition: &mast.ParenthesizedExpression{Expr: &mast.Identifier{Name: "bar"}},
									},
									// 22
									&mast.JavaWhileStatement{
										Condition: &mast.ParenthesizedExpression{
											Expr: &mast.CallExpression{
												Function:   &mast.Identifier{Name: "hasNext", Kind: mast.Method},
												Arguments:  nil,
												LangFields: nil,
											},
										},
										Body: nil,
									},
									// 23
									&mast.JavaTryStatement{
										Resources: nil,
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "t1", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										CatchClauses: []*mast.JavaCatchClause{
											{
												Parameter: &mast.JavaCatchFormalParameter{
													Modifiers: []mast.Expression{
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
															Arguments: nil,
														},
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
															Arguments: nil,
														},
													},
													Types: []mast.Expression{
														&mast.Identifier{Name: "A", Kind: mast.Typ},
														&mast.Identifier{Name: "B", Kind: mast.Typ},
													},
													Name: &mast.Identifier{Name: "ex"},
													Dimensions: []*mast.JavaDimension{
														{
															Length:      nil,
															Annotations: nil,
														},
														{
															Length:      nil,
															Annotations: nil,
														},
													},
												},
												Body: &mast.Block{
													Statements: []mast.Statement{
														&mast.ExpressionStatement{
															Expr: &mast.CallExpression{
																Function:   &mast.Identifier{Name: "t2", Kind: mast.Method},
																Arguments:  nil,
																LangFields: nil,
															},
														},
													},
												},
											},
											{
												Parameter: &mast.JavaCatchFormalParameter{
													Modifiers: []mast.Expression{
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test3", Kind: mast.Typ},
															Arguments: nil,
														},
													},
													Types: []mast.Expression{
														&mast.Identifier{Name: "C", Kind: mast.Typ},
													},
													Name:       &mast.Identifier{Name: "ex2"},
													Dimensions: nil,
												},
												Body: &mast.Block{
													Statements: []mast.Statement{
														&mast.ExpressionStatement{
															Expr: &mast.CallExpression{
																Function:   &mast.Identifier{Name: "t3", Kind: mast.Method},
																Arguments:  nil,
																LangFields: nil,
															},
														},
													},
												},
											},
										},
										Finally: &mast.JavaFinallyClause{
											Body: &mast.Block{
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "t4", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
										},
									},
									// 24
									&mast.JavaTryStatement{
										Resources: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.Identifier{Name: "file"},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													LangFields: &mast.JavaVariableDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.JavaLiteralModifier{Modifier: "private"},
														},
													},
													IsConst: false,
													Names:   []*mast.Identifier{{Name: "scanner"}},
													Type:    &mast.Identifier{Name: "Scanner", Kind: mast.Typ},
													Value: &mast.EntityCreationExpression{
														Type:  &mast.Identifier{Name: "Scanner", Kind: mast.Method},
														Value: nil,
													},
												},
											},
										},
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "t1", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										CatchClauses: []*mast.JavaCatchClause{
											{
												Parameter: &mast.JavaCatchFormalParameter{
													Modifiers: []mast.Expression{
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
															Arguments: nil,
														},
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test2", Kind: mast.Typ},
															Arguments: nil,
														},
													},
													Types: []mast.Expression{
														&mast.Identifier{Name: "A", Kind: mast.Typ},
														&mast.Identifier{Name: "B", Kind: mast.Typ},
													},
													Name: &mast.Identifier{Name: "ex"},
													Dimensions: []*mast.JavaDimension{
														{
															Length:      nil,
															Annotations: nil,
														},
														{
															Length:      nil,
															Annotations: nil,
														},
													},
												},
												Body: &mast.Block{
													Statements: []mast.Statement{
														&mast.ExpressionStatement{
															Expr: &mast.CallExpression{
																Function:   &mast.Identifier{Name: "t2", Kind: mast.Method},
																Arguments:  nil,
																LangFields: nil,
															},
														},
													},
												},
											},
											{
												Parameter: &mast.JavaCatchFormalParameter{
													Modifiers: []mast.Expression{
														&mast.Annotation{
															Name:      &mast.Identifier{Name: "Test3", Kind: mast.Typ},
															Arguments: nil,
														},
													},
													Types: []mast.Expression{
														&mast.Identifier{Name: "C", Kind: mast.Typ},
													},
													Name:       &mast.Identifier{Name: "ex2"},
													Dimensions: nil,
												},
												Body: &mast.Block{
													Statements: []mast.Statement{
														&mast.ExpressionStatement{
															Expr: &mast.CallExpression{
																Function:   &mast.Identifier{Name: "t3", Kind: mast.Method},
																Arguments:  nil,
																LangFields: nil,
															},
														},
													},
												},
											},
										},
										Finally: &mast.JavaFinallyClause{
											Body: &mast.Block{
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "t4", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
										},
									},
									// 25
									&mast.JavaTryStatement{
										Resources: nil,
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "t1", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										CatchClauses: []*mast.JavaCatchClause{
											{
												Parameter: &mast.JavaCatchFormalParameter{
													Modifiers: nil,
													Types: []mast.Expression{
														&mast.Identifier{Name: "A", Kind: mast.Typ},
													},
													Name:       &mast.Identifier{Name: "a"},
													Dimensions: nil,
												},
												Body: &mast.Block{
													Statements: []mast.Statement{
														&mast.ExpressionStatement{
															Expr: &mast.CallExpression{
																Function:   &mast.Identifier{Name: "t2", Kind: mast.Method},
																Arguments:  nil,
																LangFields: nil,
															},
														},
													},
												},
											},
										},
										Finally: nil,
									},
									// 26
									&mast.JavaTryStatement{
										Resources: nil,
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "t1", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
										CatchClauses: nil,
										Finally: &mast.JavaFinallyClause{
											Body: &mast.Block{
												Statements: []mast.Statement{
													&mast.ExpressionStatement{
														Expr: &mast.CallExpression{
															Function:   &mast.Identifier{Name: "t2", Kind: mast.Method},
															Arguments:  nil,
															LangFields: nil,
														},
													},
												},
											},
										},
									},
									// 27
									&mast.ForStatement{
										Initializers: nil,
										Condition:    nil,
										Updates:      nil,
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 28
									&mast.ForStatement{
										Initializers: nil,
										Condition:    nil,
										Updates:      nil,
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 29
									&mast.ForStatement{
										Initializers: []mast.Statement{
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													LangFields: &mast.JavaVariableDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.JavaLiteralModifier{Modifier: "private"},
														},
													},
													IsConst: false,
													Type:    &mast.Identifier{Name: "int"},
													Names:   []*mast.Identifier{{Name: "i"}},
													Value:   &mast.IntLiteral{Value: "1"},
												},
											},
										},
										Condition: nil,
										Updates:   nil,
										Body:      nil,
									},
									// 30
									&mast.ForStatement{
										Initializers: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.AssignmentExpression{
													IsShortVarDeclaration: false,
													Left: []mast.Expression{
														&mast.Identifier{Name: "i"},
													},
													Right: []mast.Expression{
														&mast.IntLiteral{Value: "2"},
													},
												},
											},
										},
										Condition: &mast.BinaryExpression{
											Left:     &mast.Identifier{Name: "i"},
											Operator: "<",
											Right:    &mast.IntLiteral{Value: "10"},
										},
										Updates: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.UpdateExpression{
													OperatorSide: mast.OperatorAfter,
													Operator:     "++",
													Operand:      &mast.Identifier{Name: "i"},
												},
											},
										},
										Body: nil,
									},
									// 31
									&mast.ForStatement{
										Initializers: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.AssignmentExpression{
													IsShortVarDeclaration: false,
													Left: []mast.Expression{
														&mast.Identifier{Name: "i"},
													},
													Right: []mast.Expression{
														&mast.IntLiteral{Value: "2"},
													},
												},
											},
											&mast.ExpressionStatement{
												Expr: &mast.AssignmentExpression{
													IsShortVarDeclaration: false,
													Left: []mast.Expression{
														&mast.Identifier{Name: "k"},
													},
													Right: []mast.Expression{
														&mast.IntLiteral{Value: "3"},
													},
												},
											},
										},
										Condition: &mast.BinaryExpression{
											Left:     &mast.Identifier{Name: "i"},
											Operator: "<",
											Right:    &mast.IntLiteral{Value: "10"},
										},
										Updates: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.UpdateExpression{
													OperatorSide: mast.OperatorAfter,
													Operator:     "++",
													Operand:      &mast.Identifier{Name: "i"},
												},
											},
											&mast.ExpressionStatement{
												Expr: &mast.UpdateExpression{
													OperatorSide: mast.OperatorAfter,
													Operator:     "--",
													Operand:      &mast.Identifier{Name: "k"},
												},
											},
										},
										Body: nil,
									},
									// 32
									&mast.ForStatement{
										Initializers: []mast.Statement{
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													LangFields: &mast.JavaVariableDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.JavaLiteralModifier{Modifier: "private"},
														},
													},
													IsConst: false,
													Type:    &mast.Identifier{Name: "int"},
													Names:   []*mast.Identifier{{Name: "i"}},
													Value:   nil,
												},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													LangFields: &mast.JavaVariableDeclarationFields{
														Modifiers: []mast.Expression{
															&mast.JavaLiteralModifier{Modifier: "private"},
														},
													},
													IsConst: false,
													Type:    &mast.Identifier{Name: "int"},
													Names:   []*mast.Identifier{{Name: "k"}},
													Value:   nil,
												},
											},
										},
										Condition: &mast.BinaryExpression{
											Left:     &mast.Identifier{Name: "i"},
											Operator: "<",
											Right:    &mast.IntLiteral{Value: "10"},
										},
										Updates: nil,
										Body:    nil,
									},
									// 33
									&mast.ForStatement{
										Initializers: nil,
										Condition: &mast.BinaryExpression{
											Left:     &mast.Identifier{Name: "i"},
											Operator: "<",
											Right:    &mast.IntLiteral{Value: "10"},
										},
										Updates: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.UpdateExpression{
													OperatorSide: mast.OperatorAfter,
													Operator:     "++",
													Operand:      &mast.Identifier{Name: "i"},
												},
											},
										},
										Body: nil,
									},
									// 34
									&mast.JavaEnhancedForStatement{
										Modifiers: []mast.Expression{
											&mast.JavaLiteralModifier{Modifier: "final"},
											&mast.Annotation{
												Name:      &mast.Identifier{Name: "Test", Kind: mast.Typ},
												Arguments: nil,
											},
										},
										Type: &mast.Identifier{Name: "String", Kind: mast.Typ},
										Name: &mast.Identifier{Name: "a"},
										Dimensions: []*mast.JavaDimension{
											{
												Length:      nil,
												Annotations: nil,
											},
											{
												Length:      nil,
												Annotations: nil,
											},
										},
										Iterable: &mast.AccessPath{
											Operand: &mast.Identifier{Name: "pkg"},
											Field:   &mast.Identifier{Name: "d"},
										},
										Body: &mast.Block{
											Statements: []mast.Statement{
												&mast.ExpressionStatement{
													Expr: &mast.CallExpression{
														Function:   &mast.Identifier{Name: "foo", Kind: mast.Method},
														Arguments:  nil,
														LangFields: nil,
													},
												},
											},
										},
									},
									// 35
									&mast.JavaEnhancedForStatement{
										Modifiers:  nil,
										Type:       &mast.Identifier{Name: "int"},
										Name:       &mast.Identifier{Name: "a"},
										Dimensions: nil,
										Iterable:   &mast.Identifier{Name: "b"},
										Body:       nil,
									},
									// 36
									&mast.ExpressionStatement{
										Expr: &mast.JavaMethodReference{
											Class: &mast.Identifier{Name: "super"},
											TypeArguments: []mast.Expression{
												&mast.Identifier{Name: "A", Kind: mast.Typ},
												&mast.Identifier{Name: "B", Kind: mast.Typ},
											},
											Method: &mast.Identifier{Name: "someMethod", Kind: mast.Method},
										},
									},
									// 37
									&mast.ExpressionStatement{
										Expr: &mast.JavaMethodReference{
											Class:         &mast.Identifier{Name: "SomeClass"},
											TypeArguments: nil,
											Method:        &mast.Identifier{Name: "new", Kind: mast.Method},
										},
									},
									// 38
									&mast.LabelStatement{
										Label: &mast.Identifier{Name: "empty_label", Kind: mast.Label},
									},
								},
								LangFields: &mast.JavaFunctionDeclarationFields{
									Modifiers:      nil,
									TypeParameters: nil,
									Annotations:    nil,
									Dimensions:     nil,
									Throws:         nil,
								},
							},
						},
					},
				},
			},
		},
	}

	// Keep track of all TS node types that we have visited
	visitedNodess := make(map[string]bool)

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			node, err := ts.ParseFile(tc.file)
			require.NoError(t, err)

			// Put all visited TS node types in to the set and check it later, this must happen
			// before translation.Run since the translator might change the node type strings to
			// share implementation logic among different languages.
			err = reflectVisit(reflect.ValueOf(node), func(node *ts.Node) {
				visitedNodess[node.Type] = true
			})
			require.NoError(t, err)

			actual, err := Run(node, ts.JavaExt)
			require.NoError(t, err)

			// We use cmp.Diff here since the diffing algorithm in require.Equal is not powerful
			// enough to give clear error messages.
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				require.FailNow(t, "mismatch (-expected +actual)", diff)
			}
		})
	}

	// Make sure all nodes have been visited and tested
	for k := range _allJavaTSNodeTypes {
		exists := visitedNodess[k]
		require.True(t, exists, "TS node %s not tested", k)
	}
}
