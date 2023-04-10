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

// JavaChecker represents a Java-specific equivalence checker.
type JavaChecker struct {
	Checker

	// LoggingOn indicates whether logging-related auto-approvals should be enabled.
	LoggingOn bool
	// importedLoggerClass stores the imported logger class name (e.g., "Logger") from supported
	// logging frameworks (e.g., slf4j and log4j) in _both_ diffs. If two diffs imported different
	// logger classes, this will be empty.
	importedLoggerClass string
}

// NewJavaChecker created and returns Java-specific equivalence
// checker.
func NewJavaChecker(checker Checker, loggingOn bool) *JavaChecker {
	return &JavaChecker{Checker: checker, LoggingOn: loggingOn}
}

// StmtIgnored checks if a given statement can be ignored when
// checking for equivalence. If it can be ignored altogether, this
// method returns true. If it can be ignored under the condition that
// some of its sub-expressions are equal, this method returns false
// and a non-nil list of sub-expressions. If it cannot be ignored
// under any circumstances, this method returns false and an empty (or
// nil) list of sub-expressions.
func (j *JavaChecker) StmtIgnored(stmt mast.Statement, symbols *symbolication.SymbolTable) (bool, []mast.Expression) {
	if !j.LoggingOn {
		return false, nil
	}
	return j.loggingStmtIgnore(stmt, symbols)
}

