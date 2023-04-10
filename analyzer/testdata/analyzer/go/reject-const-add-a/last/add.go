package add

import "time"

const limit = 11 * time.Second

// addGlobal tests adding globally defined constants to replace a
// literal use, but with a small literal value change
func addGlobal() time.Duration {
	return limit
}
