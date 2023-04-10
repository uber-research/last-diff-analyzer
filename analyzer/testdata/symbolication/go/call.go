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

type someStruct struct{}
type struct2 struct{}

// We intentionally create a name collision here with the function pointer
// variable declaration to test the ability of reasoning variable shadowing in our
// symbolication process.
type ptrfoo1 struct{}

func (*someStruct) foo1() int {
	return 7
}

func foo1() int {
	return 42
}

func bar() int {
	s := someStruct{}
	// These are type casts really, but note that tree-sitter and our system
	// will treat these as a call expression.
	a := (*struct2)(&s)
	b := (*test.struct2)(&s)

	// Note that the following uses, although syntactical similar to type casts,
	// it is actually a function call.
	foo := func(i int) int {
		return i
	}
	ptrfoo1 := &foo
	x := (*ptrfoo1)(4)
	print(x)

	return foo1() + s.foo1()
}
