package test;

// rename a default-visibility declaration
class Test {
}

private class Test2 {
  void method() {
    Test t = new Test();
  }
}
