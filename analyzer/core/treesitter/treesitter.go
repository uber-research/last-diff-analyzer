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

package treesitter

// #include "api.h"
// #include "parser.h"
// #include <stdlib.h>
// #include <string.h>
// TSLanguage *tree_sitter_java();
// TSLanguage *tree_sitter_go();
import "C"
import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unsafe"
)

// define constants of file extensions for different languages.
const (
	GoExt   = ".go"
	JavaExt = ".java"
)

// Node stores the information collected from the nodes generated by tree-sitter.
type Node struct {
	// Type is the type for the Node.
	Type string
	// Name is the name of an optional identifier (or empty string otherwise).
	Name string
	// Children contains children nodes.
	Children []*Node
	// fields is an optional map that maps field names to the node pointers in Children.
	// tree-sitter provides access to the node's children by a unique field name. This makes the children access easier
	// and sometimes essential (e.g., when all children of a node are optional and the positions in the ordered slice
	// are meaningless). However, the field names are not always available for every child, every node or every
	// language. So we make the fields map optional and create a ChildByField method to proxy the access,
	// if the access failed an error will be returned.
	// The supported field names for each node can be usually found at "src/node-types.json" in the corresponding
	// language repository of tree-sitter.
	fields map[string][]*Node
}

// ChildByField provides access to the Node's Children by a given field name.
// nil is returned if the supplied field name does not exist for the node. If multiple children exist for the given
// field, the first child is returned.
func (n *Node) ChildByField(name string) *Node {
	slice := n.fields[name]
	if slice == nil {
		return nil
	}
	return slice[0]
}

// ChildrenByField provides access to the Node's Children by a given field name. It returns a slice of children
// associated with the same field name. nil is returned if the supplied field name does not exist for the node.
func (n *Node) ChildrenByField(name string) []*Node {
	slice := n.fields[name]
	if slice == nil {
		return nil
	}
	copied := make([]*Node, len(slice))
	copy(copied, n.fields[name])
	return copied
}

