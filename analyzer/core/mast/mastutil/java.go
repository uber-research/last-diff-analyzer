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
	"analyzer/core/mast"
)

// HasJavaModifier checks if the slice of expressions contains the given Java modifier.
func HasJavaModifier(modifiers []mast.Expression, modifier string) bool {
	for _, m := range modifiers {
		if lit, ok := m.(*mast.JavaLiteralModifier); ok && lit.Modifier == modifier {
			return true
		}
	}
	return false
}
