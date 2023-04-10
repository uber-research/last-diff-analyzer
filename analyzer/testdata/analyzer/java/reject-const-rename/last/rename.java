public class rename {

    static final int c1Renamed = 42;
    
    private final int c2Renamed = c1Renamed + 7;
    
    // renameField tests (incorrect) renaming of non-private constant
    int renameField() {
        return c1Renamed + c2Renamed;
    }
}
