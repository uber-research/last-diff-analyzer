package rename

import "fmt"

type exportedField struct {
	ExportedValue int
}

// rename tests (incorrect) renaming of exported struct field
func rename(i int) int {
	fmt.Println(exportedField{ExportedValue: i})
	return i
}
