package rename

import "fmt"

// renameA tests renaming of local variables that involves nesting (variant a)
func renameA(b bool) int {
	fooRenamed := 42
	bar := 7

	if b {
		fooRenamed = 1
		bar = 0
		fmt.Println(bar)
	} else {
		fooRenamed = 2
	}

	return fooRenamed
}

// renameB tests renaming of local variables that involves nesting (variant b)
func renameB(b bool) int {
	foo := 42
	barRenamed := 7

	if b {
		foo = 1
		barRenamed = 0
		return barRenamed
	}

	return foo
}

// renameC tests renaming of local variables that involves nesting (variant c)
func renameC(b bool) int {
	foo := 42
	bar := 7
	bazRenamed := 44

	if b {
		foo = 1
		bar = 0
		return bar
	}

	return foo + bazRenamed
}

// redefineA tests renaming of a local variable that involves redefinition (variant a)
func redefineA(b bool) int {
	fooRenamed := 42

	if b {
		foo := 7
		return foo
	}

	return fooRenamed
}

// redefineB tests renaming of a local variable that involves redefinition (variant b)
func redefineB(b bool) int {
	foo := 42

	if b {
		fooRenamed := 7
		fmt.Println(fooRenamed)
		return fooRenamed
	}

	return foo
}

// redefineC tests renaming of a local variable that involves redefinition (variant c)
func redefineC(b bool) int {
	foo := 42

	if b {
		fooRenamed := 7
		fooRenamed = 1
		return fooRenamed
	}

	return foo
}

// functionA tests renaming of a local variable that involves nested functions (variant a)
func functionA() int {
	fooRenamed := 42

	f := func() int {
		foo := 7
		fmt.Println(foo)
		return foo
	}

	return f() + fooRenamed
}

// functionB tests renaming of a local variable that involves nested functions (variant b)
func functionB() int {
	foo := 42

	f := func() int {
		fooRenamed := 7
		fmt.Println(fooRenamed)
		return fooRenamed
	}

	return f() + foo
}
