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

package check

import (
	"analyzer/core/mast"
)

// GenericChecker is the generic checker for checking the functionality equivalence of two MAST forests.
type GenericChecker struct{}

// Equal method implements a generic equivalence checking process.
func (g *GenericChecker) Equal(forest []mast.Node, other []mast.Node, c LangChecker) (bool, error) {
	if len(forest) != len(other) {
		return false, nil
	}

	// iterate over the forest and check if every node is equivalent
	for i := 0; i < len(forest); i++ {
		node1, node2 := forest[i], other[i]
		if node1 == nil && node2 == nil {
			return true, nil
		}
		if node1 == nil || node2 == nil {
			return false, nil
		}
		isEqual := c.CheckNode(node1, node2, c)
		if !isEqual {
			return false, nil
		}
	}
	return true, nil
}

// The following Equal methods implement a strict equality checking
// process for each node, i.e., all fields in two nodes must be equal
// for the two nodes to be considered equal.

// CheckNode compares "generic" nodes.
func (g *GenericChecker) CheckNode(node mast.Node, other mast.Node, c LangChecker) bool {
	switch n1 := node.(type) {
	case *mast.TempGroupNode:
		if n2, ok := other.(*mast.TempGroupNode); ok {
			return g.TempGroupNodeEq(n1, n2, c)
		}
	case *mast.Root:
		if n2, ok := other.(*mast.Root); ok {
			return g.RootEq(n1, n2, c)
		}
	case *mast.Block:
		if n2, ok := other.(*mast.Block); ok {
			return g.BlockEq(n1, n2, c)
		}
	case *mast.PackageDeclaration:
		if n2, ok := other.(*mast.PackageDeclaration); ok {
			return g.PackageDeclarationEq(n1, n2, c)
		}
	case *mast.ImportDeclaration:
		if n2, ok := other.(*mast.ImportDeclaration); ok {
			return g.ImportDeclarationEq(n1, n2, c)
		}
	case *mast.ExpressionStatement:
		if n2, ok := other.(*mast.ExpressionStatement); ok {
			return g.ExpressionStatementEq(n1, n2, c)
		}
	case *mast.DeclarationStatement:
		if n2, ok := other.(*mast.DeclarationStatement); ok {
			return g.DeclarationStatementEq(n1, n2, c)
		}
	case *mast.ContinueStatement:
		if n2, ok := other.(*mast.ContinueStatement); ok {
			return g.ContinueStatementEq(n1, n2, c)
		}
	case *mast.BreakStatement:
		if n2, ok := other.(*mast.BreakStatement); ok {
			return g.BreakStatementEq(n1, n2, c)
		}
	case *mast.ReturnStatement:
		if n2, ok := other.(*mast.ReturnStatement); ok {
			return g.ReturnStatementEq(n1, n2, c)
		}
	case *mast.SwitchStatement:
		if n2, ok := other.(*mast.SwitchStatement); ok {
			return g.SwitchStatementEq(n1, n2, c)
		}
	case *mast.SwitchCase:
		if n2, ok := other.(*mast.SwitchCase); ok {
			return g.SwitchCaseEq(n1, n2, c)
		}
	case *mast.IfStatement:
		if n2, ok := other.(*mast.IfStatement); ok {
			return g.IfStatementEq(n1, n2, c)
		}
	case *mast.LabelStatement:
		if n2, ok := other.(*mast.LabelStatement); ok {
			return g.LabelStatementEq(n1, n2, c)
		}
	case *mast.Identifier:
		if n2, ok := other.(*mast.Identifier); ok {
			return g.IdentifierEq(n1, n2, c)
		}
	case *mast.ParenthesizedExpression:
		if n2, ok := other.(*mast.ParenthesizedExpression); ok {
			return g.ParenthesizedExpressionEq(n1, n2, c)
		}
	case *mast.UnaryExpression:
		if n2, ok := other.(*mast.UnaryExpression); ok {
			return g.UnaryExpressionEq(n1, n2, c)
		}
	case *mast.BinaryExpression:
		if n2, ok := other.(*mast.BinaryExpression); ok {
			return g.BinaryExpressionEq(n1, n2, c)
		}
	case *mast.IndexExpression:
		if n2, ok := other.(*mast.IndexExpression); ok {
			return g.IndexExpressionEq(n1, n2, c)
		}
	case *mast.AccessPath:
		if n2, ok := other.(*mast.AccessPath); ok {
			return g.AccessPathEq(n1, n2, c)
		}
	case *mast.CallExpression:
		if n2, ok := other.(*mast.CallExpression); ok {
			return g.CallExpressionEq(n1, n2, c)
		}
	case *mast.NullLiteral:
		if n2, ok := other.(*mast.NullLiteral); ok {
			return g.NullLiteralEq(n1, n2, c)
		}
	case *mast.BooleanLiteral:
		if n2, ok := other.(*mast.BooleanLiteral); ok {
			return g.BooleanLiteralEq(n1, n2, c)
		}
	case *mast.IntLiteral:
		if n2, ok := other.(*mast.IntLiteral); ok {
			return g.IntLiteralEq(n1, n2, c)
		}
	case *mast.FloatLiteral:
		if n2, ok := other.(*mast.FloatLiteral); ok {
			return g.FloatLiteralEq(n1, n2, c)
		}
	case *mast.StringLiteral:
		if n2, ok := other.(*mast.StringLiteral); ok {
			return g.StringLiteralEq(n1, n2, c)
		}
	case *mast.CharacterLiteral:
		if n2, ok := other.(*mast.CharacterLiteral); ok {
			return g.CharacterLiteralEq(n1, n2, c)
		}
	case *mast.UpdateExpression:
		if n2, ok := other.(*mast.UpdateExpression); ok {
			return g.UpdateExpressionEq(n1, n2, c)
		}
	case *mast.AssignmentExpression:
		if n2, ok := other.(*mast.AssignmentExpression); ok {
			return g.AssignmentExpressionEq(n1, n2, c)
		}
	case *mast.ParameterDeclaration:
		if n2, ok := other.(*mast.ParameterDeclaration); ok {
			return g.ParameterDeclarationEq(n1, n2, c)
		}
	case *mast.VariableDeclaration:
		if n2, ok := other.(*mast.VariableDeclaration); ok {
			return g.VariableDeclarationEq(n1, n2, c)
		}
	case *mast.ForStatement:
		if n2, ok := other.(*mast.ForStatement); ok {
			return g.ForStatementEq(n1, n2, c)
		}
	case *mast.KeyValuePair:
		if n2, ok := other.(*mast.KeyValuePair); ok {
			return g.KeyValuePairEq(n1, n2, c)
		}
	case *mast.LiteralValue:
		if n2, ok := other.(*mast.LiteralValue); ok {
			return g.LiteralValueEq(n1, n2, c)
		}
	case *mast.EntityCreationExpression:
		if n2, ok := other.(*mast.EntityCreationExpression); ok {
			return g.EntityCreationExpressionEq(n1, n2, c)
		}
	case *mast.FunctionLiteral:
		if n2, ok := other.(*mast.FunctionLiteral); ok {
			return g.FunctionLiteralEq(n1, n2, c)
		}
	case *mast.FieldDeclaration:
		if n2, ok := other.(*mast.FieldDeclaration); ok {
			return g.FieldDeclarationEq(n1, n2, c)
		}
	case *mast.FunctionDeclaration:
		if n2, ok := other.(*mast.FunctionDeclaration); ok {
			return g.FunctionDeclarationEq(n1, n2, c)
		}
	case *mast.CastExpression:
		if n2, ok := other.(*mast.CastExpression); ok {
			return g.CastExpressionEq(n1, n2, c)
		}
	case *mast.Annotation:
		if n2, ok := other.(*mast.Annotation); ok {
			return g.AnnotationEq(n1, n2, c)
		}
	}
	// for "untyped" nil values, for example in CheckNode(nil, nil)  call
	if node == nil && other == nil {
		return true
	}
	return false
}

