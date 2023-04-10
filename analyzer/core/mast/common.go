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

package mast

import "fmt"

// NameKind (borrowed from the JLS terminology) describes both a
// declaration kind (e.g. label decl, function/method decl) and the
// kind of an identifier that should be linked ot a given declaration
// (e.g., label, function/method call).
type NameKind int

const (
	// Blanket is for all names whose more exact kind is
	// undetermined.
	Blanket NameKind = iota
	// Label is for labels. We distinguish labels as they,
	// conceptually, live in a different namespace from other
	// identifiers, both in Java and in Go. See the following links to
	// language specs for the explanation:
	//
	// https://docs.oracle.com/javase/specs/jls/se7/html/jls-14.html#jls-14.7
	// https://golang.org/ref/spec#Label_scopes
	Label
	// Function is for function declarations and
	// function invocations.
	Function
	// Method is for method/constructor declarations and
	// method/constructor invocations.
	Method
	// Pkg is for packages or subpackages, both declarations and uses
	// (e.g., in a qualified type name) - only used in Java
	Pkg
	// Typ is for type (e.g., class) declarations and type uses
	// (e.g., var/field/param types, return types, sub-types, etc.).
	Typ
)

// Node is the interface that all MAST nodes must implement.
type Node interface {
	// node ensures that only MAST nodes can be assigned to Node.
	node()
}

//
// A note on the design of common nodes and language-specific nodes.
//

// The generic nodes defined in this file contain the common fields that are shared across different languages. However,
// if a language has a few language-specific fields that are unique to itself, we store an extra pointer pointing to a
// special construct to save the language-specific fields. The pointer has an interface type which will point to
// different language-specific constructs for different languages. Note that a new interface has to be created for every
// common node that potentially can be extended, so that the special construct for one node will not be (accidentally)
// assigned to another one.
// Take the generic FieldDeclaration node that is shared across languages for example:
// FieldDeclaration
//   - Name
//   - Type
// In Go, field declarations have an extra field "Tag" and we need to extend the generic FieldDeclaration node to store
// this. Therefore, we use a GoFieldDeclarationFields node to store the Go-specific fields:
// GoFieldDeclarationFields
//   - Tag
// Then in FieldDeclaration, an extra pointer ("LangFields") is added:
// FieldDeclaration
//   - Name
//   - Type
//   - LangFields LangFieldDeclarationFields
// which will point to a GoFieldDeclarationFields node for Go and nil for other languages. During the build of MAST,
// language-specific builders will properly assign this field to point to a language-specific "extra fields" construct.
// Note that as long as the language has extra fields to the generic ones, the LangFields pointer will never be nil even
// if the extra fields are empty.

// TempGroupNode groups together a slice of Nodes. It exists to allow returning multiple translated Nodes as a single
// node (since Translate returns only one mast.Node). This is useful for many cases, especially for "un-grouping" a
// grouped node. For example, Go allows grouping multiple import declarations into one declaration:
// import (
//
//	"package1"
//	"package2"
//	....
//
// )
// the AST would look like this:
// source_file:
//
//	import_declaration:
//	  import_spec_list:
//	    import_spec
//	    import_spec
//	  ...
//
// When translating import_spec_list, we will return a TempGroupNode that groups together the translated ImportDeclaration
// nodes from import_spec. The TempGroupNode node is then returned to source_file level and "un-grouped" there.
// In this way, the grouped import declaration is "un-grouped" into several standalone
// import declarations. The resulting MAST tree is functionally equivalent to the original code, but hides away a few
// details from Go. This is also beneficial for later analysis where the tree will automatically be identical if the
// diff simply changes multiple import declarations into a grouped one.
// For reference, the resulting MAST tree will look like this:
// Root:
//
//	ImportDeclaration
//	ImportDeclaration
//	...
//
// Similarly, we use the same technique to un-group the const declarations and var declarations as well.
type TempGroupNode struct {
	Nodes []Node
}

