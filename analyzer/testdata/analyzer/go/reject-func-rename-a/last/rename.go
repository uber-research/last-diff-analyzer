package rename

// ExportedFooRenamed tests exported function renaming
func ExportedFooRenamed() int {
	return 0
}

// rename tests (incorrect) renaming of exported function
func rename() int {
	return ExportedFooRenamed()
}
