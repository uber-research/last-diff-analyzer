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

@TestMarker package dummy;

@TestSingle(true)package dummy;

@TestMulti(a=1,b=2,c={1,2,3,@Nested(d=1),@NestedMarker})package dummy;

package dummy;

package a.b.c;

import java.util.jar.*;
import java;

@Test
@Test2
@pkg.Test3
open module test.a.b{requires transitive a;exports b;opens c to d,f;uses g;provides h with i;}

module mod{

}

@Test
class Test<A, @Test2 B extends C & D> extends E implements F, G.H {
  @Test3
  public int a = 1, b = 2;
  private String c[];
  static int d;

  // static initializer
  static {
    foo(d);
  }

  // initializer
  {
    foo();
  }

  // constructor 1
  @Test
  public <A, B> Test(int a) throws Ex {
    <A, B>super(1);
  }

  // constructor 2
  Test() {
    pkg.A.<A, B>super();
  }

  // constructor 3
  Test(String b) {
    this();
  }

  @Test public <A, B> @Test @Test2 Tuple<A, B> hello(String a) [][] throws C, D  {
    foo();
   };

  void hello2();
}

@Test
interface Test<A, @Test2 B extends C & D> extends E, F.G {
  public int a = 1;

  @Test public <A, B> @Test @Test2 Tuple<A, B> hello(String a) [][] throws C, D;

  void hello2();
}

interface Test2 {
}

enum Enum {
  ONE, TWO, THREE
}

@Test
enum T1 implements I1, I2 {
  @Test
  V1(1) {
    int a, b;
    void m1() {
    }
  },
  V2, V3;

  public void m2() {
  }

  public void m3() {
  }
}

@Test
@interface Anno {
  // annotation element declaration
  @Test2
  int A()

  default 1;

  // default annotation element declaration
  int B()[][];

  // constant declaration
  int C = 1;
}
