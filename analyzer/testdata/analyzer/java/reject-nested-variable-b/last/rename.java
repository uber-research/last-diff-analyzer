package rename;

class Test {
  int method(boolean b) {
    int fooRenamed = 42;
    int barRenamed = 7;

    if (b) {
      fooRenamed = 1;
      fooRenamed = 0;
      return barRenamed;
    }

    return barRenamed;
  }
}
