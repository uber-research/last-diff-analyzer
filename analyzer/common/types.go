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

// Package common provides support shared by various parts of the analyzer.
package common

// Analyzer defines functionality of sub-analyzers for different file
// formats.
type Analyzer interface {
	// IsAnalyzable determines if a given file name represents file
	// analyzable by this analyzer.
	IsAnalyzable(fileName string) bool
	// BaseIRBuild builds intermediate representation for relevant
	// files in the base diff.
	BaseIRBuild(filesToAnalyze []string, rootDir string) error
	// LastIRBuild builds intermediate representation for relevant
	// files in the base diff.
	LastIRBuild(filesToAnalyze []string, rootDir string) error
	// ChangesEq returns true if the changes between Go files in base
	// and last diffs are equivalent.
	ChangesEq() (bool, error)
}
