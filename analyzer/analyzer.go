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

// Package analyzer implements the main logic of an analysis
// determining if two versions of files in a modified files set
// (described by two different diffs) are semantically equivalent
// despite not being (syntactically) identical. The analyzer determines
// if a diff should be auto-approved, so the worst case scenario of
// the analyzer rejecting auto-approval or altogether failing is
// that the diff will have to be approved manually.
package analyzer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"analyzer/common"
	"analyzer/fallback"

	"go.uber.org/multierr"
)

// Status (exit) codes
const (
	Approve = 0
	Reject  = 1
	Failure = -1
)

// The two values represent prefixes of the two lines in the per-file normalized diff header.
const (
	_diffHeaderFromPrefix = "--- "
	_diffHeaderToPrefix   = "+++ "
)

// _devNullHeader represents part of the diff header describing file addition (if used to
// in the "from file" part of the header) or removal (if used in the "to file" part of the header)
const _devNullHeader = "/dev/null"

// _safeExtensions is a set of file extensions (with a dot at the beginning) that are considered
// safe to ignore no matter what changes are made to them, e.g., markdown (".md") files.
var _safeExtensions = map[string]bool{
	".md": true,
}

// Analyzer contains analyzer-specific data.
type Analyzer struct {
	// BaseDiff is location of the base diff file
	BaseDiff string
	// LastDiff is location of the last diff file
	LastDiff string

	// SubAnalyzers contains all available sub-analyzers for different
	// file formats, other than the fallback analyzer.
	SubAnalyzers []common.Analyzer

	// FallbackAnalyzer handles all files not handled by more
	// specialized analyzers
	FallbackAnalyzer fallback.Analyzer
}

// diffHeader represents the per-file diff header in the normalized diff file.
type diffHeader struct {
	from string // in unified diff it's the file name following "---"
	to   string // in unified diff it's the file name following "+++"
}

// Run is the main driver function for the analyzer assuming that
// source files representing base diff and last diff are located in
// separate directories. It compares the files reflecting changes in
// the base diff and the last diff to determine if the files are
// equivalent (in which case the function returns const value Approve)
// or not (in which case the function returns const value Reject). If
// function call result in an error, for convenience, const value
// Failure will be return as a result.
func (a *Analyzer) Run(baseDir string, lastDir string) (int, error) {
	res, filesToAnalyze := a.Setup()
	if res != Approve {
		return res, nil
	}

	for _, sub := range a.SubAnalyzers {
		err := sub.BaseIRBuild(filesToAnalyze, baseDir)
		if err != nil {
			return Failure, err
		}
		err = sub.LastIRBuild(filesToAnalyze, lastDir)
		if err != nil {
			return Failure, err
		}
	}

	err := a.FallbackAnalyzer.BaseIRBuild(filesToAnalyze, baseDir)
	if err != nil {
		return Failure, err
	}
	err = a.FallbackAnalyzer.LastIRBuild(filesToAnalyze, lastDir)
	if err != nil {
		return Failure, err
	}

	return a.Analyze()
}

// Setup sets up the analyzer before the analysis can commence.
func (a *Analyzer) Setup() (int, []string) {
	// find all files subject to analysis
	filesToAnalyze, err := a.collectFilesToAnalyze()
	if err != nil {
		return Failure, filesToAnalyze
	}
	if filesToAnalyze == nil {
		return Reject, filesToAnalyze
	}

	// add fallback analyzer - this should happen after files are
	// collected for analysis above
	a.FallbackAnalyzer = fallback.Analyzer{SubAnalyzers: a.SubAnalyzers}

	return Approve, filesToAnalyze
}

// Analyze analyzes diffs and (if necessary) source files after the
// all the intermediate representation needed for analysis is already
// constructed.
func (a *Analyzer) Analyze() (int, error) {
	// check equivalence for all supported analyzers
	for _, sub := range a.SubAnalyzers {
		eq, err := sub.ChangesEq()
		if err != nil {
			return Failure, err
		}
		if !eq {
			return Reject, nil
		}
	}
	eq, err := a.FallbackAnalyzer.ChangesEq()
	if err != nil {
		return Failure, err
	}
	if !eq {
		return Reject, nil
	}

	return Approve, nil
}

