package reject

// Foo is a dummy function
func Foo() bool {
	return true
}

// Neq tests a non semantically equivalent change
func Neq(i int) int {
	if i == 7 {
		return 1
	}
	return 0
}
