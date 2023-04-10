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
	"analyzer/core/symbolication"
)

// GoChecker represents a Go-specific equivalence checker.
type GoChecker struct {
	Checker

	// LoggingOn indicates whether logging-related auto-approvals should be enabled.
	LoggingOn bool
	// LogImportClean holds information on whether logging imports in
	// a given pair of files (for base ane last versions) were
	// "clean", that is, if they imported the right logging library
	// and were not aliased (which could make auto-approval of related
	// logging changes incorrect). This field is set by the checker
	// before comparing any nodes of the base and last forests.
	LogImportClean bool
}

// NewGoChecker created and returns Go-specific equivalence checker.
func NewGoChecker(checker Checker, loggingOn bool) *GoChecker {
	return &GoChecker{Checker: checker, LoggingOn: loggingOn}
}

// StmtIgnored checks if a given statement can be ignored when
// checking for equivalence. If it can be ignored altogether, this
// method returns true. If it can be ignored under the condition that
// some of its sub-expressions are equal, this method returns false
// and a non-nil list of sub-expressions. If it cannot be ignored
// under any circumstances, this method returns false and an empty (or
// nil) list of sub-expressions.
func (g *GoChecker) StmtIgnored(stmt mast.Statement, _ *symbolication.SymbolTable) (bool, []mast.Expression) {
	if !g.LoggingOn {
		return false, nil
	}

	exprStmt, ok := stmt.(*mast.ExpressionStatement)
	if !ok || exprStmt.Expr == nil {
		// not an expression
		return false, nil
	}
	return g.callIgnored(exprStmt.Expr, true /* outer */)
}

// CheckNode compares individual nodes.
func (g *GoChecker) CheckNode(node mast.Node, other mast.Node, c LangChecker) bool {
	switch n1 := node.(type) {
	case *mast.GoSliceExpression:
		if n2, ok := other.(*mast.GoSliceExpression); ok {
			return g.GoSliceExpressionEq(n1, n2, c)
		}
	case *mast.GoEllipsisExpression:
		if n2, ok := other.(*mast.GoEllipsisExpression); ok {
			return g.GoEllipsisExpressionEq(n1, n2, c)
		}
	case *mast.GoImaginaryLiteral:
		if n2, ok := other.(*mast.GoImaginaryLiteral); ok {
			return g.GoImaginaryLiteralEq(n1, n2, c)
		}
	case *mast.GoPointerType:
		if n2, ok := other.(*mast.GoPointerType); ok {
			return g.GoPointerTypeEq(n1, n2, c)
		}
	case *mast.GoArrayType:
		if n2, ok := other.(*mast.GoArrayType); ok {
			return g.GoArrayTypeEq(n1, n2, c)
		}
	case *mast.GoMapType:
		if n2, ok := other.(*mast.GoMapType); ok {
			return g.GoMapTypeEq(n1, n2, c)
		}
	case *mast.GoParenthesizedType:
		if n2, ok := other.(*mast.GoParenthesizedType); ok {
			return g.GoParenthesizedTypeEq(n1, n2, c)
		}
	case *mast.GoChannelType:
		if n2, ok := other.(*mast.GoChannelType); ok {
			return g.GoChannelTypeEq(n1, n2, c)
		}
	case *mast.GoFunctionType:
		if n2, ok := other.(*mast.GoFunctionType); ok {
			return g.GoFunctionTypeEq(n1, n2, c)
		}
	case *mast.GoTypeAssertionExpression:
		if n2, ok := other.(*mast.GoTypeAssertionExpression); ok {
			return g.GoTypeAssertionExpressionEq(n1, n2, c)
		}
	case *mast.GoTypeSwitchHeaderExpression:
		if n2, ok := other.(*mast.GoTypeSwitchHeaderExpression); ok {
			return g.GoTypeSwitchHeaderExpressionEq(n1, n2, c)
		}
	case *mast.GoDeferStatement:
		if n2, ok := other.(*mast.GoDeferStatement); ok {
			return g.GoDeferStatementEq(n1, n2, c)
		}
	case *mast.GoGotoStatement:
		if n2, ok := other.(*mast.GoGotoStatement); ok {
			return g.GoGotoStatementEq(n1, n2, c)
		}
	case *mast.GoFallthroughStatement:
		if n2, ok := other.(*mast.GoFallthroughStatement); ok {
			return g.GoFallthroughStatementEq(n1, n2, c)
		}
	case *mast.GoSendStatement:
		if n2, ok := other.(*mast.GoSendStatement); ok {
			return g.GoSendStatementEq(n1, n2, c)
		}
	case *mast.GoGoStatement:
		if n2, ok := other.(*mast.GoGoStatement); ok {
			return g.GoGoStatementEq(n1, n2, c)
		}
	case *mast.GoTypeDeclaration:
		if n2, ok := other.(*mast.GoTypeDeclaration); ok {
			return g.GoTypeDeclarationEq(n1, n2, c)
		}
	case *mast.GoStructType:
		if n2, ok := other.(*mast.GoStructType); ok {
			return g.GoStructTypeEq(n1, n2, c)
		}
	case *mast.GoInterfaceType:
		if n2, ok := other.(*mast.GoInterfaceType); ok {
			return g.GoInterfaceTypeEq(n1, n2, c)
		}
	case *mast.GoFieldDeclarationFields:
		if n2, ok := other.(*mast.GoFieldDeclarationFields); ok {
			return g.GoFieldDeclarationFieldsEq(n1, n2, c)
		}
	case *mast.GoForRangeStatement:
		if n2, ok := other.(*mast.GoForRangeStatement); ok {
			return g.GoForRangeStatementEq(n1, n2, c)
		}
	case *mast.GoSelectStatement:
		if n2, ok := other.(*mast.GoSelectStatement); ok {
			return g.GoSelectStatementEq(n1, n2, c)
		}
	case *mast.GoCommunicationCase:
		if n2, ok := other.(*mast.GoCommunicationCase); ok {
			return g.GoCommunicationCaseEq(n1, n2, c)
		}
	case *mast.GoFunctionDeclarationFields:
		if n2, ok := other.(*mast.GoFunctionDeclarationFields); ok {
			return g.GoFunctionDeclarationFieldsEq(n1, n2, c)
		}
	}
	return g.Checker.CheckNode(node, other, c)
}

