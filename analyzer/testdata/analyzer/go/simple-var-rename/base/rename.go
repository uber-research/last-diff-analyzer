package rename

// renameA tests simple renaming of local variable defined separately
func renameA() int {
	var foo int
	foo = 42
	return foo
}

// renameB tests simple renaming of local variable defined in-place
func renameB() int {
	foo := 42
	return foo
}
