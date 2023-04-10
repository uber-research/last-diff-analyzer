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

package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"analyzer/core/mast"
	"analyzer/core/translation"
	ts "analyzer/core/treesitter"
)

// _analyzableFileExt determines extensions for files analyzable by
// the core analyzer.
var _analyzableFileExt = []string{ts.GoExt, ts.JavaExt}

func (a *Analyzer) buildMASTForest(filesToAnalyze []string, rootDir string) ([]mast.Node, error) {
	forest := []mast.Node{}
	for _, filename := range filesToAnalyze {
		if !a.IsAnalyzable(filename) {
			continue
		}
		suffix := filepath.Ext(filename)
		if a.suffix == "" {
			a.suffix = suffix
		} else if a.suffix != suffix {
			return nil, fmt.Errorf("file %q has different suffix %q", filename, suffix)
		}

		// run tree-sitter parser and then do translation
		path := filepath.Join(rootDir, filename)
		tsNode, err := ts.ParseFile(path)
		if err != nil {
			return nil, err
		}
		mastNode, err := translation.Run(tsNode, suffix)
		if err != nil {
			return nil, err
		}
		forest = append(forest, mastNode)
	}

	return forest, nil
}

// BaseIRBuild builds intermediate representation for relevant files
// in the base diff.
func (a *Analyzer) BaseIRBuild(filesToAnalyze []string, rootDir string) error {
	forest, err := a.buildMASTForest(filesToAnalyze, rootDir)
	a.baseASTForest = forest
	return err
}

// LastIRBuild builds intermediate representation for relevant files
// in the last diff.
func (a *Analyzer) LastIRBuild(filesToAnalyze []string, rootDir string) error {
	forest, err := a.buildMASTForest(filesToAnalyze, rootDir)
	a.lastASTForest = forest
	return err
}

// IsAnalyzable determines if a given file name represents file
// analyzable by this analyzer.
func (a *Analyzer) IsAnalyzable(fileName string) bool {
	for _, ext := range _analyzableFileExt {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