// The reason why the following methods need a nil check is because
// the CheckNode method above cannot tell if a passed value is nil or
// not ("typed" nil values passed via an interface result in a non-nil
// interface value).

// TempGroupNodeEq compares TempGroupNode-s.
func (g *GenericChecker) TempGroupNodeEq(n1 *mast.TempGroupNode, n2 *mast.TempGroupNode, c LangChecker) bool {
	// alwasy return false for TempGroupNode node regardless of the
	// other node since it should never appear in the final MAST.
	return false
}

// RootEq compares Root-s.
func (g *GenericChecker) RootEq(n1 *mast.Root, n2 *mast.Root, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDeclarationList(n1.Declarations, n2.Declarations, c)
}

// BlockEq compares Block-s.
func (g *GenericChecker) BlockEq(n1 *mast.Block, n2 *mast.Block, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckStatementList(n1.Statements, n2.Statements, c)
}

// PackageDeclarationEq compares PackageDeclaration-s.
func (g *GenericChecker) PackageDeclarationEq(n1 *mast.PackageDeclaration, n2 *mast.PackageDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Annotation, n2.Annotation, c) && c.CheckNode(n1.Name, n2.Name, c)
}

// ImportDeclarationEq compares ImportDeclaration-s.
func (g *GenericChecker) ImportDeclarationEq(n1 *mast.ImportDeclaration, n2 *mast.ImportDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Alias, n2.Alias, c) && c.CheckNode(n1.Package, n2.Package, c)
}

