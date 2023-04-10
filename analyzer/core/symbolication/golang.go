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
	"unicode"
	"unicode/utf8"

	"analyzer/core/mast"
)

// GoSymbolTableBuilder is a language-specific symbol table builder for Go. All internal datastructures are shared
// with GenericSymbolTableBuilder so that both of them can modify the same objects.
type GoSymbolTableBuilder struct {
	// Here we embed the GenericSymbolTableBuilder to share the internal data structures.
	*GenericSymbolTableBuilder
	// Constructors collects all expressions representing constructors
	// as the symbolicator runs, so that they can be post-processed
	// once all symbolication information is computed.
	Constructors map[*mast.EntityCreationExpression]bool
}

// NewGoSymbolTableBuilder returns a properly initialized GoSymbolTableBuilder.
func NewGoSymbolTableBuilder() *GoSymbolTableBuilder {
	// first create a Go-specific builder with nil GenericSymbolTableBuilder field
	builder := &GoSymbolTableBuilder{}
	// assign the filed with properly initialized generic builder
	builder.GenericSymbolTableBuilder = NewGenericSymbolTableBuilder(builder)
	builder.Constructors = make(map[*mast.EntityCreationExpression]bool)
	return builder
}

// SymbolTable returns the built symbol table.
func (g *GoSymbolTableBuilder) SymbolTable() *SymbolTable {
	return g.symbolTable
}

// IsDeclarationPrivate returns if a declaration is private in Go. It simply returns whether the first letter of the
// identifier is upper-case or not.
func (g *GoSymbolTableBuilder) IsDeclarationPrivate(node mast.Node, identifier *mast.Identifier) (bool, error) {
	if identifier == nil {
		return false, nil
	}
	ch, _ := utf8.DecodeRuneInString(identifier.Name)
	// We have to use !IsUpper here since declarations such as "var _someVar int" are also private.
	return !unicode.IsUpper(ch), nil
}

