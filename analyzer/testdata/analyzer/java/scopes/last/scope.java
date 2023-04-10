package hello.world.scope;

// test renaming of private field inside public class and public field inside private class.

public class Test {

  private int aRenamed;

  private void anotherUse() {
    int aRenamed = 1;

    if (aRenamed > 5) {
      int bRenamed = 2;
    }
    this.use();
  }

  private void use() {
    this.aRenamed = 1;
  }

}

// explicit private class
private class Test2Renamed {
  public int aRenamed2;

  private void use() {
    this.aRenamed2 = 1;
  }
}

// package level class
class Test3 {
  private int a;

  private void use() {
    this.a = 1;
  }
}
