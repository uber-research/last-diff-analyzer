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

package mastutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"analyzer/core/mast"
)

func TestJoinAccessPath(t *testing.T) {
	testcases := []struct {
		path     *mast.AccessPath
		expected string
	}{
		// Foo.a
		{
			path: &mast.AccessPath{
				Operand: &mast.Identifier{Name: "Foo"},
				Field:   &mast.Identifier{Name: "a"},
			},
			expected: "Foo.a",
		},
		// Foo.a[b].c
		{
			path: &mast.AccessPath{
				Operand: &mast.IndexExpression{
					Operand: &mast.AccessPath{
						Operand: &mast.Identifier{Name: "Foo"},
						Field:   &mast.Identifier{Name: "a"},
					},
					Index: &mast.Identifier{Name: "b"},
				},
				Field: &mast.Identifier{Name: "c"},
			},
			expected: "Foo.a.c",
		},
	}

	for _, tc := range testcases {
		t.Run("Test JoinAccessPath", func(t *testing.T) {
			str, err := JoinAccessPath(tc.path)
			require.NoError(t, err)
			if diff := cmp.Diff(tc.expected, str); diff != "" {
				require.FailNow(t, "mismatch (-expected +actual)", diff)
			}
		})
	}
}
