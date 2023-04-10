package test;

private class FooRenamed {
  private String a;

  // Note: not a static class, so it's associated with an instance of Foo
  public class Bar {
    private String a; // unused.

    public void baz() {
      FooRenamed.this.a = "Set field of outer class";
    }
  }
}