// GoSliceExpressionEq compares GoSliceExpression-s.
func (g *GoChecker) GoSliceExpressionEq(n1 *mast.GoSliceExpression, n2 *mast.GoSliceExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Operand, n2.Operand, c) &&
		c.CheckNode(n1.Start, n2.Start, c) &&
		c.CheckNode(n1.End, n2.End, c) &&
		c.CheckNode(n1.Capacity, n2.Capacity, c)
}

// GoEllipsisExpressionEq compares GoEllipsisExpression-s.
func (g *GoChecker) GoEllipsisExpressionEq(n1 *mast.GoEllipsisExpression, n2 *mast.GoEllipsisExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Expr, n2.Expr, c)
}

// GoImaginaryLiteralEq compares GoImaginaryLiteral-s.
func (g *GoChecker) GoImaginaryLiteralEq(n1 *mast.GoImaginaryLiteral, n2 *mast.GoImaginaryLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Value == n2.Value
}

// GoPointerTypeEq compares GoPointerType-s.
func (g *GoChecker) GoPointerTypeEq(n1 *mast.GoPointerType, n2 *mast.GoPointerType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Type, n2.Type, c)
}

// GoArrayTypeEq compares GoArrayType-s.
func (g *GoChecker) GoArrayTypeEq(n1 *mast.GoArrayType, n2 *mast.GoArrayType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Length, n2.Length, c) && c.CheckNode(n1.Element, n2.Element, c)
}

// GoMapTypeEq compares GoMapType-s.
func (g *GoChecker) GoMapTypeEq(n1 *mast.GoMapType, n2 *mast.GoMapType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Key, n2.Key, c) && c.CheckNode(n1.Value, n2.Value, c)
}

// GoParenthesizedTypeEq compares GoParenthesizedType-s.
func (g *GoChecker) GoParenthesizedTypeEq(n1 *mast.GoParenthesizedType, n2 *mast.GoParenthesizedType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Type, n2.Type, c)
}

// GoChannelTypeEq compares GoChannelType-s.
func (g *GoChecker) GoChannelTypeEq(n1 *mast.GoChannelType, n2 *mast.GoChannelType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Direction == n2.Direction && c.CheckNode(n1.Type, n2.Type, c)
}

// GoFunctionTypeEq compares GoFunctionType-s.
func (g *GoChecker) GoFunctionTypeEq(n1 *mast.GoFunctionType, n2 *mast.GoFunctionType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDeclarationList(n1.Parameters, n2.Parameters, c) && c.CheckDeclarationList(n1.Return, n2.Return, c)
}

