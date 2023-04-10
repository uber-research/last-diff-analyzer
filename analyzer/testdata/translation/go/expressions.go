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
	nil
	true
	false
	"test\t"
	`test`
	123
	1.5
	5i
	'a'
	a*b + c
	!a
	!(a)
	a[i]
	a.b.c
	a.foo().b
	foo().b
	foo().a.b
	a[i : i+1]
	a[i : i+1 : 10]
	a[:]
	a[i:]
	foo()
	add(a, b)
	foo.bar(a, b)
	add(a, b...)
	make(int, 10)
	make([5]int, 10)
	make([]int, 10)
	make((map[string]bool), 10)
	make(*int, 10)
	make(pkg.Test, 10)
	make(chan int, 10)
	make(<-chan int, 10)
	make(chan<- int, 10)
	make(func ())
	make(func (A, B) C)
	make(func (a A, b B) (c C, d D))
	make(func () (C, D))
	make(func (A, B))
	make(func (a... A) (b... B))
	make(func (a A, ...B))
	a.(T)
	[]byte(a)
	a++
	a--

	a := []int{1, 2, 3}
	a := [...]int{1, 2, 3}
	a := Test{
		Parent: {

		},
		Key: value,
	}

	f := func(x, y int) int { return x + y }
	func () { return 1 }() // invoke directly
	func () (a A, b B) { return 1 }()
	func () (A, B) { return 1 }() // invoke directly
}
