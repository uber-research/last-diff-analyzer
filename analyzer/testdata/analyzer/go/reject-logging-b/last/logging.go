package logging

import (
	"zap"

	origzap "go.uber.org/zap"
)

var logger = origzap.NewExample()

// argMod tests change of an argument to the logging method that
// syntactically looks correct but should reject as zap is a custom
// package
func argMod() {
	logger.Debug("message", zap.Int("key", 7))
}
