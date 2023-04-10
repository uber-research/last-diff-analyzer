public class add {

    // Test adding a final modifier.
    private static final int c1 = 42;

    // Test reordering the modifiers.
    private final static int c2 = 42;

    final static int c3;
    final int c4;
    final Object c5;
    final volatile @Nullable T c6;

    int addFieldA() {
        return c1 + 7;
    }

    int addFieldB() {
        int tmp1 = 0;
        int tmp2 = 0;
        return tmp1 + 42 + tmp2;
    }

    int addFieldC(boolean p) {
        int tmp1 = 0;
        int tmp2 = 0;
        if (p) {
            tmp1 = 1;
            tmp2 = 7;
        }
        return tmp1 + 42 + tmp2;
    }

    int addLocalA() {
        int tmp1 = 0;
        int tmp2 = 0;
        return tmp1 + tmp2 + 42;
    }

    int addLocalB(boolean p) {
        int tmp1 = 0;
        int tmp2 = 0;
        if (p) {
            tmp1 = 7;
            tmp2 = 7;
        }
        return tmp1 + tmp2 + 42;
    }

}
