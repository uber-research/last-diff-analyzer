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

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const _testDataPrefix = "../../testdata/"
const _metaTestDataPrefix = _testDataPrefix + "translation/"

func TestUnsupported(t *testing.T) {
	t.Run("Test parsing unsupported yaml file", func(t *testing.T) {
		node, err := ParseFile(_testDataPrefix + "analyzer/yaml/comment/base/comment.yaml")
		require.Nil(t, node)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no available parser")
	})
}

func TestMissingFile(t *testing.T) {
	t.Run("Test parsing missing file", func(t *testing.T) {
		// add a random suffix to the filename just to be safe
		node, err := ParseFile(_testDataPrefix + "missing_file_7CE36477.go")
		require.Nil(t, node)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestParsingError(t *testing.T) {
	for _, ext := range [...]string{GoExt, JavaExt} {
		t.Run(fmt.Sprintf("Test parsing error for %q", ext), func(t *testing.T) {
			// write to a temporary file for testing parsing error
			tmpFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("temp-*%v", ext))
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			text := []byte("import a, b, c")
			_, err = tmpFile.Write(text)
			require.NoError(t, err)
			err = tmpFile.Close()
			require.NoError(t, err)

			node, err := ParseFile(tmpFile.Name())
			require.Nil(t, node)
			require.Error(t, err)
			require.Contains(t, err.Error(), "tree-sitter generated ERROR node")
		})
	}
}

func TestGoParser(t *testing.T) {
	t.Run("Test parsing go file", func(t *testing.T) {
		node, err := ParseFile(_testDataPrefix + "analyzer/go/equal/base/dummy.go")
		require.NoError(t, err)

		// the expected tree should look like this:
		// source_file:
		//   package_clause:
		//     package_identifier: dummy
		expected := &Node{
			Type: "source_file",
			Children: []*Node{
				{
					Type: "package_clause",
					Children: []*Node{
						{
							Type:     "package_identifier",
							Name:     "dummy",
							Children: []*Node{},
							fields:   map[string][]*Node{},
						},
					},
					fields: map[string][]*Node{},
				},
			},
			fields: map[string][]*Node{},
		}
		require.Equal(t, expected, node)
	})

	t.Run("Test accessing children by fields", func(t *testing.T) {
		node, err := ParseFile(_metaTestDataPrefix + "go/declarations.go")
		require.NoError(t, err)

		// the expected tree should look like this (names surrounded by [] are fields):
		// source_file:
		//   package_clause:
		//     package_identifier
		//   import_declaration:
		//     import_spec_list:
		//       import_spec
		//         [name]: dot
		//         [path]: interpreted_string_literal
		//       import_spec
		//         [name]: package_identifier
		//         [path]: interpreted_string_literal
		//       import_spec
		//         [name]: blank_identifier
		//         [path]: interpreted_string_literal
		//       import_spec
		//         [path]: interpreted_string_literal
		//    import_declaration
		//      import_spec
		//        [path]: interpreted_string_literal
		//   ...

		// now access the "path" field of the last import_spec
		expected := &Node{
			Type:     "interpreted_string_literal",
			Name:     "\"singlepackage\"",
			Children: []*Node{},
			fields:   map[string][]*Node{},
		}
		require.Equal(t, expected, node.Children[2].Children[0].ChildByField("path"))
	})
}
