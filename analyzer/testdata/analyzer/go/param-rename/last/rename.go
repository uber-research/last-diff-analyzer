package rename

import "fmt"

type test struct {
	intValue int
}

// param tests parameter renaming
func param(iRenamed int) (test, error) {
	return test{
		intValue: iRenamed,
	}, nil
}

// paramRedefineA tests parameter redefinition (variant A)
func paramRedefineA(p bool, iRenamed int) int {
	if p {
		i := 42
		return i
	}
	return iRenamed
}

// paramRedefineB tests parameter redefinition (variant A)
func paramRedefineB(p bool, i int) int {
	if p {
		iRenamed := 42
		return iRenamed
	}
	return i
}

// paramInnerA tests parameter renaming in nested function (variant a)
func paramInnerA(iRenamed int) int {

	f := func(i int) int {
		fmt.Println(i)
		return i
	}

	return iRenamed + f(7)
}

// paramInnerB tests parameter renaming in nested function (variant b)
func paramInnerB(i int) int {

	f := func(iRenamed int) int {
		fmt.Println(iRenamed)
		return iRenamed
	}

	return i + f(7)
}