// ExpressionStatementEq compares ExpressionStatement-s.
func (g *GenericChecker) ExpressionStatementEq(n1 *mast.ExpressionStatement, n2 *mast.ExpressionStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Expr, n2.Expr, c)
}

// DeclarationStatementEq compares DeclarationStatement-s.
func (g *GenericChecker) DeclarationStatementEq(n1 *mast.DeclarationStatement, n2 *mast.DeclarationStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Decl, n2.Decl, c)
}

// ContinueStatementEq compares ContinueStatement-s.
func (g *GenericChecker) ContinueStatementEq(n1 *mast.ContinueStatement, n2 *mast.ContinueStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Label, n2.Label, c)
}

// BreakStatementEq compares BreakStatement-s.
func (g *GenericChecker) BreakStatementEq(n1 *mast.BreakStatement, n2 *mast.BreakStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Label, n2.Label, c)
}

// ReturnStatementEq compares ReturnStatement-s.
func (g *GenericChecker) ReturnStatementEq(n1 *mast.ReturnStatement, n2 *mast.ReturnStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckExpressionList(n1.Exprs, n2.Exprs, c)
}

// SwitchStatementEq compares SwitchStatement-s.
func (g *GenericChecker) SwitchStatementEq(n1 *mast.SwitchStatement, n2 *mast.SwitchStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Cases) != len(n2.Cases) {
		return false
	}
	for i := 0; i < len(n1.Cases); i++ {
		isEqual := c.CheckNode(n1.Cases[i], n2.Cases[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Initializer, n2.Initializer, c) && c.CheckNode(n1.Value, n2.Value, c)
}

// SwitchCaseEq compares SwitchCase-s.
func (g *GenericChecker) SwitchCaseEq(n1 *mast.SwitchCase, n2 *mast.SwitchCase, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckExpressionList(n1.Values, n2.Values, c) && c.CheckStatementList(n1.Statements, n2.Statements, c)
}

// IfStatementEq compares IfStatement-s.
func (g *GenericChecker) IfStatementEq(n1 *mast.IfStatement, n2 *mast.IfStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Initializer, n2.Initializer, c) &&
		c.CheckNode(n1.Condition, n2.Condition, c) &&
		c.CheckNode(n1.Consequence, n2.Consequence, c) &&
		c.CheckNode(n1.Alternative, n2.Alternative, c)
}

// LabelStatementEq compares LabelStatement-s.
func (g *GenericChecker) LabelStatementEq(n1 *mast.LabelStatement, n2 *mast.LabelStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Label, n2.Label, c)
}

// IdentifierEq compares Identifier-s.
func (g *GenericChecker) IdentifierEq(n1 *mast.Identifier, n2 *mast.Identifier, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Name == n2.Name
}

// ParenthesizedExpressionEq compares ParenthesizedExpression-s.
func (g *GenericChecker) ParenthesizedExpressionEq(n1 *mast.ParenthesizedExpression, n2 *mast.ParenthesizedExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Expr, n2.Expr, c)
}

// UnaryExpressionEq compares UnaryExpression-s.
func (g *GenericChecker) UnaryExpressionEq(n1 *mast.UnaryExpression, n2 *mast.UnaryExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Operator == n2.Operator && c.CheckNode(n1.Expr, n2.Expr, c)
}

// BinaryExpressionEq compares BinaryExpression-s.
func (g *GenericChecker) BinaryExpressionEq(n1 *mast.BinaryExpression, n2 *mast.BinaryExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Operator == n2.Operator &&
		c.CheckNode(n1.Left, n2.Left, c) &&
		c.CheckNode(n1.Right, n2.Right, c)
}

