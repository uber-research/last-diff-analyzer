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

package protobuf

import (
	"bytes"
	"encoding/json"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// astEq checks if two protobuf ASTs are equivalent.
func astEq(baseAst, lastAst *parser.Proto) (bool, error) {
	// Trim the comments, inline comments, and meta (position) information from the AST nodes.
	v := &trimVisitor{}
	baseAst.Accept(v)
	lastAst.Accept(v)

	// Marshal the AST nodes to json.
	baseBytes, err := json.Marshal(baseAst)
	if err != nil {
		return false, err
	}
	lastBytes, err := json.Marshal(lastAst)
	if err != nil {
		return false, err
	}
	a, b := string(baseBytes), string(lastBytes)
	print(a, b)

	// Compare raw jsons for equality checking.
	return bytes.Equal(baseBytes, lastBytes), nil
}

// _emptyMeta is an empty meta information that will be used to clear the meta fields for all AST
// nodes during trimVisitor's visits.
var _emptyMeta = meta.Meta{}

// Implement a visitor pattern to trim comments and meta information from the AST.
type trimVisitor struct{}

// VisitComment is implemented to satisfy the Visitor interface. It does nothing.
func (*trimVisitor) VisitComment(*parser.Comment) {}

// VisitEmptyStatement clears the inline comments of the node.
func (*trimVisitor) VisitEmptyStatement(n *parser.EmptyStatement) (next bool) {
	n.InlineComment = nil
	return true
}

// VisitEnum clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitEnum(n *parser.Enum) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitEnumField clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitEnumField(n *parser.EnumField) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitExtend clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitExtend(n *parser.Extend) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitExtensions clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitExtensions(n *parser.Extensions) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitField clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitField(n *parser.Field) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitGroupField clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitGroupField(n *parser.GroupField) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitImport clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitImport(n *parser.Import) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitMapField clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitMapField(n *parser.MapField) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitMessage clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitMessage(n *parser.Message) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitOneof clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitOneof(n *parser.Oneof) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitOneofField clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitOneofField(n *parser.OneofField) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitOption clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitOption(n *parser.Option) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitPackage clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitPackage(n *parser.Package) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitReserved clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitReserved(n *parser.Reserved) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}

// VisitRPC clears the comments, inline comments, and meta information of the node.
func (v *trimVisitor) VisitRPC(n *parser.RPC) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	n.RPCRequest.Meta = _emptyMeta
	n.RPCResponse.Meta = _emptyMeta
	// The visitor driver fails to visit the options, here we manually visit the options.
	for _, opt := range n.Options {
		v.VisitOption(opt)
	}
	return true
}

// VisitService clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitService(n *parser.Service) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.InlineCommentBehindLeftCurly = nil
	n.Meta = _emptyMeta
	return true
}

// VisitSyntax clears the comments, inline comments, and meta information of the node.
func (*trimVisitor) VisitSyntax(n *parser.Syntax) (next bool) {
	n.Comments = nil
	n.InlineComment = nil
	n.Meta = _emptyMeta
	return true
}
