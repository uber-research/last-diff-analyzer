package local

// This tests an upper-case prefixed local variable renaming

func Test() int {
	LocalRenamed := 1
	var AnotherLocalRenamed int
	return LocalRenamed + AnotherLocalRenamed
}
