package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

// callAdd tests an addition of the logging method whose argument
// involves potential side-effects (to be rejected)
func callAdd(s2 *string) string {
	const s1 = "42"
	return s1
}