// CheckNode compares individual nodes.
func (j *JavaChecker) CheckNode(node mast.Node, other mast.Node, c LangChecker) bool {
	if baseRoot, ok := node.(*mast.Root); ok {
		if lastRoot, ok := other.(*mast.Root); ok {
			baseLogger, lastLogger := j.findImportedLogger(baseRoot), j.findImportedLogger(lastRoot)
			if baseLogger == lastLogger {
				j.importedLoggerClass = baseLogger
			}

		} else {
			// should not happen
			return false
		}
	}

	switch n1 := node.(type) {
	case *mast.JavaVariableDeclarator:
		if n2, ok := other.(*mast.JavaVariableDeclarator); ok {
			return j.JavaVariableDeclaratorEq(n1, n2, c)

		}
	case *mast.JavaTernaryExpression:
		if n2, ok := other.(*mast.JavaTernaryExpression); ok {
			return j.JavaTernaryExpressionEq(n1, n2, c)
		}
	case *mast.JavaAnnotatedType:
		if n2, ok := other.(*mast.JavaAnnotatedType); ok {
			return j.JavaAnnotatedTypeEq(n1, n2, c)
		}
	case *mast.JavaGenericType:
		if n2, ok := other.(*mast.JavaGenericType); ok {
			return j.JavaGenericTypeEq(n1, n2, c)
		}
	case *mast.JavaWildcard:
		if n2, ok := other.(*mast.JavaWildcard); ok {
			return j.JavaWildcardEq(n1, n2, c)
		}
	case *mast.JavaArrayType:
		if n2, ok := other.(*mast.JavaArrayType); ok {
			return j.JavaArrayTypeEq(n1, n2, c)
		}
	case *mast.JavaDimension:
		if n2, ok := other.(*mast.JavaDimension); ok {
			return j.JavaDimensionEq(n1, n2, c)
		}
	case *mast.JavaInstanceOfExpression:
		if n2, ok := other.(*mast.JavaInstanceOfExpression); ok {
			return j.JavaInstanceOfExpressionEq(n1, n2, c)
		}
	case *mast.JavaLiteralModifier:
		if n2, ok := other.(*mast.JavaLiteralModifier); ok {
			return j.JavaLiteralModifierEq(n1, n2, c)
		}
	case *mast.JavaTryStatement:
		if n2, ok := other.(*mast.JavaTryStatement); ok {
			return j.JavaTryStatementEq(n1, n2, c)
		}
	case *mast.JavaCatchClause:
		if n2, ok := other.(*mast.JavaCatchClause); ok {
			return j.JavaCatchClauseEq(n1, n2, c)
		}
	case *mast.JavaCatchFormalParameter:
		if n2, ok := other.(*mast.JavaCatchFormalParameter); ok {
			return j.JavaCatchFormalParameterEq(n1, n2, c)
		}
	case *mast.JavaFinallyClause:
		if n2, ok := other.(*mast.JavaFinallyClause); ok {
			return j.JavaFinallyClauseEq(n1, n2, c)
		}
	case *mast.JavaWhileStatement:
		if n2, ok := other.(*mast.JavaWhileStatement); ok {
			return j.JavaWhileStatementEq(n1, n2, c)
		}
	case *mast.JavaThrowStatement:
		if n2, ok := other.(*mast.JavaThrowStatement); ok {
			return j.JavaThrowStatementEq(n1, n2, c)
		}
	case *mast.JavaAssertStatement:
		if n2, ok := other.(*mast.JavaAssertStatement); ok {
			return j.JavaAssertStatementEq(n1, n2, c)
		}
	case *mast.JavaSynchronizedStatement:
		if n2, ok := other.(*mast.JavaSynchronizedStatement); ok {
			return j.JavaSynchronizedStatementEq(n1, n2, c)
		}
	case *mast.JavaDoStatement:
		if n2, ok := other.(*mast.JavaDoStatement); ok {
			return j.JavaDoStatementEq(n1, n2, c)
		}
	case *mast.JavaParameterDeclarationFields:
		if n2, ok := other.(*mast.JavaParameterDeclarationFields); ok {
			return j.JavaParameterDeclarationFieldsEq(n1, n2, c)
		}
	case *mast.JavaEnhancedForStatement:
		if n2, ok := other.(*mast.JavaEnhancedForStatement); ok {
			return j.JavaEnhancedForStatementEq(n1, n2, c)
		}
	case *mast.JavaModuleDeclaration:
		if n2, ok := other.(*mast.JavaModuleDeclaration); ok {
			return j.JavaModuleDeclarationEq(n1, n2, c)
		}
	case *mast.JavaModuleDirective:
		if n2, ok := other.(*mast.JavaModuleDirective); ok {
			return j.JavaModuleDirectiveEq(n1, n2, c)
		}
	case *mast.JavaTypeParameter:
		if n2, ok := other.(*mast.JavaTypeParameter); ok {
			return j.JavaTypeParameterEq(n1, n2, c)
		}
	case *mast.JavaClassDeclaration:
		if n2, ok := other.(*mast.JavaClassDeclaration); ok {
			return j.JavaClassDeclarationEq(n1, n2, c)
		}
	case *mast.JavaInterfaceDeclaration:
		if n2, ok := other.(*mast.JavaInterfaceDeclaration); ok {
			return j.JavaInterfaceDeclarationEq(n1, n2, c)
		}
	case *mast.JavaEnumDeclaration:
		if n2, ok := other.(*mast.JavaEnumDeclaration); ok {
			return j.JavaEnumDeclarationEq(n1, n2, c)
		}
	case *mast.JavaEnumConstantDeclaration:
		if n2, ok := other.(*mast.JavaEnumConstantDeclaration); ok {
			return j.JavaEnumConstantDeclarationEq(n1, n2, c)
		}
	case *mast.JavaClassInitializer:
		if n2, ok := other.(*mast.JavaClassInitializer); ok {
			return j.JavaClassInitializerEq(n1, n2, c)
		}
	case *mast.JavaFunctionDeclarationFields:
		if n2, ok := other.(*mast.JavaFunctionDeclarationFields); ok {
			return j.JavaFunctionDeclarationFieldsEq(n1, n2, c)
		}
	case *mast.JavaAnnotationDeclaration:
		if n2, ok := other.(*mast.JavaAnnotationDeclaration); ok {
			return j.JavaAnnotationDeclarationEq(n1, n2, c)
		}
	case *mast.JavaAnnotationElementDeclaration:
		if n2, ok := other.(*mast.JavaAnnotationElementDeclaration); ok {
			return j.JavaAnnotationElementDeclarationEq(n1, n2, c)
		}
	case *mast.JavaMethodReference:
		if n2, ok := other.(*mast.JavaMethodReference); ok {
			return j.JavaMethodReferenceEq(n1, n2, c)
		}
	case *mast.JavaClassLiteral:
		if n2, ok := other.(*mast.JavaClassLiteral); ok {
			return j.JavaClassLiteralEq(n1, n2, c)
		}
	case *mast.JavaEntityCreationExpressionFields:
		if n2, ok := other.(*mast.JavaEntityCreationExpressionFields); ok {
			return j.JavaEntityCreationExpressionFieldsEq(n1, n2, c)
		}
	case *mast.JavaVariableDeclarationFields:
		if n2, ok := other.(*mast.JavaVariableDeclarationFields); ok {
			return j.JavaVariableDeclarationFieldsEq(n1, n2, c)
		}
	case *mast.JavaCallExpressionFields:
		if n2, ok := other.(*mast.JavaCallExpressionFields); ok {
			return j.JavaCallExpressionFieldsEq(n1, n2, c)
		}
	}
	return j.Checker.CheckNode(node, other, c)
}

