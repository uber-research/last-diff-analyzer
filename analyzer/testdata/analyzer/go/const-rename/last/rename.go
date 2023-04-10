package rename

const intConstRenamed = 7

const intCompConstRenamed = 42 + intConstRenamed

// renameGlobal tests renaming of global variables
func renameGlobal() int {
	return intConstRenamed + intCompConstRenamed + packageConstARenamed
}

// renameLocalA tests renaming of local variables (variant a)
func renameLocalA(b bool) int {
	const fooRenamed = 42

	if b {
		const foo = 7
		return foo
	}

	return fooRenamed
}

// renameLocalB tests renaming of local variables (variant b)
func renameLocalB(b bool) int {
	const foo = 42

	if b {
		const fooRenamed = 7
		return fooRenamed
	}

	return foo
}

// renameLocalC tests renaming of local variables (variant c)
func renameLocalC(b bool) int {
	const fooRenamed = 42

	if b {
		const fooRenamedInner = 7
		return fooRenamedInner
	}

	return fooRenamed
}

// renameWithVarA tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarA(b bool) int {
	fooRenamed := 42

	if b {
		const foo = 7
		return foo
	}

	return fooRenamed
}

// renameWithVarB tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarB(b bool) int {
	foo := 42

	if b {
		const fooRenamed = 7
		return fooRenamed
	}

	return foo
}

// renameWithVarC tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarC(b bool) int {
	fooRenamed := 42

	if b {
		const fooRenamedInner = 7
		return fooRenamedInner
	}

	return fooRenamed
}
