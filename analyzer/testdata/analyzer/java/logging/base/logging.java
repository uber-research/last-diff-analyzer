import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        // Change string argument.
        logger.debug("Test1");
        // Change log level.
        logger.debug("Test1");
        logger.fatal("Test1");
        // Removal of safe helper function.
        logger.debug(String.format("%d %d", 1, 2));
        // Removal of a logger statement.
        logger.debug("Removal");
        // Addition of a logger statement.
        // logger.debug("Addition");
        // Binary and unary expressions.
        int x = 1;
        boolean y = false;
        logger.debug(String.format("%d %b", x + 1, !y));
        // Literals.
        logger.debug(String.format("%h %c %f %d", null, 'a', 1.1, 1));
    }

    private void test2() {
        // Call allowed method.
        logger.debug(String.format("%d", 1));
        // Call allowed method with different side-effect-free arguments.
        logger.debug(String.format("%d %d", 1, 2));
        // Call allowed method, nested.
        logger.info(String.format("%s", String.format("%s", String.format("%d", 1))));
        // Call methods with side-effect, but with same arguments.
        logger.info(Baz.sideEffect(1));
        // Nested call.
        logger.info(String.format("%d", Baz.sideEffect(1)));
    }
}
