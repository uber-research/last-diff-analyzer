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
  void root()
  {
    continue;
    continue here;
    break;
    break there;
    return;
    return a;
    a = b;
    while (a) foo();
    while (a) {
      foo();
    }
    switch (a) {
      case "1":
        foo();
        bar();
      case 2:
        test();
      case 3: case 4: case 5: foo();
      default:
      /* do nothing */
    }
    switch (a) {
      default:
        foo();
    }
    throw ex;
    assert a == b;
    assert a == b : "ERROR";
    synchronized(a) {a++;}
    if (a)
      foo();  // intentionally omit the braces to test translation of if statement without Block nodes.

    if (a) {
      t1();
    } else if (b) {
      t2();
    } else {
      t3();
    }
  hello:
    int a, b;

    do foo(); while(bar);
    do {
      foo();
    } while(bar);

    while (hasNext());

    try {
      t1();
    } catch (@Test @Test2 A|B ex[][]) {
      t2();
    } catch (@Test3 C ex2) {
      t3();
    } finally {
      t4();
    }

    // same try statement with resources
    try (file; Scanner scanner = new Scanner()) {
      t1();
    } catch (@Test @Test2 A|B ex[][]) {
      t2();
    } catch (@Test3 C ex2) {
      t3();
    } finally {
      t4();
    }

    try {
      t1();
    } catch(A a) {
      t2();
    }

    try {
      t1();
    } finally {
      t2();
    }

    for (;;) {
      foo();
    }
    for (;;) foo();
    for (int i = 1;;) {}
    for (i = 2; i < 10; i ++) {}
    for (i = 2, k = 3; i < 10; i++, k--) {}
    for (int i, k; i < 10; ) {}
    for (;i < 10; i ++) {}

    for (final @Test String a [][] : pkg.d) foo();
    for (int a : b) {}

    // method reference
    super::<A, B>someMethod;
    SomeClass::new;
    empty_label:; // an empty label without a statement
  }
}
