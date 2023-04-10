package rename

// ExportedVarRenamed tests exported global var renaming
var ExportedVarRenamed = 7

// rename tests renaming of global variables
func rename() int {
	ExportedVarRenamed += 42
	return ExportedVarRenamed
}
