package rename;

import System.out.println;

class Test {
  int redefine(boolean b) {
    int fooRenamed = 42;

    if (b) {
      int foo = 7;
      System.out.println(fooRenamed);
      return foo;
    }

    return fooRenamed;
  }
}
