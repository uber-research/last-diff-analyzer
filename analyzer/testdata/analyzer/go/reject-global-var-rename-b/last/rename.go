package rename

// replace tests (incorrect) replacement of global var defined in a
// different file (instead of having the var renamed - at the level of
// this file's AST there is no difference)
func replace() int {
	return packageVarBRenamed
}
