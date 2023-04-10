package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

func foo(s string) int {
	return len(s)
}

// callRemove tests a removal of the logging method whose argument
// involves potential side-effects (to be rejected)
func callRemove(s2 string) string {
	const s1 = "42"
	logger.Error("error", zap.Int("key1", len(s1)), zap.Int("key2", foo(s2)))
	return s1
}
