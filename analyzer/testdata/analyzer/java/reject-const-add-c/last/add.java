public class add {

    private static final add c = new add();

    // addField tests (incorrect) replacement of non-constant
    // expression with the final (but not constant) field
    add addField() {
        return c;
    }
}
