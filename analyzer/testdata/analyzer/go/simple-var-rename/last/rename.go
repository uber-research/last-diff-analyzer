package rename

// renameA tests simple renaming of local variable defined separately
func renameA() int {
	var fooRenamed int
	fooRenamed = 42
	return fooRenamed
}

// renameB tests simple renaming of local variable defined in-place
func renameB() int {
	fooRenamed := 42
	return fooRenamed
}
