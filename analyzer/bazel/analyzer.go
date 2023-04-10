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

// Package bazel implements the main logic of the analysis for Bazel-related
// files. It currently supports auto-approval of changes to
// BUILD.bazel files which involve comments, formatting, test-related
// directives, and removal (but not modification or addition) of
// dependencies.
package bazel

// Analyzer is an analyzer for Bazel files.
type Analyzer struct {
	// analyzableFileName represents file name of analyzable files (it
	// is made configurable for testing purposes, as test files have
	// different names than prod files).
	AnalyzableFileName string

	baseForest astForest
	lastForest astForest
}

// ChangesEq returns true if the changes between Bazel files in base and
// last diffs are equivalent.
func (a *Analyzer) ChangesEq() (bool, error) {
	eq, err := a.astForestEq()
	if err != nil {
		return false, err
	}

	return eq, nil
}
