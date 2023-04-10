package rename

var intVarRenamed = 7

// renameA tests renaming of global (and optionally local) variables (variant a)
func renameA() int {
	intVarRenamed += 42
	return intVarRenamed + packageVarARenamed
}

// renameB tests renaming of global (and optionally local) variables (variant b)
func renameB() int {
	intVar := 42
	return intVar + packageVarBRenamed
}

// renameC tests renaming of global (and optionally local) variables (variant c)
func renameC(b bool) int {
	if b {
		intVarRenamed = 42
		return intVarRenamed
	}
	return intVarRenamed
}

// renameD tests renaming of global (and optionally local) variables (variant d)
func renameD(b bool) int {
	if b {
		intVar := 42
		return intVar
	}
	return intVarRenamed
}
