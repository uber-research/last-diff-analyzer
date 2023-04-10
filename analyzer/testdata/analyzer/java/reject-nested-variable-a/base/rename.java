package rename;

class Test {
  int method(boolean b) {
    int foo = 42;
    int bar = 7;
    int baz = 44;

    if (b) {
      foo = 1;
      bar = 0;
      return bar;
    }

    return foo + baz;
  }
}
