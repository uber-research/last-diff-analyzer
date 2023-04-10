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

package fallback

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"

	"analyzer/common"
)

// Analyzer is the fallback analyzer for files not supported by other
// analyzers.
type Analyzer struct {
	// The following contain file contents to be compared by this
	// analyzer.
	baseFiles [][]byte
	lastFiles [][]byte

	// SubAnalyzers contains all available sub-analyzers for different
	// file formats, other than the fallback analyzer.
	SubAnalyzers []common.Analyzer
}

// ChangesEq returns true if the changes between Go files in base and
// last diffs are equivalent.
func (a *Analyzer) ChangesEq() (bool, error) {
	if len(a.baseFiles) != len(a.lastFiles) {
		return false, errors.New("lists of remaining to compare have different length")
	}
	for i, baseFile := range a.baseFiles {
		if bytes.Compare(baseFile, a.lastFiles[i]) != 0 {
			return false, nil
		}
	}
	return true, nil
}

// BaseIRBuild builds intermediate representation for relevant files
// in the base diff.
func (a *Analyzer) BaseIRBuild(filesToAnalyze []string, rootDir string) error {
	var err error
	a.baseFiles, err = a.readFiles(filesToAnalyze, rootDir)
	return err
}

// LastIRBuild builds intermediate representation for relevant files
// in the last diff.
func (a *Analyzer) LastIRBuild(filesToAnalyze []string, rootDir string) error {
	var err error
	a.lastFiles, err = a.readFiles(filesToAnalyze, rootDir)
	return err
}

// readFiles reads files from a given set of paths.
func (a *Analyzer) readFiles(filesToAnalyze []string, rootDir string) ([][]byte, error) {
	var files [][]byte
	for _, f := range filesToAnalyze {
		if !a.isAnalyzableByAny(f) {
			filePath := path.Join(rootDir, f)
			dat, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			files = append(files, dat)
		}
	}
	return files, nil
}

// isAnalyzableByAny determines if a file with a given name is analyzable
// by any of the sub-analyzers.
func (a *Analyzer) isAnalyzableByAny(fileName string) bool {
	for _, sub := range a.SubAnalyzers {
		if sub.IsAnalyzable(fileName) {
			return true
		}
	}
	return false
}
