package enum

const (
	constA = iota + 1
	constBRenamed
)

// enum tests correct handling of file-level consts that have no (explicit) value (no auto-approval but also no crash)
func enum() int {
	return constA + constBRenamed
}
