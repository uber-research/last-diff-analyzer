package rename;

import System.out.println;

class Test {
  int redefine(boolean b) {
    int foo = 42;

    if (b) {
      int foo = 7;
      System.out.println(foo);
      return foo;
    }

    return foo;
  }
}
