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

	"github.com/stretchr/testify/require"

	"analyzer/core/mast"
)

func TestHasJavaModifier(t *testing.T) {
	t.Run("Test HasJavaModifier", func(t *testing.T) {
		modifiers := []mast.Expression{
			&mast.JavaLiteralModifier{Modifier: mast.PrivateMod},
			&mast.Annotation{Name: &mast.Identifier{Name: "Foo"}},
		}
		require.True(t, HasJavaModifier(modifiers, mast.PrivateMod))
		require.False(t, HasJavaModifier(modifiers, mast.FinalMod))

		modifiers = append(modifiers, &mast.JavaLiteralModifier{Modifier: mast.FinalMod})
		require.True(t, HasJavaModifier(modifiers, mast.FinalMod))
		require.False(t, HasJavaModifier(modifiers, mast.AbstractMod))
	})
}