// IndexExpressionEq compares IndexExpression-s.
func (g *GenericChecker) IndexExpressionEq(n1 *mast.IndexExpression, n2 *mast.IndexExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Operand, n2.Operand, c) && c.CheckNode(n1.Index, n2.Index, c)
}

// AccessPathEq compares AccessPath-s.
func (g *GenericChecker) AccessPathEq(n1 *mast.AccessPath, n2 *mast.AccessPath, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Annotations) != len(n2.Annotations) {
		return false
	}
	for i := 0; i < len(n1.Annotations); i++ {
		isEqual := c.CheckNode(n1.Annotations[i], n2.Annotations[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Field, n2.Field, c) && c.CheckNode(n1.Operand, n2.Operand, c)
}

// CallExpressionEq compares CallExpression-s.
func (g *GenericChecker) CallExpressionEq(n1 *mast.CallExpression, n2 *mast.CallExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Function, n2.Function, c) && c.CheckExpressionList(n1.Arguments, n2.Arguments, c) && c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// NullLiteralEq compares NullLiteral-s.
func (g *GenericChecker) NullLiteralEq(n1 *mast.NullLiteral, n2 *mast.NullLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return true
}

// BooleanLiteralEq compares BooleanLiteral-s.
func (g *GenericChecker) BooleanLiteralEq(n1 *mast.BooleanLiteral, n2 *mast.BooleanLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value
}

// IntLiteralEq compares IntLiteral-s.
func (g *GenericChecker) IntLiteralEq(n1 *mast.IntLiteral, n2 *mast.IntLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value
}

// FloatLiteralEq compares FloatLiteral-s.
func (g *GenericChecker) FloatLiteralEq(n1 *mast.FloatLiteral, n2 *mast.FloatLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value
}

// StringLiteralEq compares StringLiteral-s.
func (g *GenericChecker) StringLiteralEq(n1 *mast.StringLiteral, n2 *mast.StringLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value && n1.IsRaw == n2.IsRaw
}

// CharacterLiteralEq compares CharacterLiteral-s.
func (g *GenericChecker) CharacterLiteralEq(n1 *mast.CharacterLiteral, n2 *mast.CharacterLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value
}

// UpdateExpressionEq compares UpdateExpression-s.
func (g *GenericChecker) UpdateExpressionEq(n1 *mast.UpdateExpression, n2 *mast.UpdateExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Operator == n2.Operator &&
		n1.OperatorSide == n2.OperatorSide &&
		c.CheckNode(n1.Operand, n2.Operand, c)
}

// AssignmentExpressionEq compares AssignmentExpression-s.
func (g *GenericChecker) AssignmentExpressionEq(n1 *mast.AssignmentExpression, n2 *mast.AssignmentExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.IsShortVarDeclaration == n2.IsShortVarDeclaration &&
		c.CheckExpressionList(n1.Left, n2.Left, c) &&
		c.CheckExpressionList(n1.Right, n2.Right, c)
}

// ParameterDeclarationEq compares ParameterDeclaration-s.
func (g *GenericChecker) ParameterDeclarationEq(n1 *mast.ParameterDeclaration, n2 *mast.ParameterDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.IsVariadic == n2.IsVariadic &&
		c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// VariableDeclarationEq compares VariableDeclaration-s.
func (g *GenericChecker) VariableDeclarationEq(n1 *mast.VariableDeclaration, n2 *mast.VariableDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Names) != len(n2.Names) {
		return false
	}
	for i := 0; i < len(n1.Names); i++ {
		isEqual := c.CheckNode(n1.Names[i], n2.Names[i], c)
		if !isEqual {
			return false
		}
	}
	// We approve modifications from non-const to const, so the only disapproval would be for
	// changing from const to non-const.
	if n1.IsConst && !n2.IsConst {
		return false
	}
	return c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.Value, n2.Value, c) &&
		c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// ForStatementEq compares ForStatement-s.
func (g *GenericChecker) ForStatementEq(n1 *mast.ForStatement, n2 *mast.ForStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Body, n2.Body, c) &&
		c.CheckNode(n1.Condition, n2.Condition, c) &&
		c.CheckStatementList(n1.Initializers, n2.Initializers, c) &&
		c.CheckStatementList(n1.Updates, n2.Updates, c)
}

