package rename

import "fmt"

var intVarRenamed = 7

// rename tests renaming of global and local variables
func rename(b bool) int {
	if b {
		intVar := 42
		fmt.Println(intVar)
		return intVarRenamed
	}
	return intVarRenamed
}
