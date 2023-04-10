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

package test1

type S1 struct{}

type S2 struct{}

func f1() error {
	return nil
}

func f2(s1 *S1, s2 S2) {
	f1()
	f3()             // used before its declaration, but should be ok
	s1, err1 := f4() // s1 is a re-declaration and err is a declaration
	if err1 != nil {
	}
	if s1 == s2 {
		if s2 == nil {
		}
		// test shadowing of type declaration
		type S2 struct{}
		// test shadowing of variable declaration
		var s2 S2 = &S2{}
		if s2 == nil {
		}
	}
}

func f3() {}

func f4() (*S1, error) {}
