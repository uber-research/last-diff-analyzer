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
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"analyzer/core/translation"
	ts "analyzer/core/treesitter"
)

const _metaTestTranslationDir = "../../testdata/translation/"
const _metaTestCheckDir = "../../testdata/check/"

func TestMASTEqual(t *testing.T) {
	checker := &GenericChecker{}
	t.Run("Test Equal methods of MAST nodes", func(t *testing.T) {
		err := filepath.WalkDir(_metaTestTranslationDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failure accessing a path %q: %v", path, err)
			}
			if d.IsDir() {
				return nil
			}
			suffix := filepath.Ext(path)
			// build tree-sitter node
			tsNode, err := ts.ParseFile(path)
			require.NoError(t, err)
			// translate to MAST node
			mastNode, err := translation.Run(tsNode, suffix)
			require.NoError(t, err)
			require.NotNil(t, mastNode)
			// check equivalence against itself
			langChecker, err := newLangChecker(checker, suffix, false /* loggingOn */)
			if err != nil {
				return err
			}
			require.True(t, checker.CheckNode(mastNode, mastNode, langChecker))
			return nil
		})
		require.NoError(t, err)
	})
}
