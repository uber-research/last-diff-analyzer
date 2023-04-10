package test;

// rename a default-visibility declaration
class TestRenamed {
}

private class Test2 {
  void method() {
    TestRenamed t = new Test();
  }
}
