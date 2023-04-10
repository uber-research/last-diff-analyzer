package add

import "time"

// Limit is used to verify rejection of a change involving global constant.
const Limit = 10 * time.Second

// addGlobal tests adding globally defined exported constants to
// replace a literal use
func addGlobal() time.Duration {
	return Limit
}
