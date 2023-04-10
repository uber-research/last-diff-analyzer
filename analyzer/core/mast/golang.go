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

//
// Go-specific node definitions, all nodes are prefixed with Go.
//

// GoSliceExpression node represents a slice expression in Go (e.g., a[i:i+5]).
type GoSliceExpression struct {
	// Operand is the expression for the slice operation.
	Operand Expression
	// Start is the begin of slice range; or nil.
	Start Expression
	// End is the end of slice range; or nil.
	End Expression
	// Capacity is the maximum capacity of slice; or nil.
	Capacity Expression
}

// node implementation for GoSliceExpression.
func (n *GoSliceExpression) node() {}

// expr implementation for GoSliceExpression.
func (n *GoSliceExpression) expr() {}

// GoEllipsisExpression node represents an ellipsis expression used in function arguments in Go (e.g., test(i, j...)).
type GoEllipsisExpression struct {
	// Expr is the expression preceded by the ellipsis (...).
	Expr Expression
}

// node implementation for GoEllipsisExpression.
func (n *GoEllipsisExpression) node() {}

// expr implementation for GoEllipsisExpression.
func (n *GoEllipsisExpression) expr() {}

// GoImaginaryLiteral node represents an imaginary literal expression (e.g., 5i).
type GoImaginaryLiteral struct {
	// Value is the value for the imaginary part, note that it includes the "i" at the end.
	Value string
}

// node implementation for GoImaginaryLiteral.
func (n *GoImaginaryLiteral) node() {}

// expr implementation for GoImaginaryLiteral.
func (n *GoImaginaryLiteral) expr() {}

// GoPointerType node represents a pointer type in Go (e.g., *Test).
type GoPointerType struct {
	// Type is the type of the struct the pointer points to.
	Type Expression
}

// node implementation for GoPointerType.
func (n *GoPointerType) node() {}

// expr implementation for GoPointerType.
func (n *GoPointerType) expr() {}

// GoArrayType node represents an array/slice type in Go (e.g., "[5]Test" and "[]Test").
type GoArrayType struct {
	// Length is the optional length of the array/slice. The value of Length indicates several variants of array type
	// in Go:
	// (1) nil indicates it is a slice type;
	// (2) StringLiteral("...") indicates that it is an array type with implicit length;
	// (3) a normal expression indicates that it is a fixed-length array type.
	Length Expression
	// Element is the type of the elements in the array.
	Element Expression
}

// node implementation for GoArrayType.
func (n *GoArrayType) node() {}

// expr implementation for GoArrayType.
func (n *GoArrayType) expr() {}

// GoMapType node represents a map type in Go (e.g., map[string]bool)).
type GoMapType struct {
	// Key is the type of key in the map.
	Key Expression
	// Value is the type of value in the map.
	Value Expression
}

// node implementation for GoMapType.
func (n *GoMapType) node() {}

// expr implementation for GoMapType.
func (n *GoMapType) expr() {}

// GoParenthesizedType node represents a parenthesized type (e.g. (int)).
type GoParenthesizedType struct {
	// Type is the inner type wrapped by the parentheses.
	Type Expression
}

// node implementation for GoParenthesizedType.
func (n *GoParenthesizedType) node() {}

// expr implementation for GoParenthesizedType.
func (n *GoParenthesizedType) expr() {}

// GoChannelTypeDirection is an enum to indicate the direction of GoChannelType.
type GoChannelTypeDirection uint8

const (
	// SendAndReceive means the channel is bidirectional (e.g., chan int).
	SendAndReceive GoChannelTypeDirection = iota
	// ReceiveOnly means the channel is receive-only (e.g., <-chan int).
	ReceiveOnly
	// SendOnly means the channel is send-only (e.g., chan<- int).
	SendOnly
)

// GoChannelType node represents a channel type (e.g. "chan int", "<- chan int" or "chan <- int").
type GoChannelType struct {
	// Direction is the direction of the channel.
	Direction GoChannelTypeDirection
	// Type is the type of the channel.
	Type Expression
}

// node implementation for GoChannelType.
func (n *GoChannelType) node() {}

// expr implementation for GoChannelType.
func (n *GoChannelType) expr() {}

