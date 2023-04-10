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

import "reflect"

func bar() (int, int) {
	return 7, 42
}

func baz() (int, interface{}) {
	return 7, 42
}

var tmp1 int = 42

// fntest tests symbolication of short declaration for functions.
//
// param1 declaration and its uses should be in the same scope to
// guarantee that short assignment to param1 is not treated as an
// actual (another) declaration of param1.
func fntest(param1 int) int {
	param1 = 0
	param1, tmp1 := bar()
	return param1 + tmp1
}

// fortestsimple tests simple symbolication of short declaration for for statements.
//
// There are three private scopes here: one for function, one for
// loop's init stmt, and one for the loop body.
func fortestsimple() int {
	acc := 0
	for i2 := 0; i2 < 42; i2++ {
		i2, tmp2 := bar()
		acc = acc + i2 + tmp2
	}
	return acc
}

// fortest tests symbolication of short declaration for for statements.
//
// There are three private scopes here: one for function, one for
// loop's init stmt, and one for the loop body.
func fortest() int {
	i3 := 0
	acc := i3
	for i3, tmp3 := bar(); i3 < 42; i3++ {
		tmptmp3 := 42
		i3, tmptmp3 := bar()
		acc = acc + i3 + tmp3 + tmptmp3
	}
	return acc
}

// nestedfntest tests symbolication of short declaration for nested functions.
//
// param4 declaration and its uses in the nested function should be in
// the same scope to guarantee that short assignment to param1 is not
// treated as an actual (another) declaration of param1.
func nestedfntest() int {

	nested := func(param4 int) int {
		param4 = 0
		param4, tmp4 := bar()
		return param4 + tmp4
	}
	return nested(42)
}

// switchtest tests symbolication of short declaration for switch statements.
func switchtest() int {
	tmp5 := 42

	switch tmp5, i5 := bar(); tmp5 + i5 {
	case 49:
		tmp5, i5 := bar()
		return tmp5 + i5
	}

	return tmp5
}

// typeswitchtest tests symbolication of short declaration for type switch statements.
func typeswitchtest() (int, interface{}) {
	tmp6 := 42

	switch tmp6, i6 := baz(); i6.(type) {
	case int:
		tmp6, i6 := baz()
		return tmp6, i6
	case string:
		return tmp6, i6
	}

	return tmp6, reflect.TypeOf(tmp6)
}

// forrangetest tests symbolication of short declaration for forrange statements.
func forrangetest() int {
	nums := []int{7, 42}

	acc := 0
	for i7, tmp7 := range nums {
		acc += i7
		i7, tmptmp7 := bar()
		acc = acc + i7 + tmp7 + tmptmp7
	}
	return acc
}

// selecttest tests symbolication of short declaration for select statements.
func selecttest(c chan int) int {
	tmp8 := 42

	select {
	case tmp8 := <-c:
		tmp8, i8 := bar()
		return tmp8 + i8
	}

	return tmp8
}

// blankIdentifiers tests assigning multiple blank identifiers
func blankIdentifiers() int {
	f := func() (int, int, int) {
		return 1, 2, 3
	}

	// ignore the first and third return values
	_, tmp9, _ := f()
	var _, tmp10, _ int = f()
	return tmp9 + tmp10
}

// ifInitializer tests the symbolication of variables in the initializer part of if statement.
func ifInitializer() {
	// The "tmp11" in the condition should be linked to the "tmp11" in the initializer part.
	if tmp11 := blankIdentifiers(); tmp11 > 5 {
	}
}
