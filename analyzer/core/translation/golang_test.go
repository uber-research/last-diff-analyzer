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

// The Go grammar is listed here: https://github.com/tree-sitter/tree-sitter-go/blob/v0.19.1/src/grammar.json
// and all node types are listed here: https://github.com/tree-sitter/tree-sitter-go/blob/v0.19.1/src/node-types.json
// a node type will appear in the final AST if it satisfies the following conditions:
// (1) it is a named type (i.e., "named": true in the node-types.json file);
// (2) it is _not_ a supertype (i.e., no "supertypes" field in the node-types.json file).
// see https://tree-sitter.github.io/tree-sitter/using-parsers#supertype-nodes for explanations on the supertypes.

// _allGoTSNodeTypes keeps track of the set of all tree-sitter nodes for verification.
var _allGoTSNodeTypes = map[string]bool{
	"argument_list":                  true,
	"array_type":                     true,
	"assignment_statement":           true,
	"binary_expression":              true,
	"block":                          true,
	"break_statement":                true,
	"call_expression":                true,
	"channel_type":                   true,
	"communication_case":             true,
	"composite_literal":              true,
	"const_declaration":              true,
	"const_spec":                     true,
	"continue_statement":             true,
	"dec_statement":                  true,
	"default_case":                   true,
	"defer_statement":                true,
	"dot":                            true,
	"element":                        true,
	"expression_case":                true,
	"expression_list":                true,
	"expression_switch_statement":    true,
	"fallthrough_statement":          true,
	"field_declaration":              true,
	"field_declaration_list":         true,
	"for_clause":                     true,
	"for_statement":                  true,
	"func_literal":                   true,
	"function_declaration":           true,
	"function_type":                  true,
	"go_statement":                   true,
	"goto_statement":                 true,
	"if_statement":                   true,
	"implicit_length_array_type":     true,
	"import_declaration":             true,
	"import_spec":                    true,
	"import_spec_list":               true,
	"inc_statement":                  true,
	"index_expression":               true,
	"interface_type":                 true,
	"interpreted_string_literal":     true,
	"keyed_element":                  true,
	"labeled_statement":              true,
	"literal_value":                  true,
	"map_type":                       true,
	"method_declaration":             true,
	"method_spec":                    true,
	"method_spec_list":               true,
	"package_clause":                 true,
	"parameter_declaration":          true,
	"parameter_list":                 true,
	"parenthesized_expression":       true,
	"parenthesized_type":             true,
	"pointer_type":                   true,
	"qualified_type":                 true,
	"range_clause":                   true,
	"receive_statement":              true,
	"return_statement":               true,
	"select_statement":               true,
	"selector_expression":            true,
	"send_statement":                 true,
	"short_var_declaration":          true,
	"slice_expression":               true,
	"slice_type":                     true,
	"source_file":                    true,
	"struct_type":                    true,
	"type_alias":                     true,
	"type_assertion_expression":      true,
	"type_case":                      true,
	"type_declaration":               true,
	"type_spec":                      true,
	"type_switch_statement":          true,
	"unary_expression":               true,
	"var_declaration":                true,
	"var_spec":                       true,
	"variadic_argument":              true,
	"variadic_parameter_declaration": true,
	"blank_identifier":               true,
	"false":                          true,
	"field_identifier":               true,
	"float_literal":                  true,
	"identifier":                     true,
	"imaginary_literal":              true,
	"int_literal":                    true,
	"label_name":                     true,
	"nil":                            true,
	"package_identifier":             true,
	"raw_string_literal":             true,
	"rune_literal":                   true,
	"true":                           true,
	"type_identifier":                true,
	"type_conversion_expression":     true,
	// The following TS nodes are _not_ included in this set:
	// "escape_sequence": dropped in tree-sitter wrapper
	// "empty_statement": dropped in tree-sitter wrapper
	// "comment": dropped in our tree-sitter wrapper
}

