import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class Foo {
    private static final Logger logger = LoggerFactory.getLogger(Foo.class);

    private void test() {
        int i = 1;
        logger.debug(String.format("%d", ++i);
    }
}