// node implementation for TempGroupNode.
func (n *TempGroupNode) node() {}

// Root is the top level node for MAST.
type Root struct {
	// Declarations is the list of declarations in the file.
	Declarations []Declaration
}

// node implementation for Root nodes.
func (n *Root) node() {}

// Block node represents a block of statements (e.g., { statement1; statement2... })
type Block struct {
	// Statements is the slice of statements in the block.
	Statements []Statement
}

// node implementation for Block nodes.
func (n *Block) node() {}

// stmt implementation for Block nodes.
func (n *Block) stmt() {}

// Annotation node represents an annotation.
type Annotation struct {
	// Name is the name of the annotation.
	Name Expression
	// Arguments is the argument slice for the list of arguments to the annotation.
	Arguments []Expression
}

// node implementation for Annotation.
func (n *Annotation) node() {}

// expr implementation for Annotation.
func (n *Annotation) expr() {}

//
// node definitions for declarations
//

// Declaration is the interface all declaration nodes must implement.
type Declaration interface {
	Node
	// decl() ensures that only declaration nodes can be assigned to Declaration.
	decl()
}

// PackageDeclaration node represents the package declaration for source code.
type PackageDeclaration struct {
	// Annotation is the optional annotation for the package declaration, only available for Java, otherwise nil.
	Annotation *Annotation
	// Name is the name of the package declaration, could either be an Identifier or an AccessPath (only for Java).
	Name Expression
}

// node implementation for PackageDeclaration.
func (n *PackageDeclaration) node() {}

// decl implementation for PackageDeclaration.
func (n *PackageDeclaration) decl() {}

// ImportDeclaration node represents an import declaration.
type ImportDeclaration struct {
	// Alias is an optional field for the imported package. This field may not apply to all languages.
	Alias *Identifier
	// Package is the expression (AccessPath / Identifier in Java and StringLiteral in Go) of the imported package.
	Package Expression
}

// node implementation for ImportDeclaration.
func (n *ImportDeclaration) node() {}

// decl implementation for ImportDeclaration.
func (n *ImportDeclaration) decl() {}

//
// node definitions for statements
//

// Statement is the interface for all statement nodes.
type Statement interface {
	Node
	// stmt() ensures that only expression nodes can be assigned to Statement.
	stmt()
}

// ExpressionStatement node represents a stand-alone expression in a statement list.
type ExpressionStatement struct {
	// Expr is the expression in the statement.
	Expr Expression
}

// node implementation for ExpressionStatement.
func (n *ExpressionStatement) node() {}

// stmt implementation for ExpressionStatement.
func (n *ExpressionStatement) stmt() {}

// DeclarationStatement node represents a declaration in a statement list.
type DeclarationStatement struct {
	// Decl is the declaration in the statement.
	Decl Declaration
}

// node implementation for DeclarationStatement.
func (n *DeclarationStatement) node() {}

// stmt implementation for DeclarationStatement.
func (n *DeclarationStatement) stmt() {}

// ContinueStatement node represents a continue statement.
type ContinueStatement struct {
	// Label is the optional label for the continue statement.
	Label *Identifier
}

// node implementation for ContinueStatement.
func (n *ContinueStatement) node() {}

// stmt implementation for ContinueStatement.
func (n *ContinueStatement) stmt() {}

// BreakStatement node represents a break statement.
type BreakStatement struct {
	// Label is the optional label for the break statement.
	Label *Identifier
}

// node implementation for BreakStatement.
func (n *BreakStatement) node() {}

// stmt implementation for BreakStatement.
func (n *BreakStatement) stmt() {}

// ReturnStatement node represents a return statement.
type ReturnStatement struct {
	// Exprs is the slice of expressions to be returned. Note that some languages like Java only permit returning
	// one expression, in that case the Exprs will only contain one element.
	Exprs []Expression
}

// node implementation for ReturnStatement.
func (n *ReturnStatement) node() {}

