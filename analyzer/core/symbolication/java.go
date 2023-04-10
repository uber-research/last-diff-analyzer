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

package symbolication

import (
	"fmt"

	"analyzer/core/mast"
	"analyzer/core/mast/mastutil"
)

// accessPathKeywords is a set of special keywords that can be present
// on access path (this, super, class), and which determine kind of
// identifiers on the class path.
var _accessPathKeywords = map[string]bool{"this": true, "super": true, "class": true}

// JavaSymbolTableBuilder is a language-specific symbol table builder for Java. All internal datastructures are shared
// with GenericSymbolTableBuilder so that both of them can modify the same objects.
type JavaSymbolTableBuilder struct {
	// generic is the pointer to the generic symbol table builder.
	*GenericSymbolTableBuilder
}

// NewJavaSymbolTableBuilder returns a properly initialized JavaSymbolTableBuilder.
func NewJavaSymbolTableBuilder() *JavaSymbolTableBuilder {
	// first create a Java-specific builder with nil GenericSymbolTableBuilder field
	builder := &JavaSymbolTableBuilder{}
	// assign the filed with properly initialized generic builder
	builder.GenericSymbolTableBuilder = NewGenericSymbolTableBuilder(builder)
	return builder
}

// SymbolTable returns the built symbol table.
func (j *JavaSymbolTableBuilder) SymbolTable() *SymbolTable {
	return j.symbolTable
}

// IsDeclarationPrivate checks whether a declaration is private in Java using the list of modifiers.
func (j *JavaSymbolTableBuilder) IsDeclarationPrivate(node mast.Node, identifier *mast.Identifier) (bool, error) {
	// We determine the privateness of each declaration based on its modifiers.
	var modifiers []mast.Expression

	// retrieve the modifiers for different declarations
	switch n := node.(type) {
	case *mast.JavaEnumDeclaration:
		modifiers = n.Modifiers

	case *mast.JavaEnumConstantDeclaration:
		modifiers = n.Modifiers

	case *mast.JavaAnnotationDeclaration:
		modifiers = n.Modifiers

	case *mast.JavaAnnotationElementDeclaration:
		modifiers = n.Modifiers

	case *mast.JavaInterfaceDeclaration:
		modifiers = n.Modifiers

	case *mast.JavaClassDeclaration:
		modifiers = n.Modifiers

	case *mast.FunctionDeclaration:
		langFields, ok := n.LangFields.(*mast.JavaFunctionDeclarationFields)
		if !ok {
			return false, langFieldsNodeNotExistsError(node)
		}
		modifiers = langFields.Modifiers

	case *mast.VariableDeclaration:
		langFields, ok := n.LangFields.(*mast.JavaVariableDeclarationFields)
		if !ok {
			return false, langFieldsNodeNotExistsError(node)
		}
		modifiers = langFields.Modifiers

	// conservatively return false for unhandled cases
	default:
		return false, nil
	}

	// find the privateness based on the modifiers
	for _, modifier := range modifiers {
		switch n := modifier.(type) {
		case *mast.JavaLiteralModifier:
			// Since Java allows multiple levels of visibility, including public, protected, private and default
			// package visibility (without visibility modifiers). Here we conservatively look for private modifier and
			// only return true if it is present.
			if n.Modifier == mast.PrivateMod {
				return true, nil
			}
		case *mast.Annotation:
			continue
		default:
			return false, fmt.Errorf("unexpected node type %T in Java modifiers", modifier)
		}
	}

	// default to public declaration
	return false, nil
}

// ProcessDeclaration handles Java-specific nodes and delegates the handlings of other nodes to the generic builder.
func (j *JavaSymbolTableBuilder) ProcessDeclaration(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if j.ignoreNodes[node] {
		return nil
	}

	var identifier *mast.Identifier

	switch n := node.(type) {
	case *mast.ImportDeclaration:
		// none of the import declaration identifiers should (at least
		// at this point) be a target of identifier resolution
		if err := j.nilifyNodeIDs(n.Package); err != nil {
			return err
		}
		if n.Alias != nil {
			return fmt.Errorf("unexpected alias in Java import declaration: %T", n)
		}
		return nil

	case *mast.JavaEnumDeclaration:
		identifier = n.Name

	case *mast.JavaEnumConstantDeclaration:
		identifier = n.Name

	case *mast.JavaAnnotationDeclaration:
		identifier = n.Name

	case *mast.JavaInterfaceDeclaration:
		identifier = n.Name

	case *mast.JavaClassDeclaration:
		identifier = n.Name

	case *mast.FieldDeclaration:
		identifier = n.Name

	// fallback to generic implementation for other declaration nodes
	default:
		return j.GenericSymbolTableBuilder.ProcessDeclaration(node)
	}

	_, err := j.createDeclarationEntry(node, identifier, true)
	return err
}

