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
// therefore any identifier with a number suffix should have a corresponding link

func TestSimpleJavaSymbolication(t *testing.T) {
	t.Run("Test creating symbol table for a simple Java program", func(t *testing.T) {

		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{
			"C1": 6, "C2": 2, "C3": 1, "f1": 2, "f2": 2, "m1": 2, "m2": 1, "p1": 3, "p2": 1, "l1": 4, "m3": 3,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/simple.java", ts.JavaExt, expectedLength)

		// now test the accuracies of the links
		// The expectedLinks maps from identifier name to a link map, which is from use identifier node to its expected
		// def identifier node. Most variable links are simple: there are only two identifiers in the test file, first
		// one is the declaration identifier and the second is the usage. So there should be two links, one from the
		// declaration to itself, and the second from the use to the declaration. The special cases will be accompanied
		// with detailed explanations.
		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			// All C1 refers to the first element, the class declaration.
			"C1": {
				names["C1"][0]: names["C1"][0], // self (top-level)
				names["C1"][1]: names["C1"][0], // first use in m2  (param)
				names["C1"][2]: names["C1"][0], // second use in m2 (local)
				names["C1"][3]: names["C1"][0], // second use in m2 (local)
				names["C1"][4]: names["C1"][0], // first use in m3 (return)
				names["C1"][5]: nil,            // second use in m3 (constructor)
			},
			"C2": {
				names["C2"][0]: names["C2"][0], // self (top-level)
				names["C2"][1]: names["C2"][0], // use in m2 (param)
			},
			"C3": {
				names["C2"][0]: names["C2"][0], // self (top-level)
			},
			"f1": {
				names["f1"][0]: names["f1"][0], // self (field in C3)
				names["f1"][1]: names["f1"][0], // use in m2
			},
			"f2": {
				names["f2"][0]: names["f2"][0], // self (field in C3)
				names["f2"][1]: names["f2"][0], // use in m2
			},
			"m1": {
				names["m1"][0]: nil, // self (method in C3)
				names["m1"][1]: nil, // use in m2
			},
			"m2": {
				names["m2"][0]: nil, // self (method in C3)
			},
			"p1": {
				names["p1"][0]: names["p1"][0], // self (param in m2)
				names["p1"][1]: names["p1"][0], // first use in m2
				names["p1"][2]: names["p1"][0], // second use in m2
			},
			"p2": {
				names["p2"][0]: names["p2"][0],
			},
			// l1 variable appears in the following locations:
			// (1) variable declaration in "C1 l1 = m3();";
			// (2) binary expression in "if (l1 == p1) {...}";
			// (3) shadowing variable declaration in "C1 l1 = m3();";
			// (4) binary expression in another "if (l1 == p1)";
			// So the links should be:
			// (1)(2) are linked to (1); (3)(4) are linked to (3).
			"l1": {
				names["l1"][0]: names["l1"][0], // self (local in m2)
				names["l1"][1]: names["l1"][0], // use in m2
				names["l1"][2]: names["l1"][2], // shadowed declaration in m2
				names["l1"][3]: names["l1"][2], // use of shadowed declaration in m2
			},
			"m3": {
				names["m3"][0]: nil, // first use before declaration in m2
				names["m3"][1]: nil, // second use before declaration in m2
				names["m3"][2]: nil, // self (method in C3)
			},
		}

		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestAccessPathJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for access path", func(t *testing.T) {
		// record the names with number suffix xand check against the number of appearances specified here
		expectedLength := map[string]int{"C1": 5, "a1": 4, "b1": 2, "arr1": 3}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/access_path.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"C1": {
				names["C1"][0]: names["C1"][0], // self (top-level)
				names["C1"][1]: names["C1"][0], // use in variable declaration
				names["C1"][2]: nil,            // use in foo (constructor)
				names["C1"][3]: names["C1"][0], // first use in bar (local var type)
				names["C1"][4]: nil,            // second use in bar (constructor)
			},
			"arr1": {
				names["arr1"][0]: names["arr1"][0], // self (variable declaration)
				names["arr1"][1]: names["arr1"][0], // use in bar
				names["arr1"][2]: nil,              // second use in bar (second element of access path)
			},
			"a1": {
				names["a1"][0]: names["a1"][0], // self (field in C1)
				names["a1"][1]: nil,            // use in foo (second element of access path)
				names["a1"][2]: names["a1"][0], // use in bar (index of second element of access path)
				names["a1"][3]: nil,            // use in bar (second element of access path)
			},
			"b1": {
				names["b1"][0]: names["b1"][0], // self in bar (local var decl)
				names["b1"][1]: names["b1"][0], // first use in bar (first element of access path)
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestLabelJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for labels", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"tmp1": 6, "tmp2": 6}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/label.java", ts.JavaExt, expectedLength)

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
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestCallJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for method calls", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"C1": 6, "foo1": 4, "value1": 2, "tmp1": 4, "foo2": 2}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/call.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"C1": {
				names["C1"][0]: names["C1"][0], // self (top-level class decl)
				names["C1"][1]: nil,            // constructor declaration
				names["C1"][2]: names["C1"][0], // method declaration (return value)
				names["C1"][3]: nil,            // entity creation expression
				names["C1"][4]: names["C1"][0], // first use in bar (local var type)
				names["C1"][5]: nil,            // second use in bar (constructor invocation)
			},
			"foo1": {
				names["foo1"][0]: nil, // self (method decl)
				names["foo1"][1]: nil, // first use in bar (call qualified with "this")
				names["foo1"][2]: nil, // second use in bar (unqualified call)
				names["foo1"][3]: nil, // third use in bar (call qualified with local var name)
			},
			"foo2": {
				names["foo2"][0]: nil, // self (method declaration)
				names["foo2"][1]: nil, // first use in bar (call qualified with "this")
			},
			"value1": {
				names["value1"][0]: nil, // self (method in annotation type declaration)
				names["value1"][1]: nil, // use in annotation (left-hand side)
			},
			"tmp1": {
				names["tmp1"][0]: nil,              // self (method decl)
				names["tmp1"][1]: names["tmp1"][1], // self (local var decl in baz)
				names["tmp1"][2]: names["tmp1"][1], // first use in baz (local var)
				names["tmp1"][3]: nil,              // second use in baz (call in method reference)
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestPkgJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for packages", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{
			"test1": 3, "C0": 3, "C1": 2, "anotherpkg1": 1,
			"yetanotherpkg1": 1, "C2": 1, "onemorepkg1": 1,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/package.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"test1": {
				names["test1"][0]: nil, // self (pkg decl)
				names["test1"][1]: nil, // id used as pkg name
				names["test1"][2]: nil, // id used as implicitly imported static field
			},
			"C0": {
				names["C0"][0]: nil,            // self (part if pkg import decl)
				names["C0"][1]: names["C0"][1], // self (class decl)
				names["C0"][2]: nil,            // use in constructor invocation
			},
			"C1": {
				names["C1"][0]: nil, // self (part if pkg import decl)
				names["C1"][0]: nil, // use as class name in implicitly imported method call
			},
			"anotherpkg1": {
				names["anotherpkg1"][0]: nil, // self (part if pkg import decl)
			},
			"yetanotherpkg1": {
				names["yetanotherpkg1"][0]: nil, // self (part if pkg import decl)
			},
			"C2": {
				names["C2"][0]: nil, // self (part if pkg import decl)
			},
			"onemorepkg1": {
				names["onemorepkg1"][0]: nil, // self (part if pkg import decl)
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestPkgMoreJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for packages on more examples", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"test1": 2}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/package_more.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"test1": {
				names["test1"][0]: nil,               // self (pkg decl)
				names["test1"][1]: names["test1"][1], // self (class decl)
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestAmbigJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for ambiguous cases", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"tmp1": 4, "tmp2": 5, "tmp3": 4, "tmp4": 3, "tmp5": 4,
			"tmp6": 3, "tmp7": 3, "tmp8": 4, "tmp9": 4, "tmp10": 3, "tmp11": 4, "tmp12": 3,
			"tmp13": 4, "tmp14": 5, "tmp15": 3, "tmp16": 3, "tmp17": 5, "tmp18": 3, "tmp19": 3, "tmp20": 3,
			"member1": 4,
		}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/ambig.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"tmp1": {
				names["tmp1"][0]: names["tmp1"][0], // self (field decl)
				names["tmp1"][1]: nil,              // self (name of method decl)
				names["tmp1"][2]: nil,              // first use in foo (method call)
				names["tmp1"][3]: names["tmp1"][0], // second use in foo (field access)
			},
			"tmp2": {
				names["tmp2"][0]: names["tmp2"][0], // self (tmp2 class decl)
				names["tmp2"][1]: names["tmp2"][1], // self (field decl name in class tmp2)
				names["tmp2"][2]: names["tmp2"][1], // use in foo (field access)
				names["tmp2"][3]: names["tmp2"][0], // first use in bar (class name to access "class" field)
				names["tmp2"][4]: names["tmp2"][0], // second use in bar (class name to access C class's "class" field)
			},
			"tmp3": {
				names["tmp3"][0]: names["tmp3"][0], // self (tmp3 class decl)
				names["tmp3"][1]: names["tmp3"][1], // self (field decl name in class tmp3)
				names["tmp3"][2]: names["tmp3"][0], // first use in foo (variable declaration)
				names["tmp3"][3]: nil,              // second use in foo (object construction)
			},
			"tmp4": {
				names["tmp4"][0]: names["tmp4"][0], // self (tmp4 class decl)
				names["tmp4"][1]: names["tmp4"][1], // self (field decl name in class tmp4)
				names["tmp4"][2]: names["tmp4"][0], // use as foo's formal parameter
			},
			"tmp5": {
				names["tmp5"][0]: names["tmp5"][0], // self (tmp5 class decl)
				names["tmp5"][1]: names["tmp5"][1], // self (field decl name in class tmp5)
				names["tmp5"][2]: names["tmp5"][0], // use in foo's throws clause
				names["tmp5"][3]: names["tmp5"][1], // use in foo (field access)
			},
			"tmp6": {
				names["tmp6"][0]: names["tmp6"][0], // self (tmp6 class decl)
				names["tmp6"][1]: names["tmp6"][1], // self (field decl name in class tmp6)
				names["tmp6"][2]: names["tmp6"][0], // use as array type in foo
			},
			"tmp7": {
				names["tmp7"][0]: names["tmp7"][0], // self (tmp7 class decl)
				names["tmp7"][1]: names["tmp7"][1], // self (field decl name in class tmp7)
				names["tmp7"][2]: names["tmp7"][0], // use as foo's "this" parameter
			},
			"tmp8": {
				names["tmp8"][0]: names["tmp8"][0], // self (tmp8 class decl)
				names["tmp8"][1]: names["tmp8"][1], // self (field decl name in class tmp8)
				names["tmp8"][2]: names["tmp8"][0], // use in field decl type in class tmp8
				names["tmp8"][3]: names["tmp8"][1], // use in foo (field access)
			},
			"tmp9": {
				names["tmp9"][0]: names["tmp9"][0], // self (tmp9 interface decl)
				names["tmp9"][1]: names["tmp9"][1], // self (field decl name in class tmp10)
				names["tmp9"][2]: names["tmp9"][0], // use in implements clause of class C
				names["tmp9"][2]: names["tmp9"][0], // use in implements clause of interface I
			},
			"tmp10": {
				names["tmp10"][0]: names["tmp10"][0], // self (tmp10 class decl)
				names["tmp10"][1]: names["tmp10"][1], // self (field decl name in class tmp10)
				names["tmp10"][2]: names["tmp10"][0], // use in extends clause of class C
			},
			"tmp11": {
				names["tmp11"][0]: names["tmp11"][0], // self (tmp11 class decl)
				names["tmp11"][1]: nil,               // self (declared method name in class tmp11)
				names["tmp11"][2]: names["tmp11"][0], // use in as method return type in class tmp11
				names["tmp11"][3]: nil,               // use in foo (method call)
			},
			"tmp12": {
				names["tmp12"][0]: names["tmp12"][0], // self (tmp12 class decl)
				names["tmp12"][1]: names["tmp12"][1], // self (field decl name in class tmp12)
				names["tmp12"][2]: names["tmp12"][0], // use in instance of clause
			},
			"tmp13": {
				names["tmp13"][0]: names["tmp13"][0], // self (tmp13 class decl)
				names["tmp13"][1]: names["tmp13"][1], // self (field decl name in class tmp13)
				names["tmp13"][2]: nil,               // use in constructor invocation
				names["tmp13"][3]: names["tmp13"][0], // use in foo's catch clause
			},
			"tmp14": {
				names["tmp14"][0]: names["tmp14"][0], // self (tmp14 class decl)
				names["tmp14"][1]: names["tmp14"][1], // self (field decl name in class tmp14)
				names["tmp14"][2]: names["tmp14"][0], // use in foo's return type
				names["tmp14"][3]: names["tmp14"][0], // use in instance of clause
				names["tmp14"][4]: names["tmp14"][0], // use in cast expression
			},
			"tmp15": {
				names["tmp15"][0]: names["tmp15"][0], // self (tmp15 class decl)
				names["tmp15"][1]: names["tmp15"][1], // self (field decl name in class tmp15)
				names["tmp15"][2]: names["tmp15"][0], // use as type parameter in call to foo
			},
			"tmp16": {
				names["tmp16"][0]: names["tmp16"][0], // self (tmp16 class decl)
				names["tmp16"][1]: names["tmp16"][1], // self (field decl name in class tmp16)
				names["tmp16"][2]: names["tmp16"][0], // use as type parameter in object construction
			},
			"tmp17": {
				names["tmp17"][0]: names["tmp17"][0], // self (tmp17 class decl)
				names["tmp17"][1]: names["tmp17"][1], // self (field decl name in class tmp17)
				names["tmp17"][2]: nil,               // first constructor name
				names["tmp17"][3]: nil,               // second constructor name
				names["tmp17"][4]: names["tmp17"][0], // use as type parameter to second constructor's invocation
			},
			"tmp18": {
				names["tmp18"][0]: names["tmp18"][0], // self (tmp18 class decl)
				names["tmp18"][1]: names["tmp18"][1], // self (field decl name in class tmp18)
				names["tmp18"][2]: names["tmp18"][0], // use in C.foo to call outer class's foo method
			},
			"tmp19": {
				names["tmp19"][0]: names["tmp19"][0], // self (tmp19 inner class decl)
				names["tmp19"][1]: names["tmp19"][1], // self (field decl name in class tmp19)
				names["tmp19"][2]: names["tmp19"][0], // use in tmp19.foo to call super class's foo method
			},
			"member1": {
				names["member1"][0]: names["member1"][0], // self (private member declaration)
				names["member1"][1]: names["member1"][1], // self (parameter declaration)
				names["member1"][2]: names["member1"][1], // use without "this" to refer to parameter declaration
				names["member1"][3]: names["member1"][0], // use with "this" to refer to the member declaration
			},
			"tmp20": {
				names["tmp20"][0]: names["tmp20"][0], // self (tmp20 inner class decl)
				names["tmp20"][1]: names["tmp20"][1], // self (field decl name in class tmp20)
				names["tmp20"][2]: names["tmp20"][0], // use in tmp20.bar to access super class's foo method reference
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}

func TestTypeJavaSymbolication(t *testing.T) {
	t.Run("Test symbolication for type names", func(t *testing.T) {
		// record the names with number suffix and check against the number of appearances specified here
		expectedLength := map[string]int{"C1": 2, "C2": 2}
		symbolTable, names := getNames(t, _metaTestDataPrefix+"java/types.java", ts.JavaExt, expectedLength)

		expectedLinks := map[string]map[*mast.Identifier]*mast.Identifier{
			"C1": {
				names["C1"][0]: names["C1"][0], // self (class decl)
				names["C1"][1]: names["C1"][0], // use in the wild card
			},
			"C2": {
				names["C2"][0]: names["C2"][0], // self (class decl)
				names["C2"][1]: names["C2"][0], // use in the generic type as the first element of access path
			},
		}
		verifySymbols(t, symbolTable, expectedLinks)
	})
}
