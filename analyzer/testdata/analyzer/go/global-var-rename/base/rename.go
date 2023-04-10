package rename

var intVar = 7

// renameA tests renaming of global (and optionally local) variables (variant a)
func renameA() int {
	intVar += 42
	return intVar + packageVarA
}

// renameB tests renaming of global (and optionally local) variables (variant b)
func renameB() int {
	intVar := 42
	return intVar + packageVarB
}

// renameC tests renaming of global (and optionally local) variables (variant c)
func renameC(b bool) int {
	if b {
		intVar = 42
		return intVar
	}
	return intVar
}

// renameD tests renaming of global (and optionally local) variables (variant d)
func renameD(b bool) int {
	if b {
		intVar := 42
		return intVar
	}
	return intVar
}
