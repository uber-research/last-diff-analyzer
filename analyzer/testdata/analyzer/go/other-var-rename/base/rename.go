package rename

// dummy function
func dummy(i int) (int, error) {
	return i, nil
}

// renameA tests renaming of variable defined in a tuple.
func renameA(i int) int {
	foo, _ := dummy(i)
	return foo
}

// renameB tests renaming variable used as function parameter.
func renameB() int {
	foo := 42
	bar := renameA(foo)
	return bar
}

// renameC tests renaming of a variable in if statement.
func renameC() {
	if a := renameB(); a > 5 {
	}
}
