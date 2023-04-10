package rename

import "fmt"

// function tests renaming of a local variable that involves nested functions
func function() int {
	foo := 42

	f := func() int {
		fooRenamed := 7
		fmt.Println(foo)
		return fooRenamed
	}

	return f() + foo
}
