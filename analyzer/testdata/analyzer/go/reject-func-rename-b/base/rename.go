package rename

// replace tests (incorrect) replacement of a call to function defined
// in a different file (instead of having the function being called
// renamed - at the level of this file's AST there is no difference)
func replace() int {
	return baz()
}
