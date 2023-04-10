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

package test;

class C1 {
}

class C2 {
}

class C3 {
  public int f1, f2;

  public void m1() {
  }

  public void m2(C1 p1, C2 p2) {
    m1();
    if (f1 == f2) {}
    C1 l1 = m3(); // used before declaration, but should be ok
    if (l1 == p1) {
      // test shadowing of variable declaration
      C1 l1 = m3();
      if (l1 == p1) {}
    }
  }

  public C1 m3() {
    return new C1();
  }
}
