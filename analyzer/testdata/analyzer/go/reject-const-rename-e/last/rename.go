package rename

import "fmt"

// rename tests renaming that involves constants, vars, and shadowing
func rename(b bool) int {
	foo := 42

	if b {
		const fooRenamed = 7
		fmt.Println(fooRenamed)
		return foo
	}

	return foo
}
