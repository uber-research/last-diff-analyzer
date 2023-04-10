package add

import "time"

const limit = 10 * time.Second

// addLocalA tests adding locally defined constants to replace a literal use (variant a)
func addLocalA() int {

	v1 := 0
	var v2 int

	const c = 42

	v3 := 0

	return v1 + v2 + v3 + c
}

// addLocalB tests adding locally defined constants to replace a literal use (variant b)
func addLocalB(b bool) int {

	const c1 = 42

	v1 := 0
	v2 := 0

	if b {
		v1 = 1
		const c2 = 7
		v2 = c2
	}

	return v1 + v2 + c1
}

// addGlobalA tests adding globally defined constants to replace a literal use (variant a)
func addGlobalA() time.Duration {
	return limit
}

// addGlobalB tests adding globally defined constants to replace a literal use (variant b)
func addGlobalB() int {

	v1 := 0
	v2 := 0

	return v1 + packageConst + v2
}

// addGlobalC tests adding globally defined constants to replace a literal use (variant c)
func addGlobalC(b bool) int {

	v1 := 0
	v2 := 0

	if b {
		v1 = 1
		v2 = packageConstA
	}

	return v1 + packageConstB + v2
}
