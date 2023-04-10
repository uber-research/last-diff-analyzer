package rename

import "fmt"

// rename tests renaming of local variables
func rename(b bool) int {
	const foo = 42

	if b {
		const fooRenamed = 7
		fmt.Println(fooRenamed)
		return foo
	}

	return foo
}
