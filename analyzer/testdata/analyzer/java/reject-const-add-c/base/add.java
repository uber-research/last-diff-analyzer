public class add {

    // addField tests (incorrect) replacement of non-constant
    // expression with the final (but not constant) field
    add addField() {
        return new add();
    }
}
