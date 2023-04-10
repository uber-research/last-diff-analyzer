package rename;

private class Test1Renamed {
}

private class Test2Renamed {

}

class Rename {
  // test field renaming
  private int mRenamed = 1;
  private int nRenamed = 2;

  // test parameter renaming
  private int someMethod(int bRenamed) {
    // test local variable renaming
    int aRenamed = 1;
    Test1Renamed t1;
    Test2Renamed t2;
    return aRenamed + bRenamed;
  }

  private int someOtherMethod() {
    // test local variable renaming
    int aRenamed2 = 1;
    Test3Renamed t3;

    // "this" keyword can actually be omitted, so we test the four possibilities
    // here.
    this.mRenamed = this.someMethod(aRenamed2);
    nRenamed = this.someMethod(aRenamed2);
    this.nRenamed = someMethod(aRenamed2);
    nRenamed = someMethod(aRenamed2);

    int nRenamed2 = 3; // shadowing the field declaration
    nRenamed2 = this.someMethod(aRenamed2);
  }
}

private class Test3Renamed {
}
