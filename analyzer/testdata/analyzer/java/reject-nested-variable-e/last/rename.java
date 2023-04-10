package rename;

class Test {
  int method(boolean b) {
    int foo = 42;

    if (b) {
      int fooRenamed = 7;
      foo = 1;
      return fooRenamed;
    }

    return foo;
  }
}
