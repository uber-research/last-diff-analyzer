package remove

// removeLocalA tests removing of locally defined constants (variant a)
func removeLocalA() int {

	v1 := 0
	var v2 int
	v3 := 0

	return v1 + v2 + v3
}

// removeLocalB tests removing of locally defined constants (variant b)
func removeLocalB(b bool) int {

	v := 0

	if b {
		v = 1
	}

	return v
}