// ProcessScope handles Java-specific scope-related nodes and delagates the handlings of other nodes to the generic
// builder.
func (j *JavaSymbolTableBuilder) ProcessScope(node mast.Node, onEnter bool) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if j.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.JavaClassDeclaration:
		// mast.JavaClassDeclaration is not handled by the
		// processTypeScope (used to handle "generic" cases) because
		// we we need to additinally handle the "this" scope here.

		// special case for exiting the scope
		if !onEnter {
			// clear the scope for "this" keyword on exit
			err := j.scopes.ClearThis()
			if err != nil {
				return err
			}
			return j.scopes.PopScope()
		}

		// first retrieve the symbol entry for the class declaration for later use
		entry, err := j.scopes.FindDeclarationEntry(n.Name, CurrentOnly, false /* activeOnly */)
		if err != nil {
			return err
		}
		if entry == nil {
			return entryNotExistError(node, n.Name)
		}

		// The privateness of the scope for the class will be determined by the privateness of the declaration itself.
		j.scopes.CreateNewScope(entry.IsPrivate /* isPrivate */)

		// mark the scope for "this" scope on enter
		return j.scopes.MarkThis()

	case *mast.JavaInterfaceDeclaration:
		return processTypeScope(j.scopes, n.Name, node, onEnter)
	case *mast.JavaAnnotationDeclaration:
		return processTypeScope(j.scopes, n.Name, node, onEnter)
	case *mast.JavaModuleDeclaration:
		// always mark module as a non-private scope
		return j.handleScope(onEnter, false /* isPrivate */)

	case *mast.EntityCreationExpression:
		if n.LangFields != nil {
			langFields, ok := n.LangFields.(*mast.JavaEntityCreationExpressionFields)
			if !ok {
				return langFieldsNodeNotExistsError(node)
			}
			// If the class body for the entity creation expression is present, for example, "new A() {...}", we should
			// create a new private scope for the declarations.
			if langFields.Body != nil {
				return j.handleScope(onEnter, true /* isPrivate */)
			}
		}
		return nil

	default:
		return j.GenericSymbolTableBuilder.ProcessScope(node, onEnter)
	}
}

// ProcessUse delegates the handling of mast.Identifier to the generic builder.
func (j *JavaSymbolTableBuilder) ProcessUse(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if j.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	// We should handle mast.AccessPath a little different from the generic one,
	// since Java allows the use of the "this" keyword. For simplicity, we only
	// try to resolve the first component _after_ the "this" keyword.
	case *mast.AccessPath:
		// _accessPathKeywords map argument is used to find out if the
		// keywords in the map are part of the access path, in which
		// case a prefix of the path preceding such keyword will be
		// returned.
		identifierChain, ignoreFirst, foundPrefix, err := mastutil.ExtractAccessPath(n, _accessPathKeywords /* searchedIDs */)
		if err != nil {
			return err
		}

		if foundPrefix != nil {
			// We found a prefix of the access path that precedes one
			// of the special keywords. This indicates that the prefix
			// represents a type (qualified such as Clazz or
			// unqualified such as tmp.Clazz) and identifiers of this
			// type should have their kinds set accordingly. The
			// translation.setExprTypeKinds function will additionally
			// do a sanity check to verify that the prefix consists
			// only of identifiers (which should be the case if it
			// precedes the special keyword).
			if err := mast.SetJavaExprTypeKinds(foundPrefix); err != nil {
				return err
			}
		}

		// unresolvedIndex marks the start index in the identifier chain that will be linked to nil declaration pointers.
		unresolvedIndex := 0

		// If the first identifier in the chain is not "this", we handle the chain in the normal way. Otherwise, we
		// resolve the first identifier after the "this" identifier, and link the rest nodes in the chain to nil.
		if identifierChain[0].Name != "this" {
			if !ignoreFirst {
				// see explanation in how values of unresolvedIndex
				// are set when handling AccessPath in the generic
				// part of the symbolicator (generic.go)
				err := j.ProcessUse(identifierChain[0])
				if err != nil {
					return err
				}
				unresolvedIndex = 1
			}
		} else if len(identifierChain) != 1 {
			if foundPrefix != nil {
				return fmt.Errorf("found a special identifer in access path %T already starting with the this identfier", n)
			}

			// Here, we handle the identifier node immediately after "this"
			// differently to support linking declarations in "this" scope.
			identifier := identifierChain[1]
			declaredIdentifier, err := j.scopes.FindDeclarationEntry(identifier, This, false /* activeOnly */)
			if err != nil {
				return err
			}
			if err := j.symbolTable.AddLink(identifier, declaredIdentifier); err != nil {
				return err
			}
			// We should ignore this identifier node in later traversals since
			// it is already handled.
			j.ignoreNodes[identifier] = true

			// The rest of the identifier nodes start at 2 (skipping "this" and the first identifier after "this")
			unresolvedIndex = 2
		}

		// link the unresolved identifiers to nil pointer
		if err := j.nilifyIDs(identifierChain, unresolvedIndex); err != nil {
			return err
		}

	default:
		return j.GenericSymbolTableBuilder.ProcessUse(node)
	}

	return nil
}

