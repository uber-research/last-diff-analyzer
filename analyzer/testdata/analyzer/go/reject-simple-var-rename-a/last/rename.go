package rename

// rename tests incorrect simple renaming of local variable
func rename() int {
	var fooARenamed int
	var fooBRenamed int

	fooBRenamed = 42
	fooARenamed = 7

	return fooARenamed + fooBRenamed
}
