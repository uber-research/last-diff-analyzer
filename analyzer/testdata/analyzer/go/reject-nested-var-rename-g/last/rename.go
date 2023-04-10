package rename

import "fmt"

// function tests renaming of a local variable that involves nested functions
func function() int {
	fooRenamed := 42

	f := func() int {
		foo := 7
		fmt.Println(fooRenamed)
		return foo
	}

	return f() + fooRenamed
}
