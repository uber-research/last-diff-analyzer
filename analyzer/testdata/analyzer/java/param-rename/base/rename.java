package rename;

private class Test {
  int method(boolean p, int i) {
    if (p) {
      int i = 1;
      return i;
    }
    return i;
  }

  int method2(boolean p, int i) {
    if (p) {
      int i = 1;
      return i;
    }
    return i;
  }

  int method3(boolean p, int i) {
    Function f = ((int i) -> i);
    return i + f(7);
  }

  int method4(boolean p, int i) {
    Function f = ((int i) -> i);
    return i + f(7);
  }
}