// JavaVariableDeclaratorEq compares JavaVariableDeclarator-s.
func (j *JavaChecker) JavaVariableDeclaratorEq(n1 *mast.JavaVariableDeclarator, n2 *mast.JavaVariableDeclarator, c LangChecker) bool {
	// always return false for JavaVariableDeclarator regardless of
	// the other nodesince it should never appear in the final MAST.
	return false
}

// JavaTernaryExpressionEq compares JavaTernaryExpression-s.
func (j *JavaChecker) JavaTernaryExpressionEq(n1 *mast.JavaTernaryExpression, n2 *mast.JavaTernaryExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Condition, n2.Condition, c) &&
		c.CheckNode(n1.Consequence, n2.Consequence, c) &&
		c.CheckNode(n1.Alternative, n2.Alternative, c)
}

// JavaAnnotatedTypeEq compares JavaAnnotatedType-s.
func (j *JavaChecker) JavaAnnotatedTypeEq(n1 *mast.JavaAnnotatedType, n2 *mast.JavaAnnotatedType, c LangChecker) bool {
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
	return c.CheckNode(n1.Type, n2.Type, c)
}

// JavaGenericTypeEq compares JavaGenericType-s.
func (j *JavaChecker) JavaGenericTypeEq(n1 *mast.JavaGenericType, n2 *mast.JavaGenericType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) && c.CheckExpressionList(n1.Arguments, n2.Arguments, c)
}

// JavaWildcardEq compares JavaWildcard-s.
func (j *JavaChecker) JavaWildcardEq(n1 *mast.JavaWildcard, n2 *mast.JavaWildcard, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Super, n2.Super, c) && c.CheckNode(n1.Extends, n2.Extends, c)
}

// JavaArrayTypeEq compares JavaArrayType-s.
func (j *JavaChecker) JavaArrayTypeEq(n1 *mast.JavaArrayType, n2 *mast.JavaArrayType, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDimensions(n1.Dimensions, n2.Dimensions, c) && c.CheckNode(n1.Name, n2.Name, c)
}

// JavaDimensionEq compares JavaDimension-s.
func (j *JavaChecker) JavaDimensionEq(n1 *mast.JavaDimension, n2 *mast.JavaDimension, c LangChecker) bool {
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
	return c.CheckNode(n1.Length, n2.Length, c)
}

// JavaInstanceOfExpressionEq compares JavaInstanceOfExpression-s.
func (j *JavaChecker) JavaInstanceOfExpressionEq(n1 *mast.JavaInstanceOfExpression, n2 *mast.JavaInstanceOfExpression, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Operand, n2.Operand, c) && c.CheckNode(n1.Type, n2.Type, c)
}

// JavaLiteralModifierEq compares JavaLiteralModifier-s.
func (j *JavaChecker) JavaLiteralModifierEq(n1 *mast.JavaLiteralModifier, n2 *mast.JavaLiteralModifier, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Modifier == n2.Modifier
}

