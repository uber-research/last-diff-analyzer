import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        // Change string argument.
        logger.debug("Test1");
        // Change log level.
        logger.error("Test1");
        logger.trace("Test1");
        // Removal of safe helper function.
        logger.debug("3 + 4");
        // Removal of a logger statement.
        // logger.debug("Removal");
        // Addition of a logger statement.
        logger.debug("Addition");
        // Binary and unary expressions.
        int x = 1;
        boolean y = false;
        logger.debug(String.format("%d %b", 2 + x, y));
        // Literals.
        logger.debug(String.format("%h %c %f %d", x, 'b', 2.2, 2));
    }

    private void test2() {
        // Call allowed method.
        logger.debug(String.format("%d", 1));
        // Call allowed method with different side-effect-free arguments.
        logger.debug(String.format("%d %d", 3, 4));
        // Call allowed method, nested.
        logger.info(String.format("%s", String.format("%s", String.format("%d", 2))));
        // Call methods with side-effect, but with same arguments.
        logger.info(Baz.sideEffect(1));
        // Nested call.
        logger.info(String.format("%d %d", Baz.sideEffect(1), 1));
    }
}
