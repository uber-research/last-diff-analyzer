package rename

import "fmt"

// rename tests renaming of local variables
func rename(b bool) int {
	const fooRenamed = 42

	if b {
		const foo = 7
		fmt.Println(foo)
		return fooRenamed
	}

	return fooRenamed
}