// JavaTryStatementEq compares JavaTryStatement-s.
func (j *JavaChecker) JavaTryStatementEq(n1 *mast.JavaTryStatement, n2 *mast.JavaTryStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.CatchClauses) != len(n2.CatchClauses) {
		return false
	}
	for i := 0; i < len(n1.CatchClauses); i++ {
		isEqual := c.CheckNode(n1.CatchClauses[i], n2.CatchClauses[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Body, n2.Body, c) &&
		c.CheckNode(n1.Finally, n2.Finally, c) &&
		c.CheckStatementList(n1.Resources, n2.Resources, c)
}

// JavaCatchClauseEq compares JavaCatchClause-s.
func (j *JavaChecker) JavaCatchClauseEq(n1 *mast.JavaCatchClause, n2 *mast.JavaCatchClause, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Parameter, n2.Parameter, c) && c.CheckNode(n1.Body, n2.Body, c)
}

// JavaCatchFormalParameterEq compares JavaCatchFormalParameter-s.
func (j *JavaChecker) JavaCatchFormalParameterEq(n1 *mast.JavaCatchFormalParameter, n2 *mast.JavaCatchFormalParameter, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Types, n2.Types, c)
}

// JavaFinallyClauseEq compares JavaFinallyClause-s.
func (j *JavaChecker) JavaFinallyClauseEq(n1 *mast.JavaFinallyClause, n2 *mast.JavaFinallyClause, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Body, n2.Body, c)
}

// JavaWhileStatementEq compares JavaWhileStatement-s.
func (j *JavaChecker) JavaWhileStatementEq(n1 *mast.JavaWhileStatement, n2 *mast.JavaWhileStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Condition, n2.Condition, c) &&
		c.CheckNode(n1.Body, n2.Body, c)
}

// JavaThrowStatementEq compares JavaThrowStatement-s.
func (j *JavaChecker) JavaThrowStatementEq(n1 *mast.JavaThrowStatement, n2 *mast.JavaThrowStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Expr, n2.Expr, c)
}

// JavaAssertStatementEq compares JavaAssertStatement-s.
func (j *JavaChecker) JavaAssertStatementEq(n1 *mast.JavaAssertStatement, n2 *mast.JavaAssertStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Condition, n2.Condition, c) && c.CheckNode(n1.ErrorString, n2.ErrorString, c)
}

// JavaSynchronizedStatementEq compares JavaSynchronizedStatement-s.
func (j *JavaChecker) JavaSynchronizedStatementEq(n1 *mast.JavaSynchronizedStatement, n2 *mast.JavaSynchronizedStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Body, n2.Body, c) && c.CheckNode(n1.Expr, n2.Expr, c)
}

// JavaDoStatementEq compares JavaDoStatement-s.
func (j *JavaChecker) JavaDoStatementEq(n1 *mast.JavaDoStatement, n2 *mast.JavaDoStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Body, n2.Body, c) && c.CheckNode(n1.Condition, n2.Condition, c)
}

// JavaParameterDeclarationFieldsEq compares JavaParameterDeclarationFields-s.
func (j *JavaChecker) JavaParameterDeclarationFieldsEq(n1 *mast.JavaParameterDeclarationFields, n2 *mast.JavaParameterDeclarationFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDimensions(n1.Dimensions, n2.Dimensions, c) &&
		n1.IsReceiver == n2.IsReceiver &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c)
}

// JavaEnhancedForStatementEq compares JavaEnhancedForStatement-s.
func (j *JavaChecker) JavaEnhancedForStatementEq(n1 *mast.JavaEnhancedForStatement, n2 *mast.JavaEnhancedForStatement, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDimensions(n1.Dimensions, n2.Dimensions, c) &&
		c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.Iterable, n2.Iterable, c) &&
		c.CheckNode(n1.Body, n2.Body, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c)
}

// JavaModuleDeclarationEq compares JavaModuleDeclaration-s.
func (j *JavaChecker) JavaModuleDeclarationEq(n1 *mast.JavaModuleDeclaration, n2 *mast.JavaModuleDeclaration, c LangChecker) bool {
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
	if len(n1.Directives) != len(n2.Directives) {
		return false
	}
	for i := 0; i < len(n1.Directives); i++ {
		isEqual := c.CheckNode(n1.Directives[i], n2.Directives[i], c)
		if !isEqual {
			return false
		}
	}
	return n1.IsOpen == n2.IsOpen && c.CheckNode(n1.Name, n2.Name, c)
}

