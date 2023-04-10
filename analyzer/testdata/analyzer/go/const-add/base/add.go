package add

import "time"

// addLocalA tests adding locally defined constants to replace a literal use (variant a)
func addLocalA() int {

	v1 := 0
	var v2 int
	v3 := 0

	return v1 + v2 + v3 + 42
}

// addLocalB tests adding locally defined constants to replace a literal use (variant b)
func addLocalB(b bool) int {

	v1 := 0
	v2 := 0

	if b {
		v1 = 1
		v2 = 7
	}

	return v1 + v2 + 42
}

// addGlobalA tests adding globally defined constants to replace a literal use (variant a)
func addGlobalA() time.Duration {
	return 10 * time.Second
}

// addGlobalB tests adding globally defined constants to replace a literal use (variant b)
func addGlobalB() int {

	v1 := 0
	v2 := 0

	return v1 + 42 + v2
}

// addGlobalC tests adding globally defined constants to replace a literal use (variant c)
func addGlobalC(b bool) int {

	v1 := 0
	v2 := 0

	if b {
		v1 = 1
		v2 = 7
	}

	return v1 + 42 + v2
}
