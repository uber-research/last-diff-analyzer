public class remove {

    int v1 = 0;

    private static final int c1 = 42;

    int v2 = 1;
    
    private final int c2 = 7;

    int v3 = 2;

    private final int c3 = c1 + 7;

    int removeLocalA() {
        int tmp1 = 0;
        final int c4 = 42;
        int tmp2 = 0;
        return tmp1 + tmp2;
    }
    
    int removeLocalB(boolean p) {
        final int c4 = 42;
        int tmp1 = 0;
        if (p) {
            int tmp2 = 1;
            final int c5 = 7;
        }        
        return tmp1;
    }
}
