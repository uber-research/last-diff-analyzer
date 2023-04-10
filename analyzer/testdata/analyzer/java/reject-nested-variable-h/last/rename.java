package rename;

import System.out.println;

class Test {
  static int f() {
    return 1;
  }

  int method(boolean b) {
    int foo = 42;

    if (b) {
      int fooRenamed = 7;
      System.out.println(foo);
      return fooRenamed;
    }

    return Test.f() + foo;
  }
}
