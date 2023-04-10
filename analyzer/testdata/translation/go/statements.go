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

func root() {
	defer foo(bar)
	continue
	continue here
	break
	break there
	goto label1
	return
	return a
	return a, b
	c <- v
	go add(1, 2)
	a = b
	a, b = c, d
	a := 1
	a, b := c, d

	switch a = foo(); a {
	case "1", 2:
		foo()
		bar()
		var a, b = 1, 2
		fallthrough
	case 3:
		test()
	default:
		/* do nothing */
	}

	switch {
	case a < b:
		return 1
	}

	switch a := 1; n := c.(type) {
	}

	switch (&c).(type) {
	case *Test:
		/* do nothing */
	default:
		foo()
	}

	switch a++; a {}

	if a {
		foo()
	}

	if a {
		t1()
	} else if b {
		t2()
	} else {
		t3()
	}

	if a := 1; a {
	}
	if a = 1; a {

	}
hello:
	var a, b int
	type T = map[string]bool
	type T int

	var (
		a, b int = 1, 2
		c        = 3
	)
	var d = 4
	var e, f int
	var h, i int = foo()

	const (
		a, b int = 1, 2
		c        = 3
		d
	)
	const e = 4

	type Test struct {
		*A
		B    T
		C    int "tag"
		D, E int `t`
	}

	for i := 1; i < 10; i++ {
		foo()
	}

	for {  // no for_clause at all
		foo()
	}

	for ;; { // empty for_clause
		foo()
	}

	for hasNext() {
		foo()
	}

	for range lst {
	}

	for i, k := range lst {
	}

	for pkg.A = range lst {
	}

	select {
	case a, b := <-ch:
		foo()
	var x, y int = 1, 2
	case c = <-x:
		bar()
	case <-quit:
	default:
		test()
	}

	select {
	default: // intentional empty default case
	}

	; // empty_statement, should not appear in the final MAST
	empty_label: // an empty label without a statement
}
