package nouse

var tmp = 7

// noUse tests the case when variable is used in one diff and unused in the other
func noUse() string {
	const s1 = "42"
	return s1
}
