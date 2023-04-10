package rename

import (
	fmt_import_renamed "fmt"
	my_fmt_import_renamed "fmt"
)

// replace tests (incorrect) replacement of the package alias (instead
// of the package alias being renamed)
func replace(i int) int {
	my_fmt_import_renamed.Println(i)
	fmt_import_renamed.Println(i)
	return i
}