// stmt implementation for ReturnStatement.
func (n *ReturnStatement) stmt() {}

// SwitchStatement node represents a switch statement.
type SwitchStatement struct {
	// Initializer is the optional initializer for the switch statement. Some languages (such as Go) support writing an
	// initializer for the switch statement (e.g., "switch a := 1; a {...}" ). Otherwise it is set to nil.
	Initializer Statement
	// Value is the expression for the switch statement.
	Value Expression
	// Cases is the list of switch case expressions.
	Cases []*SwitchCase
}

// node implementation for SwitchStatement.
func (n *SwitchStatement) node() {}

// stmt implementation for SwitchStatement.
func (n *SwitchStatement) stmt() {}

// SwitchCase node represents a case expression in the switch body.
type SwitchCase struct {
	// Values is the list of expressions for the switch case. For default case,
	// this is set to nil.
	Values []Expression
	// Statements is the list of expressions for the switch case, can be nil if
	// no statements are present. We do not use mast.Block here for consistency
	// with mast.FunctionDeclaration and mast.GoCommunicationCase.
	Statements []Statement
}

// node implementation for SwitchCase.
func (n *SwitchCase) node() {}

// stmt() implementation for SwitchCase.
func (n *SwitchCase) stmt() {}

// IfStatement node represents an if statement (e.g., "if (a) {...} else (b) {...}")
type IfStatement struct {
	// Initializer is an optional initializer for the if statement, only available in Go (e.g., "if a := 1; a {...}").
	Initializer Statement
	// Condition is the condition for the if statement.
	Condition Expression
	// Consequence is the true branch for the if statement, could be a Block node if multiple statements are present.
	Consequence Statement
	// Alternative is the false branch for the if statement, could be a Block node if multiple statements are present.
	// Also, it could another IfStatement for "if ... else if ..." structure.
	Alternative Statement
}

// node implementation for IfStatement.
func (n *IfStatement) node() {}

// stmt() implementation for IfStatement.
func (n *IfStatement) stmt() {}

// LabelStatement node represents a label statement (e.g., "label_name: ").
type LabelStatement struct {
	// Label is the label name of the label statement.
	Label *Identifier
}

// node implementation for LabelStatement.
func (n *LabelStatement) node() {}

// stmt() implementation for LabelStatement.
func (n *LabelStatement) stmt() {}

//
// node definitions for expressions
//

// Expression is the interface for all expression nodes.
type Expression interface {
	Node
	// expr() ensures that only expression nodes can be assigned to Expression.
	expr()
}

// Identifier node represents an identifier.
type Identifier struct {
	Name string
	Kind NameKind
}

// node implementation for Identifier.
func (n *Identifier) node() {}

// expr implementation for Identifier.
func (n *Identifier) expr() {}

// ParenthesizedExpression node represents a parenthesized expression (e.g., (a+b)).
type ParenthesizedExpression struct {
	// Expr is the parenthesized expression.
	Expr Expression
}

// node implementation for ParenthesizedExpression.
func (n *ParenthesizedExpression) node() {}

// expr implementation for ParenthesizedExpression.
func (n *ParenthesizedExpression) expr() {}

// UnaryExpression node represents an unary expression in the language.
type UnaryExpression struct {
	// Operator is the operator for the unary expression.
	Operator string
	// Expr is the expression for the unary operation.
	Expr Expression
}

// node implementation for UnaryExpression.
func (n *UnaryExpression) node() {}

// expr implementation for UnaryExpression.
func (n *UnaryExpression) expr() {}

// BinaryExpression node represents a binary expression in the language.
type BinaryExpression struct {
	// Operator is the operator for the binary expression.
	Operator string
	// Left is the left-hand side expression of the binary expression.
	Left Expression
	// Right is the right-hand side expression of the binary expression.
	Right Expression
}

// node implementation for BinaryExpression.
func (n *BinaryExpression) node() {}

// expr implementation for BinaryExpression.
func (n *BinaryExpression) expr() {}

