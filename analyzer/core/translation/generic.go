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
	"errors"
	"fmt"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"
)

// GenericTranslator handles the common translations of many common AST nodes that are shared across different languages.
// It is served as a fallback translator for the language-specific ones, so it will return an error if
// a node type is not handled.
type GenericTranslator struct {
	// Lang is a language-specific translator that generic translator will call back for recursion.
	Lang Translator
}

// Translate handles the translations of common nodes.
// It will fall back to language-specific translator Lang.Translate for recursion.
func (t *GenericTranslator) Translate(node *ts.Node) (mast.Node, error) {
	// delegate the translations to unexported methods
	switch nodeType := node.Type; nodeType {
	/* a few artificial nodes that are designed to share the logic across languages */
	case "list":
		return t.list(node)
	case "access_path":
		return t.accessPath(node)
	case "key_value_pair":
		return t.keyValuePair(node)
	case "literal_value":
		return t.literalValue(node)
	/* common nodes shared across languages */
	case "identifier":
		return t.identifier(node)
	case "root":
		return t.root(node)
	case "block":
		return t.block(node)
	case "parenthesized_expression":
		return t.parenthesizedExpression(node)
	case "unary_expression":
		return t.unaryExpression(node)
	case "binary_expression":
		return t.binaryExpression(node)
	case "index_expression":
		return t.indexExpression(node)
	case "assignment_expression":
		return t.assignmentExpression(node)
	case "null_literal":
		return t.nullLiteral(node)
	case "true", "false":
		return t.booleanLiteral(node)
	case "string_literal":
		return t.stringLiteral(node)
	case "int_literal":
		return t.intLiteral(node)
	case "float_literal":
		return t.floatLiteral(node)
	case "character_literal":
		return t.characterLiteral(node)
	case "update_expression":
		return t.updateExpression(node)
	case "continue_statement":
		return t.continueStatement(node)
	case "break_statement":
		return t.breakStatement(node)
	case "return_statement":
		return t.returnStatement(node)
	case "if_statement":
		return t.ifStatement(node)
	case "labeled_statement":
		return t.labeledStatement(node)
	default:
		return nil, fmt.Errorf("unsupported Node type %q", nodeType)
	}
}

