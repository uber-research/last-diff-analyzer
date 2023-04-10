import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        // Change string argument.
        logger.debug(String.format("%d", sideEffect(1));
    }
}
