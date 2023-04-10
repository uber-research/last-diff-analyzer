package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

// rename tests renaming of the logging method (should reject as we
// are renaming to "incompatible" logging method)
func rename() {
	logger.Debug("message")
}
