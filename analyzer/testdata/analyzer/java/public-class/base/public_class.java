package test;

// This tests renaming _inside_ a public class

public class SomeThing {
  private int m;

  private int method() {
    return this.m;
  }
}
