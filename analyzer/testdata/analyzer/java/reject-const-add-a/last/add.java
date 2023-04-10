public class add {

    private static final int c1 = 42;
    
    private final int c3 = c1 + 8;

    // addField tests adding globally defined constants to replace a
    // literal use, but with a small literal value change
    int addField() {
        return c3;
    }
}
