package rename

// rename tests incorrect simple renaming of local variable
func rename() int {
	var fooARenamed int
	var fooBRenamed int

	fooARenamed = 42
	fooBRenamed = 7

	return fooBRenamed - fooARenamed
}