// buildNode recursively builds a Node using a *C.TSTreeCursor. TSTreeCursor is a cursor struct in tree-sitter that
// can be moved around using tree-sitter APIs and it allows walking a syntax tree with maximum efficiency.
// See https://tree-sitter.github.io/tree-sitter/using-parsers#walking-trees-with-tree-cursors for more details.
// More specifically, it
// (1) populates the Name field using the position information from the current TSNode and source string;
// (2) recursively builds children and populates fields map if a field name exists for each child.
// It is guaranteed that the cursor will move back to its original state.
func buildNode(cursor *C.TSTreeCursor, source string, parentType string) (*Node, error) {
	// retrieve the current node from the cursor
	tsNode := C.ts_tree_cursor_current_node(cursor)
	nodeType := C.GoString(C.ts_node_type(tsNode))

	if nodeType == "ERROR" {
		return nil, errors.New("tree-sitter generated ERROR node")
	}

	// note that tree-sitter constructs concrete syntax trees, where terminals are included as unnamed nodes, so to
	// build an AST we skip the unnamed nodes. However, some anonymous terminal nodes are also important, so a few
	// exceptions are made below, otherwise we will simply return nil to skip the node generation.
	// see https://tree-sitter.github.io/tree-sitter/using-parsers#named-vs-anonymous-nodes for more details.
	if !C.ts_node_is_named(tsNode) {
		switch parentType {
		case "dimensions", "binary_expression", "unary_expression", "channel_type", "inc_statement", "dec_statement",
			"update_expression", "modifiers", "range_clause", "receive_statement", "module_declaration",
			"module_directive":
			// (1) we only collect named nodes, so "dimensions" node in Java (String "[][][]") will not contain
			//     "[" or "]" since they are terminal nodes. This will be reflected here where child will be nil.
			//     However, this is an important piece of information (otherwise we cannot tell apart "@NotNull[][]" vs
			//     "[] @NotNull []"). So here, we also collect "[" and "]";
			// (2) similarly for unary and binary expression, the operator node ("+", "-" etc.) should also be
			//     preserved;
			// (3) channel type also needs "chan" and "<-" to indicate the direction of the channel.
			// (4) inc_statement / dec_statement nodes in Go and update_expression nodes in Java require the
			//     position/value of the operator;
			// (5) "modifiers" nodes in Java contain multiple terminal nodes such as "public", "private" etc.
			// (6) "range_clause" nodes in Go need the terminal operator node "=" vs ":=" to determine whether it is
			//     an assignment or a short variable declaration;
			// (7) "receive_statement" nodes in Go need the terminal operator node "=" vs ":=" to determine whether it
			//     is an assignment or a short variable declaration;
			// (8) "module_declaration" nodes in Java need the terminal "open" node to determine the property of the
			//     module;
			// (9) "module_directive" nodes in Java need the terminal keyword (such as "requires" and "exports") for
			//     the module directive.
			return &Node{Type: nodeType}, nil
		default:
			return nil, nil
		}
	}

	// skip node generation for the following nodes:
	// (1) comment;
	// (2) escape_sequence represents escape sequences in string nodes (e.g., "\t" in "hello\tworld"). However, we
	//     preserve the entire string verbatim, so it is not necessary to keep this node;
	// (3) empty_statement represents an empty statement (";") in Go.
	if nodeType == "comment" || nodeType == "escape_sequence" {
		if C.ts_tree_cursor_goto_first_child(cursor) {
			return nil, errors.New("unexpected children for comment/escape_sequence node")
		}
		return nil, nil
	}

	// empty_statement will have a ";" terminal child, so here we do a special check
	if nodeType == "empty_statement" {
		if count := C.uint(C.ts_node_child_count(tsNode)); count != 1 {
			return nil, fmt.Errorf("unexpected children count for empty_statement: %v", count)
		}
		if childType := C.GoString(C.ts_node_type(C.ts_node_child(tsNode, 0))); childType != ";" {
			return nil, fmt.Errorf("unexpected child type for empty_statement: %q", childType)
		}
		return nil, nil
	}

	node := &Node{
		Type:     nodeType,
		Children: []*Node{},
		fields:   map[string][]*Node{},
	}

	// we use the position to retrieve missing name information from source.
	if strings.Contains(nodeType, "identifier") || strings.Contains(nodeType, "literal") {
		start, end := uint32(C.ts_node_start_byte(tsNode)), uint32(C.ts_node_end_byte(tsNode))
		node.Name = source[start:end]
	}

	// a few cases that require special handling to fill in the Name field
	switch nodeType {
	case "true", "false", "null_literal", "nil", "this", "super",
		"void_type", "boolean_type", "floating_point_type", "integral_type", "label_name",
		"requires_modifier":
		start, end := uint32(C.ts_node_start_byte(tsNode)), uint32(C.ts_node_end_byte(tsNode))
		node.Name = source[start:end]
	case "dot":
		node.Name = "."
	case "asterisk":
		node.Name = "*"
	}

	// recursively build children, the cursor movement APIs (ts_tree_cursor_*) will return true if the operation was
	// successfull and false otherwise.
	if !C.ts_tree_cursor_goto_first_child(cursor) {
		return node, nil
	}

	for {
		child, err := buildNode(cursor, source, nodeType)
		if err != nil {
			return nil, err
		}

		if child != nil {
			// populate the fields map if field name exists
			if cFieldName := C.ts_tree_cursor_current_field_name(cursor); cFieldName != nil {
				fieldName := C.GoString(cFieldName)
				node.fields[fieldName] = append(node.fields[fieldName], child)
			}
			node.Children = append(node.Children, child)
		}
		// move on to the next sibling, and break the loop if no sibling is found
		if !C.ts_tree_cursor_goto_next_sibling(cursor) {
			break
		}
	}
	// finally go back to parent node
	C.ts_tree_cursor_goto_parent(cursor)
	return node, nil
}

// ParseFile reads and parses the file and returns a root Node
// An empty Node is returned if any error occurs.
func ParseFile(path string) (*Node, error) {
	// create a parser and set the language according to the file's extension
	parser := C.ts_parser_new()
	defer C.ts_parser_delete(parser)

	var language *C.TSLanguage
	switch suffix := filepath.Ext(path); suffix {
	case GoExt:
		language = C.tree_sitter_go()
	case JavaExt:
		language = C.tree_sitter_java()
	default:
		return nil, fmt.Errorf("no available parser for file %q", path)
	}

	// set language for parser
	C.ts_parser_set_language(parser, language)

	// open the file and read content
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	source := string(bytes)

	// allocate a CString in CGO, note that it is our responsibility to deallocate it
	cSource := C.CString(source)
	tsTree := C.ts_parser_parse_string(parser, nil, cSource, C.uint(C.strlen(cSource)))
	defer C.free(unsafe.Pointer(cSource))
	defer C.ts_tree_delete(tsTree)

	// create cursor for the root node and pass it to buildNode
	root := C.ts_tree_root_node(tsTree)
	cursorStruct := C.ts_tree_cursor_new(root)
	cursor := &cursorStruct
	defer C.ts_tree_cursor_delete(cursor)

	// build the tree in go
	goTree, err := buildNode(cursor, source, "" /* parentType */)
	if err != nil {
		return nil, err
	}

	return goTree, nil
}
