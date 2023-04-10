package format

// Foo is a dummy function
func Foo() bool {
	return true
}

// LineSplit tests a change where a line is split to two
func LineSplit(i int) int {
	if i == 7 || i == 42 {
		return 1
	}
	return 42
}
