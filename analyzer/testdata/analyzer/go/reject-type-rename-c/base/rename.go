package rename

import "fmt"

// replace tests (incorrect) replacement of a type defined in a
// different file (instead of having the type renamed - at the level
// of this file's AST there is no difference)
func replace(i int) int {
	t := testA{i}
	fmt.Println(t)
	return i
}
