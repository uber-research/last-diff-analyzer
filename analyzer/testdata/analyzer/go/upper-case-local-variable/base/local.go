package local

// This tests an upper-case prefixed local variable renaming

func Test() int {
	Local := 1
	var AnotherLocal int
	return Local + AnotherLocal
}
