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

package yaml

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// astEq checks if two yaml ASTs are equivalent.
func astEq(baseAst *yaml.Node, lastAst *yaml.Node) (bool, error) {
	removeCommentAndFormat(baseAst)
	removeCommentAndFormat(lastAst)

	baseBytes, err := yaml.Marshal(baseAst)
	if err != nil {
		return false, err
	}
	lastBytes, err := yaml.Marshal(lastAst)
	if err != nil {
		return false, err
	}
	return bytes.Compare(baseBytes, lastBytes) == 0, nil
}

// removeCommentAndFormat removes comments and formatting-related data
// from the AST.
func removeCommentAndFormat(ast *yaml.Node) {
	if ast != nil {
		ast.HeadComment = ""
		ast.LineComment = ""
		ast.FootComment = ""
		ast.Line = 0
		ast.Column = 0
		removeCommentAndFormat(ast.Alias)
		for _, n := range ast.Content {
			removeCommentAndFormat(n)
		}
	}
}
