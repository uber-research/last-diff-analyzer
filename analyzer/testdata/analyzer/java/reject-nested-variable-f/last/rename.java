package rename;

class Test {
  int method(boolean b) {
    int foo = 42;

    if (b) {
      int foo = 7;
      foo = 1;
      return foo;
    }

    return foo;
  }
}