// GoFunctionType node represents a function type (e.g. "func (a A) B").
type GoFunctionType struct {
	// Parameters is the slice of parameters.
	Parameters []Declaration
	// Return is the slice of result parameters.
	Return []Declaration
}

// node implementation for GoFunctionType.
func (n *GoFunctionType) node() {}

// expr implementation for GoFunctionType.
func (n *GoFunctionType) expr() {}

// GoTypeAssertionExpression node represents a type assertion expression in Go (e.g., "a.(T)").
type GoTypeAssertionExpression struct {
	// Operand is the operand for the type assertion operation.
	Operand Expression
	// Type is the type for the type assertion operation.
	Type Expression
}

// node implementation for GoTypeAssertionExpression.
func (n *GoTypeAssertionExpression) node() {}

// expr implementation for GoTypeAssertionExpression.
func (n *GoTypeAssertionExpression) expr() {}

// GoTypeSwitchHeaderExpression node represents the type switch expression in Go (e.g., "n := a.(type)").
// This is only intended to be used in SwitchStatement node for Go AST, where it is possible to switch on the type of
// an expression.
type GoTypeSwitchHeaderExpression struct {
	// Alias is an optional alias for the type.
	Alias *Identifier
	// Operand is the operand for the type switch expression.
	Operand Expression
}

// node implementation for GoTypeSwitchHeaderExpression.
func (n *GoTypeSwitchHeaderExpression) node() {}

// expr implementation for GoTypeSwitchHeaderExpression.
func (n *GoTypeSwitchHeaderExpression) expr() {}

// GoDeferStatement node represents a defer statement in Go (e.g., "defer foo(bar)")).
type GoDeferStatement struct {
	// Expr is the expression of the defer statement.
	Expr Expression
}

// node implementation for GoDeferStatement.
func (n *GoDeferStatement) node() {}

// stmt implementation for GoDeferStatement.
func (n *GoDeferStatement) stmt() {}

// GoGotoStatement node represents a goto statement in Go (e.g., "goto label_name")).
type GoGotoStatement struct {
	// Label is the label for the goto statement.
	Label *Identifier
}

// node implementation for GoGotoStatement.
func (n *GoGotoStatement) node() {}

// stmt implementation for GoGotoStatement.
func (n *GoGotoStatement) stmt() {}

// GoFallthroughStatement node represents a fallthrough statement in Go (e.g., "fallthrough")).
type GoFallthroughStatement struct{}

// node implementation for GoFallthroughStatement.
func (n *GoFallthroughStatement) node() {}

// stmt implementation for GoFallthroughStatement.
func (n *GoFallthroughStatement) stmt() {}

// GoSendStatement node represents a send statement in Go (e.g., "c <- v")).
type GoSendStatement struct {
	// Channel is the channel expression for the send statement.
	Channel Expression
	// Value is the value to be sent on the Channel.
	Value Expression
}

// node implementation for GoSendStatement.
func (n *GoSendStatement) node() {}

// stmt implementation for GoSendStatement.
func (n *GoSendStatement) stmt() {}

// GoGoStatement node represents a go statement in Go (e.g., "go add(1, 2)")).
type GoGoStatement struct {
	// Call is the function call expression for the go statement.
	Call *CallExpression
}

// node implementation for GoGoStatement.
func (n *GoGoStatement) node() {}

// stmt implementation for GoGoStatement.
func (n *GoGoStatement) stmt() {}

// GoTypeDeclaration node represents a type declaration in Go (e.g., "type T map[string]bool" or "type T = int").
type GoTypeDeclaration struct {
	// IsAlias indicates whether the type declaration is a type alias or not.
	// For example, "type T int" is a type declaration and "type T = int" is a type alias.
	// See https://golang.org/ref/spec#Type_declarations for additional explanations.
	IsAlias bool
	// Name is the alias for the type.
	Name *Identifier
	// Type is the type for the type alias statement.
	Type Expression
}

// node implementation for GoTypeDeclaration.
func (n *GoTypeDeclaration) node() {}

// decl implementation for GoTypeDeclaration.
func (n *GoTypeDeclaration) decl() {}

