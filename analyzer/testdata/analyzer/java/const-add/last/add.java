public class add {

    private static final int c1 = 42;
    
    private final int c2 = 7;
    
    private final int c3 = c1 + 7;

    int addFieldA() {
        return c3;
    }
    
    int addFieldB() {
        int tmp1 = 0;
        int tmp2 = 0;
        return tmp1 + c1 + tmp2;
    }
    
    int addFieldC(boolean p) {
        int tmp1 = 0;
        int tmp2 = 0;
        if (p) {
            tmp1 = 1;
            tmp2 = c2;
        }
        return tmp1 + c1 + tmp2;
    }

    int addLocalA() {
        int tmp1 = 0;
        final int c4 = 42;
        int tmp2 = 0;
        return tmp1 + tmp2 + c4;
    }

    int addLocalB(boolean p) {        
        final int c4 = 42;
        int tmp1 = 0;
        int tmp2 = 0;
        if (p) {
            tmp1 = 7;
            final int c5 = 7;
            tmp2 = c5;
        }
        return tmp1 + tmp2 + c4;
    }

}