// ProcessDeclaration adds the declared identifiers to the current scope for Go-specific nodes and delegates other nodes
// to the generic builder.
func (g *GoSymbolTableBuilder) ProcessDeclaration(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.GoTypeDeclaration:
		_, err := g.createDeclarationEntry(node, n.Name, true)
		return err

	case *mast.AssignmentExpression:
		if !n.IsShortVarDeclaration {
			// skip processing if it is a normal assignment statement
			return nil
		}
		// Otherwise, it is a short variable declaration and requires special handling.
		// Go's short variable declaration is tricky, consider the following example:
		// var a int     <--- "a" is declared here
		// a, b := foo() <--- a short variable declaration that _redeclares_ (which is the same as assignment) "a"
		//                    and _declares_ "b"
		// Here, the LHS of the short variable declaration contains two identifiers, where "a" is already declared
		// before. According to the Golang language specification [1]:
		// (1) "a short variable declaration may redeclare variables provided they were originally declared earlier in
		//      the _same_ block", i.e., we should not look up the entire stack to determine if it is a redeclaration;
		// (2) "Redeclaration does not introduce a new variable; it just assigns a new value to the original.";
		// (3) The LHS must all be identifiers, i.e., no AccessPath nodes are allowed.
		// [1] https://golang.org/ref/spec#Short_variable_declarations.

		// Therefore, we first iterate through the identifier list on LHS and find which
		// identifiers are re-declarations and which are true declarations, and properly set the
		// links.

		// We create a stack entry to contain the inactive declarations (refer to the comments in
		// handling of VariableDeclaration in GenericSymbolTableBuilder for detailed explanations).
		g.inactiveDecls = append(g.inactiveDecls, []*SymbolTableEntry{})

		for _, expression := range n.Left {
			// According to Golang language specifications [1], the LHS of short variable
			// declarations can only consist of identifiers.
			// [1] https://golang.org/ref/spec#Short_variable_declarations
			expr, ok := expression.(*mast.Identifier)
			if !ok {
				return fmt.Errorf("unexpected node in LHS of short variable declaration: %T", expr)
			}

			// first find if it is already defined in the _current_ scope only
			// Here we need to look for inactive declarations as well (please refer to the comments
			// of "SymbolTableEntry.IsActive" for understanding what inactive means) since we could
			// have declaration and assignment in the _same_ short variable declaration for blank
			// identifiers, e.g.,
			// _, e, _ := foo() // the second blank identifier is permitted by Go
			// In this case, the first "_" will be treated as a declaration by our system (an
			// inactive declaration entry will be added), and the second should be treated as just
			// a normal assignment (by finding the inactive declaration entry).
			decl, err := g.scopes.FindDeclarationEntry(expr, CurrentOnly /* option */, false /* activeOnly */)
			if err != nil {
				return err
			}
			// An identifier is only a declaration if we cannot find the declaration in the current scope.
			if decl == nil {
				// The short variable declaration can only appear inside functions, therefore, it is always a private
				// declaration.

				// We need a special case here for creating inactive declaration entries, please refer to
				// the comments for "SymbolTableEntry.IsActive" for more explanations.
				entry := &SymbolTableEntry{
					Identifier: expr,
					Node:       n,
					IsPrivate:  true,
					IsActive:   false,
				}
				err := g.scopes.AddDeclarationEntry(entry, mast.Blanket /* kind */)
				if err != nil {
					return err
				}
				// add inactive declaration to the stack entry
				l := len(g.inactiveDecls)
				g.inactiveDecls[l-1] = append(g.inactiveDecls[l-1], entry)
				// add a "self" link for the declared identifier and mark the identifier node as
				// already handled
				if err := g.symbolTable.AddLink(expr, entry); err != nil {
					return err
				}
				g.ignoreNodes[expr] = true
			}
		}
		return nil

	case *mast.FunctionDeclaration:
		langFields, ok := n.LangFields.(*mast.GoFunctionDeclarationFields)
		if !ok {
			return langFieldsNodeNotExistsError(node)
		}
		if langFields.Receiver == nil {
			// add only functions and not add method declarations
			if _, err := g.createDeclarationEntry(node, n.Name, true /* active */); err != nil {
				return err
			}
		}

	case *mast.ImportDeclaration:
		// The Alias field in mast.ImportDeclaration is optional, so we skip if it is nil.
		if n.Alias == nil {
			return nil
		}
		// The alias for the import declaration will always be private to the current file.
		entry := &SymbolTableEntry{Identifier: n.Alias, Node: n, IsPrivate: true, IsActive: true}
		return g.scopes.AddDeclarationEntry(entry, mast.Blanket /* kind */)

	case *mast.FieldDeclaration:
		// in Go we can have public fields in private scopes so we
		// have to check the private-ness of a field declaration
		// explicitly
		// TODO: handle embedded structs (struct declarations without names)
		if n.Name != nil {
			isPrivate, err := g.lang.IsDeclarationPrivate(n, n.Name)
			if err != nil {
				return err
			}
			entry := &SymbolTableEntry{
				Identifier: n.Name,
				Node:       n,
				IsPrivate:  isPrivate,
				IsActive:   true,
			}
			return g.scopes.AddDeclarationEntry(entry, n.Name.Kind)
		}

	// fallback to generic implementation for other declaration nodes
	default:
		if err := g.GenericSymbolTableBuilder.ProcessDeclaration(node); err != nil {
			return err
		}
	}

	return nil
}

