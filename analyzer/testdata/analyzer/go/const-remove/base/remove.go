package remove

// removeLocalA tests removing of locally defined constants (variant a)
func removeLocalA() int {

	v1 := 0
	var v2 int

	const c = 42

	v3 := 0

	return v1 + v2 + v3
}

// removeLocalB tests removing of locally defined constants (variant b)
func removeLocalB(b bool) int {

	const c1 = 42

	v := 0

	if b {
		v = 1
		const c2 = 7
	}

	return v
}
