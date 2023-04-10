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

class Root {
  void root () {
    null;
    true;
    false;
    "test\t";
    123;
    0x5;
    0o5;
    0b11;
    1.5;
    0x0.C90FDAP2f;
    'a';
    a*b + c;
    !a;
    !(a);
    a[i];
    a.b.c;
    a.foo().b;
    foo().b;
    foo().a.b;
    a ? b : c;
    add(a, b);
    foo.bar(a, b);
    this.foo(a);
    super.bar(a);
    a.b.c.foo(a);
    (void)a;
    (int)a;
    (float)a;
    (a.b.c.Map<@Even int, boolean>)a;
    (Test<@NotNull ? extends A>)a;
    (Test<@Foo @Bar ? super @Foo A>)a;
    (Test<?>)a;
    (@NotNull String [] @Foo [] @Bar @Test [])a;
    (T1 & T2)a;
    a instanceof B;
    a++;
    a--;
    ++a;
    --a;
    () -> {
      foo();
    }
    () -> foo();
    (@T1 @T2 ClassName name. this) -> a;
    (@T1 @T2 ClassName name, Test1 t1, final Test2 t2[]@Anno[], Test3 ...a) -> a;
    (a, b) -> a;
    int a;
    final int a;
    int a = 2, b, c = foo();
    String [][] a, b[];
    @Test1 @Test2 String a[], b[][];
    int[] a = new int[]{1, 2};
    int[] a = new int[5];
    int[][] a = new int[b][c];
    int[][] a = new int[][]{{1,2},{3,4}};
    String a = new String("test");
    foo(String[].class); // the argument will be translated to JavaClassLiteral.
    foo(String.class); // due to quirks from tree-sitter, the argument will be translated to AccessPath.
    new ConcurrentHashMap<
    MethodDeclaration,
    Annotations.@NotNull Measure
    >();
    outer.new T();
    T a = outer.new T() {
      public void hello() {
        return;
      }
    };
    Clazz.<String>foo();
    Clazz.<String, Integer>foo();
    Clazz.<String, Integer>foo(a, b);
    a.this.bar;
    super.bar;
    a.super.bar;
    a.b.super.bar;
    a.this.bar();
    super.bar();
    a.super.bar();
    a.super.bar(a, b);
  }
}