// KeyValuePairEq compares KeyValuePair-s.
func (g *GenericChecker) KeyValuePairEq(n1 *mast.KeyValuePair, n2 *mast.KeyValuePair, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Key, n2.Key, c) && c.CheckNode(n1.Value, n2.Value, c)
}

// LiteralValueEq compares LiteralValue-s.
func (g *GenericChecker) LiteralValueEq(n1 *mast.LiteralValue, n2 *mast.LiteralValue, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckExpressionList(n1.Values, n2.Values, c)
}

// EntityCreationExpressionEq compares EntityCreationExpression-s.
func (g *GenericChecker) EntityCreationExpressionEq(n1 *mast.EntityCreationExpression, n2 *mast.EntityCreationExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.Value, n2.Value, c) &&
		c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// FunctionLiteralEq compares FunctionLiteral-s.
func (g *GenericChecker) FunctionLiteralEq(n1 *mast.FunctionLiteral, n2 *mast.FunctionLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckStatementList(n1.Statements, n2.Statements, c) &&
		c.CheckDeclarationList(n1.Parameters, n2.Parameters, c) &&
		c.CheckDeclarationList(n1.Returns, n2.Returns, c)
}

// FieldDeclarationEq compares FieldDeclaration-s.
func (g *GenericChecker) FieldDeclarationEq(n1 *mast.FieldDeclaration, n2 *mast.FieldDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// FunctionDeclarationEq compares FunctionDeclaration-s.
func (g *GenericChecker) FunctionDeclarationEq(n1 *mast.FunctionDeclaration, n2 *mast.FunctionDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckStatementList(n1.Statements, n2.Statements, c) &&
		c.CheckDeclarationList(n1.Parameters, n2.Parameters, c) &&
		c.CheckDeclarationList(n1.Returns, n2.Returns, c) &&
		c.CheckNode(n1.LangFields, n2.LangFields, c)
}

// CastExpressionEq compares CastExpression-s.
func (g *GenericChecker) CastExpressionEq(n1 *mast.CastExpression, n2 *mast.CastExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckExpressionList(n1.Types, n2.Types, c) && c.CheckNode(n1.Operand, n2.Operand, c)
}

// AnnotationEq compares Annotation-s.
func (g *GenericChecker) AnnotationEq(n1 *mast.Annotation, n2 *mast.Annotation, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Arguments) != len(n2.Arguments) {
		return false
	}
	for i := 0; i < len(n1.Arguments); i++ {
		isEqual := c.CheckNode(n1.Arguments[i], n2.Arguments[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Name, n2.Name, c)
}

// CheckDeclarationList is a helper function that check if every
// declaration in the given two lists is equal. All elements in the
// list must be non-nil.
func (g *GenericChecker) CheckDeclarationList(declarations []mast.Declaration, other []mast.Declaration, c LangChecker) bool {
	if len(declarations) != len(other) {
		return false
	}
	for i := 0; i < len(declarations); i++ {
		isEqual := c.CheckNode(declarations[i], other[i], c)
		if !isEqual {
			return false
		}
	}
	return true
}

// CheckStatementList is a helper function that check if every
// statement in the given two lists is equal. All elements in the list
// must be non-nil.
func (g *GenericChecker) CheckStatementList(statements []mast.Statement, other []mast.Statement, c LangChecker) bool {
	if len(statements) != len(other) {
		return false
	}
	for i := 0; i < len(statements); i++ {
		isEqual := c.CheckNode(statements[i], other[i], c)
		if !isEqual {
			return false
		}
	}
	return true
}

// CheckExpressionList is a helper function that check if every
// expression in the given two lists is equal. All elements in the
// list must be non-nil.
func (g *GenericChecker) CheckExpressionList(expressions []mast.Expression, other []mast.Expression, c LangChecker) bool {
	if len(expressions) != len(other) {
		return false
	}
	for i := 0; i < len(expressions); i++ {
		isEqual := c.CheckNode(expressions[i], other[i], c)
		if !isEqual {
			return false
		}
	}
	return true
}

// CheckDimensions is a helper function that check if two arrays of
// dimensions are equal. All elements in the list must be non-nil.
func (g *GenericChecker) CheckDimensions(dimensions []*mast.JavaDimension, other []*mast.JavaDimension, c LangChecker) bool {
	if len(dimensions) != len(other) {
		return false
	}
	for i := 0; i < len(dimensions); i++ {
		isEqual := c.CheckNode(dimensions[i], other[i], c)
		if !isEqual {
			return false
		}
	}
	return true

}