// GoTypeAssertionExpressionEq compares GoTypeAssertionExpression-s.
func (g *GoChecker) GoTypeAssertionExpressionEq(n1 *mast.GoTypeAssertionExpression, n2 *mast.GoTypeAssertionExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Operand, n2.Operand, c) && c.CheckNode(n1.Type, n2.Type, c)
}

// GoTypeSwitchHeaderExpressionEq compares GoTypeSwitchHeaderExpression-s.
func (g *GoChecker) GoTypeSwitchHeaderExpressionEq(n1 *mast.GoTypeSwitchHeaderExpression, n2 *mast.GoTypeSwitchHeaderExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Alias, n2.Alias, c) && c.CheckNode(n1.Operand, n2.Operand, c)
}

// GoDeferStatementEq compares GoDeferStatement-s.
func (g *GoChecker) GoDeferStatementEq(n1 *mast.GoDeferStatement, n2 *mast.GoDeferStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Expr, n2.Expr, c)
}

// GoGotoStatementEq compares GoGotoStatement-s.
func (g *GoChecker) GoGotoStatementEq(n1 *mast.GoGotoStatement, n2 *mast.GoGotoStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Label, n2.Label, c)
}

// GoFallthroughStatementEq compares GoFallthroughStatement-s.
func (g *GoChecker) GoFallthroughStatementEq(n1 *mast.GoFallthroughStatement, n2 *mast.GoFallthroughStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return true
}

// GoSendStatementEq compares GoSendStatement-s.
func (g *GoChecker) GoSendStatementEq(n1 *mast.GoSendStatement, n2 *mast.GoSendStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Channel, n2.Channel, c) && c.CheckNode(n1.Value, n2.Value, c)
}

// GoGoStatementEq compares GoGoStatement-s.
func (g *GoChecker) GoGoStatementEq(n1 *mast.GoGoStatement, n2 *mast.GoGoStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Call, n2.Call, c)
}

// GoTypeDeclarationEq compares GoTypeDeclaration-s.
func (g *GoChecker) GoTypeDeclarationEq(n1 *mast.GoTypeDeclaration, n2 *mast.GoTypeDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.IsAlias == n2.IsAlias &&
		c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.Type, n2.Type, c)
}

// GoStructTypeEq compares GoStructType-s.
func (g *GoChecker) GoStructTypeEq(n1 *mast.GoStructType, n2 *mast.GoStructType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Declarations) != len(n2.Declarations) {
		return false
	}
	for i := 0; i < len(n1.Declarations); i++ {
		isEqual := c.CheckNode(n1.Declarations[i], n2.Declarations[i], c)
		if !isEqual {
			return false
		}
	}
	return true
}

// GoInterfaceTypeEq compares GoInterfaceType-s.
func (g *GoChecker) GoInterfaceTypeEq(n1 *mast.GoInterfaceType, n2 *mast.GoInterfaceType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDeclarationList(n1.Declarations, n2.Declarations, c)
}

// GoFieldDeclarationFieldsEq compares GoFieldDeclarationFields-s.
func (g *GoChecker) GoFieldDeclarationFieldsEq(n1 *mast.GoFieldDeclarationFields, n2 *mast.GoFieldDeclarationFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Tag, n2.Tag, c)
}

// GoForRangeStatementEq compares GoForRangeStatement-s.
func (g *GoChecker) GoForRangeStatementEq(n1 *mast.GoForRangeStatement, n2 *mast.GoForRangeStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Assignment, n2.Assignment, c) &&
		c.CheckNode(n1.Body, n2.Body, c) &&
		c.CheckNode(n1.Iterable, n2.Iterable, c)
}

// GoSelectStatementEq compares GoSelectStatement-s.
func (g *GoChecker) GoSelectStatementEq(n1 *mast.GoSelectStatement, n2 *mast.GoSelectStatement, c LangChecker) bool {
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
	return true
}

// GoCommunicationCaseEq compares GoCommunicationCase-s.
func (g *GoChecker) GoCommunicationCaseEq(n1 *mast.GoCommunicationCase, n2 *mast.GoCommunicationCase, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Communication, n2.Communication, c) && c.CheckStatementList(n1.Statements, n2.Statements, c)
}

// GoFunctionDeclarationFieldsEq compares GoFunctionDeclarationFields-s.
func (g *GoChecker) GoFunctionDeclarationFieldsEq(n1 *mast.GoFunctionDeclarationFields, n2 *mast.GoFunctionDeclarationFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Receiver, n2.Receiver, c)
}
