package rename

import "fmt"

// rename tests renaming that involves constants, vars, and shadowing
func rename(b bool) int {
	fooRenamed := 42

	if b {
		const foo = 7
		fmt.Println(foo)
		return fooRenamed
	}

	return fooRenamed
}