// collectFilesToAnalyze returns names of files that need to be
// analyzed, nil if analysis is infeasible (e.g. because files got
// removed), or empty list of files if changes are the same in the
// base diff and in the last diff for all files.
func (a *Analyzer) collectFilesToAnalyze() ([]string, error) {
	baseDiffFileLines, err := splitDiff(a.BaseDiff)
	if err != nil {
		return nil, err
	}

	lastDiffFileLines, err := splitDiff(a.LastDiff)
	if err != nil {
		return nil, err
	}

	// headers for equal normalized per-modified-file diffs
	equalDiffs := make(map[diffHeader]bool)
	// find per-file diffs that have the same header in both the base
	// diff and the last diff, and whose per-file (line) content is
	// identical (in other words, find per-file changes that are
	// identical in both diffs)
OUTER:
	for baseHeader, baseLines := range baseDiffFileLines {
		if lastLines, ok := lastDiffFileLines[baseHeader]; !ok {
			// either base diff or last diff is missing a given header
			continue
		} else if len(baseLines) != len(lastLines) {
			// number of lines in per-file diffs is different
			continue
		} else {
			for i, l := range baseLines {
				if l != (lastLines)[i] {
					// actual line difference between per-file diffs
					continue OUTER
				}
			}
		}
		equalDiffs[baseHeader] = true
	}

	// find out if the remaining (non-identical) per-file diffs can be
	// successfully compared; in order for this to happen, they should
	// represent a file modification or a file addition in the base diff
	// (we currently do not support other types of changes)

	baseModFiles := a.collectModifiedFiles(equalDiffs, baseDiffFileLines)
	lastModFiles := a.collectModifiedFiles(equalDiffs, lastDiffFileLines)

	baseAddedFiles := a.collectAddedFiles(equalDiffs, baseDiffFileLines)
	lastAddedFiles := a.collectAddedFiles(equalDiffs, lastDiffFileLines)

	// for all non-equal per-file diffs (not captured above), both in
	// the base diff and in the last diff, record file name if
	// per-file diff represents a file modification or a file addition
	// in the base diff, otherwise fail immediately

	affectedFiles, onlyAffectedFiles := collectAffectedFiles(equalDiffs,
		baseDiffFileLines, lastDiffFileLines,
		baseAddedFiles, lastAddedFiles,
		baseModFiles, lastModFiles)

	if !onlyAffectedFiles {
		return nil, nil
	}
	// we actually need an empty array to distinguish from nil value
	// being returned
	filesToAnalyze := []string{}
	// return an array so that the future iteration order is fixed as
	// it will be used to construct two intermediate representations
	// for comparison
	for name := range affectedFiles {
		if _, ok := _safeExtensions[filepath.Ext(name)]; ok {
			continue
		}
		filesToAnalyze = append(filesToAnalyze, name)
	}
	return filesToAnalyze, nil
}

// parseDiff analyzes normalized diff file and returns its content
// split into per-file lines. It returns a mapping from the modified file
// pair (from and to files) to the lines in the diff file representing
// modifications for this file.
func splitDiff(diffFileName string) (map[diffHeader][]string, error) {
	f, err := os.Open(diffFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open diff file %q: %v", diffFileName, err)
	}
	defer f.Close()

	diffFiles, err := splitDiffImpl(f)
	if err != nil {
		return nil, fmt.Errorf("cannot parse normalized diff file %q: %v", diffFileName, err)
	}

	return diffFiles, nil
}