// GoStructType node represents a struct type declaration in Go.
type GoStructType struct {
	// Declarations is a slice of field declarations in the struct.
	Declarations []*FieldDeclaration
}

// node implementation for GoStructType.
func (n *GoStructType) node() {}

// expr implementation for GoStructType.
func (n *GoStructType) expr() {}

// GoInterfaceType node represents an interface type in Go (e.g., "T interface {...}").
type GoInterfaceType struct {
	// Declarations is a list of either:
	// (1) FunctionDeclaration (without Body) for method specifications;
	// (2) FieldDeclaration (without name and tag) for interface embeddings.
	Declarations []Declaration
}

// node implementation for GoInterfaceType.
func (n *GoInterfaceType) node() {}

// expr implementation for GoInterfaceType.
func (n *GoInterfaceType) expr() {}

// GoForRangeStatement node represents a for range statement in Go (e.g., "for i, n := range list {...}").
type GoForRangeStatement struct {
	// Assignment is the assignment expression containing the left and right of the range clause. The range clause
	// behaves very similar to an (imbalanced) AssignmentExpression, therefore here we reuse the structure. If there is
	// no assignment / declaration (e.g., "for range lst {...}"), this field is set to nil and the following Iterable
	// field is set instead. Note that only one of Assignment and Iterable can be non-nil.
	Assignment *AssignmentExpression
	// Iterable is the expression for the range clause.
	Iterable Expression
	// Body is the body of the for statement.
	Body *Block
}

// node implementation for GoForRangeStatement.
func (n *GoForRangeStatement) node() {}

// stmt implementation for GoForRangeStatement.
func (n *GoForRangeStatement) stmt() {}

// GoSelectStatement node represents a select statement in Go (e.g., "select {...}").
type GoSelectStatement struct {
	// Cases is the slice of all communication cases.
	Cases []*GoCommunicationCase
}

// node implementation for GoSelectStatement.
func (n *GoSelectStatement) node() {}

// stmt implementation for GoSelectStatement.
func (n *GoSelectStatement) stmt() {}

// GoCommunicationCase node represents a special case statement that is used in select statement
// (e.g., "case i1 = <-c1"). It is similar to SwitchCase node, but instead of Expression, one Statement node is required
// for the case condition. If a bare expression is given, it has to be wrapped inside a ExpressionStatement.
type GoCommunicationCase struct {
	// Communication is the statement for the communication case. For default case, it is set to nil.
	Communication Statement
	// Statements is the list of expressions for the communication case, can be
	// nil if no statements are present.
	// Please see the comment by FunctionDeclaration to see why we cannot use
	// mast.Block as a child node here.
	Statements []Statement
}

// node implementation for GoCommunicationCase.
func (n *GoCommunicationCase) node() {}

// stmt implementation for GoCommunicationCase.
func (n *GoCommunicationCase) stmt() {}

//
// Definitions for language-specific fields that extend the generic MAST nodes.
//

// GoFunctionDeclarationFields node stores the Go-specific extra fields for the generic FunctionDeclaration node, to
// represent a method declaration in Go (e.g., "func (a *A) test () {...}").
type GoFunctionDeclarationFields struct {
	// Receivers is an optional receiver for the method declarations. Note that for function declarations, it is set to
	// nil.
	Receiver Declaration
}

// node implementation for GoFunctionDeclarationFields.
func (n *GoFunctionDeclarationFields) node() {}

// langFunctionDeclarationFields implementation for GoFunctionDeclarationFields.
func (n *GoFunctionDeclarationFields) langFunctionDeclarationFields() {}

// GoFieldDeclarationFields node stores the Go-specific fields for the generic FieldDeclaration node, to represent a
// field declaration in Go (e.g., "a int").
type GoFieldDeclarationFields struct {
	// Tag is an optional tag for the field, only availabe for Go.
	Tag *StringLiteral
}

// node implementation for GoFieldDeclarationFields.
func (n *GoFieldDeclarationFields) node() {}

// langFieldDeclarationFields implementation for GoFieldDeclarationFields.
func (n *GoFieldDeclarationFields) langFieldDeclarationFields() {}
