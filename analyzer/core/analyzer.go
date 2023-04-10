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

// Package core implements the main logic of the analysis for "core"
// source files, that is source files that can be translated into the
// same meta-IR and are therefore subject to uniform equivalence
// analysis.
package core

import (
	"analyzer/core/check"
	"analyzer/core/mast"
	"analyzer/core/symbolication"
	"analyzer/core/transformation"
)

// Analyzer is an analyzer for core files.
type Analyzer struct {
	// RenamingOn indicates whether renaming is enabled or not.
	RenamingOn bool
	// LoggingOn indicates whether logging support is enabled or not.
	LoggingOn bool

	// suffix is the suffix of the files to analyze, needed to determine the language.
	suffix string
	// baseASTForest is the forest of all MAST nodes for base diff.
	baseASTForest []mast.Node
	// baseASTForest is the forest of all MAST nodes for last diff.
	lastASTForest []mast.Node
}

// ChangesEq returns true if the changes between core files in base and last diffs are equivalent.
func (a *Analyzer) ChangesEq() (bool, error) {
	if len(a.baseASTForest) != len(a.lastASTForest) {
		return false, nil
	}

	// early return if no files are actually processed
	if len(a.baseASTForest) == 0 {
		return true, nil
	}

	// run symbolication, and then transformation on the base and last forest
	forests := [...][]mast.Node{a.baseASTForest, a.lastASTForest}
	symbolTables := make([]*symbolication.SymbolTable, len(forests))

	for i, forest := range forests {
		table, err := symbolication.Run(forest, a.suffix)
		if err != nil {
			return false, err
		}
		symbolTables[i] = table

		// setup transformers
		var transformers []transformation.Transformer

		// add renamer if it is enabled
		if a.RenamingOn {
			transformers = append(transformers, transformation.NewRenamer(table))
		}

		// run the transformers
		for _, transformer := range transformers {
			var err error
			forests[i], err = transformer.Transform(forest)
			if err != nil {
				return false, err
			}
		}
	}

	// run check on the base and last forest
	baseForest := forests[0]
	lastForest := forests[1]
	baseTable := symbolTables[0]
	lastTable := symbolTables[1]
	return check.Run(baseForest, lastForest, baseTable, lastTable, a.suffix, a.LoggingOn)
}
