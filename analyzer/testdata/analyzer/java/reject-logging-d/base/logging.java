import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        // Call functions with side-effects in nested.
        logger.debug(String.format("%d %s", sideEffect(1), String.format("%d", sideEffect(1))));
    }
}
