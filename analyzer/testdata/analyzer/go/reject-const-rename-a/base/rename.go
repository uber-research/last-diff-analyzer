package rename

// ExportedConst tests exported const renaming
const ExportedConst = 42

// rename tests (incorrect) renaming of exported constant
func rename() int {
	return ExportedConst
}
