package rename

// replace tests (incorrect) replacement of constant defined in a
// different file (instead of having the constant renamed - at the
// level of this file's AST there is no difference)
func replace() int {
	return packageConstA
}
