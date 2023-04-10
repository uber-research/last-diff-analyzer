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

@interface MyAnnotation {
    int value1() default 0;
}

@MyAnnotation(value1 = 42)
class C1 {
    public C1() {}

    int foo1() {
        return 7;
    }

    C1 foo2() {
        return new C1();
    }
    
    
    int bar() {
        C1 tmp = new C1();
        return this.foo1() + foo1() + tmp.foo1() + this.foo2().bar();
    }
    
    void tmp1(Object o) {
        System.out.println(o);
    }

    void baz()  {
        int tmp1 = 7;
        List<Integer> nums = Arrays.asList(tmp1, 42);
        nums.forEach(this::tmp1);
    }

}
