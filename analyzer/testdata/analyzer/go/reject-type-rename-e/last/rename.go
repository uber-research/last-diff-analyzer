package rename

type test struct {
	varAccessedValueRenamed int
}

// rename tests (incorrect, at least for now) renaming of a field that
// is the used via a variable
func rename(i int) int {
	t := test{varAccessedValueRenamed: i}
	return t.varAccessedValueRenamed
}
