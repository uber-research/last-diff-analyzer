package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

func foo(s string) int {
	return len(s)
}

// argRemove tests a removal of the argument to the logging method
// that involves potential side-effects (to be rejected)
func argRemove(s1 string, s2 string) {
	logger.Error("error", zap.Int("key1", len(s1)), zap.Int("key2", foo(s2)))
}
