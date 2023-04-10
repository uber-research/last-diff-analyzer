import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        // Call functions with side effects in unary expression.
        logger.debug(String.format("%d", !sideEffect(2));
    }
}
