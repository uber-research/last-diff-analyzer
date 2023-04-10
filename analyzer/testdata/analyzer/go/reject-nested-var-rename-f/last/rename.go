package rename

// redefine tests renaming of a local variable that involves redefinition
func redefine(b bool) int {
	fooRenamed := 42

	if b {
		foo := 7
		fooRenamed = 1
		return foo
	}

	return fooRenamed
}
