package rename

import "fmt"

// rename tests parameter renaming in nested function
func rename(iRenamed int) int {

	f := func(i int) int {
		fmt.Println(iRenamed)
		return i
	}

	return iRenamed + f(7)
}
