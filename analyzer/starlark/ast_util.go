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

package starlark

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"go.starlark.net/syntax"
)

// astForest stores starlark ASTs for a given change set in a form
// ensuring consistent iteration order.
type astForest struct {
	// fileNames contains names of modified files.
	fileNames []string
	// asts contains ASTs for modified files.
	asts []*syntax.File
}

// BaseIRBuild builds intermediate representation for relevant files
// in the base diff.
func (a *Analyzer) BaseIRBuild(filesToAnalyze []string, rootDir string) error {
	var err error
	a.baseForest, err = a.astForestBuild(filesToAnalyze, rootDir)
	return err
}

// LastIRBuild builds intermediate representation for relevant files
// in the last diff.
func (a *Analyzer) LastIRBuild(filesToAnalyze []string, rootDir string) error {
	var err error
	a.lastForest, err = a.astForestBuild(filesToAnalyze, rootDir)
	return err
}

// astForestBuild builds a forest of ASTs for all specified starlark
// files.
func (a *Analyzer) astForestBuild(filesToAnalyze []string, rootDir string) (astForest, error) {
	forest := astForest{}
	for _, f := range filesToAnalyze {
		if a.IsAnalyzable(f) {
			srcFilePath := path.Join(rootDir, f)
			ast, err := astBuild(srcFilePath)
			if err != nil {
				return forest, fmt.Errorf("failed to create a starlark AST for %q: %v", srcFilePath, err)
			}
			forest.fileNames = append(forest.fileNames, f)
			forest.asts = append(forest.asts, ast)
		}
	}
	return forest, nil
}

// IsAnalyzable determines if a given file name represents file
// analyzable by this analyzer.
func (a *Analyzer) IsAnalyzable(fileName string) bool {
	return strings.HasSuffix(fileName, "star")
}

// astBuild construct and returns an AST for a starlark file at a given path.
func astBuild(filePath string) (*syntax.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Disable RetainComments so that the resulting AST will not contain any comments.
	// See more at https://pkg.go.dev/go.starlark.net/syntax.
	return syntax.Parse(filePath, nil /* src */, ^syntax.RetainComments /* mode */)
}

// astForestEq determines equivalence of the starlark AST forests.
func (a *Analyzer) astForestEq() (bool, error) {
	if len(a.baseForest.fileNames) != len(a.lastForest.fileNames) ||
		len(a.baseForest.asts) != len(a.lastForest.asts) ||
		len(a.baseForest.fileNames) != len(a.lastForest.asts) {
		// We build both starlark AST forests based on the same (single) list of files, so this should
		// never happen.
		return false, errors.New("lists of starlark changes to compare have different length")
	}

	for i, baseAst := range a.baseForest.asts {
		if a.baseForest.fileNames[i] != a.lastForest.fileNames[i] {
			return false, errors.New("lists of starlark files to compare are different")
		}
		lastAst := a.lastForest.asts[i]

		eq, err := astEq(baseAst, lastAst)
		if err != nil {
			return false, err
		}
		if !eq {
			return false, nil
		}

	}
	return true, nil
}