// ProcessOther handles the JavaClassDeclaration node and pre-register the declaration nodes in it for visibility
// throughout the entire class. The handling of other nodes is delegated to the generic builder.
func (j *JavaSymbolTableBuilder) ProcessOther(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if j.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.JavaClassDeclaration:
		// Similar to mast.Root, the declarations in a Java class should be visible in the entire class. So we iterate
		// through the declarations and put them into the current scope.
		for _, decl := range n.Body {
			if err := j.ProcessDeclaration(decl); err != nil {
				return err
			}
			if err := j.PostProcessDeclaration(decl); err != nil {
				return err
			}
		}

	// fallback to generic builder
	default:
		return j.GenericSymbolTableBuilder.ProcessOther(node)
	}
	return nil
}

// PostSymbolicationFixup is responsible for any actions that need to be
// performed after the whole symbolication process is finished.
func (j *JavaSymbolTableBuilder) PostSymbolicationFixup() error {
	// In Java, at the point of MAST creation, the only thing we know
	// is wheter a given declaration is final, but final does not
	// necessarily mean constant. For example, this is a final
	// non-constant declaration as we can only asssume that the
	// variable will not be re-assigned, but not that the internal
	// state of SomeObject will not change over time:
	//
	// final SomeObject c = new SomeObject()
	//
	// According to the JLS
	// (https://docs.oracle.com/javase/specs/jls/se11/html/jls-4.html#jls-4.12.4)
	// for the value to be considered constant, it has to be of
	// primitive or String value and must be initialized with a
	// constant expression. The ConstantFixup method for Java is
	// responsible for updating IsConst value for variable
	// declarations to be true only if the declaration's value is
	// truly constant and not just final (before this method runs,
	// IsConst only means final).
	fixedDecls := make(map[*mast.VariableDeclaration]bool)
	symbols := j.symbolTable.OrderedSymbols()
	for _, symbol := range symbols {
		if _, err := j.fixDeclEntry(symbol, fixedDecls); err != nil {
			return err
		}
	}
	return nil
}

// fixDeclEntry is a helper method to fix declaration entry for a
// single identifier.
func (j *JavaSymbolTableBuilder) fixDeclEntry(symbol *mast.Identifier, fixedDecls map[*mast.VariableDeclaration]bool) (isConst bool, err error) {
	declEntry, err := j.symbolTable.DeclarationEntry(symbol)
	if err != nil {
		return false, err
	}
	if declEntry == nil {
		return false, nil
	}
	if declEntry.Node == nil {
		return false, fmt.Errorf("a symbol table entry missing for the declaration of %q", symbol.Name)

	}
	vdecl, ok := declEntry.Node.(*mast.VariableDeclaration)
	if !ok {
		return false, nil
	}

	if fixedDecls[vdecl] {
		return vdecl.IsConst, nil
	}

	fixedDecls[vdecl] = true

	if vdecl.Value == nil || vdecl.Type == nil {
		return false, nil
	}

	// make sure that declaration type is a String or a primitive type
	typ, ok := vdecl.Type.(*mast.Identifier)
	if !ok {
		return false, nil
	}
	switch typ.Name {
	case "String", "byte", "short", "int", "long", "float", "double", "boolean", "char":
		vdecl.IsConst, err = j.isConstValue(vdecl.Value, fixedDecls)
		return vdecl.IsConst, err
	}

	// unless explicitly determined as such, a variable does not
	// represent a constant
	vdecl.IsConst = false

	return false, nil
}

// isConstValue is a helper function used to determine if a given
// value is constant.
func (j *JavaSymbolTableBuilder) isConstValue(expr mast.Expression, fixedDecls map[*mast.VariableDeclaration]bool) (bool, error) {
	// determination of whether a value is constant is based on the JLS:
	// https://docs.oracle.com/javase/specs/jls/se11/html/jls-15.html#jls-15.28
	switch val := expr.(type) {
	case *mast.NullLiteral, *mast.IntLiteral, *mast.FloatLiteral, *mast.StringLiteral, *mast.CharacterLiteral:
		return true, nil
	case *mast.BinaryExpression:
		switch val.Operator {
		case "*", "/", "%", "+", "-", "<<", ">>", ">>>", "<", "<=", ">", ">=", "==", "!=", "&", "^", "|", "&&", "||":
			isConstLeft, err := j.isConstValue(val.Left, fixedDecls)
			if err != nil {
				return false, err
			}
			isConstRight, err := j.isConstValue(val.Right, fixedDecls)
			if err != nil {
				return false, err
			}
			if isConstLeft && isConstRight {
				return true, nil
			}
		}
	case *mast.ParenthesizedExpression:
		return j.isConstValue(val.Expr, fixedDecls)
	case *mast.Identifier:
		return j.fixDeclEntry(val, fixedDecls)
	}
	return false, nil
}
