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

package test

// breakLabelTest tests symbolication of "break" labels named the same
// as other language-level constructs.
func breakLabelTest() int {
	tmp1 := 42

tmp1:
	for i := 0; i < 7; i++ {
		tmp1 = tmp1 + 1
		break tmp1
	}

	return tmp1
}

// contLabelTest tests symbolication of "continue" labels named the same
// as other language-level constructs.
func contLabelTest() int {
	tmp2 := 42

tmp2:
	for i := 0; i < 7; i++ {
		tmp2 = tmp2 + 1
		continue tmp2
	}

	return tmp2
}

// gotoLabelTest tests symbolication of "goto" labels named the same
// as other language-level constructs.
func gotoLabelTest() int {
	tmp3 := 42

tmp3:
	for i := 0; i < 7; i++ {
		tmp3 = tmp3 + 1
		if tmp3 == 1 {
			goto tmp3
		}
	}

	return tmp3
}