// ProcessScope handles Go-specific scope-related nodes and delagates handling of other nodes to the generic builder.
func (g *GoSymbolTableBuilder) ProcessScope(node mast.Node, onEnter bool) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.GoTypeDeclaration:
		// Go's type declaration can either be a struct declaration, interface declaration or a type alias. If it is
		// actually a struct or interface declaration, we should create a separate scope for it. The privateness of the
		// scope depends on the declaration itself.
		switch n.Type.(type) {
		case *mast.GoStructType, *mast.GoInterfaceType:
			// Since ProcessScope will be called _after_ ProcessDeclaration, here the declaration should already be
			// registered in the scope. So we retrieve it first to find the privateness of the declaration, the scope
			// will simply inherit the privateness. Note that the Go language specification [1] is a little vague on
			// this, which does not mention the exported field names in unexported struct declarations. However, the
			// exported identifiers should not be visible inside an unexported identifiers.
			// [1] https://golang.org/ref/spec#Exported_identifiers
			if err := processTypeScope(g.scopes, n.Name, node, onEnter); err != nil {
				return err
			}

			// The scope of the struct type or interface type has been processed, so we ignore them here.
			g.ignoreNodes[n.Type] = true

			return nil

		default:
			// We do nothing if the type declaration is not a struct or interface declaration, e.g., a type alias.
			return nil
		}

	case *mast.GoStructType, *mast.GoInterfaceType:
		// Named struct and interface declarations will be handled when handling mast.GoTypeDeclaration, here we handle
		// the anonymous ones, e.g., "var a struct{x int; y int}".

		// The scope for the anonymous struct and interface types are always private.
		return g.handleScope(onEnter, true /* isPrivate */)

	case *mast.GoCommunicationCase:
		// A case clause of the select statement needs its own scope
		// to include both optional assignment in the receive
		// statement and the body of the clause. For example, below
		// tmp8 is declared only once (in the receive statement) and
		// then used twice in the body of the clause (also in
		// symbolication/go/short_decl.go test file):
		//
		//	select {
		//	case tmp8 := <-c:
		//		tmp8, i8 := bar()
		//		return tmp8 + i8
		//	}

		// For communication case nodes, they can only appear in private scopes,
		// therefore they must have private scopes as well.
		return g.handleScope(onEnter, true /* isPrivate */)

	default:
		return g.GenericSymbolTableBuilder.ProcessScope(node, onEnter)
	}
}

// ProcessUse delegates the handling of uses to generic builder.
func (g *GoSymbolTableBuilder) ProcessUse(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.EntityCreationExpression:
		g.Constructors[n] = true
	}

	return g.GenericSymbolTableBuilder.ProcessUse(node)
}

// ProcessOther delegates the handling of other nodes to the generic builder.
func (g *GoSymbolTableBuilder) ProcessOther(node mast.Node) error {
	// If the node is in the ignore set, it means that it has already been properly processed in upper-level and should
	// be ignored.
	if g.ignoreNodes[node] {
		return nil
	}

	return g.GenericSymbolTableBuilder.ProcessOther(node)
}

// PostProcessDeclaration handles declaration node after they have been visited.
func (g *GoSymbolTableBuilder) PostProcessDeclaration(node mast.Node) error {
	if g.ignoreNodes[node] {
		return nil
	}

	switch n := node.(type) {
	case *mast.AssignmentExpression:
		if !n.IsShortVarDeclaration {
			// skip processing if it is a normal assignment statement
			return nil
		}
		// We need to "activate" this declaration so that it's
		// available for all statements following this one (please see
		// mast.AssignmentExpression handling in ProcessDeclaration for
		// additional explanation).
		g.activateDeclarations()
	default:
		return g.GenericSymbolTableBuilder.PostProcessDeclaration(node)
	}
	return nil
}

// PostSymbolicationFixup is responsible for any actions that need to be
// performed after the whole symbolication process is finished.
func (g *GoSymbolTableBuilder) PostSymbolicationFixup() error {
	for c := range g.Constructors {
		values, err := getKeyValuePairs(c)
		if err != nil {
			return err
		}
		if values == nil {
			// we don't do special symbolication for constructors that do not contain key-value pairs
			continue
		}
		isMap, declType, err := g.isMapType(c.Type)
		if err != nil {
			return err
		}
		if isMap {
			// we don't do special symbolication for map constructors
			continue
		}
		if err := g.symbolicateFieldKeys(values, declType); err != nil {
			return err
		}
	}
	return nil
}