// IndexExpression node represents a index expression (e.g., a[i]).
type IndexExpression struct {
	// Operand is the operand for the index operation.
	Operand Expression
	// Index is the index expression for the index operation.
	Index Expression
}

// node implementation for IndexExpression.
func (n *IndexExpression) node() {}

// expr implementation for IndexExpression.
func (n *IndexExpression) expr() {}

// AccessPath node represents an access path. It generally represents the "dot" operation (e.g., "expr.expr.expr...").
// For example, the following expressions will be represented by this node:
// (1) selector expression ("a.b.c") in Go;
// (2) field access ("this.field") in Java;
// (3) scoped identifier (package "java.lang.Boolean";);
// (4) scoped type identifier ("@UI Activity.@Safe Callback").
type AccessPath struct {
	// Operand is the expression for the select operation.
	Operand Expression
	// Annotations is an optional slice of annotations for annotations. In particular, it is used in scoped type identifiers in Java to add
	// annotations to the _field_ instead of the operands (The annotation "@UI" in "@UI Activity.@Safe @NotNull Callback"
	// annotates "Activity" and the @Safe and @NotNull annotate the "CallBack".
	// See https://docs.oracle.com/javase/specs/jls/se11/html/jls-4.html#jls-4.3 for a detailed explanation.
	// Note that @UI will not appear in the slice of annotations since it will be grouped together with rest of the
	// expression to form an JavaAnnotatedType. Therefore, the "@UI Activity.@Safe @NotNull Callback" expression will
	// have the following MAST node structure (simplified for brevity):
	// JavaAnnotatedType
	//   - Annotations: ["@UI"]
	//   - Type:
	//       AccessPath
	//         - Operand: "Activity"
	//         - Annotations: ["@Safe", "@NotNull"]
	//         - Field: "Callback"
	Annotations []*Annotation
	// Field is the field to access.
	Field *Identifier
}

// node implementation for AccessPath.
func (n *AccessPath) node() {}

// expr implementation for AccessPath.
func (n *AccessPath) expr() {}

// LangCallExpressionFields is the interface for the language-specific fields in CallExpression node.
type LangCallExpressionFields interface {
	Node
	// langCallExpressionFields ensures only language-specific fields for
	// CallExpression can be assigned to LangCallExpressionFields.
	langCallExpressionFields()
}

// CallExpression node represents a call expression or method invocation in Java (e.g., add(a, b)).
type CallExpression struct {
	// LangFields stores the language-specific fields.
	LangFields LangCallExpressionFields
	// Function is the expression for the function.
	Function Expression
	// Arguments is the argument slice for the list of arguments to the function, or nil if no arguments are given.
	Arguments []Expression
}

// node implementation for CallExpression.
func (n *CallExpression) node() {}

// expr implementation for CallExpression.
func (n *CallExpression) expr() {}

// NullLiteral node represents a Null literal expression (e.g., nil in Go and null in Java).
type NullLiteral struct{}

// node implementation for NullLiteral.
func (n *NullLiteral) node() {}

// expr implementation for NullLiteral.
func (n *NullLiteral) expr() {}

// BooleanLiteral node represents a boolean literal expression (e.g., true or false).
type BooleanLiteral struct {
	// Value is the value for the literal.
	Value bool
}

// node implementation for BooleanLiteral.
func (n *BooleanLiteral) node() {}

// expr implementation for BooleanLiteral.
func (n *BooleanLiteral) expr() {}

// IntLiteral node represents an int literal expression.
type IntLiteral struct {
	// Value is the value for the literal. Due to different bases of int literals (e.g., hexadecimal and octal) across
	// languages, we use strings to store the int literal.
	Value string
}

// node implementation for IntLiteral.
func (n *IntLiteral) node() {}

// expr implementation for IntLiteral.
func (n *IntLiteral) expr() {}

// FloatLiteral node represents a float literal expression.
type FloatLiteral struct {
	// Value is the value for the literal. Due to imprecision of floats and different bases of float literals, we use
	// strings to store the float literal.
	Value string
}

