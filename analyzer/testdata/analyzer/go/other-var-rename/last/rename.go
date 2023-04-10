package rename

// dummy function
func dummy(i int) (int, error) {
	return i, nil
}

// renameA tests renaming of variable defined in a tuple.
func renameA(i int) int {
	fooRenamed, _ := dummy(i)
	return fooRenamed
}

// renameB tests renaming variable used as function parameter.
func renameB() int {
	fooRenamed := 42
	bar := renameA(fooRenamed)
	return bar
}

// renameC tests renaming of a variable in if statement.
func renameC() {
	if aRenamed := renameB(); aRenamed > 5 {
	}
}
