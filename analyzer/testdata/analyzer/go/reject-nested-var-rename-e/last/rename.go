package rename

// redefine tests renaming of a local variable that involves redefinition
func redefine(b bool) int {
	foo := 42

	if b {
		fooRenamed := 7
		foo = 1
		return fooRenamed
	}

	return foo
}