// root translates root to mast.Root.
func (t *GenericTranslator) root(node *ts.Node) (mast.Node, error) {
	// translate children and return a mast.Root node
	// The root node may contain declaration nodes such as variable_declaration nodes that are translated to multiple
	// mast.VariableDeclaration nodes grouped by mast.TempGroupNode. Therefore here we un-group them.
	translated, err := translateNodes(t.Lang, node.Children, true /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	cast, err := toDeclarations(translated)
	if err != nil {
		return nil, err
	}
	result := &mast.Root{
		Declarations: cast,
	}
	return result, nil
}

// block translates block to mast.Block.
func (t *GenericTranslator) block(node *ts.Node) (mast.Node, error) {
	// block can have any number of children so we skip the check here.

	// translate and cast statements
	translated, err := translateNodes(t.Lang, node.Children, true /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	statements, err := toStatements(translated)
	if err != nil {
		return nil, err
	}

	result := &mast.Block{
		Statements: statements,
	}
	return result, nil
}

// identifier translates identifier / literal to mast.Identifier.
func (t *GenericTranslator) identifier(node *ts.Node) (mast.Node, error) {
	result := &mast.Identifier{
		Name: node.Name,
		Kind: mast.Blanket,
	}
	return result, nil
}

// parenthesizedExpression translates parenthesized_expression to mast.
func (t *GenericTranslator) parenthesizedExpression(node *ts.Node) (mast.Node, error) {
	// parenthesized_expression must have one child: expression (anonymous)
	if len(node.Children) != 1 {
		return nil, childrenNumberError(node)
	}
	translated, err := t.Lang.Translate(node.Children[0])
	if err != nil {
		return nil, err
	}
	expression, ok := translated.(mast.Expression)
	if !ok {
		return nil, nodeTypeError(translated)
	}
	result := &mast.ParenthesizedExpression{
		Expr: expression,
	}
	return result, nil
}

// unaryExpression translates unary_expression to mast.UnaryExpression.
func (t *GenericTranslator) unaryExpression(node *ts.Node) (mast.Node, error) {
	// unary_expression must have two child: operator, operand
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	operator := node.Children[0].Type

	translated, err := t.Lang.Translate(node.Children[1])
	if err != nil {
		return nil, err
	}
	expression, ok := translated.(mast.Expression)
	if !ok {
		return nil, err
	}

	result := &mast.UnaryExpression{
		Operator: operator,
		Expr:     expression,
	}
	return result, nil
}

// binaryExpression translates binary_expression to mast.BinaryExpression.
func (t *GenericTranslator) binaryExpression(node *ts.Node) (mast.Node, error) {
	// binary_expression must have three children: left, operator, right.
	if len(node.Children) != 3 {
		return nil, childrenNumberError(node)
	}
	operator := node.Children[1].Type

	nodes := []*ts.Node{node.Children[0], node.Children[2]}

	translated, err := translateNodes(t.Lang, nodes, false /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	expressions, err := toExpressions(translated)
	if err != nil {
		return nil, err
	}

	result := &mast.BinaryExpression{
		Operator: operator,
		Left:     expressions[0],
		Right:    expressions[1],
	}
	return result, nil
}

// indexExpression translates index_expression to mast.IndexExpression.
func (t *GenericTranslator) indexExpression(node *ts.Node) (mast.Node, error) {
	// index_expression must have two children: operand and index.
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	translated, err := translateNodes(t.Lang, node.Children, false /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	expressions, err := toExpressions(translated)
	if err != nil {
		return nil, err
	}

	result := &mast.IndexExpression{
		Operand: expressions[0],
		Index:   expressions[1],
	}
	return result, nil
}

// accessPath translates access_path to mast.accessPath.
func (t *GenericTranslator) accessPath(node *ts.Node) (mast.Node, error) {
	// access_path must have two children: operand and field.
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	translated, err := translateNodes(t.Lang, node.Children, false /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	expressions, err := toExpressions(translated)
	if err != nil {
		return nil, err
	}

	// the second child (field) must be an identifier node.
	field, ok := expressions[1].(*mast.Identifier)
	if !ok {
		return nil, nodeTypeError(expressions[1])
	}

	result := &mast.AccessPath{
		Operand: expressions[0],
		Field:   field,
	}
	return result, nil
}

// assignmentExpression translates assignment_expression to mast.AssignmentExpression.
func (t *GenericTranslator) assignmentExpression(node *ts.Node) (mast.Node, error) {
	// assignment_expression must have two children: left, right.
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	// store the generated slice of expressions for left and right
	expressions := [2][]mast.Expression{}
	for i, child := range node.Children {
		// left or right could be an expression_list which will generate a slice of expressions grouped by
		// mast.TempGroupNode, so here we un-group them.
		translated, err := translateNodes(t.Lang, []*ts.Node{child}, true /* shouldUngroup */)
		if err != nil {
			return nil, err
		}
		expressions[i], err = toExpressions(translated)
		if err != nil {
			return nil, err
		}
	}

	result := &mast.AssignmentExpression{
		Left:  expressions[0],
		Right: expressions[1],
	}

	return result, nil
}

// list translates *_list nodes to multiple nodes grouped by mast.TempGroupNode.
func (t *GenericTranslator) list(node *ts.Node) (mast.Node, error) {
	// *_list can have any number of children, so we skip the check here.
	// each element inside the list node could also be translated to multiple nodes (e.g., a field_declaration
	// could be translated to multiple GoFieldDeclaration nodes), therefore here we un-group them.
	translated, err := translateNodes(t.Lang, node.Children, true /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	result := &mast.TempGroupNode{
		Nodes: translated,
	}
	return result, nil
}

// nullLiteral translates nil / null_literal to mast.NullLiteral.
func (t *GenericTranslator) nullLiteral(node *ts.Node) (mast.Node, error) {
	if node.Name != "null" && node.Name != "nil" {
		return nil, fmt.Errorf("unexpected value for null literal node: %q", node.Name)
	}
	return &mast.NullLiteral{}, nil
}

// booleanLiteral translates true / false to mast.BooleanLiteral.
func (t *GenericTranslator) booleanLiteral(node *ts.Node) (mast.Node, error) {
	if node.Name != "true" && node.Name != "false" {
		return nil, fmt.Errorf("unexpected value for boolean literal node: %q", node.Name)
	}

	result := &mast.BooleanLiteral{
		Value: node.Name == "true",
	}

	return result, nil
}

// stringLiteral translates string_literal to mast.StringLiteral.
func (t *GenericTranslator) stringLiteral(node *ts.Node) (mast.Node, error) {
	result := &mast.StringLiteral{
		IsRaw: false,
		Value: node.Name,
	}
	return result, nil
}

// intLiteral translates int_literal to mast.IntLiteral.
func (t *GenericTranslator) intLiteral(node *ts.Node) (mast.Node, error) {
	result := &mast.IntLiteral{
		Value: node.Name,
	}
	return result, nil
}

// floatLiteral translates float_literal to mast.FloatLiteral.
func (t *GenericTranslator) floatLiteral(node *ts.Node) (mast.Node, error) {
	result := &mast.FloatLiteral{
		Value: node.Name,
	}
	return result, nil
}

// characterLiteral translates character_literal to mast.JavaCharacterLiteral.
func (t *GenericTranslator) characterLiteral(node *ts.Node) (mast.Node, error) {
	result := &mast.CharacterLiteral{
		Value: node.Name,
	}
	return result, nil
}

// updateExpression translates update_expression to mast.UpdateExpression.
func (t *GenericTranslator) updateExpression(node *ts.Node) (mast.Node, error) {
	// update_expression must have two children: operator, expression, but the position of the operator is not fixed.
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	// find the position of the operator to determine the operator side.
	var operatorSide mast.UpdateExpressionOperatorSide
	var operator, expression *ts.Node
	if node.Children[0].Type == "++" || node.Children[0].Type == "--" {
		operatorSide = mast.OperatorBefore
		operator, expression = node.Children[0], node.Children[1]
	} else {
		if node.Children[1].Type != "++" && node.Children[1].Type != "--" {
			return nil, errors.New("operator not found for update_expression node")
		}
		operatorSide = mast.OperatorAfter
		operator, expression = node.Children[1], node.Children[0]
	}

	// translate and cast the expression node
	translated, err := t.Lang.Translate(expression)
	if err != nil {
		return nil, err
	}
	castExpression, ok := translated.(mast.Expression)
	if !ok {
		return nil, nodeTypeError(translated)
	}

	result := &mast.UpdateExpression{
		OperatorSide: operatorSide,
		Operator:     operator.Type,
		Operand:      castExpression,
	}

	return result, nil
}

// continueStatement translates continue_statement to mast.ContinueStatement.
func (t *GenericTranslator) continueStatement(node *ts.Node) (mast.Node, error) {
	// continue_statement can have 0-1 children: label (optional, anonymous).
	if len(node.Children) > 1 {
		return nil, childrenNumberError(node)
	}

	result := &mast.ContinueStatement{}

	// return early if no label is present
	if len(node.Children) == 0 {
		return result, nil
	}

	// translate and cast the label
	translated, err := t.Lang.Translate(node.Children[0])
	if err != nil {
		return nil, err
	}
	label, ok := translated.(*mast.Identifier)
	if !ok {
		return nil, nodeTypeError(translated)
	}
	label.Kind = mast.Label
	result.Label = label

	return result, nil
}

// breakStatement translates break_statement to mast.BreakStatement.
func (t *GenericTranslator) breakStatement(node *ts.Node) (mast.Node, error) {
	// breakStatement can have 0-1 children: label (optional, anonymous).
	if len(node.Children) > 1 {
		return nil, childrenNumberError(node)
	}

	result := &mast.BreakStatement{}

	// return early if no label is present
	if len(node.Children) == 0 {
		return result, nil
	}

	// translate and cast the label
	translated, err := t.Lang.Translate(node.Children[0])
	if err != nil {
		return nil, err
	}
	label, ok := translated.(*mast.Identifier)
	if !ok {
		return nil, nodeTypeError(translated)
	}
	label.Kind = mast.Label
	result.Label = label

	return result, nil
}

// returnStatement translates return_statement to mast.ReturnStatement.
func (t *GenericTranslator) returnStatement(node *ts.Node) (mast.Node, error) {
	// return_statement can have 0-1 children: expression_list (or expression in Java, optional).
	if len(node.Children) > 1 {
		return nil, childrenNumberError(node)
	}

	result := &mast.ReturnStatement{}

	// early return if no expressions are present in the return statement
	if len(node.Children) == 0 {
		return result, nil
	}

	// we could have a list of expressions (e.g., expression_list node in Go) in Children that requires un-grouping.
	translated, err := translateNodes(t.Lang, node.Children, true /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	expressions, err := toExpressions(translated)
	if err != nil {
		return nil, err
	}

	result.Exprs = expressions

	return result, nil
}

// ifStatement translates if_statement to mast.IfStatement.
func (t *GenericTranslator) ifStatement(node *ts.Node) (mast.Node, error) {
	// if_statement must have 2-4 children: initializer (optional), condition, consequence, alternative (optional).
	if len(node.Children) < 2 || len(node.Children) > 4 {
		return nil, childrenNumberError(node)
	}

	initializer, condition := node.ChildByField("initializer"), node.ChildByField("condition")
	consequence, alternative := node.ChildByField("consequence"), node.ChildByField("alternative")
	// only initializer and alternative are optional
	if condition == nil || consequence == nil {
		return nil, nilChildError(node)
	}

	// If consequence is an empty body, we can directly set it to nil for better performance and simpler MAST structure.
	// That is, "if a {}" will have a "nil" Consequence field instead of mast.Block with empty Statements slice.
	if len(consequence.Children) == 0 {
		consequence = nil
	}

	// translate and cast the condition
	translatedCondition, err := t.Lang.Translate(condition)
	if err != nil {
		return nil, err
	}
	castCondition, ok := translatedCondition.(mast.Expression)
	if !ok {
		return nil, nodeTypeError(translatedCondition)
	}

	// translate and cast initializer, consequence and alternative
	translatedStatements, err := translateNodes(t.Lang, []*ts.Node{initializer, consequence, alternative}, false /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	castStatements, err := toStatements(translatedStatements)
	if err != nil {
		return nil, err
	}

	result := &mast.IfStatement{
		Initializer: castStatements[0],
		Condition:   castCondition,
		Consequence: castStatements[1],
		Alternative: castStatements[2],
	}

	return result, nil
}

// labeledStatement translates labeled_statement to a mast.LabelStatement or a mast.TempGroupNode grouping together (
// mast.LabelStatement, multiple other statements...).
func (t *GenericTranslator) labeledStatement(node *ts.Node) (mast.Node, error) {
	// labeled_statement must have 1 or 2 children: label, statement (optional, anonymous).
	// Tree-sitter attaches the nearest statement with the label itself in the labeled_statement node, for simpler tree
	// structure, we will split them and translate them to a mast.TempGroupNode grouping a mast.LabelStatement and the
	// statement for the upper-level translation to un-group. This will give us a flattened block of statements.
	// For example, consider the following code:
	//
	// L1: int a = 1;
	// L2: label: int b = 2, c = 3;
	// L3: int d = 4;
	//
	// Here tree-sitter will generate a labeled_statement representing the entire L2, which we translate it to
	// mast.TempGroupNode
	//   - mast.LabelStatement      ("label:")
	//   - mast.VariableDeclaration ("int b = 2;")
	//   - mast.VariableDeclaration ("int c = 3;")
	// Later in the upper level when the mast.TempGroupNode is un-grouped, we will have:
	// mast.Block
	//   - mast.VariableDeclaration ("int a = 1;")
	//   - mast.LabelStatement      ("label:")
	//   - mast.VariableDeclaration ("int b = 2;")
	//   - mast.VariableDeclaration ("int c = 3;")
	//   - mast.VariableDeclaration ("int d = 4;")
	if len(node.Children) != 1 && len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	// The attached statement might be translated to multiple statements (e.g., "int a, b;" will be translated to
	// "int a; int b;"), therefore we un-group it here.
	translated, err := translateNodes(t.Lang, node.Children, true /* shouldUngroup */)
	if err != nil {
		return nil, err
	}

	// The first node is the label.
	castLabel, ok := translated[0].(*mast.Identifier)
	if !ok {
		return nil, nodeTypeError(translated[0])
	}
	castLabel.Kind = mast.Label
	label := &mast.LabelStatement{Label: castLabel}

	// early return if it is an empty label without a statement
	if len(translated) == 1 {
		return label, nil
	}

	result := &mast.TempGroupNode{
		Nodes: []mast.Node{label},
	}
	// Go AST could generate an expression node instead of statement nodes, so here we call toStatements to
	// automatically handle the transformations to mast.ExpressionStatement. See toStatements for additional
	// explanations.
	castStatements, err := toStatements(translated[1:])
	if err != nil {
		return nil, err
	}
	for _, stmt := range castStatements {
		result.Nodes = append(result.Nodes, stmt)
	}

	return result, nil
}

// keyValuePair translates key_value_pair to mast.KeyValuePair.
func (t *GenericTranslator) keyValuePair(node *ts.Node) (mast.Node, error) {
	// key_value_pair must have two children: key, value.
	if len(node.Children) != 2 {
		return nil, childrenNumberError(node)
	}

	// translate and cast key and value
	translated, err := translateNodes(t.Lang, node.Children, false /* shouldUngroup */)
	if err != nil {
		return nil, err
	}
	expressions, err := toExpressions(translated)
	if err != nil {
		return nil, err
	}

	result := &mast.KeyValuePair{
		Key:   expressions[0],
		Value: expressions[1],
	}
	return result, nil
}

// literalValue translates literal_value to mast.LiteralValue.
func (t *GenericTranslator) literalValue(node *ts.Node) (mast.Node, error) {
	// literal_value can have any number of children, so we skip the check here.

	// early return if it is an empty literal_value
	if len(node.Children) == 0 {
		result := &mast.LiteralValue{}
		return result, nil
	}

	// The translations is the same as list: translate all children and get result in mast.TempGroupNode.
	node.Type = "list"
	// Here we do not need to call language-specific translator since "list" is handled by GenericTranslator already.
	translated, err := t.Translate(node)
	if err != nil {
		return nil, err
	}
	cast, ok := translated.(*mast.TempGroupNode)
	if !ok {
		return nil, nodeTypeError(translated)
	}
	expressions, err := toExpressions(cast.Nodes)
	if err != nil {
		return nil, err
	}

	result := &mast.LiteralValue{
		Values: expressions,
	}

	return result, nil
}