// splitDiffImpl implements the actual splitting logic for splitDiff function
func splitDiffImpl(diffReader io.Reader) (map[diffHeader][]string, error) {
	// the general idea here is to do the following when reading the
	// normalized diff file line-by-line:
	// - find a start of the per-file diff header (subsequent lines starting with "---" and "+++")
	// - store all the following lines (until encountering another header) in a per-header map
	// normalized diffs do not conform to "standard" unified diff file format (e.g. no chunks)
	// which prevents the use of off-the-shelf parsing solutions

	diffFileLines := make(map[diffHeader][]string)
	var currentMapKey diffHeader

	scanner := bufio.NewScanner(diffReader)
	if !scanner.Scan() {
		return diffFileLines, scanner.Err()
	}

	str := scanner.Text()
	header, err := attemptHeaderRead(str, scanner)
	if err != nil {
		return nil, fmt.Errorf("bad header: %v", err)
	}
	if header == nil {
		// normalized diff file should start with a file header
		return nil, errors.New("no starting diff header")
	}
	currentMapKey = *header
	for scanner.Scan() {
		str := scanner.Text()
		header, err := attemptHeaderRead(str, scanner)
		if err != nil {
			return nil, fmt.Errorf("bad header: %v", err)
		}
		if header != nil {
			// a new header has been found
			currentMapKey = *header
		} else {
			// a non-header line - add it to the map
			diffFileLines[currentMapKey] = append(diffFileLines[currentMapKey], str)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return diffFileLines, nil
}

// attemptHeaderRead checks if a string passed as an argument is a start of a
// header and either parses and returns the rest of the header, or
// returns nil (to indicate that the string represents a non-header
// line).
func attemptHeaderRead(str string, scanner *bufio.Scanner) (*diffHeader, error) {
	if !strings.HasPrefix(str, _diffHeaderFromPrefix) {
		// the argument string  does not represent start of a header
		return nil, nil
	}
	fromFile := str[len(_diffHeaderFromPrefix):]
	if !scanner.Scan() {
		return nil, multierr.Append(errors.New("unexpected end of diff file"), scanner.Err())

	}
	str = scanner.Text()
	if !strings.HasPrefix(str, _diffHeaderToPrefix) {
		// per-file diff header should consist of two lines, one
		// starting with "---" and the other one starting with "+++"
		return nil, errors.New("incomplete diff header")
	}
	toFile := str[len(_diffHeaderToPrefix):]
	return &diffHeader{fromFile, toFile}, nil
}

// collectModifiedFiles determines which changes in a given diff
// represent modified files, and returns a set of such modified files.
func (a *Analyzer) collectModifiedFiles(equalDiffs map[diffHeader]bool,
	diffFileLines map[diffHeader][]string) map[string]bool {
	modifiedFiles := make(map[string]bool)
	for header := range diffFileLines {
		if equalDiffs[header] {
			// change represented by this diff has already been compared
			continue
		}
		if header.from == header.to {
			modifiedFiles[header.to] = true
		}
	}
	return modifiedFiles
}

// collectAddedFiles determines which changes in a given diff
// represent added files, and returns a set of such added files.
func (a *Analyzer) collectAddedFiles(equalDiffs map[diffHeader]bool,
	diffFileLines map[diffHeader][]string) map[string]bool {
	addedFiles := make(map[string]bool)
	for header := range diffFileLines {
		if equalDiffs[header] {
			// change represented by this diff has already been compared
			continue
		}
		if header.from == _devNullHeader {
			addedFiles[header.to] = true
		}
	}
	return addedFiles
}

// collectAffectedFiles determines which changes in the diffs reflect
// file modifications or additions in the base diff (as opposed to
// file removal, renaming, addition in the last diff only), and
// returns a set of the files affected this way.  If and only if all
// files are affected this way, the second returned value is true,
// otherwise (e.g., a file was removed) it is false.
func collectAffectedFiles(equalDiffs map[diffHeader]bool,
	baseDiffFilesMap map[diffHeader][]string,
	lastDiffFilesMap map[diffHeader][]string,
	baseAddedFiles map[string]bool,
	lastAddedFiles map[string]bool,
	baseModFiles map[string]bool,
	lastModFiles map[string]bool) (map[string]bool, bool) {
	affectedFiles := make(map[string]bool)
	for header := range baseDiffFilesMap {
		if equalDiffs[header] {
			// change represented by this diff has already been compared
			continue
		}
		if baseModFiles[header.from] {
			// base diff file modification
			affectedFiles[header.from] = true
			continue
		}
		if baseAddedFiles[header.to] && lastAddedFiles[header.to] {
			// base diff file addition (use "to" file as "from" file
			// is /dev/null (must also show up as addition in last
			// diff, even if it's modified in the last diff, otherwise
			// it's a removal)
			affectedFiles[header.to] = true
			continue
		}
		// everything else (removals, renamings) files is not OK
		return affectedFiles, false
	}
	for header := range lastDiffFilesMap {
		if equalDiffs[header] {
			// change represented by this diff has already been compared
			continue
		}
		if lastModFiles[header.from] {
			// last diff file modification
			affectedFiles[header.from] = true
			continue
		}
		if lastAddedFiles[header.to] && baseAddedFiles[header.to] {
			// last diff file addition (use "to" file as "from" file is /dev/null
			// file must be also added in the base diff, but it's already captured
			// when analyzing base diff files
			continue
		}
		// everything else (removals, renamings, additions only in the
		// last diff) is not OK
		return affectedFiles, false
	}
	return affectedFiles, true
}
