package rename

import "fmt"

// ExportedTest tests exported function renaming
type ExportedTest struct {
	intValue int
}

// rename tests (incorrect) renaming of exported type
func rename(i int) int {
	fmt.Println(ExportedTest{})
	return i
}
