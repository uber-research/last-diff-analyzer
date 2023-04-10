package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

func foo(s string) int {
	return len(s)
}

// argSideEffect tests the case when non-equivalent arguments to the
// logging method involve potential side-effects (to be rejected)
func argSideEffect() {
	logger.Error("error", zap.Int("key3", foo("42")), zap.String("key2", "error"), zap.Int("key1", foo("7")))
}
