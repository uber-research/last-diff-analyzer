package rename;

private class Test {
  int method(boolean p, int iRenamed) {
    if (p) {
      int i = 1;
      return i;
    }
    return iRenamed;
  }

  int method2(boolean p, int i) {
    if (p) {
      int iRenamed = 1;
      return iRenamed;
    }
    return i;
  }

  int method3(boolean p, int iRenamed) {
    Function f = ((int i) -> i);
    return iRenamed + f(7);
  }

  int method4(boolean p, int i) {
    Function f = ((int iRenamed) -> iRenamed);
    return i + f(7);
  }
}
