package rename

const intConst = 7

const intCompConst = 42 + intConst

// renameGlobal tests renaming of global variables
func renameGlobal() int {
	return intConst + intCompConst + packageConstA
}

// renameLocalA tests renaming of local variables (variant a)
func renameLocalA(b bool) int {
	const foo = 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}

// renameLocalB tests renaming of local variables (variant b)
func renameLocalB(b bool) int {
	const foo = 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}

// renameLocalC tests renaming of local variables (variant c)
func renameLocalC(b bool) int {
	const foo = 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}

// renameWithVarA tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarA(b bool) int {
	foo := 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}

// renameWithVarB tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarB(b bool) int {
	foo := 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}

// renameWithVarC tests renaming that involves constants, vars, and shadowing (variant a)
func renameWithVarC(b bool) int {
	foo := 42

	if b {
		const foo = 7
		return foo
	}

	return foo
}
