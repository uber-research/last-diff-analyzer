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

// this test includes single and multiple imports and a comment to be ignored during translation
package rename

import (
	. "example/package1"
	t "package2"
	_ "package3"
	"package4"
)

import "singlepackage"

import () //empty declaration, will simply be dropped

func (a *A) test(a, b int) (c, d string) {
	foo()
}

func test(a, b int) (c, d string) {
	foo()
}

type Test interface {
	Embedded
	hello(a int) (b int)
}

type (
	A = B // type alias
	C D   // type spec
)

const () // empty declaration, will simply be dropped
var ()   // same as above
type ()  // same as above