// node implementation for IntLiteral.
func (n *FloatLiteral) node() {}

// expr implementation for IntLiteral.
func (n *FloatLiteral) expr() {}

// StringLiteral node represents a string literal expression.
type StringLiteral struct {
	// IsRaw indicates whether the string literal is a raw string or not.
	IsRaw bool
	// Value is the value for the literal.
	Value string
}

// node implementation for StringLiteral.
func (n *StringLiteral) node() {}

// expr implementation for StringLiteral.
func (n *StringLiteral) expr() {}

// CharacterLiteral node represents a character literal expression.
type CharacterLiteral struct {
	// Value is the value for the literal.
	Value string
}

// node implementation for CharacterLiteral.
func (n *CharacterLiteral) node() {}

// expr implementation for CharacterLiteral.
func (n *CharacterLiteral) expr() {}

// UpdateExpressionOperatorSide indicates the side of operator in UpdateExpression (e.g., ++a vs a++).
type UpdateExpressionOperatorSide uint8

const (
	// OperatorBefore indicates the operator is placed _before_ the expression (e.g., ++a).
	OperatorBefore UpdateExpressionOperatorSide = iota
	// OperatorAfter indicates the operator is placed _after_ the expression (e.g., a++).
	OperatorAfter
)

// UpdateExpression node represents an update expression (e.g., a++).
type UpdateExpression struct {
	// OperatorSide indicates the side of the operator.
	OperatorSide UpdateExpressionOperatorSide
	// Operator is the operator of the update expression.
	Operator string
	// Operand is the operand of the update expression.
	Operand Expression
}

// node implementation for UpdateExpression.
func (n *UpdateExpression) node() {}

// expr implementation for UpdateExpression.
func (n *UpdateExpression) expr() {}

const (
	// AssignmentEqualOperator means the equal operator "=".
	AssignmentEqualOperator string = "="
	// AssignmentDeclareOperator means the operator ":=" in Go.
	AssignmentDeclareOperator string = ":="
)

// AssignmentExpression node represents an assignment expression (e.g., "a = b" or "a, b = c, d").
// Note that some languages only permit assigning one element in an assignment expression, in which case the Left and
// Right will only contain one element.
// Specially in Go, short variable declaration shares the same structure as AssignmentExpression. We use the same
// MAST node to represent short variable declarations, though tree-sitter uses a different node type for it. Our design
// follows that of the ast package in Go.
type AssignmentExpression struct {
	// IsShortVarDeclaration is a flag that indicates this assignment expression is a short variable declaration (e.g.,
	// "a, b := c, d") in Go.
	IsShortVarDeclaration bool
	// Left is the expression list on the left hand side of an assignment expression.
	Left []Expression
	// Right is the expression list on the right hand side of an assignment expression.
	Right []Expression
}

// node implementation for AssignmentExpression.
func (n *AssignmentExpression) node() {}

// expr implementation for AssignmentExpression.
func (n *AssignmentExpression) expr() {}

// LangParameterDeclarationFields is the interface for the language-specific fields in ParameterDeclaration node.
type LangParameterDeclarationFields interface {
	Node
	// langParameterDeclarationFields ensures only language-specific fields for
	// ParameterDeclaration can be assigned to LangParameterDeclarationFields.
	langParameterDeclarationFields()
}

// ParameterDeclaration node represents a parameter declaration.
type ParameterDeclaration struct {
	// LangFields stores the language-specific fields.
	LangFields LangParameterDeclarationFields
	// IsVariadic indicates whether this paramater is a variadic parameter declaration or not (e.g., "(a... int)" in
	// Go and "(int ...a)" in Java).
	IsVariadic bool
	// Type is the type of the parameter.
	Type Expression
	// Name is the name of the parameter.
	Name *Identifier
}

// node implementation for ParameterDeclaration.
func (n *ParameterDeclaration) node() {}

