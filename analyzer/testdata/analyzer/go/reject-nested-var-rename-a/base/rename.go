package rename

// rename tests renaming of local variables that involves nesting
func rename(b bool) int {
	foo := 42
	bar := 7
	baz := 44

	if b {
		foo = 1
		bar = 0
		return bar
	}

	return foo + baz
}
