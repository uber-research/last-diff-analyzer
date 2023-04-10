package rename

import "fmt"

type test struct {
	intValue int
}

// param tests parameter renaming
func param(i int) (test, error) {
	return test{
		intValue: i,
	}, nil
}

// paramRedefineA tests parameter redefinition (variant A)
func paramRedefineA(p bool, i int) int {
	if p {
		i := 42
		return i
	}
	return i
}

// paramRedefineB tests parameter redefinition (variant A)
func paramRedefineB(p bool, i int) int {
	if p {
		i := 42
		return i
	}
	return i
}

// paramInnerA tests parameter renaming in nested function (variant a)
func paramInnerA(i int) int {

	f := func(i int) int {
		fmt.Println(i)
		return i
	}

	return i + f(7)
}

// paramInnerB tests parameter renaming in nested function (variant b)
func paramInnerB(i int) int {

	f := func(i int) int {
		fmt.Println(i)
		return i
	}

	return i + f(7)
}
