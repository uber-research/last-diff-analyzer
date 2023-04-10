package logging

import "go.uber.org/zap"

var logger = zap.NewExample()

func getLogger() *zap.Logger {
	return logger
}

// callRemove tests a removal of the logging method with access path containing a function call
func callRemove(s2 string) string {
	const s1 = "42"
	return s1
}
