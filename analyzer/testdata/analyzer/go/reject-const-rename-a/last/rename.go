package rename

// ExportedConstRenamed tests exported const renaming
const ExportedConstRenamed = 42

// rename tests (incorrect) renaming of exported constant
func rename() int {
	return ExportedConstRenamed
}
