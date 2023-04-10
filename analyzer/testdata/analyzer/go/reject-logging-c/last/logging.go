package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

// modDeref tests changing an argument to the logging method to
// dereference a pointer (should be rejected)
func modDeref(s *string) {
	logger.Error("error", zap.Any("key", *s))
}