// decl implementation for ParameterDeclaration.
func (n *ParameterDeclaration) decl() {}

// LangVariableDeclarationFields is the interface for the language-specific fields in VariableDeclaration node.
type LangVariableDeclarationFields interface {
	Node
	// langVariableDeclarationFields ensures only language-specific fields for
	// VariableDeclaration can be assigned to LangVariableDeclarationFields.
	langVariableDeclarationFields()
}

// VariableDeclaration node represents a variable declaration (e.g., "var a int = 2" in Go or "int a = 2;" in Java).
type VariableDeclaration struct {
	// LangFields stores the language-specific fields.
	LangFields LangVariableDeclarationFields
	// IsConst indicates whether this variable declaration is a const (or final) declaration.
	IsConst bool
	// Names is a slice of variables being declared. Although we try to split multi-variable declaration into multiple
	// single-variable declarations (e.g., "int x = 1, y = 2;" -> "int x = 1; int y = 2;") for simpler tree structure,
	// there are cases for some languages such as Go that allows imbalanced variable declarations (e.g.,
	// "var x, y int = foo()", where "foo()" returns two values). Currently, the length of this slice must be 1 for
	// other languages except for Go. Note that this is a special case and the RHS must be a "a single multi-valued
	// expression such as a function call, a channel or map operation, or a type assertion" from [1]. So the Value field
	// remains a single value.
	// [1] https://golang.org/ref/spec#Assignments
	Names []*Identifier
	// Type is the optional type for the variable declaration.
	Type Expression
	// Value is the optional initial value for the variable declaration.
	Value Expression
}

// node implementation for VariableDeclaration.
func (n *VariableDeclaration) node() {}

// decl implementation for VariableDeclaration.
func (n *VariableDeclaration) decl() {}

// ForStatement node represents a plain for statement that consists of "init; condition; update" (e.g.,
// "for (i = 1; i < 10; i++) {...}"). This node is designed to be general enough to support many language features:
// (1) Initializers and Updates are slices of statements to support multiple statements in them. For example, Java
//
//	allows writing multiple statements in "init" and "update";
//
// (2) All Initializers, Condition and Updates are optional to support omissions of them. For example, Go allows writing
//
//	for statements in many ways, including "for {...}" and "for condition {...}".
type ForStatement struct {
	// Initializers is an optional slice of initializers.
	Initializers []Statement
	// Condition is an optional condition of the for statement.
	Condition Expression
	// Updates is an optional slice of update statements.
	Updates []Statement
	// Body is the body of the for statement. Some languages such as Java allows writing a single statement in a for
	// statement, for better unification we will wrap it inside a *Block node instead.
	Body *Block
}

// node implementation for ForStatement.
func (n *ForStatement) node() {}

// stmt implementation for ForStatement.
func (n *ForStatement) stmt() {}

// KeyValuePair node represents a common structure in many other nodes, such as keyed parameter (e.g., "(a=1, b=3)") or
// struct literals (e.g., "{A: a, B: b}") etc. It is designed to be a general node structure that can be shared across
// languages.
type KeyValuePair struct {
	// Key is the key for the key-value pair.
	Key Expression
	// Value is the Value for the key-value pair.
	Value Expression
}

// node implementation for KeyValuePair.
func (n *KeyValuePair) node() {}

// expr implementation for KeyValuePair.
func (n *KeyValuePair) expr() {}

// LiteralValue node represents an array initializer (e.g, "{1, 2, 3}") or a struct initializer (e.g., "{A: a, B: b}").
type LiteralValue struct {
	// Values is the slice of values in the literal value, each element could either be a KeyValuePair or a single
	// Expression.
	Values []Expression
}

// node implementation for LiteralValue.
func (n *LiteralValue) node() {}

// expr implementation for LiteralValue.
func (n *LiteralValue) expr() {}

