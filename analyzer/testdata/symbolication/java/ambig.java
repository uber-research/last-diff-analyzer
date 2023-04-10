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

import java.util.List;
import java.util.Arrays;

class C {    
  static int tmp1 = 7;
    
  static double tmp1() {
      return 42;
  }
  
  static void foo() {
    System.out.println(tmp1() + tmp1);
  }
}

class tmp2 {

    static int tmp2 = 42;

    static void foo() {
        System.out.println(tmp2);
    }

    class C {
    }

    private static void bar() {
        System.out.println(tmp2.class);
        System.out.println(tmp2.C.class);
    }

}

class tmp3 {

    static int tmp3 = 42;

    static void foo() {
        tmp3 tmp = new tmp3();
        System.out.println(tmp);
    }
}

class tmp4 {

    static int tmp4 = 42;

    static void foo(tmp4 tmp) {
        System.out.println(tmp);
    }
}

class tmp5 extends Exception {

    static int tmp5 = 42;

    static void foo() throws tmp5 {
        System.out.println(tmp5);
    }
}

class tmp6 {

    static int tmp6 = 42;

    static void foo() {
        tmp6[] tmp;
        System.out.println(tmp);
    }
}

class tmp7 {

    static int tmp7 = 42;

    static void foo(tmp7 this) {
        System.out.println(this);
    }
}

class tmp8 {

    static tmp8 tmp8;

    static void foo() {
        System.out.println(tmp8);
    }
}

interface tmp9 {
}

class tmp10 {

    static int tmp9;
    static int tmp10;

    class C extends tmp10 implements tmp9 {
    }

    interface I extends tmp9 {
    }
}

class tmp11 {

    static tmp11 tmp11() {
        return null;
    }

    static void foo() {
        System.out.println(tmp11());
    }
}

class tmp12 {

    static int tmp12 = 42;

    static void foo(Object o) {
        if (o instanceof tmp12) {
            System.out.println("HelloWorld!");
        }
    }
}

class tmp13 extends Exception {

    static int tmp13 = 42;

    static void foo() {
        try {
            throw new tmp13();
        } catch (tmp13 ex) {
            System.out.println("HelloWorld!");
        }
    }
}

class tmp14 {

    static int tmp14 = 42;

    static tmp14 foo(Object o) {
        if (o instanceof tmp14) {
            return (tmp14)o;
        }
        return null;
    }
}

class tmp15 {

    class C {
        static <T> String bar(T o) {
            return o.toString();
        }
    }

    static int tmp15 = 42;

    static String foo() {
        return C.<tmp15>bar("HelloWorld!");
    }
}

class tmp16 {

    class C<T> {
    }

    static int tmp16 = 42;

    static void foo() {
        System.out.println(new C<tmp16>());
    }
}

class tmp17 {

    static int tmp17 = 42;

    public <T> tmp17(T o) {
    }

    public tmp17() {
        <tmp17>this(null);
    }
}

public class tmp18 {

    static int tmp18 = 42;

    int foo() {
        return 7;
    }

    class C {
        int foo() {
            return tmp18.this.foo();
        }
    }
}

public class AnotherC {

    int foo() {
        return 42;
    }

    class tmp19 extends AnotherC {

        static int tmp19 = 42;

        int foo() {
            return tmp19.super.foo();
        }
    }
}

public class YetAnotherC
{
    private int member1;
    public YetAnotherC(int member1) {
        // Method parameter shadows a private member, but we should still be
        // able to access this member1 via "this".
        int k = member1;  // refers to the parameter declaration
        this.member1 = 1; // refers to the member declaration
    }

    void foo(Object o) {
        System.out.println(o);
    }

    class tmp20 extends YetAnotherC {

        static int tmp20 = 42;

        void foo(Object o) {
            System.out.print(o);
        }

        void bar() {
            List<Integer> nums = Arrays.asList(7, 42);
            nums.forEach(tmp20.super::foo);
        }
    }
}
