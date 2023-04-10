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

var global1 int = 1

type someStruct struct {
	b1 int
}

// simple test a very simple access path
func simple() int {
	a1 := someStruct{42}
	return a1.b1
}

type anotherStruct struct {
	tmp1 int
}

func foo1() anotherStruct {
	return anotherStruct{42}
}

// funcStartPath tests access path starting with a function call
func funcStartPath() int {
	tmp1 := 7
	return foo1().tmp1 + (foo1()).tmp1 + tmp1
}

type yetAnotherStruct struct {
	tmp2 int
}

func (s *yetAnotherStruct) bar1() *yetAnotherStruct {
	return s
}

// funcMidPath tests access path a function (method) call in the middle
func funcMidPath() int {
	c1 := yetAnotherStruct{42}
	return c1.bar1().tmp2
}

type oneMoreStruct struct {
	tmp3 int
}

func foo3() oneMoreStruct {
	return oneMoreStruct{42}
}

func bar3() func() oneMoreStruct {
	return foo3
}

// funcRetPath tests access path with a function returning another function
func funcRetPath() int {
	return bar3()().tmp3
}

type yetOneMoreStruct struct {
	arr1 []oneMoreStruct
}

func foo4() yetOneMoreStruct {
	return yetOneMoreStruct{nil}
}

func (*yetOneMoreStruct) bar4() func() yetOneMoreStruct {
	return foo4
}

// funcRetMidPath tests access path with a function returning another function in the middle of the path
func funcMidRetPath() int {
	c4 := yetOneMoreStruct{make([]oneMoreStruct, 3)}
	return c4.bar4()().arr1[global1].tmp3
}
