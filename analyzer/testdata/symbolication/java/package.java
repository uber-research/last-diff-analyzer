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

package test1; // package name identifier should be linked to null

// all identifiers used for package imports should be linked to null
import C0.C1;
import anotherpkg1.*;
import static yetanotherpkg1.C2;
import static onemorepkg1.*;

class C0 {

  public int foo() {
      // bar is a method in class pkg1.C1 (not defined here) - C1
      // should be linked to null as it comes from import declaration
      return C1.bar();
  }

  public static int baz() {
      // since we do not symbolicate package names, test1 should be
      // linked to null rather than the current package declaration
      return (new test1.C0()).foo();
  }

  public static int bak() {
      // test1 can be implicitly imported static field from
      // onemorepkg1 and should not be linked with the current package
      // declaration
      return test1;
  }
}

