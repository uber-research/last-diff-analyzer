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

// Visitor is the interface that all MAST node visitors must implement. It contains a Pre(Node) and a Post(Node) method
// that is called before and after the traversal of each node.
type Visitor interface {
	// Pre takes in a MAST node for processing _before_ traversing its children and could return an error.
	Pre(Node) error
	// Post takes in a MAST node processing _after_ traversing its children and could return an error.
	Post(Node) error
}

// Walk takes a Visitor and walks the MAST. The input node must be non-nil.
func Walk(v Visitor, node Node) error {
	// Nilness checking for interface types in Golang is not very intuitive: "an interface value that holds a nil
	// concrete value is itself non-nil" -- A Tour of Go. An simple way of looking at this is that each interface type
	// will have a tuple of type and value (T, V) under the hood, and an interface value will be nil if and only if
	// (T=nil, V=nil). https://golang.org/doc/faq#nil_error provides more detailed explanations.
	// Due to this, testing "node == nil" here will not prevent us from nil dereferencing since node could be, for
	// example, (T=*Identifier, V=nil). Instead, we should make sure the nil fields in each node never get passed down.
	// Therefore, when we are traversing the children, we put nilness guard around each field and only do recursion if
	// it is not nil. This design follows the Go "ast" package. Note that although some fields in MAST nodes will never
	// be nil, we still wrap the recursions for the fields inside nilness guards just to be safe.
	// P.S. reflect package can actually be used to test the nilness of an interface (e.g., "reflect.ValueOf(i).IsNil()"),
	// but its usage is generally discouraged due to performance issues.
	// P.P.S. Checking nilness of the node fields that have interface types is fine. For example, we can check the
	// nilness of the field "Type" in "ParameterDeclaration", which has an interface type "Expression". This is because
	// during the translations an empty field will be assigned a pure nil pointer (i.e., (T=nil, V=nil)).

	if err := v.Pre(node); err != nil {
		return err
	}
	switch n := node.(type) {
	case *Root:
		if err := walkSlice(v, n.Declarations); err != nil {
			return err
		}

	case *Block:
		if err := walkSlice(v, n.Statements); err != nil {
			return err
		}

	case *PackageDeclaration:
		if n.Annotation != nil {
			if err := Walk(v, n.Annotation); err != nil {
				return err
			}
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}

	case *ImportDeclaration:
		if n.Alias != nil {
			if err := Walk(v, n.Alias); err != nil {
				return err
			}
		}
		if n.Package != nil {
			if err := Walk(v, n.Package); err != nil {
				return err
			}
		}

	case *ExpressionStatement:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *DeclarationStatement:
		if n.Decl != nil {
			if err := Walk(v, n.Decl); err != nil {
				return err
			}
		}

	case *ContinueStatement:
		if n.Label != nil {
			if err := Walk(v, n.Label); err != nil {
				return err
			}
		}

	case *BreakStatement:
		if n.Label != nil {
			if err := Walk(v, n.Label); err != nil {
				return err
			}
		}

	case *ReturnStatement:
		if err := walkSlice(v, n.Exprs); err != nil {
			return err
		}

	case *SwitchStatement:
		if n.Initializer != nil {
			if err := Walk(v, n.Initializer); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}
		for _, c := range n.Cases {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *SwitchCase:
		if err := walkSlice(v, n.Values); err != nil {
			return err
		}
		if err := walkSlice(v, n.Statements); err != nil {
			return err
		}

	case *IfStatement:
		if n.Initializer != nil {
			if err := Walk(v, n.Initializer); err != nil {
				return err
			}
		}
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}
		if n.Consequence != nil {
			if err := Walk(v, n.Consequence); err != nil {
				return err
			}
		}
		if n.Alternative != nil {
			if err := Walk(v, n.Alternative); err != nil {
				return err
			}
		}

	case *LabelStatement:
		if n.Label != nil {
			if err := Walk(v, n.Label); err != nil {
				return err
			}
		}

	case *ParenthesizedExpression:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *UnaryExpression:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *BinaryExpression:
		if n.Left != nil {
			if err := Walk(v, n.Left); err != nil {
				return err
			}
		}
		if n.Right != nil {
			if err := Walk(v, n.Right); err != nil {
				return err
			}
		}

	case *IndexExpression:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}
		if n.Index != nil {
			if err := Walk(v, n.Index); err != nil {
				return err
			}
		}

	case *AccessPath:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}
		if n.Annotations != nil {
			for _, c := range n.Annotations {
				if err := Walk(v, c); err != nil {
					return err
				}
			}
		}
		if n.Field != nil {
			if err := Walk(v, n.Field); err != nil {
				return err
			}
		}

	case *CallExpression:
		if n.Function != nil {
			if err := Walk(v, n.Function); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Arguments); err != nil {
			return err
		}
		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

	case *UpdateExpression:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}

	case *AssignmentExpression:
		if err := walkSlice(v, n.Left); err != nil {
			return err
		}
		if err := walkSlice(v, n.Right); err != nil {
			return err
		}

	case *ParameterDeclaration:
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

	case *VariableDeclaration:
		for _, c := range n.Names {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}
		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

	case *ForStatement:
		if err := walkSlice(v, n.Initializers); err != nil {
			return err
		}
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Updates); err != nil {
			return err
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *KeyValuePair:
		if n.Key != nil {
			if err := Walk(v, n.Key); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}

	case *LiteralValue:
		if err := walkSlice(v, n.Values); err != nil {
			return err
		}

	case *EntityCreationExpression:
		if n.Object != nil {
			if err := Walk(v, n.Object); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}
		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

	case *FunctionLiteral:
		if err := walkSlice(v, n.Parameters); err != nil {
			return err
		}
		if err := walkSlice(v, n.Returns); err != nil {
			return err
		}
		if n.Statements != nil {
			if err := walkSlice(v, n.Statements); err != nil {
				return err
			}
		}

	case *FieldDeclaration:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

	case *FunctionDeclaration:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Parameters); err != nil {
			return err
		}

		if err := walkSlice(v, n.Returns); err != nil {
			return err
		}

		if n.LangFields != nil {
			if err := Walk(v, n.LangFields); err != nil {
				return err
			}
		}

		if n.Statements != nil {
			if err := walkSlice(v, n.Statements); err != nil {
				return err
			}
		}

	case *GoSliceExpression:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}
		if n.Start != nil {
			if err := Walk(v, n.Start); err != nil {
				return err
			}
		}
		if n.End != nil {
			if err := Walk(v, n.End); err != nil {
				return err
			}
		}
		if n.Capacity != nil {
			if err := Walk(v, n.Capacity); err != nil {
				return err
			}
		}

	case *GoEllipsisExpression:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *GoPointerType:
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *GoArrayType:
		if n.Length != nil {
			if err := Walk(v, n.Length); err != nil {
				return err
			}
		}
		if n.Element != nil {
			if err := Walk(v, n.Element); err != nil {
				return err
			}
		}

	case *GoMapType:
		if n.Key != nil {
			if err := Walk(v, n.Key); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}

	case *GoParenthesizedType:
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *GoChannelType:
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *GoFunctionType:
		if err := walkSlice(v, n.Parameters); err != nil {
			return err
		}
		if err := walkSlice(v, n.Return); err != nil {
			return err
		}

	case *GoTypeAssertionExpression:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *GoTypeSwitchHeaderExpression:
		if n.Alias != nil {
			if err := Walk(v, n.Alias); err != nil {
				return err
			}
		}
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}

	case *GoDeferStatement:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *GoGotoStatement:
		if n.Label != nil {
			if err := Walk(v, n.Label); err != nil {
				return err
			}
		}

	case *GoSendStatement:
		if n.Channel != nil {
			if err := Walk(v, n.Channel); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}

	case *GoGoStatement:
		if n.Call != nil {
			if err := Walk(v, n.Call); err != nil {
				return err
			}
		}

	case *GoTypeDeclaration:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *GoStructType:
		for _, c := range n.Declarations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *GoInterfaceType:
		if err := walkSlice(v, n.Declarations); err != nil {
			return err
		}

	case *GoFieldDeclarationFields:
		if n.Tag != nil {
			if err := Walk(v, n.Tag); err != nil {
				return err
			}
		}

	case *GoForRangeStatement:
		if n.Assignment != nil {
			if err := Walk(v, n.Assignment); err != nil {
				return err
			}
		}
		if n.Iterable != nil {
			if err := Walk(v, n.Iterable); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *GoSelectStatement:
		for _, c := range n.Cases {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *GoCommunicationCase:
		if n.Communication != nil {
			if err := Walk(v, n.Communication); err != nil {
				return err
			}
		}
		if n.Statements != nil {
			if err := walkSlice(v, n.Statements); err != nil {
				return err
			}
		}

	case *GoFunctionDeclarationFields:
		if n.Receiver != nil {
			if err := Walk(v, n.Receiver); err != nil {
				return err
			}
		}

	case *JavaTernaryExpression:
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}
		if n.Consequence != nil {
			if err := Walk(v, n.Consequence); err != nil {
				return err
			}
		}
		if n.Alternative != nil {
			if err := Walk(v, n.Alternative); err != nil {
				return err
			}
		}

	case *JavaClassLiteral:
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *CastExpression:
		if n.Types != nil {
			if err := walkSlice(v, n.Types); err != nil {
				return err
			}
		}
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}

	case *Annotation:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.Arguments {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaAnnotatedType:
		for _, c := range n.Annotations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *JavaGenericType:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Arguments); err != nil {
			return err
		}

	case *JavaWildcard:
		if n.Super != nil {
			if err := Walk(v, n.Super); err != nil {
				return err
			}
		}
		if n.Extends != nil {
			if err := Walk(v, n.Extends); err != nil {
				return err
			}
		}

	case *JavaArrayType:
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaDimension:
		if n.Length != nil {
			if err := Walk(v, n.Length); err != nil {
				return err
			}
		}
		for _, c := range n.Annotations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaInstanceOfExpression:
		if n.Operand != nil {
			if err := Walk(v, n.Operand); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}

	case *JavaTryStatement:
		if err := walkSlice(v, n.Resources); err != nil {
			return err
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}
		for _, c := range n.CatchClauses {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Finally != nil {
			if err := Walk(v, n.Finally); err != nil {
				return err
			}
		}

	case *JavaCatchClause:
		if n.Parameter != nil {
			if err := Walk(v, n.Parameter); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaCatchFormalParameter:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if err := walkSlice(v, n.Types); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaFinallyClause:
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaWhileStatement:
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaThrowStatement:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}

	case *JavaAssertStatement:
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}
		if n.ErrorString != nil {
			if err := Walk(v, n.ErrorString); err != nil {
				return err
			}
		}

	case *JavaSynchronizedStatement:
		if n.Expr != nil {
			if err := Walk(v, n.Expr); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaDoStatement:
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}
		if n.Condition != nil {
			if err := Walk(v, n.Condition); err != nil {
				return err
			}
		}

	case *JavaParameterDeclarationFields:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaCallExpressionFields:
		if err := walkSlice(v, n.TypeArguments); err != nil {
			return err
		}

	case *JavaEnhancedForStatement:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Iterable != nil {
			if err := Walk(v, n.Iterable); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := Walk(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaModuleDeclaration:
		for _, c := range n.Annotations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.Directives {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaModuleDirective:
		if err := walkSlice(v, n.Exprs); err != nil {
			return err
		}

	case *JavaTypeParameter:
		for _, c := range n.Annotations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Extends); err != nil {
			return err
		}

	case *JavaClassDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.TypeParameters {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.SuperClass != nil {
			if err := Walk(v, n.SuperClass); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Interfaces); err != nil {
			return err
		}
		if err := walkSlice(v, n.Body); err != nil {
			return err
		}

	case *JavaInterfaceDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		for _, c := range n.TypeParameters {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Extends); err != nil {
			return err
		}
		if err := walkSlice(v, n.Body); err != nil {
			return err
		}

	case *JavaEnumDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Interfaces); err != nil {
			return err
		}
		if err := walkSlice(v, n.Body); err != nil {
			return err
		}

	case *JavaEnumConstantDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Arguments); err != nil {
			return err
		}
		if err := walkSlice(v, n.Body); err != nil {
			return err
		}

	case *JavaClassInitializer:
		if n.Block != nil {
			if err := Walk(v, n.Block); err != nil {
				return err
			}
		}

	case *JavaFunctionDeclarationFields:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		for _, c := range n.TypeParameters {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		for _, c := range n.Annotations {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Throws); err != nil {
			return err
		}

	case *JavaAnnotationDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.Body); err != nil {
			return err
		}

	case *JavaAnnotationElementDeclaration:
		if err := walkSlice(v, n.Modifiers); err != nil {
			return err
		}
		if n.Type != nil {
			if err := Walk(v, n.Type); err != nil {
				return err
			}
		}
		if n.Name != nil {
			if err := Walk(v, n.Name); err != nil {
				return err
			}
		}
		if n.Value != nil {
			if err := Walk(v, n.Value); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *JavaMethodReference:
		if n.Class != nil {
			if err := Walk(v, n.Class); err != nil {
				return err
			}
		}
		if err := walkSlice(v, n.TypeArguments); err != nil {
			return err
		}
		if n.Method != nil {
			if err := Walk(v, n.Method); err != nil {
				return err
			}
		}

	case *JavaEntityCreationExpressionFields:
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := walkSlice(v, n.Body); err != nil {
				return err
			}
		}

	case *JavaVariableDeclarationFields:
		for _, c := range n.Modifiers {
			if err := Walk(v, c); err != nil {
				return err
			}
		}
		for _, c := range n.Dimensions {
			if err := Walk(v, c); err != nil {
				return err
			}
		}

	case *CharacterLiteral, *JavaLiteralModifier, *FloatLiteral, *StringLiteral, *GoFallthroughStatement,
		*BooleanLiteral, *IntLiteral, *NullLiteral, *Identifier, *GoImaginaryLiteral:
		// The nodes above do not have children, so there is nothing to do.

	default:
		// *TempGroupNode and *JavaVariableDeclarator nodes are intentionally left out to be caught here, since they
		// are intermediate nodes that should never appear in the final tree.
		return fmt.Errorf("unexpected node %T during MAST traversal", node)
	}

	if err := v.Post(node); err != nil {
		return err
	}

	return nil
}

// walkSlice walks a slice of MAST nodes and calls Visitor on each of the element.
func walkSlice[T Node](v Visitor, slice []T) error {
	for _, c := range slice {
		if err := Walk(v, c); err != nil {
			return err
		}
	}
	return nil
}

// inspector implements the Visitor interface in a functional way, so that caller can simply
// supply an anonymous function instead of defining a complete struct for simple traversals. It is
// meant to be used in Inspect.
type inspector func(Node)

// Pre implements the required Visitor interface. It calls f on the node.
func (f inspector) Pre(node Node) error {
	f(node)
	return nil
}

// Post implements the required Visitor interface. It is a no-op.
func (f inspector) Post(Node) error {
	// no-op
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling f(node), and then Inspect
// invokes f recursively for each of the non-nil children.
func Inspect(node Node, f func(Node)) error {
	return Walk(inspector(f), node)
}
