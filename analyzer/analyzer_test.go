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

package analyzer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"analyzer/bazel"
	"analyzer/core"
	"analyzer/gomod"
	"analyzer/protobuf"
	"analyzer/sql"
	"analyzer/starlark"
	"analyzer/thrift"
	"analyzer/yaml"
)

const _metaTestDataPrefix = "./testdata/analyzer/"

// featureFlags stores the set of feature flags for the analyzer.
type featureFlags uint32

const (
	// _logging indicates whether auto-approvals for logging-related changes are enabled.
	_logging featureFlags = 1 << iota
)

func TestAnalyzer(t *testing.T) {
	// Default set of feature flags.
	defaultFlags := _logging

	testCases := []struct {
		testDir  string
		expected int
		flags    featureFlags
	}{
		{
			// test equal files
			testDir:  _metaTestDataPrefix + "go/equal/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test arbitrary (comparable only based on their "raw"
			// content) equal files when the patches themselved are
			// different (emulates a rebase case)
			testDir:  _metaTestDataPrefix + "equal-arbitrary/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test a non-equivalent change in the last diff
			testDir:  _metaTestDataPrefix + "go/reject/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test a non-equivalent change in the last diff where all
			// modified files are arbitrary (not supported) files
			testDir:  _metaTestDataPrefix + "reject-arbitrary-all/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test a non-equivalent change in the last diff
			// where some modified files are non-go files
			// and the changes in Go files are not semantically equivalent
			testDir:  _metaTestDataPrefix + "go/reject-nogo-some/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test a non-equivalent change in the last diff
			// where some modified files are non-go files
			// and the changes in Go files are semantically equivalent
			testDir:  _metaTestDataPrefix + "go/reject-nogo-some-rest-equal/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of comments added in the last diff
			testDir:  _metaTestDataPrefix + "go/comment/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of formatting changes in the last diff
			testDir:  _metaTestDataPrefix + "go/format/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence in the case of added (instead of modfied) files in the base diff
			// (using comments as an example of such change)
			testDir:  _metaTestDataPrefix + "go/added-base/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence in the case of added (instead of modfied) files in the last diff
			// (using comments as an example of such change)
			testDir:  _metaTestDataPrefix + "go/reject-added-last/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of simple variable renaming
			testDir:  _metaTestDataPrefix + "go/simple-var-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) simple variable renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-simple-var-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) simple variable renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-simple-var-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of nested variable renaming
			testDir:  _metaTestDataPrefix + "go/nested-var-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant d)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant e)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant f)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant g)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-g/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) nested variable renaming (variant h)
			testDir:  _metaTestDataPrefix + "go/reject-nested-var-rename-h/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of code where a global variable is used in base diff but not in last diff
			testDir:  _metaTestDataPrefix + "go/global-var-nouse/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of global variable renaming
			testDir:  _metaTestDataPrefix + "go/global-var-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) global variable renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-global-var-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) global variable renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-global-var-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) global variable renaming (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-global-var-rename-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of various other variable renamings
			testDir:  _metaTestDataPrefix + "go/other-var-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of parameter renaming
			testDir:  _metaTestDataPrefix + "go/param-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) parameter renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-param-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) parameter renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-param-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of import alias renaming
			testDir:  _metaTestDataPrefix + "go/import-alias-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) import alias renaming
			testDir:  _metaTestDataPrefix + "go/reject-import-alias-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of function renaming
			testDir:  _metaTestDataPrefix + "go/func-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) function renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-func-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) function renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-func-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of constant renaming
			testDir:  _metaTestDataPrefix + "go/const-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant d)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant e)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming (variant f)
			testDir:  _metaTestDataPrefix + "go/reject-const-rename-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming for enums (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-const-enum-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming for enums (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-const-enum-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) constant renaming for enums (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-const-enum-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of type renaming
			testDir:  _metaTestDataPrefix + "go/type-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant d)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant e)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) type renaming (variant f)
			testDir:  _metaTestDataPrefix + "go/reject-type-rename-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of comments-related change in BUILD.bazel file
			testDir:  _metaTestDataPrefix + "bazel/approve/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of formatting-related change in BUILD.bazel file
			testDir:  _metaTestDataPrefix + "bazel/reject/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of comments-related change in yaml file
			testDir:  _metaTestDataPrefix + "yaml/approve/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of formatting-related change in yaml file
			testDir:  _metaTestDataPrefix + "yaml/reject/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of const removal
			testDir:  _metaTestDataPrefix + "go/const-remove/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of const addition (to replace a literal use)
			testDir:  _metaTestDataPrefix + "go/const-add/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) global const removal
			testDir:  _metaTestDataPrefix + "go/reject-const-remove/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const addition (to replace a literal use but with a different value)
			testDir:  _metaTestDataPrefix + "go/reject-const-add-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) exported const addition (to replace a literal use)
			testDir:  _metaTestDataPrefix + "go/reject-const-add-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const modification
			// that does not have uses in any of the diffs
			testDir:  _metaTestDataPrefix + "go/reject-const-mod-same-name/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) approval of stmt list of different length
			testDir:  _metaTestDataPrefix + "go/reject-more-stmt/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			testDir:  _metaTestDataPrefix + "go/struct-create/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) struct creation
			testDir:  _metaTestDataPrefix + "go/reject-struct-create-g/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) passing of variadic argument with and without the trailing ...
			testDir:  _metaTestDataPrefix + "go/reject-variadic/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of logging-related changes
			testDir:  _metaTestDataPrefix + "go/logging/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of logging-related changes when logging support is off
			testDir:  _metaTestDataPrefix + "go/logging/",
			expected: Reject,
			flags:    defaultFlags & ^_logging,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant a)
			testDir:  _metaTestDataPrefix + "go/reject-logging-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant b)
			testDir:  _metaTestDataPrefix + "go/reject-logging-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant c)
			testDir:  _metaTestDataPrefix + "go/reject-logging-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant d)
			testDir:  _metaTestDataPrefix + "go/reject-logging-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant e)
			testDir:  _metaTestDataPrefix + "go/reject-logging-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant f)
			testDir:  _metaTestDataPrefix + "go/reject-logging-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant g)
			testDir:  _metaTestDataPrefix + "go/reject-logging-g/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant h)
			testDir:  _metaTestDataPrefix + "go/reject-logging-h/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant i)
			testDir:  _metaTestDataPrefix + "go/reject-logging-i/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant j)
			testDir:  _metaTestDataPrefix + "go/reject-logging-j/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant k)
			testDir:  _metaTestDataPrefix + "go/reject-logging-k/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant l)
			testDir:  _metaTestDataPrefix + "go/reject-logging-l/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) logging-related changes (variant m)
			testDir:  _metaTestDataPrefix + "go/reject-logging-m/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of comments-related change in go.mod file
			testDir:  _metaTestDataPrefix + "gomod/approve/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of formatting-related change in go.mod file
			testDir:  _metaTestDataPrefix + "gomod/reject/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of last statements in blocks
			testDir:  _metaTestDataPrefix + "go/reject-last-stmt/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incompatible) access path change
			testDir:  _metaTestDataPrefix + "go/reject-access-path/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of files residing in different packages
			testDir:  _metaTestDataPrefix + "go/rename-different-package/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of changes to the first components of access paths
			testDir:  _metaTestDataPrefix + "go/access-path/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of changes to a local upper-case prefixed variable
			testDir:  _metaTestDataPrefix + "go/upper-case-local-variable/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of renames for Java
			testDir:  _metaTestDataPrefix + "java/rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test a incompatible change of changing the behavior while renaming the variable
			testDir:  _metaTestDataPrefix + "java/reject-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test renaming inside nested classes
			testDir:  _metaTestDataPrefix + "java/nested-classes/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// TODO: remove this test when proper symbolications for access paths are implemented
			// test rejection of renaming an identifier within an access path.
			testDir:  _metaTestDataPrefix + "java/reject-access-path/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of changes _inside_ a public class
			testDir:  _metaTestDataPrefix + "java/public-class/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test incompatible changes of renaming the public class
			testDir:  _metaTestDataPrefix + "java/reject-public-class/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test mixed privateness of scopes
			testDir:  _metaTestDataPrefix + "java/scopes/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of completely identical files
			testDir:  _metaTestDataPrefix + "java/equal/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of changes to modifiers (e.g., adding a final modifier)
			testDir:  _metaTestDataPrefix + "java/modifier/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test removing a final modifier
			testDir:  _metaTestDataPrefix + "java/reject-modifier-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test adding a private modifier
			testDir:  _metaTestDataPrefix + "java/reject-modifier-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test removing a private modifier
			testDir:  _metaTestDataPrefix + "java/reject-modifier-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test adding a static modifier
			testDir:  _metaTestDataPrefix + "java/reject-modifier-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test removing a static modifier
			testDir:  _metaTestDataPrefix + "java/reject-modifier-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test adding a final modifier with an annotation
			testDir:  _metaTestDataPrefix + "java/reject-modifier-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of parameter renaming
			testDir:  _metaTestDataPrefix + "java/param-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of adding comments
			testDir:  _metaTestDataPrefix + "java/comment/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test renaming of a public field inside a public class
			testDir:  _metaTestDataPrefix + "java/reject-scopes-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test renaming of package-visible class
			testDir:  _metaTestDataPrefix + "java/reject-scopes-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant a
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant b
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant c
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant d
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant e
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant f
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant g
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-g/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test rejection of renaming a different nested variable - variant h
			testDir:  _metaTestDataPrefix + "java/reject-nested-variable-h/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of renaming constants
			testDir:  _metaTestDataPrefix + "java/const-rename/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of adding constants
			testDir:  _metaTestDataPrefix + "java/const-add/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test equivalence of removing constants
			testDir:  _metaTestDataPrefix + "java/const-remove/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const renaming
			testDir:  _metaTestDataPrefix + "java/reject-const-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const addition - variant a
			testDir:  _metaTestDataPrefix + "java/reject-const-add-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const addition - variant b
			testDir:  _metaTestDataPrefix + "java/reject-const-add-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const addition - variant c
			testDir:  _metaTestDataPrefix + "java/reject-const-add-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) const removal
			testDir:  _metaTestDataPrefix + "java/reject-const-remove/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test equivalence of logging-related changes for Java
			testDir:  _metaTestDataPrefix + "java/logging/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of logging-related changes for Java when logging support is off
			testDir:  _metaTestDataPrefix + "java/logging/",
			expected: Reject,
			flags:    defaultFlags & ^_logging,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant a
			testDir:  _metaTestDataPrefix + "java/reject-logging-a/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant b
			testDir:  _metaTestDataPrefix + "java/reject-logging-b/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant c
			testDir:  _metaTestDataPrefix + "java/reject-logging-c/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant d
			testDir:  _metaTestDataPrefix + "java/reject-logging-d/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant e
			testDir:  _metaTestDataPrefix + "java/reject-logging-e/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) logging-related changes for Java - variant f
			testDir:  _metaTestDataPrefix + "java/reject-logging-f/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of (incorrect) change in initializer part of an if statement
			testDir:  _metaTestDataPrefix + "go/reject-if-initializer/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test comment-only changes to sql files
			testDir:  _metaTestDataPrefix + "sql/approve/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of statement addition to a sql file
			testDir:  _metaTestDataPrefix + "sql/reject/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test comment-only changes to protobuf files
			testDir:  _metaTestDataPrefix + "protobuf/comment/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of field renaming to a protobuf file
			testDir:  _metaTestDataPrefix + "protobuf/reject-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of value changing to a protobuf file
			testDir:  _metaTestDataPrefix + "protobuf/reject-value/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test comment-only changes to starlark files
			testDir:  _metaTestDataPrefix + "starlark/comment/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of statement addition to a starlark file
			testDir:  _metaTestDataPrefix + "starlark/reject-add/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of statement removal to a starlark file
			testDir:  _metaTestDataPrefix + "starlark/reject-remove/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of variable rename to a starlark file (currently unsupported)
			testDir:  _metaTestDataPrefix + "starlark/reject-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test comment-only changes to thrift files
			testDir:  _metaTestDataPrefix + "thrift/comment/",
			expected: Approve,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of field addition to a thrift file
			testDir:  _metaTestDataPrefix + "thrift/reject-add/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of field removal to a thrift file
			testDir:  _metaTestDataPrefix + "thrift/reject-remove/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test non-equivalence of key rename to a thrift file
			testDir:  _metaTestDataPrefix + "thrift/reject-rename/",
			expected: Reject,
			flags:    defaultFlags,
		},
		{
			// test ignore of markdown files
			testDir:  _metaTestDataPrefix + "markdown/modify/",
			expected: Approve,
			flags:    defaultFlags,
		},
	}
	for _, tc := range testCases {
		baseDiff := tc.testDir + "base.diff"
		lastDiff := tc.testDir + "last.diff"
		baseDir := tc.testDir + "base"
		lastDir := tc.testDir + "last"
		t.Run(fmt.Sprintf("%s == %s", strings.TrimPrefix(baseDiff, _metaTestDataPrefix),
			strings.TrimPrefix(lastDiff, _metaTestDataPrefix)), func(t *testing.T) {
			a := Analyzer{
				BaseDiff: baseDiff,
				LastDiff: lastDiff,
			}
			// add core analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &core.Analyzer{RenamingOn: true, LoggingOn: tc.flags&_logging != 0})
			// add Bazel files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &bazel.Analyzer{AnalyzableFileName: "BUILD.bazel.test"})
			// add yaml files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &yaml.Analyzer{})
			// add go.mod files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &gomod.Analyzer{})
			// add SQL files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &sql.Analyzer{})
			// add protobuf files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &protobuf.Analyzer{})
			// add starlark files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &starlark.Analyzer{})
			// add thrift files analyzer
			a.SubAnalyzers = append(a.SubAnalyzers, &thrift.Analyzer{})

			status, err := a.Run(baseDir, lastDir)
			require.NoError(t, err)
			require.Equal(t, tc.expected, status)
		})
	}
}
