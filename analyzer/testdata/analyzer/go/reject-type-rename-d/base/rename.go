package rename

import "fmt"

// replace tests (incorrect) replacement of a struct field defined in
// a different file (instead of having the struct field renamed - at
// the level of this file's AST there is no difference)
func replace(i int) int {
	t := testC{bar: i}
	fmt.Println(t)
	return t.bar
}