// isMapType determines if a given expression represents a map
// declaration type or not and additionally returns the declared type
// (whether it's a map type or not).
func (g *GoSymbolTableBuilder) isMapType(e mast.Expression) (bool, mast.Expression, error) {
	switch typ := e.(type) {
	case *mast.GoMapType:
		return true, typ, nil
	case *mast.Identifier:
		declEntry, err := g.symbolTable.DeclarationEntry(typ)
		if err != nil {
			return false, nil, err
		}
		if declEntry == nil {
			return false, nil, nil
		}
		decl, ok := declEntry.Node.(*mast.GoTypeDeclaration)
		if !ok {
			return false, nil, nil
		}
		if _, ok := decl.Type.(*mast.GoMapType); ok {
			return true, decl.Type, nil
		}
		return false, decl.Type, nil
	}

	return false, nil, nil
}

// getKeyValuePairs returns a list of key-value pairs from the
// struct/map constructor if this constructor contains such a list
// (and not just a list of identifiers). Otherwise it returns nil.
func getKeyValuePairs(cons *mast.EntityCreationExpression) ([]*mast.KeyValuePair, error) {
	if cons.Value == nil {
		return nil, nil
	}
	var l []*mast.KeyValuePair
	keyValuePairs := false
	others := false
	for _, e := range cons.Value.Values {
		if pair, ok := e.(*mast.KeyValuePair); ok {
			l = append(l, pair)
			keyValuePairs = true
			continue
		}
		others = true
	}
	if keyValuePairs && others {
		// a mix of key-value pairs and other expressions (should not happen)
		return nil, fmt.Errorf("a mix of key-value expressions and other expressions in the constructor: %T", cons)
	}
	return l, nil
}

// symbolicateFieldKeys symbolicates keys used in a constructor (e.g.,
// a struct constructor).
func (g *GoSymbolTableBuilder) symbolicateFieldKeys(values []*mast.KeyValuePair, declType mast.Expression) error {
	var fieldDeclarations []*mast.FieldDeclaration
	if declType != nil {
		if sdecl, ok := declType.(*mast.GoStructType); ok {
			fieldDeclarations = sdecl.Declarations
		}
	}
	for _, v := range values {
		key, ok := v.Key.(*mast.Identifier)
		if !ok {
			// can't do any symbolication as key is not an identifier
			continue
		}
		fieldDecl := getDeclaredField(key.Name, fieldDeclarations)
		if fieldDecl == nil {
			// field not found so nothing to symbolicate (possible,
			// for example, with embedded structs - handling these is
			// a TODO, see reject-struct-create-g test for an example)
			continue
		}
		entry, err := g.symbolTable.DeclarationEntry(fieldDecl.Name)
		if err != nil {
			return err
		}
		if entry != nil {
			// we found a connection between a key in the struct
			// constructor and the struct declaration - establish it,
			// even if it means replacing existing link
			keyEntry, err := g.symbolTable.DeclarationEntry(key)
			if err != nil {
				return err
			}
			if keyEntry == entry {
				// indentifier representing they key already points to
				// the right field declaration
				continue

			}
			if err := g.symbolTable.ReplaceLink(key, entry); err != nil {
				return err
			}
			g.ignoreNodes[key] = true
			continue
		}
		// We are not sure what this key represents, but it's not a
		// key representing a known struct field and it's not a key
		// used to initialize a map. Consequently (and
		// conservatively), we do not want to associate it with some
		// random identifier that happens to be in one of the
		// reachable scopes, so that we can avoid spurious renaming
		// (and spurious-auto approval) such as in the following where
		// someStruct is both the name of the function parameter and
		// of a struct field that cannot be located within the code we
		// have at our disposal during the analysis process.
		//
		// func foo(someField int) someStruct {
		//   return someStruct{someField: 42}
		// }
		//
		if err := g.symbolTable.RemoveLink(key); err != nil {
			return err
		}
	}
	return nil
}

// getDeclaredField returns a field with a given name declared in the
// struct type
func getDeclaredField(name string, fieldDeclarations []*mast.FieldDeclaration) *mast.FieldDeclaration {
	for _, decl := range fieldDeclarations {
		if decl.Name == nil {
			// TODO: handle embedded structs (struct declarations without names)
			// see similar comment in ProcessDeclaration
			continue
		}
		if decl.Name.Name == name {
			return decl
		}
	}
	return nil
}