// JavaModuleDirectiveEq compares JavaModuleDirective-s.
func (j *JavaChecker) JavaModuleDirectiveEq(n1 *mast.JavaModuleDirective, n2 *mast.JavaModuleDirective, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.Keyword == n2.Keyword && c.CheckExpressionList(n1.Exprs, n2.Exprs, c)
}

// JavaTypeParameterEq compares JavaTypeParameter-s.
func (j *JavaChecker) JavaTypeParameterEq(n1 *mast.JavaTypeParameter, n2 *mast.JavaTypeParameter, c LangChecker) bool {
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
	return c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckExpressionList(n1.Extends, n2.Extends, c)
}

// JavaClassDeclarationEq compares JavaClassDeclaration-s.
func (j *JavaChecker) JavaClassDeclarationEq(n1 *mast.JavaClassDeclaration, n2 *mast.JavaClassDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.TypeParameters) != len(n2.TypeParameters) {
		return false
	}
	for i := 0; i < len(n1.TypeParameters); i++ {
		isEqual := c.CheckNode(n1.TypeParameters[i], n2.TypeParameters[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.SuperClass, n2.SuperClass, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Interfaces, n2.Interfaces, c) &&
		c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaInterfaceDeclarationEq compares JavaInterfaceDeclaration-s.
func (j *JavaChecker) JavaInterfaceDeclarationEq(n1 *mast.JavaInterfaceDeclaration, n2 *mast.JavaInterfaceDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.TypeParameters) != len(n2.TypeParameters) {
		return false
	}
	for i := 0; i < len(n1.TypeParameters); i++ {
		isEqual := c.CheckNode(n1.TypeParameters[i], n2.TypeParameters[i], c)
		if !isEqual {
			return false
		}
	}
	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Extends, n2.Extends, c) &&
		c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaEnumDeclarationEq compares JavaEnumDeclaration-s.
func (j *JavaChecker) JavaEnumDeclarationEq(n1 *mast.JavaEnumDeclaration, n2 *mast.JavaEnumDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Interfaces, n2.Interfaces, c) &&
		c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaEnumConstantDeclarationEq compares JavaEnumConstantDeclaration-s.
func (j *JavaChecker) JavaEnumConstantDeclarationEq(n1 *mast.JavaEnumConstantDeclaration, n2 *mast.JavaEnumConstantDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Arguments, n2.Arguments, c) &&
		c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaClassInitializerEq compares JavaClassInitializer-s.
func (j *JavaChecker) JavaClassInitializerEq(n1 *mast.JavaClassInitializer, n2 *mast.JavaClassInitializer, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return n1.IsStatic == n2.IsStatic && c.CheckNode(n1.Block, n2.Block, c)
}

// JavaFunctionDeclarationFieldsEq compares JavaFunctionDeclarationFields-s.
func (j *JavaChecker) JavaFunctionDeclarationFieldsEq(n1 *mast.JavaFunctionDeclarationFields, n2 *mast.JavaFunctionDeclarationFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.TypeParameters) != len(n2.TypeParameters) {
		return false
	}
	for i := 0; i < len(n1.TypeParameters); i++ {
		isEqual := c.CheckNode(n1.TypeParameters[i], n2.TypeParameters[i], c)
		if !isEqual {
			return false
		}
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
	return c.CheckDimensions(n1.Dimensions, n2.Dimensions, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckExpressionList(n1.Throws, n2.Throws, c)
}

// JavaCallExpressionFieldsEq compares JavaCallExpressionFields nodes.
func (j *JavaChecker) JavaCallExpressionFieldsEq(n1 *mast.JavaCallExpressionFields, n2 *mast.JavaCallExpressionFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckExpressionList(n1.TypeArguments, n2.TypeArguments, c)
}

// JavaAnnotationDeclarationEq compares JavaAnnotationDeclaration-s.
func (j *JavaChecker) JavaAnnotationDeclarationEq(n1 *mast.JavaAnnotationDeclaration, n2 *mast.JavaAnnotationDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c) &&
		c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaAnnotationElementDeclarationEq compares JavaAnnotationElementDeclaration-s.
func (j *JavaChecker) JavaAnnotationElementDeclarationEq(n1 *mast.JavaAnnotationElementDeclaration, n2 *mast.JavaAnnotationElementDeclaration, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckDimensions(n1.Dimensions, n2.Dimensions, c) &&
		c.CheckNode(n1.Name, n2.Name, c) &&
		c.CheckNode(n1.Type, n2.Type, c) &&
		c.CheckNode(n1.Value, n2.Value, c) &&
		c.CheckExpressionList(n1.Modifiers, n2.Modifiers, c)
}

// JavaMethodReferenceEq compares JavaMethodReference-s.
func (j *JavaChecker) JavaMethodReferenceEq(n1 *mast.JavaMethodReference, n2 *mast.JavaMethodReference, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Method, n2.Method, c) &&
		c.CheckNode(n1.Class, n2.Class, c) &&
		c.CheckExpressionList(n1.TypeArguments, n2.TypeArguments, c)
}

// JavaClassLiteralEq compares JavaClassLiteral-s.
func (j *JavaChecker) JavaClassLiteralEq(n1 *mast.JavaClassLiteral, n2 *mast.JavaClassLiteral, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	return c.CheckNode(n1.Type, n2.Type, c)
}

// JavaEntityCreationExpressionFieldsEq compares JavaEntityCreationExpressionFields-s.
func (j *JavaChecker) JavaEntityCreationExpressionFieldsEq(n1 *mast.JavaEntityCreationExpressionFields, n2 *mast.JavaEntityCreationExpressionFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	if len(n1.Dimensions) != len(n2.Dimensions) {
		return false
	}
	for i := 0; i < len(n1.Dimensions); i++ {
		isEqual := c.CheckNode(n1.Dimensions[i], n2.Dimensions[i], c)
		if !isEqual {
			return false
		}
	}

	return c.CheckDeclarationList(n1.Body, n2.Body, c)
}

// JavaVariableDeclarationFieldsEq compares JavaVariableDeclarationFields-s.
func (j *JavaChecker) JavaVariableDeclarationFieldsEq(n1 *mast.JavaVariableDeclarationFields, n2 *mast.JavaVariableDeclarationFields, c LangChecker) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}

	// We allow _additions_ of the "final" modifiers (removals should be considered unsafe).
	// Moreover, the order of the modifiers actually do not matter in terms of functionalities.
	// So here we do some additional processing to approve such cases.

	// We will first build a set of all literal modifiers in n1. During the iteration, we store
	// other types of modifier nodes in a slice to be checked later.
	var others1 []mast.Expression
	literalModifiers := make(map[string]bool)
	for _, modifier := range n1.Modifiers {
		m, ok := modifier.(*mast.JavaLiteralModifier)
		if !ok {
			others1 = append(others1, modifier)
			continue
		}
		literalModifiers[m.Modifier] = true
	}

	// When we iterate all modifiers in n2, we remove the found modifiers in the set. If a modifier
	// is not found in the set, n2 is adding literal modifiers. We should report unsafe unless it
	// is adding a "final" modifier. Again, we store other modifiers in a slice to be checked later.
	var others2 []mast.Expression
	for _, modifier := range n2.Modifiers {
		m, ok := modifier.(*mast.JavaLiteralModifier)
		if !ok {
			others2 = append(others2, modifier)
			continue
		}
		// Remove the modifier if found in the set.
		if _, ok := literalModifiers[m.Modifier]; ok {
			delete(literalModifiers, m.Modifier)
			continue
		}
		// Otherwise, it is not ok to add modifiers other than "final".
		if m.Modifier != mast.FinalMod {
			return false
		}
	}
	// If there are still literal modifiers in the set, it means n2 is missing some literal
	// modifiers that appear in n1. This should be considered unsafe.
	if len(literalModifiers) != 0 {
		return false
	}

	// Literal modifiers have been checked above, so here we only check the remainders.
	return c.CheckExpressionList(others1, others2, c) && c.CheckDimensions(n1.Dimensions, n2.Dimensions, c)
}
