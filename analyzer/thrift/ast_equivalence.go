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

package thrift

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.uber.org/thriftrw/ast"
)

// astEq checks if two thrift ASTs are equivalent.
func astEq(baseAst, lastAst *ast.Program) (bool, error) {
	visitor := &trimVisitor{}
	ast.Walk(visitor, baseAst)
	if visitor.err != nil {
		return false, visitor.err
	}
	ast.Walk(visitor, lastAst)
	if visitor.err != nil {
		return false, visitor.err
	}

	// Since the thriftw library does not provide an easy way to unparse the AST nodes, we simply
	// fall back to using json.Marshal to compare the equality of the nodes.
	baseBytes, err := json.Marshal(baseAst)
	if err != nil {
		return false, err
	}
	lastBytes, err := json.Marshal(lastAst)
	if err != nil {
		return false, err
	}

	return bytes.Equal(baseBytes, lastBytes), nil
}

// trimVisitor trims the position information (i.e., Line and Column) and comments (i.e., Doc).
type trimVisitor struct {
	// err keeps track of the error occurred during the visit since the library-defined interface
	// does not allow us to return an error.
	err error
}

func (v *trimVisitor) Visit(_ ast.Walker, node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Program:
		// no-op for top-level node.
	case ast.ConstantBoolean, ast.ConstantDouble, ast.ConstantInteger, ast.ConstantString:
		// These nodes are type aliases for builtin types, no-op.
	case ast.ConstantList, ast.ConstantMap, ast.ConstantMapItem, ast.ConstantReference:
		// These nodes are not pointer types under the ast.Node interface. So we cannot modify them
		// here since they are copies. Instead, we will modify them at their parent nodes instead.
		// Specifically, we will call stripConstantValue that strips the nodes and returns a
		// stripped copy, then the field with these types gets swapped at the parent node.
	case ast.BaseType, ast.ListType, ast.MapType, ast.SetType, ast.TypeReference:
		// Similar to the ast.ConstantValue types above, we will modify the type nodes at their
		// parent nodes by calling stripType to get a stripped version of the type nodes.
	case *ast.Annotation:
		n.Line, n.Column = 0, 0
	case *ast.Constant:
		n.Line, n.Column, n.Doc = 0, 0, ""
		n.Type = v.stripType(n.Type)
		n.Value = v.stripConstantValue(n.Value)
	case *ast.CppInclude:
		n.Line, n.Column = 0, 0
	case *ast.Enum:
		n.Line, n.Column, n.Doc = 0, 0, ""
	case *ast.EnumItem:
		n.Line, n.Column, n.Doc = 0, 0, ""
	case *ast.Field:
		n.Line, n.Column, n.Doc = 0, 0, ""
		n.Default = v.stripConstantValue(n.Default)
		n.Type = v.stripType(n.Type)
	case *ast.Function:
		n.Line, n.Column, n.Doc = 0, 0, ""
		n.ReturnType = v.stripType(n.ReturnType)
	case *ast.Include:
		n.Line, n.Column = 0, 0
	case *ast.Namespace:
		n.Line, n.Column = 0, 0
	case *ast.Service:
		n.Line, n.Column, n.Doc = 0, 0, ""
	case *ast.Struct:
		n.Line, n.Column, n.Doc = 0, 0, ""
	case *ast.Typedef:
		n.Line, n.Column, n.Doc = 0, 0, ""
		n.Type = v.stripType(n.Type)
	default:
		v.err = fmt.Errorf("unhandled node type %T", node)
	}

	// If any error occurred when handling the node, we stop further traversal by returning a nil visitor.
	if v.err != nil {
		return nil
	}

	return v
}

// stripType removes the position information (i.e., Line and Column) for the ast.Type nodes
// recursively and returns the stripped version.
func (v *trimVisitor) stripType(typ ast.Type) ast.Type {
	if typ == nil {
		return nil
	}
	switch t := typ.(type) {
	case ast.BaseType:
		t.Line, t.Column = 0, 0
		return t
	case ast.ListType:
		t.Line, t.Column = 0, 0
		t.ValueType = v.stripType(t.ValueType)
		return t
	case ast.MapType:
		t.Line, t.Column = 0, 0
		t.KeyType = v.stripType(t.KeyType)
		t.ValueType = v.stripType(t.ValueType)
		return t
	case ast.SetType:
		t.Line, t.Column = 0, 0
		t.ValueType = v.stripType(t.ValueType)
		return t
	case ast.TypeReference:
		t.Line, t.Column = 0, 0
		return t
	default:
		v.err = fmt.Errorf("unhandled type node type %T", typ)
		return nil
	}
}

// stripConstantValue removes the position information (i.e., Line and Column) for the
// ast.ConstantValue nodes recursively and returns the stripped version.
func (v *trimVisitor) stripConstantValue(constant ast.ConstantValue) ast.ConstantValue {
	if constant == nil {
		return nil
	}
	switch c := constant.(type) {
	case ast.ConstantBoolean, ast.ConstantDouble, ast.ConstantInteger, ast.ConstantString:
		// These nodes are type aliases for builtin types, no-op.
		return c
	case ast.ConstantList:
		c.Line, c.Column = 0, 0
		for i := range c.Items {
			c.Items[i] = v.stripConstantValue(c.Items[i])
		}
		return c
	case ast.ConstantMap:
		c.Line, c.Column = 0, 0
		for i := range c.Items {
			c.Items[i].Line, c.Items[i].Column = 0, 0
			c.Items[i].Key = v.stripConstantValue(c.Items[i].Key)
			c.Items[i].Value = v.stripConstantValue(c.Items[i].Value)
		}
		return c
	case ast.ConstantReference:
		c.Line, c.Column = 0, 0
		return c
	default:
		v.err = fmt.Errorf("unhandled constant value type %T", constant)
		return nil
	}
}