// LangEntityCreationExpressionFields is the interface for the
// language-specific fields in EntityCreationExpression node.
type LangEntityCreationExpressionFields interface {
	Node
	// langEntityCreationExpressionFields ensures only
	// language-specific fields for EntityCreationExpression can be
	// assigned to langEntityCreationExpressionFields.
	langEntityCreationExpressionFields()
}

// EntityCreationExpression node represents a literal object creation expression, such as composite literal nodes
// (e.g., "[]int{1, 2, 3}", "Test{A: a}") in Go or array creation and object creation in Java
// (e.g., "new int[]{1, 2, 3}" and "new Test(1, 2, 3)").
type EntityCreationExpression struct {
	// LangFields stores the language-specific fields.
	LangFields LangEntityCreationExpressionFields
	// Object is the outer object expression (e.g., "outer.new Inner();") for the entity creation expression. Only
	// available in Java.
	Object Expression
	// Type is the type of the entity creation expression.
	Type Expression
	// Value is the value for the entity creation expression.
	Value *LiteralValue
}

// node implementation for EntityCreationExpression.
func (n *EntityCreationExpression) node() {}

// expr implementation for EntityCreationExpression.
func (n *EntityCreationExpression) expr() {}

// FunctionLiteral node represents an anonymous function literal. Specifically, it represents lambda expression
// ("(o -> foo())") in Java and func literal ("func (a A) (b B) {...}") in Go.
type FunctionLiteral struct {
	// Parameters is the optional slice of parameter declaration. Each element could either be a ParameterDeclaration or
	// a JavaParameterDeclaration.
	Parameters []Declaration
	// Returns is the optional slice of result parameter declaration, only available on Go. Each element could only be
	// a ParameterDeclaration node, but for better compatibility (so that we do not need two sets of helper functions
	// for normal parameters and return parameters), we still use Declaration type here.
	Returns []Declaration
	// Statements is the body of the function literal. This could be either an actual list of statements (representing a block
	// statements in the function, e.g., "(o -> {...})" in Java) or a single expression wrapped in mast.ExpressionStatement
	// (for returning a single expression, e.g., "(o -> foo(o))" in Java).
	// For Go, the Body would always be a list of statements representing a block ("func (a A) B {...}").
	// Please see the comment by  FunctionDeclaration to see why we cannot use the actual mast.Block as a child node here.
	Statements []Statement
}

// node implementation for FunctionLiteral.
func (n *FunctionLiteral) node() {}

// expr implementation for FunctionLiteral.
func (n *FunctionLiteral) expr() {}

// CastExpression node represents a type cast expression (e.g., "(NewType) a" in Java and "[]byte(json)" in Go). Note
// that cast expressions are represented by different tree-sitter nodes in Go ("type_conversion_expression") and Java (
// "cast_expression").
type CastExpression struct {
	// Types is the slice of types to cast for the operation. This is designed for Java which allows multiple types
	// being present in the cast expression (e.g., "(T1 & T2 & T3 ...)expr"). Note that in Go the length of the slice
	// should always be 1.
	Types []Expression
	// Operand is the expression to perform the cast operation on.
	Operand Expression
}

// node implementation for CastExpression.
func (n *CastExpression) node() {}

// expr implementation for CastExpression.
func (n *CastExpression) expr() {}

// LangFieldDeclarationFields is the interface for the language-specific fields in FieldDeclaration node.
type LangFieldDeclarationFields interface {
	Node
	// langFieldDeclarationFields ensures only language-specific fields for FunctionDeclaration can be assigned to
	// LangFieldDeclarationFields.
	langFieldDeclarationFields()
}

// FieldDeclaration node represents a common field declaration.
type FieldDeclaration struct {
	// LangFields stores the extra fields for language-specific fields.
	LangFields LangFieldDeclarationFields
	// Name is the name of the field.
	Name *Identifier
	// Type is the type for the field.
	Type Expression
}

// node implementation for FieldDeclaration.
func (n *FieldDeclaration) node() {}

// decl implementation for FieldDeclaration.
func (n *FieldDeclaration) decl() {}

