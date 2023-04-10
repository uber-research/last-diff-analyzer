package rename

func rename() int {
	return fooRenamed() + bazRenamed()
}

// declare out of order to make sure that it works for non-linear AST
// (where a function can be declaredas part of another function
// declaration)
func fooRenamed() int {
	return 7
}

// verifies that renaming a function does not cause a problem with param renaming
func paramRenamed(i int) int {
	return i + 42
}

// verifies that renaming a function does not cause a problem with variable renaming
func variableRenamed() int {
	res := 42
	return res
}
