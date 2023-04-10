package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

// argAdd tests an addition of the argument to the logging method that
// involves potential side-effects (to be rejected)
func argAdd(s1 string, s2 *string) {
	logger.Error("error", zap.Int("key1", len(s1)))
}
