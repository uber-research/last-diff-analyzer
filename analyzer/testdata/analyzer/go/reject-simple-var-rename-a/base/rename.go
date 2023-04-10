package rename

// rename tests incorrect simple renaming of local variable
func rename() int {
	var fooA int
	var fooB int

	fooA = 42
	fooB = 7

	return fooA + fooB
}