func TestGoTranslation(t *testing.T) {
	// unit tests for different translation rules
	testCases := []struct {
		description string
		file        string
		expected    mast.Node
	}{
		{
			description: "Test translating declarations",
			file:        _metaTestDataPrefix + "go/declarations.go",
			expected: &mast.Root{
				Declarations: []mast.Declaration{
					&mast.PackageDeclaration{
						Name: &mast.Identifier{Name: "rename"},
					},
					&mast.ImportDeclaration{
						Alias:   &mast.Identifier{Name: "."},
						Package: &mast.StringLiteral{Value: `"example/package1"`},
					},
					&mast.ImportDeclaration{
						Alias:   &mast.Identifier{Name: "t"},
						Package: &mast.StringLiteral{Value: `"package2"`},
					},
					&mast.ImportDeclaration{
						Alias:   &mast.Identifier{Name: "_"},
						Package: &mast.StringLiteral{Value: `"package3"`},
					},
					&mast.ImportDeclaration{
						Alias:   nil,
						Package: &mast.StringLiteral{Value: `"package4"`},
					},
					&mast.ImportDeclaration{
						Alias:   nil,
						Package: &mast.StringLiteral{Value: `"singlepackage"`},
					},
					&mast.FunctionDeclaration{
						Name: &mast.Identifier{Name: "test", Kind: mast.Method},
						Parameters: []mast.Declaration{
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "int"},
								Name:       &mast.Identifier{Name: "a"},
							},
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "int"},
								Name:       &mast.Identifier{Name: "b"},
							},
						},
						Returns: []mast.Declaration{
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "string"},
								Name:       &mast.Identifier{Name: "c"},
							},
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "string"},
								Name:       &mast.Identifier{Name: "d"},
							},
						},
						LangFields: &mast.GoFunctionDeclarationFields{
							Receiver: &mast.ParameterDeclaration{
								IsVariadic: false,
								Type: &mast.GoPointerType{
									Type: &mast.Identifier{Name: "A"},
								},
								Name: &mast.Identifier{Name: "a"},
							},
						},
						Statements: []mast.Statement{
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
									Arguments: nil,
								},
							},
						},
					},
					&mast.FunctionDeclaration{
						Name: &mast.Identifier{Name: "test", Kind: mast.Function},
						Parameters: []mast.Declaration{
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "int"},
								Name:       &mast.Identifier{Name: "a"},
							},
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "int"},
								Name:       &mast.Identifier{Name: "b"},
							},
						},
						Returns: []mast.Declaration{
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "string"},
								Name:       &mast.Identifier{Name: "c"},
							},
							&mast.ParameterDeclaration{
								IsVariadic: false,
								Type:       &mast.Identifier{Name: "string"},
								Name:       &mast.Identifier{Name: "d"},
							},
						},
						LangFields: &mast.GoFunctionDeclarationFields{},
						Statements: []mast.Statement{
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
									Arguments: nil,
								},
							},
						},
					},
					&mast.GoTypeDeclaration{
						IsAlias: false,
						Name:    &mast.Identifier{Name: "Test"},
						Type: &mast.GoInterfaceType{
							Declarations: []mast.Declaration{
								&mast.FieldDeclaration{
									Name:       nil,
									Type:       &mast.Identifier{Name: "Embedded"},
									LangFields: &mast.GoFieldDeclarationFields{},
								},
								&mast.FunctionDeclaration{
									Name: &mast.Identifier{Name: "hello", Kind: mast.Function},
									Parameters: []mast.Declaration{
										&mast.ParameterDeclaration{
											IsVariadic: false,
											Type:       &mast.Identifier{Name: "int"},
											Name:       &mast.Identifier{Name: "a"},
										},
									},
									Returns: []mast.Declaration{
										&mast.ParameterDeclaration{
											IsVariadic: false,
											Type:       &mast.Identifier{Name: "int"},
											Name:       &mast.Identifier{Name: "b"},
										},
									},
									LangFields: &mast.GoFunctionDeclarationFields{},
									Statements: nil,
								},
							},
						},
					},
					&mast.GoTypeDeclaration{
						IsAlias: true,
						Name:    &mast.Identifier{Name: "A"},
						Type:    &mast.Identifier{Name: "B"},
					},
					&mast.GoTypeDeclaration{
						IsAlias: false,
						Name:    &mast.Identifier{Name: "C"},
						Type:    &mast.Identifier{Name: "D"},
					},
				},
			},
		},
		{
			description: "Test translating expressions",
			file:        _metaTestDataPrefix + "go/expressions.go",
			expected: &mast.Root{
				Declarations: []mast.Declaration{
					&mast.FunctionDeclaration{
						Name:       &mast.Identifier{Name: "root", Kind: mast.Function},
						Parameters: nil,
						Returns:    nil,
						LangFields: &mast.GoFunctionDeclarationFields{},
						Statements: []mast.Statement{
							&mast.ExpressionStatement{Expr: &mast.NullLiteral{}},
							&mast.ExpressionStatement{Expr: &mast.BooleanLiteral{Value: true}},
							&mast.ExpressionStatement{Expr: &mast.BooleanLiteral{Value: false}},
							&mast.ExpressionStatement{
								Expr: &mast.StringLiteral{
									IsRaw: false,
									Value: `"test\t"`,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.StringLiteral{
									IsRaw: true,
									Value: "`test`",
								},
							},
							&mast.ExpressionStatement{Expr: &mast.IntLiteral{Value: "123"}},
							&mast.ExpressionStatement{Expr: &mast.FloatLiteral{Value: "1.5"}},
							&mast.ExpressionStatement{Expr: &mast.GoImaginaryLiteral{Value: "5i"}},
							&mast.ExpressionStatement{Expr: &mast.CharacterLiteral{Value: "'a'"}},
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
							&mast.ExpressionStatement{
								Expr: &mast.UnaryExpression{
									Operator: "!",
									Expr:     &mast.Identifier{Name: "a"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.UnaryExpression{
									Operator: "!",
									Expr: &mast.ParenthesizedExpression{
										Expr: &mast.Identifier{Name: "a"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.IndexExpression{
									Operand: &mast.Identifier{Name: "a"},
									Index:   &mast.Identifier{Name: "i"},
								},
							},
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
							&mast.ExpressionStatement{
								Expr: &mast.AccessPath{
									Operand: &mast.CallExpression{
										Function: &mast.AccessPath{
											Operand:     &mast.Identifier{Name: "a"},
											Annotations: nil,
											Field:       &mast.Identifier{Name: "foo"},
										},
										Arguments: nil,
									},
									Field: &mast.Identifier{Name: "b"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AccessPath{
									Operand: &mast.CallExpression{
										Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
										Arguments: nil,
									},
									Field: &mast.Identifier{Name: "b"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AccessPath{
									Operand: &mast.AccessPath{
										Operand: &mast.CallExpression{
											Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
											Arguments: nil,
										},
										Annotations: nil,
										Field:       &mast.Identifier{Name: "a"},
									},
									Field: &mast.Identifier{Name: "b"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.GoSliceExpression{
									Operand: &mast.Identifier{Name: "a"},
									Start:   &mast.Identifier{Name: "i"},
									End: &mast.BinaryExpression{
										Operator: "+",
										Left:     &mast.Identifier{Name: "i"},
										Right:    &mast.IntLiteral{Value: "1"},
									},
									Capacity: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.GoSliceExpression{
									Operand: &mast.Identifier{Name: "a"},
									Start:   &mast.Identifier{Name: "i"},
									End: &mast.BinaryExpression{
										Operator: "+",
										Left:     &mast.Identifier{Name: "i"},
										Right:    &mast.IntLiteral{Value: "1"},
									},
									Capacity: &mast.IntLiteral{Value: "10"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.GoSliceExpression{
									Operand:  &mast.Identifier{Name: "a"},
									Start:    nil,
									End:      nil,
									Capacity: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.GoSliceExpression{
									Operand:  &mast.Identifier{Name: "a"},
									Start:    &mast.Identifier{Name: "i"},
									End:      nil,
									Capacity: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
									Arguments: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "add", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.Identifier{Name: "a"},
										&mast.Identifier{Name: "b"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.AccessPath{
										Operand:     &mast.Identifier{Name: "foo"},
										Annotations: nil,
										Field:       &mast.Identifier{Name: "bar"},
									},
									Arguments: []mast.Expression{
										&mast.Identifier{Name: "a"},
										&mast.Identifier{Name: "b"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "add", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.Identifier{Name: "a"},
										&mast.GoEllipsisExpression{
											Expr: &mast.Identifier{Name: "b"},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.Identifier{Name: "int"},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoArrayType{
											Length:  &mast.IntLiteral{Value: "5"},
											Element: &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoArrayType{
											Length:  nil,
											Element: &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoParenthesizedType{
											Type: &mast.GoMapType{
												Key:   &mast.Identifier{Name: "string"},
												Value: &mast.Identifier{Name: "bool"},
											},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoPointerType{
											Type: &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.AccessPath{
											Operand:     &mast.Identifier{Name: "pkg"},
											Annotations: nil,
											Field:       &mast.Identifier{Name: "Test"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoChannelType{
											Direction: mast.SendAndReceive,
											Type:      &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoChannelType{
											Direction: mast.ReceiveOnly,
											Type:      &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoChannelType{
											Direction: mast.SendOnly,
											Type:      &mast.Identifier{Name: "int"},
										},
										&mast.IntLiteral{Value: "10"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: nil,
											Return:     nil,
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "A"},
													Name:       nil,
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "B"},
													Name:       nil,
												},
											},
											Return: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "C"},
													Name:       nil,
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "A"},
													Name:       &mast.Identifier{Name: "a"},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "B"},
													Name:       &mast.Identifier{Name: "b"},
												},
											},
											Return: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "C"},
													Name:       &mast.Identifier{Name: "c"},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "D"},
													Name:       &mast.Identifier{Name: "d"},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: nil,
											Return: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "C"},
													Name:       nil,
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "D"},
													Name:       nil,
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "A"},
													Name:       nil,
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "B"},
													Name:       nil,
												},
											},
											Return: nil,
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: true,
													Type:       &mast.Identifier{Name: "A"},
													Name:       &mast.Identifier{Name: "a"},
												},
											},
											Return: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: true,
													Type:       &mast.Identifier{Name: "B"},
													Name:       &mast.Identifier{Name: "b"},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "make", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.GoFunctionType{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "A"},
													Name:       &mast.Identifier{Name: "a"},
												},
												&mast.ParameterDeclaration{
													IsVariadic: true,
													Type:       &mast.Identifier{Name: "B"},
													Name:       nil,
												},
											},
											Return: nil,
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.GoTypeAssertionExpression{
									Operand: &mast.Identifier{Name: "a"},
									Type:    &mast.Identifier{Name: "T"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CastExpression{
									Types: []mast.Expression{
										&mast.GoArrayType{
											Length:  nil,
											Element: &mast.Identifier{Name: "byte"},
										},
									},
									Operand: &mast.Identifier{Name: "a"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.UpdateExpression{
									OperatorSide: mast.OperatorAfter,
									Operator:     "++",
									Operand:      &mast.Identifier{Name: "a"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.UpdateExpression{
									OperatorSide: mast.OperatorAfter,
									Operator:     "--",
									Operand:      &mast.Identifier{Name: "a"},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
									},
									Right: []mast.Expression{
										&mast.EntityCreationExpression{
											Object: nil,
											Type: &mast.GoArrayType{
												Length:  nil,
												Element: &mast.Identifier{Name: "int"},
											},
											Value: &mast.LiteralValue{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
													&mast.IntLiteral{Value: "2"},
													&mast.IntLiteral{Value: "3"},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
									},
									Right: []mast.Expression{
										&mast.EntityCreationExpression{
											Object: nil,
											Type: &mast.GoArrayType{
												Length:  &mast.StringLiteral{Value: "..."},
												Element: &mast.Identifier{Name: "int"},
											},
											Value: &mast.LiteralValue{
												Values: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
													&mast.IntLiteral{Value: "2"},
													&mast.IntLiteral{Value: "3"},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
									},
									Right: []mast.Expression{
										&mast.EntityCreationExpression{
											Object: nil,
											Type:   &mast.Identifier{Name: "Test"},
											Value: &mast.LiteralValue{
												Values: []mast.Expression{
													&mast.KeyValuePair{
														Key: &mast.Identifier{Name: "Parent"},
														Value: &mast.LiteralValue{
															Values: nil,
														},
													},
													&mast.KeyValuePair{
														Key:   &mast.Identifier{Name: "Key"},
														Value: &mast.Identifier{Name: "value"},
													},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "f"},
									},
									Right: []mast.Expression{
										&mast.FunctionLiteral{
											Parameters: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "int"},
													Name:       &mast.Identifier{Name: "x"},
												},
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "int"},
													Name:       &mast.Identifier{Name: "y"},
												},
											},
											Returns: []mast.Declaration{
												&mast.ParameterDeclaration{
													IsVariadic: false,
													Type:       &mast.Identifier{Name: "int"},
													Name:       nil,
												},
											},
											Statements: []mast.Statement{
												&mast.ReturnStatement{
													Exprs: []mast.Expression{
														&mast.BinaryExpression{
															Left:     &mast.Identifier{Name: "x"},
															Operator: "+",
															Right:    &mast.Identifier{Name: "y"},
														},
													},
												},
											},
										},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.FunctionLiteral{
										Parameters: nil,
										Returns:    nil,
										Statements: []mast.Statement{
											&mast.ReturnStatement{
												Exprs: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
												},
											},
										},
									},
									Arguments: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.FunctionLiteral{
										Parameters: nil,
										Returns: []mast.Declaration{
											&mast.ParameterDeclaration{
												IsVariadic: false,
												Type:       &mast.Identifier{Name: "A"},
												Name:       &mast.Identifier{Name: "a"},
											},
											&mast.ParameterDeclaration{
												IsVariadic: false,
												Type:       &mast.Identifier{Name: "B"},
												Name:       &mast.Identifier{Name: "b"},
											},
										},
										Statements: []mast.Statement{
											&mast.ReturnStatement{
												Exprs: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
												},
											},
										},
									},
									Arguments: nil,
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.CallExpression{
									Function: &mast.FunctionLiteral{
										Parameters: nil,
										Returns: []mast.Declaration{
											&mast.ParameterDeclaration{
												IsVariadic: false,
												Type:       &mast.Identifier{Name: "A"},
												Name:       nil,
											},
											&mast.ParameterDeclaration{
												IsVariadic: false,
												Type:       &mast.Identifier{Name: "B"},
												Name:       nil,
											},
										},
										Statements: []mast.Statement{
											&mast.ReturnStatement{
												Exprs: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
												},
											},
										},
									},
									Arguments: nil,
								},
							},
						},
					},
				},
			},
		},
		{
			description: "Test translating statements",
			file:        _metaTestDataPrefix + "go/statements.go",
			expected: &mast.Root{
				Declarations: []mast.Declaration{
					&mast.FunctionDeclaration{
						Name:       &mast.Identifier{Name: "root", Kind: mast.Function},
						Parameters: nil,
						Returns:    nil,
						LangFields: &mast.GoFunctionDeclarationFields{},
						Statements: []mast.Statement{
							&mast.GoDeferStatement{
								Expr: &mast.CallExpression{
									Function: &mast.Identifier{Name: "foo", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.Identifier{Name: "bar"},
									},
								},
							},
							&mast.ContinueStatement{
								Label: nil,
							},
							&mast.ContinueStatement{
								Label: &mast.Identifier{Name: "here", Kind: mast.Label},
							},
							&mast.BreakStatement{
								Label: nil,
							},
							&mast.BreakStatement{
								Label: &mast.Identifier{Name: "there", Kind: mast.Label},
							},
							&mast.GoGotoStatement{
								Label: &mast.Identifier{Name: "label1", Kind: mast.Label},
							},
							&mast.ReturnStatement{
								Exprs: nil,
							},
							&mast.ReturnStatement{
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "a"},
								},
							},
							&mast.ReturnStatement{
								Exprs: []mast.Expression{
									&mast.Identifier{Name: "a"},
									&mast.Identifier{Name: "b"},
								},
							},
							&mast.GoSendStatement{
								Channel: &mast.Identifier{Name: "c"},
								Value:   &mast.Identifier{Name: "v"},
							},
							&mast.GoGoStatement{
								Call: &mast.CallExpression{
									Function: &mast.Identifier{Name: "add", Kind: mast.Function},
									Arguments: []mast.Expression{
										&mast.IntLiteral{Value: "1"},
										&mast.IntLiteral{Value: "2"},
									},
								},
							},
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
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: false,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
										&mast.Identifier{Name: "b"},
									},
									Right: []mast.Expression{
										&mast.Identifier{Name: "c"},
										&mast.Identifier{Name: "d"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
									},
									Right: []mast.Expression{
										&mast.IntLiteral{Value: "1"},
									},
								},
							},
							&mast.ExpressionStatement{
								Expr: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "a"},
										&mast.Identifier{Name: "b"},
									},
									Right: []mast.Expression{
										&mast.Identifier{Name: "c"},
										&mast.Identifier{Name: "d"},
									},
								},
							},
							&mast.SwitchStatement{
								Initializer: &mast.ExpressionStatement{
									Expr: &mast.AssignmentExpression{
										IsShortVarDeclaration: false,
										Left: []mast.Expression{
											&mast.Identifier{Name: "a"},
										},
										Right: []mast.Expression{
											&mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
								Value: &mast.Identifier{Name: "a"},
								Cases: []*mast.SwitchCase{
									{
										Values: []mast.Expression{
											&mast.StringLiteral{Value: `"1"`},
											&mast.IntLiteral{Value: "2"},
										},
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
													Arguments: nil,
												},
											},
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "bar", Kind: mast.Function},
													Arguments: nil,
												},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													IsConst: false,
													Names:   []*mast.Identifier{{Name: "a"}},
													Type:    nil,
													Value:   &mast.IntLiteral{Value: "1"},
												},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													IsConst: false,
													Names:   []*mast.Identifier{{Name: "b"}},
													Type:    nil,
													Value:   &mast.IntLiteral{Value: "2"},
												},
											},
											&mast.GoFallthroughStatement{},
										},
									},
									{
										Values: []mast.Expression{
											&mast.IntLiteral{Value: "3"},
										},
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "test", Kind: mast.Function},
													Arguments: nil,
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
							&mast.SwitchStatement{
								Initializer: nil,
								Value:       nil,
								Cases: []*mast.SwitchCase{
									{
										Values: []mast.Expression{
											&mast.BinaryExpression{
												Operator: "<",
												Left:     &mast.Identifier{Name: "a"},
												Right:    &mast.Identifier{Name: "b"},
											},
										},
										Statements: []mast.Statement{
											&mast.ReturnStatement{
												Exprs: []mast.Expression{
													&mast.IntLiteral{Value: "1"},
												},
											},
										},
									},
								},
							},
							&mast.SwitchStatement{
								Initializer: &mast.ExpressionStatement{
									Expr: &mast.AssignmentExpression{
										IsShortVarDeclaration: true,
										Left: []mast.Expression{
											&mast.Identifier{Name: "a"},
										},
										Right: []mast.Expression{
											&mast.IntLiteral{Value: "1"},
										},
									},
								},
								Value: &mast.GoTypeSwitchHeaderExpression{
									Alias:   &mast.Identifier{Name: "n"},
									Operand: &mast.Identifier{Name: "c"},
								},
								Cases: nil,
							},
							&mast.SwitchStatement{
								Initializer: nil,
								Value: &mast.GoTypeSwitchHeaderExpression{
									Alias: nil,
									Operand: &mast.ParenthesizedExpression{
										Expr: &mast.UnaryExpression{
											Operator: "&",
											Expr:     &mast.Identifier{Name: "c"},
										},
									},
								},
								Cases: []*mast.SwitchCase{
									{
										Values: []mast.Expression{
											&mast.GoPointerType{
												Type: &mast.Identifier{Name: "Test"},
											},
										},
										Statements: nil,
									},
									{
										Values: nil,
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
													Arguments: nil,
												},
											},
										},
									},
								},
							},
							&mast.SwitchStatement{
								Initializer: &mast.ExpressionStatement{
									Expr: &mast.UpdateExpression{
										OperatorSide: mast.OperatorAfter,
										Operator:     "++",
										Operand:      &mast.Identifier{Name: "a"},
									},
								},
								Value: &mast.Identifier{Name: "a"},
								Cases: nil,
							},
							&mast.IfStatement{
								Initializer: nil,
								Condition:   &mast.Identifier{Name: "a"},
								Consequence: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
								Alternative: nil,
							},
							&mast.IfStatement{
								Initializer: nil,
								Condition:   &mast.Identifier{Name: "a"},
								Consequence: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "t1", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
								Alternative: &mast.IfStatement{
									Condition: &mast.Identifier{Name: "b"},
									Consequence: &mast.Block{
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "t2", Kind: mast.Function},
													Arguments: nil,
												},
											},
										},
									},
									Alternative: &mast.Block{
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "t3", Kind: mast.Function},
													Arguments: nil,
												},
											},
										},
									},
								},
							},
							&mast.IfStatement{
								Initializer: &mast.ExpressionStatement{
									Expr: &mast.AssignmentExpression{
										IsShortVarDeclaration: true,
										Left: []mast.Expression{
											&mast.Identifier{Name: "a"},
										},
										Right: []mast.Expression{
											&mast.IntLiteral{Value: "1"},
										},
									},
								},
								Condition:   &mast.Identifier{Name: "a"},
								Consequence: nil,
								Alternative: nil,
							},
							&mast.IfStatement{
								Initializer: &mast.ExpressionStatement{
									Expr: &mast.AssignmentExpression{
										IsShortVarDeclaration: false,
										Left: []mast.Expression{
											&mast.Identifier{Name: "a"},
										},
										Right: []mast.Expression{
											&mast.IntLiteral{Value: "1"},
										},
									},
								},
								Condition:   &mast.Identifier{Name: "a"},
								Consequence: nil,
								Alternative: nil,
							},
							&mast.LabelStatement{
								Label: &mast.Identifier{Name: "hello", Kind: mast.Label},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "a"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   nil,
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "b"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   nil,
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.GoTypeDeclaration{
									IsAlias: true,
									Name:    &mast.Identifier{Name: "T"},
									Type: &mast.GoMapType{
										Key:   &mast.Identifier{Name: "string"},
										Value: &mast.Identifier{Name: "bool"},
									},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.GoTypeDeclaration{
									IsAlias: false,
									Name:    &mast.Identifier{Name: "T"},
									Type:    &mast.Identifier{Name: "int"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "a"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   &mast.IntLiteral{Value: "1"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "b"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   &mast.IntLiteral{Value: "2"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "c"}},
									Type:    nil,
									Value:   &mast.IntLiteral{Value: "3"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "d"}},
									Type:    nil,
									Value:   &mast.IntLiteral{Value: "4"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "e"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   nil,
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names:   []*mast.Identifier{{Name: "f"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   nil,
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: false,
									Names: []*mast.Identifier{
										{Name: "h"},
										{Name: "i"},
									},
									Type: &mast.Identifier{Name: "int"},
									Value: &mast.CallExpression{
										Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
										Arguments: nil,
									},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: true,
									Names:   []*mast.Identifier{{Name: "a"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   &mast.IntLiteral{Value: "1"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: true,
									Names:   []*mast.Identifier{{Name: "b"}},
									Type:    &mast.Identifier{Name: "int"},
									Value:   &mast.IntLiteral{Value: "2"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: true,
									Names:   []*mast.Identifier{{Name: "c"}},
									Type:    nil,
									Value:   &mast.IntLiteral{Value: "3"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: true,
									Names:   []*mast.Identifier{{Name: "d"}},
									Type:    nil,
									Value:   nil,
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.VariableDeclaration{
									IsConst: true,
									Names:   []*mast.Identifier{{Name: "e"}},
									Type:    nil,
									Value:   &mast.IntLiteral{Value: "4"},
								},
							},
							&mast.DeclarationStatement{
								Decl: &mast.GoTypeDeclaration{
									IsAlias: false,
									Name:    &mast.Identifier{Name: "Test"},
									Type: &mast.GoStructType{
										Declarations: []*mast.FieldDeclaration{
											{
												Name:       nil,
												Type:       &mast.Identifier{Name: "A"},
												LangFields: &mast.GoFieldDeclarationFields{Tag: nil},
											},
											{
												Name:       &mast.Identifier{Name: "B"},
												Type:       &mast.Identifier{Name: "T"},
												LangFields: &mast.GoFieldDeclarationFields{Tag: nil},
											},
											{
												Name: &mast.Identifier{Name: "C"},
												Type: &mast.Identifier{Name: "int"},
												LangFields: &mast.GoFieldDeclarationFields{
													Tag: &mast.StringLiteral{
														IsRaw: false,
														Value: `"tag"`,
													},
												},
											},
											{
												Name: &mast.Identifier{Name: "D"},
												Type: &mast.Identifier{Name: "int"},
												LangFields: &mast.GoFieldDeclarationFields{
													Tag: &mast.StringLiteral{
														IsRaw: true,
														Value: "`t`",
													},
												},
											},
											{
												Name: &mast.Identifier{Name: "E"},
												Type: &mast.Identifier{Name: "int"},
												LangFields: &mast.GoFieldDeclarationFields{
													Tag: &mast.StringLiteral{
														IsRaw: true,
														Value: "`t`",
													},
												},
											},
										},
									},
								},
							},
							&mast.ForStatement{
								Initializers: []mast.Statement{
									&mast.ExpressionStatement{
										Expr: &mast.AssignmentExpression{
											IsShortVarDeclaration: true,
											Left: []mast.Expression{
												&mast.Identifier{Name: "i"},
											},
											Right: []mast.Expression{
												&mast.IntLiteral{Value: "1"},
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
								Body: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
							},
							&mast.ForStatement{
								Initializers: nil,
								Condition:    nil,
								Updates:      nil,
								Body: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
							},
							&mast.ForStatement{
								Initializers: nil,
								Condition:    nil,
								Updates:      nil,
								Body: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
							},
							&mast.ForStatement{
								Initializers: nil,
								Condition: &mast.CallExpression{
									Function:  &mast.Identifier{Name: "hasNext", Kind: mast.Function},
									Arguments: nil,
								},
								Updates: nil,
								Body: &mast.Block{
									Statements: []mast.Statement{
										&mast.ExpressionStatement{
											Expr: &mast.CallExpression{
												Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
												Arguments: nil,
											},
										},
									},
								},
							},
							&mast.GoForRangeStatement{
								Assignment: nil,
								Iterable:   &mast.Identifier{Name: "lst"},
								Body:       nil,
							},
							&mast.GoForRangeStatement{
								Assignment: &mast.AssignmentExpression{
									IsShortVarDeclaration: true,
									Left: []mast.Expression{
										&mast.Identifier{Name: "i"},
										&mast.Identifier{Name: "k"},
									},
									Right: []mast.Expression{
										&mast.Identifier{Name: "lst"},
									},
								},
								Iterable: nil,
								Body:     nil,
							},
							&mast.GoForRangeStatement{
								Assignment: &mast.AssignmentExpression{
									IsShortVarDeclaration: false,
									Left: []mast.Expression{
										&mast.AccessPath{
											Operand:     &mast.Identifier{Name: "pkg"},
											Annotations: nil,
											Field:       &mast.Identifier{Name: "A"},
										},
									},
									Right: []mast.Expression{
										&mast.Identifier{Name: "lst"},
									},
								},
								Iterable: nil,
								Body:     nil,
							},
							&mast.GoSelectStatement{
								Cases: []*mast.GoCommunicationCase{
									{
										Communication: &mast.ExpressionStatement{
											Expr: &mast.AssignmentExpression{
												IsShortVarDeclaration: true,
												Left: []mast.Expression{
													&mast.Identifier{Name: "a"},
													&mast.Identifier{Name: "b"},
												},
												Right: []mast.Expression{
													&mast.UnaryExpression{
														Operator: "<-",
														Expr:     &mast.Identifier{Name: "ch"},
													},
												},
											},
										},
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "foo", Kind: mast.Function},
													Arguments: nil,
												},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													IsConst: false,
													Names:   []*mast.Identifier{{Name: "x"}},
													Type:    &mast.Identifier{Name: "int"},
													Value:   &mast.IntLiteral{Value: "1"},
												},
											},
											&mast.DeclarationStatement{
												Decl: &mast.VariableDeclaration{
													IsConst: false,
													Names:   []*mast.Identifier{{Name: "y"}},
													Type:    &mast.Identifier{Name: "int"},
													Value:   &mast.IntLiteral{Value: "2"},
												},
											},
										},
									},
									{
										Communication: &mast.ExpressionStatement{
											Expr: &mast.AssignmentExpression{
												IsShortVarDeclaration: false,
												Left: []mast.Expression{
													&mast.Identifier{Name: "c"},
												},
												Right: []mast.Expression{
													&mast.UnaryExpression{
														Operator: "<-",
														Expr:     &mast.Identifier{Name: "x"},
													},
												},
											},
										},
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "bar", Kind: mast.Function},
													Arguments: nil,
												},
											},
										},
									},
									{
										Communication: &mast.ExpressionStatement{
											Expr: &mast.UnaryExpression{
												Operator: "<-",
												Expr:     &mast.Identifier{Name: "quit"},
											},
										},
										Statements: nil,
									},
									// the default case
									{
										Communication: nil,
										Statements: []mast.Statement{
											&mast.ExpressionStatement{
												Expr: &mast.CallExpression{
													Function:  &mast.Identifier{Name: "test", Kind: mast.Function},
													Arguments: nil,
												},
											},
										},
									},
								},
							},
							&mast.GoSelectStatement{
								Cases: []*mast.GoCommunicationCase{
									// the empty default case
									{
										Communication: nil,
										Statements:    nil,
									},
								},
							},
							&mast.LabelStatement{
								Label: &mast.Identifier{Name: "empty_label", Kind: mast.Label},
							},
						},
					},
				},
			},
		},
	}

	// Keep track of all TS node types that we have visited
	visitedNodes := make(map[string]bool)

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			node, err := ts.ParseFile(tc.file)
			require.NoError(t, err)

			// Put all visited TS node types in to the set and check it later, this must happen
			// before translation.Run since the translator might change the node type strings to
			// share implementation logic among different languages.
			err = reflectVisit(reflect.ValueOf(node), func(node *ts.Node) {
				visitedNodes[node.Type] = true
			})
			require.NoError(t, err)

			actual, err := Run(node, ts.GoExt)
			require.NoError(t, err)

			// We use cmp.Diff here since the diffing algorithm in require.Equal is not powerful
			// enough to give clear error messages.
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				require.FailNow(t, "mismatch (-expected +actual)", diff)
			}
		})
	}

	// Make sure all node types have been visited and tested
	for k := range _allGoTSNodeTypes {
		exists := visitedNodes[k]
		require.True(t, exists, "TS node %s not tested", k)
	}
}
