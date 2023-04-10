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

package symbolication

import (
	"testing"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"
)

// We specifically designed the test file to have a unique number attached to each identifier ("f1, f2..." etc),
// therefore any identifier with a number suffix should have a corresponding link.

func TestSimpleSymbolication(t *testing.T) {
	t.Run("Test symbolication for a simple case", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{
			"test1": 1, "S1": 3, "S2": 5, "f1": 2, "f2": 1, "s1": 3, "s2": 5, "f3": 2, "f4": 2, "err1": 2,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/simple.go", ts.GoExt, expectedLength)

		// now test the accuracies of the links
		// The expectedLinks maps from identifier name to a link map, which is from use identifier node to its expected
		// def identifier node. Most variable links are simple: there are only two identifiers in the test file, first
		// one is the declaration identifier and the second is the usage. So there should be two links, one from the
		// declaration to itself, and the second from the use to the declaration. The special cases will be accompanied
		// by detailed explanations.
		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"test1": {
				names["test1"][0]: nil, // self (pkg decl)
			},
			"S1": {
				names["S1"][0]: names["S1"][0], // self (top-level)
				names["S1"][1]: names["S1"][0], // first use (param type in f2)
				names["S1"][2]: names["S1"][0], // second use (param type in f4)
			},
			// The S2 is actually a type identifier, there are five S2 variables:
			// (1) type declaration in "type S2 struct{}";
			// (2) parameter declaration in "f2(s1 *S1, s2 S2)";
			// (3) shadowing type declaration in "type S2 struct{}";
			// (4)(5) variable declaration and entity creation expression in "var s2 S2 = &S2{}".
			// So the links should be:
			// (1)(2) should be linked to (1); (3)(4)(5) should be linked to (3) due to shadowing.
			"S2": {
				names["S2"][0]: names["S2"][0], // self (top-level)
				names["S2"][1]: names["S2"][0], // first use (param type in f2)
				names["S2"][2]: names["S2"][2], // shadowed declaration in f2
				names["S2"][3]: names["S2"][2], // first use of shadowed declaration in f2
				names["S2"][4]: names["S2"][2], // second use of shadowed declaration in f2
			},
			// Similar to S2, there are in total five s2 identifiers:
			// (1) parameter declaration in "f2(s1 *S1, s2 S2)";
			// (2) binary expression in "if s1 == s2 {...}";
			// (3) another binary expression in "if s2 == nil {}" inside the if above;
			// (4) shadowing variable declaration in "var s2 S2 = &S2{}";
			// (5) binary expression in "if s2 == nil {}" which refers to the shadow declaration in (4).
			// So the links should be:
			// (1)(2)(3) should be linked to (1); (4)(5) should be linked to (4) due to shadowing.
			"s2": {
				names["s2"][0]: names["s2"][0], // self (param)
				names["s2"][1]: names["s2"][0], // first use in f2
				names["s2"][2]: names["s2"][0], // second use in f2
				names["s2"][3]: names["s2"][3], // shadowed declaration in f2
				names["s2"][4]: names["s2"][3], // first use of shadowed declaration in f2
			},
			"f1": {
				names["f1"][0]: names["f1"][0], // self (top-level)
				names["f1"][1]: names["f1"][0], // first use in f2
			},
			"f2": {
				names["f1"][0]: names["f1"][0], // self (top-level)
			},
			"s1": {
				names["s1"][0]: names["s1"][0], // self (param)
				names["s1"][1]: names["s1"][0], // redeclaration (use) in s2
				names["s1"][2]: names["s1"][0], // second use in f2
			},
			"f3": {
				names["f3"][0]: names["f3"][1], // use before decl in f2
				names["f3"][1]: names["f3"][1], // self (top-level)
			},
			"f4": {
				names["f4"][0]: names["f4"][1], // use  before decl in f2
				names["f4"][1]: names["f4"][1], // self (top-level)
			},
			"err1": {
				names["err1"][0]: names["err1"][0], // self (in f2)
				names["err1"][1]: names["err1"][0], // use in f2
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestOrderSymbolication(t *testing.T) {
	t.Run("Test symbolication for unusual use/def order", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"tmp1": 4}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/order.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (top-level)
				names["tmp1"][1]: names["tmp1"][0], // first func use to top-level
				names["tmp1"][2]: names["tmp1"][2], // self (in-func)
				names["tmp1"][3]: names["tmp1"][2], // last func use to in-func decl
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestShortDecSymbolication(t *testing.T) {
	t.Run("Test symbolication for short variable declarations", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{
			"param1": 4, "tmp1": 3, "i2": 5, "tmp2": 2, "i3": 7, "tmp3": 2, "tmptmp3": 3,
			"param4": 4, "tmp4": 2, "tmp5": 6, "i5": 4, "tmp6": 7, "i6": 5,
			"i7": 4, "tmp7": 2, "tmptmp7": 2, "tmp8": 5, "tmp9": 2, "tmp10": 2, "i8": 2,
			"tmp11": 2,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/short_decl.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"param1": {
				names["param1"][0]: names["param1"][0], // self (param)
				names["param1"][1]: names["param1"][0], // first func use
				names["param1"][2]: names["param1"][0], // short assignment (== second func use)
				names["param1"][3]: names["param1"][0], // third func use
			},
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (top-level)
				names["tmp1"][1]: names["tmp1"][1], // self (short assignment == declaration)
				names["tmp1"][2]: names["tmp1"][1], // first use of redeclared var
			},
			"i2": {
				names["i2"][0]: names["i2"][0], // self (loop init stmt)
				names["i2"][1]: names["i2"][0], // use in loop cond
				names["i2"][2]: names["i2"][0], // use in loop post-iter stmt
				names["i2"][3]: names["i2"][3], // self (short assignment == declaration) in loop body
				names["i2"][4]: names["i2"][3], // first use of redeclared var in loop body
			},
			"tmp2": {
				names["tmp2"][0]: names["tmp2"][0], // self (short assignment == declaration) in loop body
				names["tmp2"][1]: names["tmp2"][0], // use of redeclared var in loop body
			},
			"i3": {
				names["i3"][0]: names["i3"][0], // self (function body)
				names["i3"][1]: names["i3"][0], // use in function body
				names["i3"][2]: names["i3"][2], // self (short assignment == declaration) in loop init stmt
				names["i3"][3]: names["i3"][2], // use in loop cond
				names["i3"][4]: names["i3"][2], // use in loop post-iter stmt
				names["i3"][5]: names["i3"][5], // self (short assignment == declaration) in loop body
				names["i3"][6]: names["i3"][5], // first use of redeclared var in loop body
			},
			"tmp3": {
				names["tmp3"][0]: names["tmp3"][0], // self (short assignment == declaration) in loop body
				names["tmp3"][1]: names["tmp3"][0], // use of redeclared var in loop body
			},
			"tmptmp3": {
				names["tmptmp3"][0]: names["tmptmp3"][0], // self (loop body)
				names["tmptmp3"][1]: names["tmptmp3"][0], // short assignment (== second func use) in loop body
				names["tmptmp3"][2]: names["tmptmp3"][0], // second use in loop body
			},
			"param4": {
				names["param4"][0]: names["param4"][0], // self (param)
				names["param4"][1]: names["param4"][0], // first func use
				names["param4"][2]: names["param4"][0], // short assignment (== second func use)
				names["param4"][3]: names["param4"][0], // third func use
			},
			"tmp4": {
				names["tmp4"][0]: names["tmp4"][0], // self (short assignment == declaration)
				names["tmp4"][1]: names["tmp4"][0], // first use of redeclared var
			},
			"tmp5": {
				names["tmp5"][0]: names["tmp5"][0], // self (short assignment == declaration) in function
				names["tmp5"][1]: names["tmp5"][1], // self (short assignment == declaration) in switch init
				names["tmp5"][2]: names["tmp5"][1], // use in switch tag
				names["tmp5"][3]: names["tmp5"][3], // self (short assignment == declaration) in case clause
				names["tmp5"][4]: names["tmp5"][3], // use in case clause
				names["tmp5"][5]: names["tmp5"][0], // use in function
			},
			"i5": {
				names["i5"][0]: names["i5"][0], // self (short assignment == declaration) in switch init
				names["i5"][1]: names["i5"][0], // use in switch tag
				names["i5"][2]: names["i5"][2], // self (short assignment == declaration) in case clause
				names["i5"][3]: names["i5"][2], // use in case clause
			},
			"tmp6": {
				names["tmp6"][0]: names["tmp6"][0], // self (short assignment == declaration) in function
				names["tmp6"][1]: names["tmp6"][1], // self (short assignment == declaration) in switch init
				names["tmp6"][2]: names["tmp6"][2], // self (short assignment == declaration) in case clause
				names["tmp6"][3]: names["tmp6"][2], // use in case clause
				names["tmp6"][4]: names["tmp6"][1], // use in default clause
				names["tmp6"][5]: names["tmp6"][0], // first use in function
				names["tmp6"][6]: names["tmp6"][0], // second use in function
			},
			"i6": {
				names["i6"][0]: names["i6"][0], // self (short assignment == declaration) in switch init
				names["i6"][1]: names["i6"][0], // use in switch tag
				names["i6"][2]: names["i6"][2], // self (short assignment == declaration) in case clause
				names["i6"][3]: names["i6"][2], // use in case clause
				names["i6"][3]: names["i6"][2], // use in default clause
			},
			"i7": {
				names["i7"][0]: names["i7"][0], // self (short assignment == declaration) in loop init stmt
				names["i7"][1]: names["i7"][0], // use in loop body
				names["i7"][2]: names["i7"][2], // self (short assignment == declaration) in loop body
				names["i7"][3]: names["i7"][2], // use of redeclared var in loop body
			},
			"tmp7": {
				names["tmp7"][0]: names["tmp7"][0], // self (short assignment == declaration) in loop init stmt
				names["tmp7"][1]: names["tmp7"][0], // use in loop body
			},
			"tmptmp7": {
				names["tmptmp7"][0]: names["tmptmp7"][0], // self (short assignment == declaration) in loop body
				names["tmptmp7"][1]: names["tmptmp7"][0], // use in loop body
			},
			"tmp8": {
				names["tmp8"][0]: names["tmp8"][0], // self (short assignment == declaration) in function
				names["tmp8"][1]: names["tmp8"][1], // self (short assignment == declaration) in select case clause
				names["tmp8"][2]: names["tmp8"][1], // short assignment (== first use in select case clause)
				names["tmp8"][3]: names["tmp8"][1], // second use in select case clause
				names["tmp8"][4]: names["tmp8"][0], // use in function
			},
			"tmp9": {
				names["tmp9"][0]: names["tmp9"][0], // self (short assignment == declaration)
				names["tmp9"][1]: names["tmp9"][0], // use in return statement
			},
			"tmp10": {
				names["tmp10"][0]: names["tmp10"][0], // self (variable declaration)
				names["tmp10"][1]: names["tmp10"][0], // use in return statement
			},
			"i8": {
				names["i8"][0]: names["i8"][0], // self (short assignment == declaration) in select case clause
				names["i8"][1]: names["i8"][0], // use in select case clause
			},
			"tmp11": {
				names["tmp11"][0]: names["tmp11"][0], // self (short variable declaration in if initializer)
				names["tmp11"][1]: names["tmp11"][0], // use in if condition
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestAmbigSymbolication(t *testing.T) {
	t.Run("Test symbolication for various ambiguous use/def chains", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"tmp1": 5, "tmp2": 5, "tmp3": 4, "a1": 2}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/ambig_decl.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (top-level)
				names["tmp1"][1]: names["tmp1"][0], // use in func as func return type
				names["tmp1"][2]: names["tmp1"][2], // self (in-func)
				names["tmp1"][3]: names["tmp1"][0], // use in func to top-level
				names["tmp1"][4]: names["tmp1"][2], // use in func to in-func decl
			},
			"tmp2": {
				names["tmp2"][0]: names["tmp2"][0], // self (top-level)
				names["tmp2"][1]: names["tmp2"][0], // first use in func as func return type
				names["tmp2"][2]: names["tmp2"][2], // self (in-func)
				names["tmp2"][3]: names["tmp2"][0], // use in func to top-level
				names["tmp2"][4]: names["tmp2"][2], // use in func to in-func decl
			},
			"tmp3": {
				names["tmp3"][0]: names["tmp3"][0], // self (top-level)
				names["tmp3"][1]: names["tmp3"][1], // self (for loop init)
				names["tmp3"][2]: names["tmp3"][0], // use of top-level in for loop's rang
				names["tmp3"][3]: names["tmp3"][1], // use of for loop init decl in loop body
			},
			"a1": {
				names["a1"][0]: names["a1"][0], // self (top level field)
				names["a1"][1]: nil,            // second element of access path
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestAccessPathSymbolication(t *testing.T) {
	t.Run("Test symbolication for access path", func(t *testing.T) {
		// record the names with number suffix xand check against the number of appearances specified here
		expectedLength := map[string]int{
			"global1": 2, "a1": 2, "b1": 2, "tmp1": 5, "foo1": 3, "tmp2": 2,
			"bar1": 2, "c1": 2, "tmp3": 3, "foo3": 2, "bar3": 2, "foo4": 2,
			"bar4": 2, "c4": 2, "arr1": 2,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/access_path.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"global1": {
				names["global1"][0]: names["global1"][0], // self (global variable declaration)
				names["global1"][1]: names["global1"][0], // index of third element of access path in funcMidRetPath
			},
			"a1": {
				names["a1"][0]: names["a1"][0], // self (in-func)
				names["a1"][1]: names["a1"][0], // use in in-func (first element of access path)
			},
			"b1": {
				names["b1"][0]: names["b1"][0], // self (top-level field)
				names["b1"][1]: nil,            // second element of access path
			},
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (top-level field)
				names["tmp1"][1]: names["tmp1"][1], // self (in-func)
				names["tmp1"][2]: nil,              // second element of access path
				names["tmp1"][3]: nil,              // second element of access path (parenthesized)
				names["tmp1"][4]: names["tmp1"][1], // func use to in-func decl
			},
			"foo1": {
				names["foo1"][0]: names["foo1"][0], // self (top-level func)
				names["foo1"][1]: names["foo1"][0], // first use in in-func (first element of access path)
				names["foo1"][2]: names["foo1"][0], // second use in in-func (first element of access path)
			},
			"tmp2": {
				names["tmp2"][0]: names["tmp2"][0], // self (top-level field)
				names["tmp2"][1]: nil,              // third element of access path
			},
			"bar1": {
				names["bar1"][0]: nil, // self (method decl)
				names["bar1"][1]: nil, // second element of access path
			},
			"arr1": {
				names["arr1"][0]: names["arr1"][0], // self (field declaration)
				names["bar1"][1]: nil,              // third element of access path
			},
			"c1": {
				names["c1"][0]: names["c1"][0], // self (in-func)
				names["c1"][1]: names["c1"][0], // use in in-func (first element of access path)
			},
			"tmp3": {
				names["tmp3"][0]: names["tmp3"][0], // self (top-level field)
				names["tmp3"][1]: nil,              // second element of access path in funcRetPath
				names["tmp3"][2]: nil,              // forth element of access path in funcMidRetPath
			},
			"foo3": {
				names["foo3"][0]: names["foo3"][0], // self (top-level func)
				names["foo3"][1]: names["foo3"][0], // return value in bar3 func
			},
			"bar3": {
				names["bar3"][0]: names["bar3"][0], // self (top-level func)
				names["bar3"][1]: names["bar3"][0], // use in func (first element of access path)
			},
			"foo4": {
				names["foo4"][0]: names["foo4"][0], // self (top-level func)
				names["foo4"][1]: names["foo4"][0], // return value in bar4 func
			},
			"bar4": {
				names["bar4"][0]: nil, // self (method decl)
				names["bar4"][1]: nil, // use in func (second element of access path)
			},
			"c4": {
				names["c4"][0]: names["c4"][0], // self (in-func)
				names["c4"][1]: names["c4"][0], // use in in-func (first element of access path)
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestLabelSymbolication(t *testing.T) {
	t.Run("Test symbolication for labels", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"tmp1": 6, "tmp2": 6, "tmp3": 7}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/label.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (var decl in func)
				names["tmp1"][1]: names["tmp1"][1], // self (label decl in func)
				names["tmp1"][2]: names["tmp1"][0], // first use of var in func
				names["tmp1"][3]: names["tmp1"][0], // second use of var in func
				names["tmp1"][4]: names["tmp1"][1], // use of label in func
				names["tmp1"][5]: names["tmp1"][0], // third use of var in func
			},
			"tmp2": {
				names["tmp2"][0]: names["tmp2"][0], // self (var decl in func)
				names["tmp2"][1]: names["tmp2"][1], // self (label decl in func)
				names["tmp2"][2]: names["tmp2"][0], // first use of var in func
				names["tmp2"][3]: names["tmp2"][0], // second use of var in func
				names["tmp2"][4]: names["tmp2"][1], // use of label in func
				names["tmp2"][5]: names["tmp2"][0], // third use of var in func
			},
			"tmp3": {
				names["tmp3"][0]: names["tmp3"][0], // self (var decl in func)
				names["tmp3"][1]: names["tmp3"][1], // self (label decl in func)
				names["tmp3"][2]: names["tmp3"][0], // first use of var in func
				names["tmp3"][3]: names["tmp3"][0], // second use of var in func
				names["tmp3"][4]: names["tmp3"][0], // third use of var in func
				names["tmp3"][5]: names["tmp3"][1], // use of label in func
				names["tmp3"][6]: names["tmp3"][0], // third use of var in func
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestCallSymbolication(t *testing.T) {
	t.Run("Test symbolication for function and method calls", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"foo1": 4, "struct2": 3, "ptrfoo1": 3}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"go/call.go", ts.GoExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"foo1": {
				names["foo1"][0]: nil,              // self (method decl)
				names["foo1"][1]: names["foo1"][1], // self (func decl)
				names["foo1"][2]: names["foo1"][1], // func call
				names["foo1"][3]: nil,              // method call
			},
			"struct2": {
				names["struct2"][0]: names["struct2"][0], // self (struct decl)
				names["struct2"][1]: names["struct2"][0], // type conversion
				// TODO: enable this when symbolication for access paths is implemented
				// names["struct2"][2]: names["struct2"][0], // type conversion in access path
			},
			"ptrfoo1": {
				names["ptrfoo1"][0]: names["ptrfoo1"][0], // self (struct decl)
				// The function pointer variable declaration below intentionally
				// shadows the struct decl above.
				names["ptrfoo1"][1]: names["ptrfoo1"][1], // self (lambda func)
				names["ptrfoo1"][2]: names["ptrfoo1"][1], // lambda function call
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}
