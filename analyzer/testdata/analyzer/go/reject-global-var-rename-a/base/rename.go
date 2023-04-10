package rename

// ExportedVar tests exported global var renaming
var ExportedVar = 7

// rename tests renaming of global variables
func rename() int {
	ExportedVar += 42
	return ExportedVar
}
