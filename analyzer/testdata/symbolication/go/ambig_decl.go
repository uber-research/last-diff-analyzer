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

type tmp1 struct{}

// sameLineDecl tests if use/def link for a variable used and declared
// on the same line is created correctly.
func sameLineDecl() tmp1 {
	var tmp1 = tmp1{}
	return tmp1
}

type tmp2 struct{}

// sameLineShortDecl tests if use/def link for a variable used and
// declared on the same line is created correctly.
func sameLineShortDecl() tmp2 {
	tmp2 := tmp2{}
	return tmp2
}

type tmp3 struct {
	a1 []int
}

// rangeDecl tests if use/def link for a variable used and
// declared in the range loop is created correctly.
func rangeDecl() int {
	acc := 0
	for _, tmp3 := range (tmp3{[]int{7, 42}}).a1 {
		acc = acc + tmp3
	}
	return acc
}
