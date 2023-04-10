package rename

import "fmt"

// renameA tests renaming of local variables that involves nesting (variant a)
func renameA(b bool) int {
	foo := 42
	bar := 7

	if b {
		foo = 1
		bar = 0
		fmt.Println(bar)
	} else {
		foo = 2
	}

	return foo
}

// renameB tests renaming of local variables that involves nesting (variant b)
func renameB(b bool) int {
	foo := 42
	bar := 7

	if b {
		foo = 1
		bar = 0
		return bar
	}

	return foo
}

// renameC tests renaming of local variables that involves nesting (variant c)
func renameC(b bool) int {
	foo := 42
	bar := 7
	baz := 44

	if b {
		foo = 1
		bar = 0
		return bar
	}

	return foo + baz
}

// redefineA tests renaming of a local variable that involves redefinition (variant a)
func redefineA(b bool) int {
	foo := 42

	if b {
		foo := 7
		return foo
	}

	return foo
}

// redefineB tests renaming of a local variable that involves redefinition (variant b)
func redefineB(b bool) int {
	foo := 42

	if b {
		foo := 7
		fmt.Println(foo)
		return foo
	}

	return foo
}

// redefineC tests renaming of a local variable that involves redefinition (variant c)
func redefineC(b bool) int {
	foo := 42

	if b {
		foo := 7
		foo = 1
		return foo
	}

	return foo
}

// functionA tests renaming of a local variable that involves nested functions (variant a)
func functionA() int {
	foo := 42

	f := func() int {
		foo := 7
		fmt.Println(foo)
		return foo
	}

	return f() + foo
}

// functionB tests renaming of a local variable that involves nested functions (variant b)
func functionB() int {
	foo := 42

	f := func() int {
		foo := 7
		fmt.Println(foo)
		return foo
	}

	return f() + foo
}
