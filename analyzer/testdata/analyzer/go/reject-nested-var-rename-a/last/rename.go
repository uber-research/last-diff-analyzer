package rename

// rename tests renaming of local variables that involves nesting
func rename(b bool) int {
	foo := 42
	bar := 7
	bazRenamed := 44

	if b {
		foo = 1
		bazRenamed = 0
		return bar
	}

	return foo + bazRenamed
}
