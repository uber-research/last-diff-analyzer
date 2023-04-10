package rename

import "fmt"

// redefine tests renaming of a local variable that involves redefinition
func redefine(b bool) int {
	foo := 42

	if b {
		fooRenamed := 7
		fmt.Println(fooRenamed)
		return foo
	}

	return foo
}
