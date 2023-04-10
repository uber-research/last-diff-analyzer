package enum

const (
	constA = iota + 1
	constB
)

// enum tests correct handling of file-level consts that have no (explicit) value (no auto-approval but also no crash)
func enum() int {
	return constA + constB
}
