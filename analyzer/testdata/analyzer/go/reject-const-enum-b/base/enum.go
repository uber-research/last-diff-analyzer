package enum

// enum tests correct handling of package-level consts that have no (explicit) value (no auto-approval but also no crash)
func enum() int {
	return packageConstA + packageConstB
}
