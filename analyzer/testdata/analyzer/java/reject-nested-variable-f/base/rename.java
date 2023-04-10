package rename;

class Test {
  int method(boolean b) {
    int fooRenamed = 42;

    if (b) {
      int foo = 7;
      fooRenamed = 1;
      return foo;
    }

    return fooRenamed;
  }
}