// LangFunctionDeclarationFields is the interface for the language-specific fields in FunctionDeclaration node.
type LangFunctionDeclarationFields interface {
	Node
	// langFunctionDeclarationFields ensures only language-specific fields for FunctionDeclaration can be assigned to
	// LangFunctionDeclarationFields.
	langFunctionDeclarationFields()
}

// FunctionDeclaration node represents a function declaration.
type FunctionDeclaration struct {
	// LangFields stores the language-specific fields for FunctionDeclaration.
	LangFields LangFunctionDeclarationFields
	// Name is the name of the function declaration.
	Name *Identifier
	// Parameters is an optional list of formal parameters for the function declaration.
	Parameters []Declaration
	// Returns is an optional list of return parameters for the function declaration.
	Returns []Declaration
	// Statements represent the optional body of the method declaration.
	// This field can be:
	// (1) nil to indicate that this is a function declaration without
	//     definition / implementation, e.g., "public int f1(T t1, T t2);"
	//     in Java;
	// (2) an empty list to indicate a function declaration with an empty body,
	//     e.g., "func (a int) int {}" in Go.
	// We cannot use a mast.Block here as block will create its own (nested)
	// scope, which does not include the function parameters. This is acceptable
	// in Java but is not acceptable in Go
	// (https://golang.org/ref/spec#Declarations_and_scope : "The scope of an
	// identifier denoting a method receiver, function parameter, or result
	// variable is the function body").
	// Note that this is required by further analyses (e.g., symbolication for
	// short assignments in Go) where the analysis requires the function
	// parameters and variable declarations in Statements to be in the same
	// scope.
	Statements []Statement
}

// node implementation for MethodDeclaration.
func (n *FunctionDeclaration) node() {}

// decl implementation for MethodDeclaration.
func (n *FunctionDeclaration) decl() {}

// SetCallKind sets the kind of the standalone identifier representing
// function/method call name for the generic (language-agnostic) case.
//
// This method does some additional verification to make sure that we
// do not miss setting the kind for other relevant cases, such as for
// calls via an access path, for example a.b.foo().
// TODO: NameKind should be decoupled from the node itself and should be handled
// in the symbolication process.
func SetCallKind(expr Expression, kind NameKind) error {
	switch e := expr.(type) {
	case *Identifier:
		// simple call: foo()
		e.Kind = kind
	case *AccessPath:
		// a call via access path: a.b.foo()
		//
		// don't do anything (as we don't know if foo() is, for
		// example, a method name or a variable name), which is safe
		// since we do not symbolicate identifiers on access path
		// beyond the first one anyway
	case *UnaryExpression:
		// The type conversions (not type assertion) in Golang is interpreted
		// as call expressions in tree-sitter, e.g., converting variable "a" to
		// int would be "int(a)" (with no syntactical difference to a normal
		// function call). Moreover, you can convert to pointer types such as
		// "(*int)(a)". We do nothing here to assign the identifier a Blanket
		// kind to be matched against any kind of declarations in the
		// symbolication process in order to support both type casts and
		// lambda function calls.
		if e.Operator != "*" {
			return fmt.Errorf("unexpected unary operator in call expression %s", e.Operator)
		}
	case *CallExpression:
		// a call via another call (don't do anything): foo()()
	case *FunctionLiteral:
		// Go function literal or Java lambda expression (don't do
		// anything): func() {...} or () -> {...}
	case *ParenthesizedExpression:
		// For parenthesized expression, we recursively call this function to
		// unwrap the parentheses.
		return SetCallKind(e.Expr, kind)
	default:
		// we want to make sure that we mark all identifiers
		// representing calls so if we missed a type of call
		// expression, we should throw an error
		//
		// Note that we intentionally leave out handling of BinaryExpression
		// above to be caught here since it is almost never legal in call
		// expressions in languages we support.
		return fmt.Errorf("unknown call expression %T", expr)
	}
	return nil
}
